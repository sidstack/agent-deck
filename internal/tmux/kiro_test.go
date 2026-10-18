package tmux

import "testing"

func TestDetectToolFromCommandKiro(t *testing.T) {
	for _, cmd := range []string{"kiro-cli chat --v2 --tui", "/opt/bin/kiro-cli chat --v3", `"/opt/AI Tools/kiro-cli" chat`, "KIRO-CLI chat"} {
		if got := detectToolFromCommand(cmd); got != "kiro" {
			t.Errorf("detect(%q) = %q", cmd, got)
		}
	}
	for _, cmd := range []string{"echo kiro-cli", "cat /opt/bin/kiro-cli", "git -C /tmp/kiro status", "my-kiro-cli-wrapper"} {
		if got := detectToolFromCommand(cmd); got == "kiro" {
			t.Errorf("false positive: %q", cmd)
		}
	}
	if DefaultRawPatterns("kiro") != nil {
		t.Fatal("Kiro turn patterns have not been verified")
	}
}
