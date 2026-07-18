// jind-ai-hunk is a jin plugin that hands the user's in-progress hunk
// review off to the target jin session's agent.
//
// The plugin does not render diffs or collect comments itself — hunk
// already does that. Its whole job is to send a short message telling
// the agent to load `hunk skill path` and pick up the user's inline
// notes for this workdir.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/takaaki-s/jind-ai-hunk/internal/jin"
	"github.com/takaaki-s/jind-ai-hunk/internal/prompt"
)

func main() {
	var (
		sessionID = flag.String("session", os.Getenv("JIN_SESSION_ID"), "target jin session ID")
		workdir   = flag.String("workdir", os.Getenv("JIN_WORKDIR"), "hunk session repo root (defaults to the jin session's workdir)")
		dryRun    = flag.Bool("dry-run", false, "print the assembled message instead of sending it")
	)
	flag.Parse()

	if *sessionID == "" && !*dryRun {
		fmt.Fprintln(os.Stderr, "jind-ai-hunk: no session ID (pass --session <id> or set JIN_SESSION_ID)")
		os.Exit(1)
	}

	client := jin.New(os.Getenv("JIN_BIN"))

	// If the caller did not tell us a workdir, ask jin for the session's own
	// workdir. Non-fatal — the prompt template falls back to "." otherwise.
	if *workdir == "" && *sessionID != "" {
		if s, err := client.Info(*sessionID); err == nil && s.Workdir != "" {
			*workdir = s.Workdir
		}
	}

	msg := prompt.Build(*workdir)

	if *dryRun {
		fmt.Print(msg)
		return
	}

	if err := client.Send(*sessionID, msg); err != nil {
		fmt.Fprintln(os.Stderr, "jind-ai-hunk: send failed:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "jind-ai-hunk: sent hunk-review handoff to session", *sessionID)
}
