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
	routes := RouteMap{
		"/<id:int>/info": func(d *Dialogue) *Dialogue {
			return WriteResponseString(d, http.StatusOK, "OK")
		},
	}

	handler := Chain(Switch(routes), NotFoundHandler())

	tests := []struct {
		path string
		want bool
	}{
		{path: "/123/info", want: true},
		{path: "/abc/info", want: false},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		d := NewDialogue(req, rr)
		handler(d)

		if rr.Code == http.StatusOK && !test.want {
			t.Errorf("handler returned wrong status code for path %s: got %v want %v",
				test.path, rr.Code, http.StatusNotFound)
		} else if rr.Code == http.StatusNotFound && test.want {
			t.Errorf("handler returned wrong status code for path %s: got %v want %v",
				test.path, rr.Code, http.StatusOK)
		}
	}
}

func TestValidateAndExtractParams(t *testing.T) {
	tests := []struct {
		pattern string
		path    string
		want    bool
	}{
		{pattern: "/<id:int>/info", path: "/123/info", want: true},
		{pattern: "/<id:int>/info", path: "/abc/info", want: false},
	}

	for _, test := range tests {
		got := validateAndExtractParams(test.pattern, test.path)
		if got != test.want {
			t.Errorf("validateAndExtractParams(%q, %q) = %v, want %v", test.pattern, test.path, got, test.want)
		}
	}
}
