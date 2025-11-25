package config

import (
	"net/url"
)

type URL struct {
	u url.URL
}

func (u *URL) UnmarshalText(text []byte) error {
	parsed, err := url.Parse(string(text))
	if err != nil {
		return err
	}
	u.u = *parsed
	return nil
}

func (u *URL) URL() *url.URL {
	copy := new(url.URL)
	*copy = u.u
	return copy
}
