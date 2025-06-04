package httpdebug

import (
	"net/http"
)

// Middleware enables debug mode in the context when an incoming request
// contains the given debug header.
func Middleware(debugHeader Header) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(debugHeader.Key) == debugHeader.Value {
				r = r.WithContext(enableDebugMode(r.Context()))
			}

			next.ServeHTTP(w, r)
		})
	}
}
