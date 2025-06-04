package httpdebug

import (
	"context"
	"net/http"
)

// Middleware creates an HTTP middleware function that adds debug to context if the "debug" header has the expected value.
func (c *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the debug header
		if r.Header.Get(c.cfg.Header) == c.cfg.Password {
			// Add the expected key to the context
			ctx := context.WithValue(r.Context(), ctxKey{}, struct{}{})
			r = r.WithContext(ctx)
		}

		// Continue with the next handler
		next.ServeHTTP(w, r)
	})
}
