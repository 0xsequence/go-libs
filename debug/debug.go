package debug

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/0xsequence/go-libs/httpdebug"
)

type Debug struct {
	BasicAuth BasicAuth        `toml:"basic_auth"`
	Header    httpdebug.Header `toml:"header"`
}

type BasicAuth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// Return middleware for basic auth on passed realm
func (b BasicAuth) Middleware(realm string) func(next http.Handler) http.Handler {
	return middleware.BasicAuth(realm, map[string]string{
		b.Username: b.Password,
	})
}
