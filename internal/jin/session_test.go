package jin

import (
	"errors"
	"reflect"
	"testing"
)

type call struct {
	name string
	args []string
}

type fakeRunner struct {
	calls []call
	out   []byte
	err   error
}

func (f *fakeRunner) run(name string, args ...string) ([]byte, error) {
	f.calls = append(f.calls, call{name: name, args: append([]string(nil), args...)})
	return f.out, f.err
}

func TestClient_Info_ParsesJSON(t *testing.T) {
	fr := &fakeRunner{out: []byte(`{"id":"abc-123","workdir":"/w","tmux_pane":"%42"}`)}
	c := &Client{Bin: "jin", Run: fr.run}

	sess, err := c.Info("some-selector")
	if err != nil {
		t.Fatalf("Info: %v", err)
	}
	if sess.ID != "abc-123" || sess.Workdir != "/w" || sess.TmuxPane != "%42" {
		t.Errorf("session = %+v", sess)
	}
	want := []call{{name: "jin", args: []string{"session", "info", "some-selector", "--json"}}}
	if !reflect.DeepEqual(fr.calls, want) {
		t.Errorf("calls = %+v, want %+v", fr.calls, want)
	}
}

func TestClient_Info_MissingID(t *testing.T) {
	fr := &fakeRunner{out: []byte(`{"workdir":"/w"}`)}
	c := &Client{Bin: "jin", Run: fr.run}
	if _, err := c.Info("x"); err == nil {
		t.Errorf("expected error on empty id")
	}
}

func TestClient_Info_MalformedJSON(t *testing.T) {
	fr := &fakeRunner{out: []byte(`not json`)}
	c := &Client{Bin: "jin", Run: fr.run}
	if _, err := c.Info("x"); err == nil {
		t.Errorf("expected parse error")
	}
}

func TestClient_Info_RunnerError(t *testing.T) {
	fr := &fakeRunner{err: errors.New("boom")}
	c := &Client{Bin: "jin", Run: fr.run}
	if _, err := c.Info("x"); err == nil {
		t.Errorf("expected propagated error")
	}
}

func TestClient_Send_ArgsAndErrors(t *testing.T) {
	fr := &fakeRunner{out: []byte("ok\n")}
	c := &Client{Bin: "/opt/jin/jin", Run: fr.run}

	if err := c.Send("session-1", "hello\nworld"); err != nil {
		t.Fatalf("Send: %v", err)
	}
	want := []call{{name: "/opt/jin/jin", args: []string{"session", "send", "session-1", "hello\nworld"}}}
	if !reflect.DeepEqual(fr.calls, want) {
		t.Errorf("calls = %+v, want %+v", fr.calls, want)
	}

	// Error propagation.
	fr2 := &fakeRunner{err: errors.New("daemon down")}
	c2 := &Client{Bin: "jin", Run: fr2.run}
	if err := c2.Send("s", "p"); err == nil {
		t.Errorf("expected error propagation")
	}
}

func TestNew_DefaultBin(t *testing.T) {
	if c := New(""); c.Bin != "jin" {
		t.Errorf("empty bin should default to 'jin', got %q", c.Bin)
	}
	if c := New("/usr/local/bin/jin"); c.Bin != "/usr/local/bin/jin" {
		t.Errorf("explicit bin not preserved: %q", c.Bin)
	}
}
