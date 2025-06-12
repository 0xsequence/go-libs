package httpdebug

// Header represents a special header that enables debug mode.
type Header struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}
