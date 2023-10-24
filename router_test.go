package dialogue

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
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
	uuid_value, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name           string
		pattern        string
		path           string
		expected       bool
		expectedParams map[string]Param
	}{
		{
			name:     "Valid UUID",
			pattern:  "/user/<id:uuid4>",
			path:     "/user/123e4567-e89b-12d3-a456-426614174000",
			expected: true,
			expectedParams: map[string]Param{
				"id": {Type: "uuid4", Value: uuid_value},
			},
		},
		{
			name:     "Invalid UUID",
			pattern:  "/user/<id:uuid4>",
			path:     "/user/invalid-uuid",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			d := NewDialogue(req, w)

			got := validateAndExtractParams(d, tt.pattern)
			if got != tt.expected {
				t.Errorf("validateAndExtractParams() = %v, want %v", got, tt.expected)
			}

			if got && tt.expected { // Only check params if expected and got are true
				for param, value := range tt.expectedParams {
					if d.PathParams[param].Type != value.Type || d.PathParams[param].Value != value.Value {
						t.Errorf("Param %v: got %v, want %v", param, d.PathParams[param], value)
					}
				}
			}
		})
	}
}
