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
