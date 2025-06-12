package httpdebug

import (
	"net/http"
)

// Middleware enables debug mode in the context when an incoming request
// contains the given debug header.
func Middleware(debugHeader Header) func(next http.Handler) http.Handler {
	// If the debug header is not fully defined, just return a passthrough handler.
	if debugHeader.Key == "" || debugHeader.Value == "" {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// Otherwise, wrap the handler and check for the header.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(debugHeader.Key) == debugHeader.Value {
				r = r.WithContext(enableDebugMode(r.Context()))
			}
			next.ServeHTTP(w, r)
		})
	}
}
