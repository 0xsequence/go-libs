package config

type Service struct {
	URL           string `toml:"url"`            // Service URL
	JWTSecret     string `toml:"jwt_secret"`     // Secret for signing S2S JWT tokens
	JWTToken      string `toml:"jwt_token"`      // Custom static JWT token for S2S comms. Mutually exclusive with JWTSecret.
	DebugRequests bool   `toml:"debug_requests"` // Enables HTTP request logging in CURL format
}
