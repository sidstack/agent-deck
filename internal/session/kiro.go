package session

import "strings"

const defaultKiroCommand = "kiro-cli chat --v2 --tui"

// buildKiroCommand launches Kiro's native TUI in the session's working directory.
// Only the picker alias expands to our default. An explicit command keeps its
// engine and profile arguments; no resume or approval flags are inferred.
func (i *Instance) buildKiroCommand(baseCommand string) string {
	if cmd := strings.TrimSpace(baseCommand); cmd == "" || cmd == "kiro" {
		baseCommand = defaultKiroCommand
	}
	return i.buildEnvSourceCommand() + baseCommand
}
