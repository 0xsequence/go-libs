package debug

import (
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
