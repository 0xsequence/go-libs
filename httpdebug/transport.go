package httpdebug

import (
	"net/http"

	"github.com/go-chi/transport"
)

// Transport propagates debug mode by adding the given debug header to outgoing
// requests, if enabled in the context.
func Transport(debugHeader Header) func(next http.RoundTripper) http.RoundTripper {
	// If the debug header is not fully defined, just return a passthrough transport.
	if debugHeader.Key == "" || debugHeader.Value == "" {
		return func(next http.RoundTripper) http.RoundTripper {
			return transport.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
				return next.RoundTrip(r)
			})
		}
	}

	// Wrap the transport to inject the header if debug mode is enabled.
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
