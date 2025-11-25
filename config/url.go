package config

import (
	"fmt"
	"net/url"
)

type BaseURL struct {
	u url.URL
}

func (u *BaseURL) UnmarshalText(text []byte) error {
	str := string(text)
	parsed, err := url.Parse(str)
	if err != nil {
		return err //nolint:wrapcheck
	}
	if parsed.Host == "" {
		return fmt.Errorf("host is required: %q", str)
	}
	u.u = *parsed
	return nil
}

func (u *BaseURL) URL() *url.URL {
	if u == nil || u.u.Host == "" {
		return nil
	}
	copy := new(url.URL)
	*copy = u.u
	return copy
}
