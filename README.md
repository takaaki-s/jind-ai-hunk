# jind-ai-hunk

A [jin](https://github.com/takaaki-s/jind-ai) plugin that bridges an
interactive [hunk](https://github.com/modem-dev/hunk) review to the
target jin session's agent.

Two actions, deliberately separate:

- **`open`** — opens `hunk diff` in a new tmux window pointed at the
  session's workdir. If a hunk session is already running for that
  workdir, does nothing and points you at it.
- **`send`** — sends a short handoff message to the target jin session
  telling its agent to load hunk's own review skill and address your
  inline notes. Refuses if there is no live hunk session for that
  workdir yet.

Typical flow: run `open` → add notes with `c` in the hunk window → run
`send`. Notes live in hunk's local daemon and vanish when hunk exits,
so `send` before closing the hunk window.

## Requirements

- jin ≥ 0.8.0 (schema-v2 `actions[]`, `JIN_BIN`, `JIN_SESSION_ID`, `JIN_WORKDIR`).
- Go ≥ 1.24 to build the plugin binary.
- [hunk](https://github.com/modem-dev/hunk) on `PATH` (`npm i -g hunkdiff`
  or `brew install hunk`).
- `tmux` — `open` spawns a new tmux window.
- `git` on `PATH` (hunk needs it).

## Install

### From a git source

```
jin plugin install github.com/takaaki-s/jind-ai-hunk
```

`install.source.build` runs `go build -o bin/jind-ai-hunk ./cmd/jind-ai-hunk`
during install; that binary is what the `send` action invokes.

### Local development

```
git clone https://github.com/takaaki-s/jind-ai-hunk ~/dev/jind-ai-hunk
cd ~/dev/jind-ai-hunk
go build -o bin/jind-ai-hunk ./cmd/jind-ai-hunk     # --link installs skip build
jin plugin install --link ~/dev/jind-ai-hunk
```

Rebuild `bin/jind-ai-hunk` after any Go source change; the symlinked
plugin directory picks it up on the next invocation.

## Usage

Both actions surface in jin's action panel with the labels above. From
a shell, invoke them explicitly:

```
jin plugin run jind-ai-hunk open --session <session-id>
# ... add notes in the hunk window ...
jin plugin run jind-ai-hunk send --session <session-id>
```

`open` is the default — `jin plugin run jind-ai-hunk` runs it.

### Flags on `send`

The `send` action calls the Go binary, which accepts:

| flag | meaning |
|------|---------|
| `--session <id>` | target jin session; falls back to `JIN_SESSION_ID` |
| `--workdir <path>` | hunk session repo root; falls back to `JIN_WORKDIR`, then to `jin session info --json` on the target session |
| `--dry-run` | print the assembled message and exit — nothing is sent |

### Prompt language

The handoff message is written in English so the target agent sees a
consistent, model-friendly instruction. The agent's *response* language
still follows whatever the session config specifies — English
instructions do not force English replies.

## What this plugin does NOT do

- It does not read hunk's JSON directly. Deferring to `hunk skill path`
  keeps the plugin insensitive to schema changes and lets hunk own the
  navigate / list / comment vocabulary.
- It does not render a diff. Use `hunk diff` (or any `hunk show …`
  variant) for that.
- It does not persist notes. Hunk's local daemon holds your live notes;
  quit hunk only after you have run `send`.
- It does not resolve the caller's jin session automatically. Pass
  `--session` explicitly (or invoke from jin's action panel where the
  session context is injected).

## License

MIT — see [LICENSE](./LICENSE).
