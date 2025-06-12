package httpdebug

import (
	"net/http"

	"github.com/go-chi/transport"
)

// Transport propagates debug mode by adding the given debug header to outgoing
// requests, if enabled in the context.
func Transport(debugHeader Header) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return transport.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			if IsDebugModeEnabled(r.Context()) {
				r = transport.CloneRequest(r)
				r.Header.Set(debugHeader.Key, debugHeader.Value)
			}

			return next.RoundTrip(r)
		})
	}
}
