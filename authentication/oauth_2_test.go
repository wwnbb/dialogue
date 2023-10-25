package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"

	D "github.com/wwnbb/dialogue"
	"golang.org/x/oauth2"
)

func TestOAuth2(t *testing.T) {
	// Set up the OAuth2 config and the next handler function.
	config := OAuth2Config{
		Config: oauth2.Config{
			ClientID: "client-id",
			Endpoint: oauth2.Endpoint{},
		},
	}
	next := func(d *D.Dialogue) *D.Dialogue {
		return d
	}

	// Create a new request with a bearer token.
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer mock-token")

	// Create a new response recorder.
	rr := httptest.NewRecorder()

	// Create a new Dialogue with the request and response recorder.
	d := D.NewDialogue(req, rr)

	// Call the OAuth2 function.
	handler := OAuth2(config, next)
	handler(d)

	// Check the response status code.
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
