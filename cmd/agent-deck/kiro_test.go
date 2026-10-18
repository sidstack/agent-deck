package main

import (
	"testing"

	"github.com/asheshgoplani/agent-deck/internal/session"
)

func TestKiroCreationCatalog(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("XDG_DATA_HOME", "")
	session.ClearUserConfigCache()
	t.Cleanup(session.ClearUserConfigCache)
	catalog, err := buildCreationCatalog("default")
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range catalog.Tools {
		if tool.Name == "kiro" {
			if tool.Kind != "kiro" || len(tool.Models) != 0 || len(tool.ReasoningEfforts) != 0 {
				t.Fatalf("Kiro must not inherit unsupported model/effort controls: %+v", tool)
			}
			return
		}
	}
	t.Fatal("Kiro missing from remote creation catalog")
}
