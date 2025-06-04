package debugger

import (
	"cmp"
	"context"
)

const (
	DefaultHeader   string = "debug"
	DefaultPassword string = "pass"
)

type Client struct {
	cfg Config
}

type Config struct {
	Header   string `toml:"header"`
	Password string `toml:"password"`
}

func New(opts Config) *Client {
	opts.Header = cmp.Or(opts.Header, DefaultHeader)
	opts.Password = cmp.Or(opts.Password, DefaultPassword)

	return &Client{
		cfg: opts,
	}
}

type ctxKey struct{}

func IsDebug(ctx context.Context) bool {
	return ctx.Value(ctxKey{}) != nil
}
