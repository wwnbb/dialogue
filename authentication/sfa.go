package authentication

import (
	"encoding/base64"
	D "github.com/wwnbb/dialogue"
	"net/http"
	"strings"
)

type BasicAuthConfig struct {
	Username string
	Password string
}

func BasicAuth(config BasicAuthConfig) D.DialogueFunc {
	return func(d *D.Dialogue) *D.Dialogue {
		auth := d.Request.Header.Get("Authorization")
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			d = D.SetHeader(d, "WWW-Authenticate", `Basic realm="Restricted"`)
			return D.WriteResponseString(d, http.StatusUnauthorized, "Unauthorized")
		}

		payload, _ := base64.StdEncoding.DecodeString(parts[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || pair[0] != config.Username || pair[1] != config.Password {
			return D.WriteResponseString(d, http.StatusUnauthorized, "Unauthorized")
		}

		return d
	}
}
