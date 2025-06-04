package httpdebug

import (
	"net/http"

	"github.com/go-chi/transport"
)

func (c *Client) Transport(next http.RoundTripper) http.RoundTripper {
	return transport.RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		r := transport.CloneRequest(req)

		if val := r.Context().Value(ctxKey{}); val != nil {
			r.Header.Set(c.cfg.Header, c.cfg.Password)
		}

		return next.RoundTrip(r)
	})
}
