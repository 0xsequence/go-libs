package middleware

import (
	"net/http"

	"github.com/0xsequence/go-libs/config"
	"github.com/go-chi/chi/v5/middleware"
)

// BasicAuth protects routes with basic authentication. Returns 404 if credentials are not configured.
func BasicAuth(creds config.BasicAuth) func(next http.Handler) http.Handler {
	if creds.Username == "" || creds.Password == "" {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Missing or misconfigured credentials.
				// Return HTTP 404.
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			})
		}
	}
	return middleware.BasicAuth("sequence", map[string]string{
		creds.Username: creds.Password,
	})
}
