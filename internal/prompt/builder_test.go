package prompt

import (
	"strings"
	"testing"
)

func TestBuild_ContainsWorkdirAndSkillHint(t *testing.T) {
	got := Build("/home/me/proj")

	for _, want := range []string{
		"hunk skill path",
		"/home/me/proj",
		"--type user",
		"--include-patch",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt missing %q\n---\n%s", want, got)
		}
	}
}

func TestBuild_EmptyWorkdirFallback(t *testing.T) {
	got := Build("   ")
	// Falls back to "." so the agent still has a concrete target string.
	if !strings.Contains(got, "workdir: .\n") {
		t.Errorf("expected fallback workdir '.' in prompt, got:\n%s", got)
	}
}
