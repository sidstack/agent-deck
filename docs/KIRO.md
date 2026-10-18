# Kiro CLI support in this fork

This fork adds Kiro as a built-in tool so its native terminal interface can use
Agent Deck's existing session dashboard and persistent tmux sessions.

```sh
agent-deck add /path/to/project -c kiro
agent-deck
```

Select **kiro** in the new-session tool picker, or create a session with the CLI.
Kiro is also available in first-run setup, Settings' default-tool selection,
and the tool-visibility checklist. Use Up/Down to reach entries below the fold.
Attach and detach through Agent Deck normally. Closing the dashboard leaves the
live Kiro process in its tmux session.

The default launch is:

```sh
kiro-cli chat --v2 --tui
```

Kiro runs in the project's working directory and loads its own profiles,
steering, MCP configuration, credentials, and permissions. This integration
does not rewrite Kiro configuration or add approval-bypass flags.

For an existing native profile, supply the full command:

```sh
agent-deck add /path/to/project -c 'kiro-cli chat --v2 --tui --agent reviewer'
```

Explicit commands keep their arguments. The `kiro` picker alias alone expands
to the V2/TUI default. Installation detection checks for `kiro-cli` on PATH;
the existing `show_only_installed_tools` and `hidden_tools` settings apply.
The name `kiro` is now reserved for this built-in; rename an existing
`[tools.kiro]` custom wrapper to another name to keep that wrapper.

## Support boundary

- Supported: built-in identity, icon, tool picker, CLI creation, remote creation
  catalog, launch in the selected directory, and existing native terminal
  attachment/persistence.
- This is live-process persistence. Restarting an exited default Kiro process
  starts a **fresh conversation**; this fork does not discover or resume Kiro
  conversation IDs, import history, or fork conversations.
- No Kiro hook installation or verified thinking/approval status patterns are
  included. Treat dashboard turn-state indicators as coarse existing behavior,
  not authoritative Kiro telemetry.
- No Kiro model/effort picker, MCP manager, account migration, or conductor
  integration is claimed. Explicit engine/profile options belong to the native
  command; V3 compatibility has not been validated by this patch.

## Validation

Tests run only in a disposable Docker container with an isolated home and tmux
server. A fake `kiro-cli` verifies the exact V2/TUI arguments and project cwd
through the real session launcher. Command detection covers quoted paths,
profile names containing other agents' names, and false-positive rejection.
The patch makes no authenticated Kiro calls. Real Kiro rendering and approvals
through the user's terminal remain a manual acceptance check.

To reproduce the affected-package checks using upstream's disposable containers:

```sh
scripts/verify-docker.sh --stages build,vet,test --pkgs './internal/session ./internal/tmux ./internal/ui ./cmd/agent-deck'
```
