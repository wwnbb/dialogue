package dialogue

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
	d := &Dialogue{}
	handler1 := func(d *Dialogue) *Dialogue {
		d.Logs = append(d.Logs, "handler1 called")
		return d
	}
	handler2 := func(d *Dialogue) *Dialogue {
		d.Logs = append(d.Logs, "handler2 called")
		return d
	}

	chain := Chain(handler1, handler2)
	chain(d)

	if len(d.Logs) != 2 {
		t.Errorf("Chain did not call all handlers: %v", d.Logs)
	}
}

func TestMap(t *testing.T) {
	d := &Dialogue{}
	handler := func(d *Dialogue) *Dialogue {
		d.Logs = append(d.Logs, "handler called")
		return d
	}

	d.Map(handler)

	if len(d.Logs) != 1 || d.Logs[0] != "handler called" {
		t.Errorf("Map did not call the handler: %v", d.Logs)
	}
}

func TestSwitch(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	d := NewDialogue(req, rec)

	routes := RouteMap{
		"/test": func(d *Dialogue) *Dialogue {
			d.Logs = append(d.Logs, "test route hit")
			return d
		},
	}

	switchHandler := Switch(routes)
	switchHandler(d)

	if len(d.Logs) != 1 || d.Logs[0] != "test route hit" {
		t.Errorf("Switch did not route correctly: %v", d.Logs)
	}
}
