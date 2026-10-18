package ui

import (
	"slices"
	"testing"

	"github.com/asheshgoplani/agent-deck/internal/session"
)

func TestKiroInNewSessionPicker(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	session.ClearUserConfigCache()
	t.Cleanup(session.ClearUserConfigCache)
	if !slices.Contains(buildPresetCommands(), "kiro") {
		t.Fatal("Kiro missing from new-session picker")
	}
}
