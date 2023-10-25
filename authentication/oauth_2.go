package authentication

import (
	"context"
	"errors"
	"net/http"
	"strings"

	D "github.com/wwnbb/dialogue"
	"golang.org/x/oauth2"
)

type OAuth2Config struct {
	Config oauth2.Config
}

func OAuth2(config OAuth2Config, next D.DialogueFunc) D.DialogueFunc {
	return func(d *D.Dialogue) *D.Dialogue {
		// Extract the token from the Authorization header
		authHeader := d.Request.Header.Get("Authorization")
		if authHeader == "" {
			// No Authorization header; return 401 Unauthorized
			return D.WriteResponseString(d, http.StatusUnauthorized, "Unauthorized")
		}

		tokenType, token, err := parseAuthHeader(authHeader)
		if err != nil || tokenType != "Bearer" {
			// Malformed Authorization header; return 401 Unauthorized
			return D.WriteResponseString(d, http.StatusUnauthorized, "Unauthorized")
		}

		// Create a new context with the token
		ctx := context.WithValue(d.Request.Context(), oauth2.HTTPClient, &http.Client{})
		tokenSource := config.Config.TokenSource(ctx, &oauth2.Token{AccessToken: token})

		// Verify the token
		_, err = tokenSource.Token()
		if err != nil {
			// Invalid token; return 401 Unauthorized
			return D.WriteResponseString(d, http.StatusUnauthorized, "Unauthorized")
		}

		// Token is valid; proceed to the next handler
		return d.Map(next)
	}
}

func parseAuthHeader(authHeader string) (string, string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", "", errors.New("malformed Authorization header")
	}
	return parts[0], parts[1], nil
}
