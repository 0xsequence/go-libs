package debug

import (
	"github.com/0xsequence/go-libs/httpdebug"
)

// Deprecated: Use github.com/0xsequence/go-libs/config.Debug instead.
type Debug struct {
	Enabled   bool             `toml:"enabled"`
	BasicAuth BasicAuth        `toml:"basic_auth"`
	Header    httpdebug.Header `toml:"header"`
}

// Deprecated: Use github.com/0xsequence/go-libs/config.BasicAuth instead.
type BasicAuth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
