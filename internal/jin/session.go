// Package jin is a thin wrapper around the `jin` CLI that this plugin
// invokes to look up session metadata and to send the assembled review
// prompt back into the target session.
package jin

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Runner executes a command and returns its combined stdout. Injecting
// this via Client makes the wrapper testable without spawning real jin.
type Runner func(name string, args ...string) ([]byte, error)

// DefaultRunner runs the command and returns stdout only; stderr is
// folded into the returned error on non-zero exit.
func DefaultRunner(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w (stderr: %s)",
			name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return []byte(stdout.String()), nil
}

// Client talks to the jin daemon via the CLI. Bin resolves to
// $JIN_BIN when the dispatcher injects it, or "jin" from PATH otherwise.
type Client struct {
	Bin string
	Run Runner
}

// New returns a Client using the real exec runner.
func New(bin string) *Client {
	if bin == "" {
		bin = "jin"
	}
	return &Client{Bin: bin, Run: DefaultRunner}
}

// Session is the subset of `jin session info --json` output the plugin
// needs. Additional fields the CLI emits are ignored.
type Session struct {
	ID       string `json:"id"`
	Workdir  string `json:"workdir"`
	TmuxPane string `json:"tmux_pane"`
}

// Info resolves selector (a session ID or name) to a Session.
func (c *Client) Info(selector string) (Session, error) {
	out, err := c.Run(c.Bin, "session", "info", selector, "--json")
	if err != nil {
		return Session{}, err
	}
	var s Session
	if err := json.Unmarshal(out, &s); err != nil {
		return Session{}, fmt.Errorf("parse session info: %w", err)
	}
	if s.ID == "" {
		return Session{}, fmt.Errorf("session info returned empty id for selector %q", selector)
	}
	return s, nil
}

// Send delivers prompt to the session identified by id. prompt is a
// single argument to the jin CLI — no shell is interposed, so quoting
// is not this caller's concern.
func (c *Client) Send(id, prompt string) error {
	if _, err := c.Run(c.Bin, "session", "send", id, prompt); err != nil {
		return err
	}
	return nil
}
