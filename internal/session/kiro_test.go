package session

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"al.essio.dev/pkg/shellescape"
)

func TestKiroIdentityAndLaunch(t *testing.T) {
	r := Init(nil)
	if !r.IsBuiltin("kiro") {
		t.Fatal("Kiro must be a built-in tool")
	}
	for _, cmd := range []string{"kiro", "kiro-cli chat --v2 --tui", "kiro-cli chat --v2 --tui --agent claude-reviewer", "kiro-cli chat --v2 --tui --agent codex-reviewer", "/opt/bin/kiro-cli chat --v3", `"/opt/AI Tools/kiro-cli" chat`} {
		if got := r.Match(cmd); got != "kiro" {
			t.Errorf("Match(%q) = %q", cmd, got)
		}
	}
	for _, cmd := range []string{"echo kiro-cli", "git -C /tmp/kiro status", "my-kiro-cli-wrapper", "cat /opt/bin/kiro-cli"} {
		if got := r.Match(cmd); got == "kiro" {
			t.Errorf("false match: %q", cmd)
		}
	}
	if got := GetToolCommand("kiro"); got != "kiro-cli chat --v2 --tui" {
		t.Errorf("command = %q", got)
	}
	if GetToolIcon("kiro") != r.Icon("kiro") || r.Icon("kiro") == "" {
		t.Error("Kiro icon missing or inconsistent")
	}
	if IsClaudeCompatible("kiro") || IsCodexCompatible("kiro") {
		t.Fatal("Kiro must not inherit another agent's permissions or resume protocol")
	}
	inst := &Instance{Tool: "kiro"}
	for _, cmd := range []string{"", "kiro"} {
		if got := inst.buildKiroCommand(cmd); !strings.HasSuffix(got, "kiro-cli chat --v2 --tui") {
			t.Errorf("launch %q = %q", cmd, got)
		}
	}
	explicit := "kiro-cli chat --v3 --agent reviewer"
	if got := inst.buildKiroCommand(explicit); !strings.HasSuffix(got, explicit) {
		t.Errorf("explicit command changed: %q", got)
	}
	if inst.CanRestartGeneric() {
		t.Error("Kiro has no verified generic resume support")
	}
}

func TestKiroInstalledProbeUsesNativeBinary(t *testing.T) {
	old := lookPathFn
	t.Cleanup(func() { lookPathFn = old })
	lookPathFn = func(cmd string) (string, error) {
		if cmd == "kiro-cli" {
			return "/fake/kiro-cli", nil
		}
		return "", os.ErrNotExist
	}
	r := InitFiltered(nil, true, nil)
	if !r.IsVisible("kiro") || r.FilterFallback() {
		t.Fatal("installed kiro-cli should show Kiro without enabling fallback")
	}
	lookPathFn = func(string) (string, error) { return "", os.ErrNotExist }
	if r = InitFiltered(nil, true, []string{"kiro"}); r.IsVisible("kiro") {
		t.Fatal("explicitly hidden Kiro should stay hidden")
	}
}

func TestKiroPickerHasBuiltin(t *testing.T) {
	if !slices.Contains(pickerPresetOrder, "kiro") {
		t.Fatal("Kiro missing from picker catalog")
	}
}

func TestKiroExitToShellOptIn(t *testing.T) {
	exitToShellTestEnv(t)
	inst := NewInstanceWithTool("kiro-exit", t.TempDir(), "kiro")
	raw := inst.buildKiroCommand("kiro")
	if got := inst.wrapExitToShell(raw); got != raw {
		t.Fatal("disabled exit_to_shell changed the Kiro command")
	}
	inst.ExitToShell = boolPtr(true)
	wrapped := inst.wrapExitToShell(raw)
	out := runPaneCommand(t, wrapped, inst.ProjectPath, exitToShellFakeEnv(t, "kiro-cli"))
	if !strings.Contains(out, "AGENT_RAN chat --v2 --tui") || !strings.Contains(out, "SHELL_FALLBACK -i") {
		t.Fatalf("Kiro did not launch and fall back to the requested shell: %s", out)
	}
}

func TestKiroNativeProcessStartsInProject(t *testing.T) {
	skipIfNoTmuxBinary(t)
	dir := t.TempDir()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	ClearUserConfigCache()
	t.Cleanup(ClearUserConfigCache)
	bin := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(bin, 0o700); err != nil {
		t.Fatal(err)
	}
	receipt := filepath.Join(dir, "launch.txt")
	fake := "#!/bin/sh\nprintf '%s\\n' \"$PWD\" \"$@\" > " + shellescape.Quote(receipt) + "\nexec sleep 60\n"
	if err := os.WriteFile(filepath.Join(bin, "kiro-cli"), []byte(fake), 0o700); err != nil {
		t.Fatal(err)
	}

	inst := NewInstanceWithTool("kiro-native-test", dir, "kiro")
	inst.Command = "kiro"
	t.Cleanup(func() { _ = inst.Kill() })
	if err := inst.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(receipt)
		if err == nil {
			want := dir + "\nchat\n--v2\n--tui\n"
			if string(data) != want {
				t.Fatalf("native launch = %q, want %q", data, want)
			}
			if !inst.tmuxSession.Exists() {
				t.Fatal("Kiro process did not remain in its persistent tmux session")
			}
			if strings.Contains(string(data), "trust") || strings.Contains(string(data), "resume") {
				t.Fatal("unexpected approval or resume flag")
			}
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("native Kiro launch did not write its receipt")
}
