package config

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestURLUnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantHost string
		wantPort string
	}{
		{
			name:     "http://localhost:4242",
			input:    "http://localhost:4242",
			wantHost: "localhost",
			wantPort: "4242",
		},
		{
			name:     "https with port",
			input:    "https://example.com:8080",
			wantHost: "example.com",
			wantPort: "8080",
		},
		{
			name:     "no port",
			input:    "http://example.com",
			wantHost: "example.com",
			wantPort: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u BaseURL
			err := u.UnmarshalText([]byte(tt.input))
			assert.NoError(t, err)

			parsedURL := u.URL()
			assert.Equal(t, tt.wantHost, parsedURL.Hostname())
			assert.Equal(t, tt.wantPort, parsedURL.Port())
		})
	}
}

func TestURLCopy(t *testing.T) {
	var u BaseURL
	err := u.UnmarshalText([]byte("http://localhost:4242"))
	assert.NoError(t, err)

	url1 := u.URL()
	url2 := u.URL()

	assert.Equal(t, url1.String(), url2.String())
	if url1 == url2 {
		t.Error("URL() should return different pointers")
	}

	// Make sure the original URL is not modified and we get a new copy every time.
	url1.Scheme = "https"
	assert.Equal(t, "http", url2.Scheme)
	url3 := u.URL()
	assert.Equal(t, "http", url3.Scheme)
}

func TestURLUninitialized(t *testing.T) {
	// Test case where UnmarshalText was never called (url not defined in config)
	var u BaseURL
	parsedURL := u.URL()
	assert.Nil(t, parsedURL, "URL() should return nil when BaseURL is uninitialized")
}

func TestURLEmptyText(t *testing.T) {
	// Test case where UnmarshalText is called with empty text
	var u BaseURL
	err := u.UnmarshalText([]byte(""))
	assert.Error(t, err, "UnmarshalText should error on empty text")

	// Even if UnmarshalText fails, URL() should return nil
	parsedURL := u.URL()
	assert.Nil(t, parsedURL, "URL() should return nil when UnmarshalText failed")
}

func TestURLNilReceiver(t *testing.T) {
	// Test case where receiver is nil
	var u *BaseURL
	parsedURL := u.URL()
	assert.Nil(t, parsedURL, "URL() should return nil when receiver is nil")
}
