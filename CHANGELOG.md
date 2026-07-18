# Changelog

All notable changes to this project will be documented here. The format
loosely follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.0] — 2026-07-18

Initial release. A [jin](https://github.com/takaaki-s/jind-ai) plugin
that bridges an interactive [hunk](https://github.com/modem-dev/hunk)
review to the target jin session's agent.

### Added

Two manual actions on the `schema_version: 2` manifest, wired via
`actions[]`:

- `open` — spawns `hunk diff` in a new tmux window pointed at the jin
  session's workdir. If a hunk session is already registered for that
  workdir, no-op with a hint instead of duplicating.
- `send` — sends a short handoff message to the target jin session
  telling its agent to load `hunk skill path` and walk the user's
  inline notes (`hunk session comment list --repo <workdir> --type user`).
  Refuses if no live hunk session for that workdir is registered.

Other:

- `--dry-run` on the `send` binary prints the assembled message and
  exits without sending.
- Session workdir is auto-resolved via `jin session info --json` when
  `JIN_WORKDIR` is not injected by the dispatcher.
- Handoff prompt is written in English. Agent responses still follow
  the session's configured language.
- `jin plugin install --link` support via `install.source.build`.

### Requirements

- jin ≥ 0.8.0 (schema-v2 `actions[]`)
- Go ≥ 1.24 to build
- [hunk](https://github.com/modem-dev/hunk) (`npm i -g hunkdiff` or
  `brew install hunk`) on `PATH`
- `tmux` — the `open` action spawns a new tmux window
- `git` on `PATH` (hunk needs it)

[0.1.0]: https://github.com/takaaki-s/jind-ai-hunk/releases/tag/v0.1.0
