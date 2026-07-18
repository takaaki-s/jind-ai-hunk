#!/usr/bin/env bash
# jind-ai-hunk — jin plugin entrypoint. Two verbs dispatched via
# manifest actions[]:
#
#   open  — Split the target jin session's tmux pane and run `hunk diff`
#           in the new half so the user can review side-by-side with the
#           agent. Fire-and-forget: this script exits immediately, hunk
#           stays running until the user quits it. If a hunk session for
#           this workdir is already live, prints a hint and exits without
#           spawning a duplicate.
#   send  — Send a short handoff message to the target jin session telling
#           its agent to load `hunk skill path` and address the user's
#           inline notes. Requires a live hunk session for the workdir.
#
# Notes live in hunk's local daemon and vanish when hunk exits, so send
# BEFORE closing the hunk window.

set -u

verb="${1:-}"
if [ $# -gt 0 ]; then shift; fi

if [ -z "${JIN_SESSION_ID:-}" ]; then
  echo "jind-ai-hunk: no session context. Invoke from jin's action panel or 'jin plugin run jind-ai-hunk <action> --session <id>'." >&2
  exit 1
fi

if ! command -v hunk >/dev/null 2>&1; then
  echo "jind-ai-hunk: 'hunk' not found on PATH. Install it (npm i -g hunkdiff or brew install hunk)." >&2
  exit 1
fi

WORKDIR="${JIN_WORKDIR:-}"
if [ -z "$WORKDIR" ]; then
  echo "jind-ai-hunk: JIN_WORKDIR not set — the plugin needs the target workdir." >&2
  exit 1
fi

JIN="${JIN_BIN:-jin}"
DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
BIN="$DIR/bin/jind-ai-hunk"

case "$verb" in
  open)
    if hunk session get --repo "$WORKDIR" >/dev/null 2>&1; then
      echo "jind-ai-hunk: hunk is already reviewing $WORKDIR — switch to that pane, or run the 'send' action when you are done." >&2
      exit 0
    fi
    # jin owns the tmux plumbing: `jin pane split <session-id>` splits the
    # session's own pane in its own tmux server, inheriting the session's
    # working directory (which is $JIN_WORKDIR). No $TMUX or
    # JIN_CALLER_TMUX_* handling needed here — jin routes to the right
    # server internally.
    exec "$JIN" pane split "$JIN_SESSION_ID" --direction right --size 50% --no-focus -- hunk diff
    ;;

  send)
    if [ ! -x "$BIN" ]; then
      echo "jind-ai-hunk: binary not found at $BIN. Build it: go build -o bin/jind-ai-hunk ./cmd/jind-ai-hunk" >&2
      exit 1
    fi
    if ! hunk session get --repo "$WORKDIR" >/dev/null 2>&1; then
      echo "jind-ai-hunk: no live hunk session for $WORKDIR — run the 'open' action first, add notes, then 'send'." >&2
      exit 1
    fi
    exec "$BIN" "$@"
    ;;

  "")
    echo "jind-ai-hunk: no action specified. Use 'open' to start a review or 'send' to hand it off." >&2
    exit 1
    ;;

  *)
    echo "jind-ai-hunk: unknown action '$verb' (expected 'open' or 'send')." >&2
    exit 1
    ;;
esac
