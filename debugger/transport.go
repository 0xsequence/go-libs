package debugger

import (
	"net/http"
	"slices"
)

func (c *Client) Transport(next http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		r := cloneRequest(req)

		if val := r.Context().Value(ctxKey{}); val != nil {
			r.Header.Set(c.cfg.Header, c.cfg.Password)
		}

		return next.RoundTrip(r)
	})
}

// RoundTripFunc, similar to http.HandlerFunc, is an adapter
// to allow the use of ordinary functions as http.RoundTrippers.
type RoundTripFunc func(r *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// cloneRequest creates a shallow copy of a given request
// to comply with stdlib's http.RoundTripper contract:
//
// RoundTrip should not modify the request, except for
// consuming and closing the Request's Body. RoundTrip may
// read fields of the request in a separate goroutine. Callers
// should not mutate or reuse the request until the Response's
// Body has been closed.
func cloneRequest(orig *http.Request) *http.Request {
	clone := &http.Request{}
	*clone = *orig

	clone.Header = make(http.Header, len(orig.Header))
	for key, value := range orig.Header {
		clone.Header[key] = slices.Clone(value)
	}

	return clone
}
