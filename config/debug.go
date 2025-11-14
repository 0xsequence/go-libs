package config

import (
	"github.com/0xsequence/go-libs/httpdebug"
)

type Debug struct {
	Enabled   bool             `toml:"enabled"`
	BasicAuth BasicAuth        `toml:"basic_auth"`
	Header    httpdebug.Header `toml:"header"`
}
