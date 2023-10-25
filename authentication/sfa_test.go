package authentication

import (
	"encoding/base64"
	D "github.com/wwnbb/dialogue"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	var hello_world = func(d *D.Dialogue) *D.Dialogue {
		return D.WriteResponseString(d, http.StatusOK, "Hello World")
	}
	config := BasicAuthConfig{
		Username: "testuser",
		Password: "testpass",
	}

	basicAuthMiddleware := BasicAuth(config)

	tests := []struct {
		name               string
		givenAuthorization string
		expectedStatusCode int
	}{
		{
			name:               "Valid Credentials",
			givenAuthorization: "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpass")),
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Credentials",
			givenAuthorization: "Basic " + base64.StdEncoding.EncodeToString([]byte("invaliduser:invalidpass")),
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Missing Authorization Header",
			givenAuthorization: "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.givenAuthorization != "" {
				req.Header.Set("Authorization", tt.givenAuthorization)
			}

			rec := httptest.NewRecorder()
			d := D.NewDialogue(req, rec)

			D.Chain(basicAuthMiddleware, hello_world, D.NotFoundHandler())(d)

			res := rec.Result()

			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}
		})
	}
}
