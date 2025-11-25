package config

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/test-go/testify/assert"
)

func TestServicesTOML(t *testing.T) {
	var config struct {
		Services struct {
			Metadata Service `toml:"metadata"`
			Indexer  Service `toml:"indexer"`
		} `toml:"services"`
	}

	tt := []struct {
		testName    string
		tomlData    string
		expectError bool
	}{
		{
			testName: "valid services",
			tomlData: `
			[services.metadata]
				url = "http://localhost:4242"
				jwt_secret = "secret"
				debug_requests = false

			[services.indexer]
				url = "https://indexer.sequence.app"
				access_key = "key"
				debug_requests = true
`,
			expectError: false,
		},
	}

	for _, tt := range tt {
		t.Run(tt.testName, func(t *testing.T) {
			err := toml.Unmarshal([]byte(tt.tomlData), &config)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, "http://localhost:4242", config.Services.Metadata.URL().String())
			assert.Equal(t, "https://indexer.sequence.app", config.Services.Indexer.URL().String())
		})
	}
}

func TestServiceTOMLFields(t *testing.T) {
	var config struct {
		Service Service `toml:"service"`
		Indexer Service `toml:"indexer"`
	}

	tt := []struct {
		testName    string
		tomlData    string
		expectError bool
	}{
		{
			testName: "disabled service - URL can be empty",
			tomlData: `
			[service]
				disabled = true
				url = ""
				jwt_token = "token"
				debug_requests = false
`,
			expectError: false,
		},
		{
			testName: "enabled service - URL can't be empty",
			tomlData: `
			[service]
				url = ""
				jwt_token = "token"
				debug_requests = false
`,
			expectError: true,
		},
		{
			testName: "jwt_secret and jwt_token",
			tomlData: `
			[service]
				url = "http://localhost:4242"
				jwt_secret = "secret"
				jwt_token = "token"
				debug_requests = false
`,
			expectError: true,
		},
		{
			testName: "jwt_token and access_key",
			tomlData: `
			[service]
				url = "http://localhost:4242"
				jwt_token = "token"
				access_key = "key"
				debug_requests = false
`,
			expectError: true,
		},
		{
			testName: "jwt_secret and access_key",
			tomlData: `
			[service]
				url = "http://localhost:4242"
				jwt_secret = "secret"
				access_key = "key"
				debug_requests = false
`,
			expectError: true,
		},
	}

	for _, tt := range tt {
		t.Run(tt.testName, func(t *testing.T) {
			err := toml.Unmarshal([]byte(tt.tomlData), &config)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
