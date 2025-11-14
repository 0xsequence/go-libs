package config

type BasicAuth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
