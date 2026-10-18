package ui

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"al.essio.dev/pkg/shellescape"

	"github.com/asheshgoplani/agent-deck/internal/session"
)

func TestKiroCreateDialogLaunch(t *testing.T) {
	home := t.TempDir()
	project := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	session.ClearUserConfigCache()
	t.Cleanup(session.ClearUserConfigCache)
	bin := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(bin, 0o700); err != nil {
		t.Fatal(err)
	}
	receipt := filepath.Join(project, "kiro-launch.txt")
	fake := "#!/bin/sh\nprintf '%s\\n' \"$PWD\" \"$@\" > " + shellescape.Quote(receipt) + "\nexec sleep 60\n"
	if err := os.WriteFile(filepath.Join(bin, "kiro-cli"), []byte(fake), 0o700); err != nil {
		t.Fatal(err)
	}
	h := &Home{}
	msg := h.createSessionInGroupWithWorktreeAndOptions(
		"kiro-dialog-test", project, "", "kiro", "test",
		"", "", "", false, false, nil, nil, "", "", "", false, nil, "", "", "", false,
	)().(sessionCreatedMsg)
	if msg.err != nil {
		t.Fatal(msg.err)
	}
	inst := msg.instance
	t.Cleanup(func() { _ = inst.KillAndWait() })
	if inst.Tool != "kiro" {
		t.Fatalf("dialog created tool %q, want kiro", inst.Tool)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if data, err := os.ReadFile(receipt); err == nil {
			if want := project + "\nchat\n--v2\n--tui\n"; string(data) != want {
				t.Fatalf("Kiro launch = %q, want %q", data, want)
			}
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("dialog failed to launch native Kiro")
}

func TestKiroInNewSessionPicker(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	session.ClearUserConfigCache()
	t.Cleanup(session.ClearUserConfigCache)
	if !slices.Contains(buildPresetCommands(), "kiro") {
		t.Fatal("Kiro missing from new-session picker")
	}
}

func TestKiroVisibilitySetting(t *testing.T) {
	panel := NewToolVisibilityPanel()
	panel.LoadConfig(&session.UserConfig{
		Tools: map[string]session.ToolDef{"kiro": {}},
		UI:    session.UISettings{HiddenTools: []string{"kiro"}},
	})
	count := 0
	for index, name := range panel.toolNames {
		if name == "kiro" {
			count++
			panel.cursor = index
		}
	}
	if count != 1 || !slices.Contains(panel.HiddenTools(), "kiro") {
		t.Fatal("Kiro must have exactly one row that respects hidden_tools")
	}
	panel.visibleTools["kiro"] = true
	if slices.Contains(panel.HiddenTools(), "kiro") {
		t.Fatal("showing Kiro should remove it from hidden_tools")
	}
}
