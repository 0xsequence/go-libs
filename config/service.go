package config

import (
	"fmt"
	"net/url"
)

type Service struct {
	Disabled bool    `toml:"disabled"` // Disables the service.
	url      BaseURL `toml:"url"`      // Service BaseURL. Use URL() to get copy of *url.URL.

	// Mutually exclusive fields.
	JWTSecret string `toml:"jwt_secret"` // Secret for signing JWT tokens for S2S comms. Mutually exclusive with JWTToken and AccessKey.
	JWTToken  string `toml:"jwt_token"`  // Custom static JWT token for S2S comms. Mutually exclusive with JWTSecret and AccessKey.
	AccessKey string `toml:"access_key"` // Access key used as X-Access-Key header. Mutually exclusive with JWTSecret and JWTToken.

	DebugRequests bool `toml:"debug_requests"` // Enables HTTP request logging in CURL format.
}

// UnmarshalTOML implements custom TOML unmarshaling with validation.
func (s *Service) UnmarshalTOML(v any) error {
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map[string]any, got %T", v)
	}

	if val, ok := m["disabled"].(bool); ok {
		s.Disabled = val
		if val {
			return nil // Disabled. Stop validating other fields.
		}
	}

	if val, ok := m["url"].(string); ok {
		if err := s.url.UnmarshalText([]byte(val)); err != nil {
			return fmt.Errorf("failed to unmarshal url: %w", err)
		}
	}
	if val, ok := m["jwt_secret"].(string); ok {
		s.JWTSecret = val
	}
	if val, ok := m["jwt_token"].(string); ok {
		s.JWTToken = val
	}
	if val, ok := m["access_key"].(string); ok {
		s.AccessKey = val
	}
	if val, ok := m["debug_requests"].(bool); ok {
		s.DebugRequests = val
	}

	return s.Validate() //nolint:wrapcheck
}

func (s *Service) Validate() error {
	// Validate mutually exclusive auth fields
	switch {
	case
		s.JWTSecret != "" && s.JWTToken != "",
		s.JWTSecret != "" && s.AccessKey != "",
		s.JWTToken != "" && s.AccessKey != "":
		return fmt.Errorf("mutually exclusive auth fields: only one of jwt_secret, jwt_token, or access_key can be set")
	}
	return nil
}

func (s *Service) URL() *url.URL {
	return s.url.URL()
}
