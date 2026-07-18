// Package prompt renders the review-handoff message that this plugin
// sends to `jin session send`. The message tells the target agent to
// load hunk's own review skill and pick up the user's inline notes.
//
// The plugin does not parse hunk's JSON on its own: the SKILL.md at
// `hunk skill path` already teaches the agent the full navigate /
// list / act flow. Keeping this file as a template lets us tune the
// wording without touching the CLI layer.
//
// The instruction is written in English on purpose — LLM responses
// still follow whatever language the session config specifies, and
// English lands more reliably across models.
package prompt

import (
	"fmt"
	"strings"
)

// Build returns the message body sent to the jin session. workdir is
// the repo root the hunk session is loaded against; it is embedded so
// the agent can address the right session without extra discovery.
func Build(workdir string) string {
	workdir = strings.TrimSpace(workdir)
	if workdir == "" {
		workdir = "."
	}
	var b strings.Builder
	b.WriteString("Please act on the human-authored review notes I left in Hunk.\n\n")
	b.WriteString("Load the hunk-review skill first: run `hunk skill path` and read the file it points to. Then follow that skill for the review on the workdir below.\n\n")
	fmt.Fprintf(&b, "- workdir: %s\n\n", workdir)
	b.WriteString("Important:\n")
	b.WriteString("- Fetch my notes with `hunk session comment list --repo <workdir> --type user --json`. The `--type user` flag is required — without it you get the legacy agent-authored view, not my notes.\n")
	b.WriteString("- Pull the corresponding hunk bodies from `hunk session review --repo <workdir> --include-patch --json` (the `patch` field per file).\n")
	b.WriteString("- Each note carries filePath / hunkIndex / newRange / body. Summarize your plan for addressing each note before you start editing.\n")
	return b.String()
}
