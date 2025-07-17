package config_test

import (
	"testing"

	"github.com/Ricky004/watchdata/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUri(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantScheme string
		wantValue  string
		wantErr    bool
	}{
		{
			name:       "valid http URI",
			input:      "http://example.com",
			wantScheme: "http",
			wantValue:  "//example.com",
			wantErr:    false,
		},
		{
			name:       "valid https URI",
			input:      "https://example.com/path",
			wantScheme: "https",
			wantValue:  "//example.com/path",
			wantErr:    false,
		},
		{
			name:       "valid ftp URI",
			input:      "ftp://files.example.com",
			wantScheme: "ftp",
			wantValue:  "//files.example.com",
			wantErr:    false,
		},
		{
			name:       "valid file URI",
			input:      "file:///path/to/file",
			wantScheme: "file",
			wantValue:  "///path/to/file",
			wantErr:    false,
		},
		{
			name:       "valid custom scheme",
			input:      "custom-scheme://value",
			wantScheme: "custom-scheme",
			wantValue:  "//value",
			wantErr:    false,
		},
		{
			name:       "scheme with numbers",
			input:      "http2://example.com",
			wantScheme: "http2",
			wantValue:  "//example.com",
			wantErr:    false,
		},
		{
			name:       "scheme with plus",
			input:      "git+ssh://git@github.com/user/repo.git",
			wantScheme: "git+ssh",
			wantValue:  "//git@github.com/user/repo.git",
			wantErr:    false,
		},
		{
			name:       "scheme with dot",
			input:      "x.y://value",
			wantScheme: "x.y",
			wantValue:  "//value",
			wantErr:    false,
		},
		{
			name:       "scheme with hyphen",
			input:      "x-y://value",
			wantScheme: "x-y",
			wantValue:  "//value",
			wantErr:    false,
		},
		{
			name:       "simple value without authority",
			input:      "mailto:user@example.com",
			wantScheme: "mailto",
			wantValue:  "user@example.com",
			wantErr:    false,
		},
		{
			name:       "empty value",
			input:      "scheme:",
			wantScheme: "scheme",
			wantValue:  "",
			wantErr:    false,
		},
		{
			name:    "missing scheme",
			input:   "//example.com",
			wantErr: true,
		},
		{
			name:    "missing colon",
			input:   "httpexample.com",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "scheme starting with number",
			input:   "2http://example.com",
			wantErr: true,
		},
		{
			name:    "scheme with invalid character",
			input:   "ht_tp://example.com",
			wantErr: true,
		},
		{
			name:       "complex URI with query and fragment",
			input:      "https://example.com:8080/path?query=value&other=test#fragment",
			wantScheme: "https",
			wantValue:  "//example.com:8080/path?query=value&other=test#fragment",
			wantErr:    false,
		},
		{
			name:       "URI with special characters in value",
			input:      "custom://user:pass@host:port/path?q=a&b=c#frag",
			wantScheme: "custom",
			wantValue:  "//user:pass@host:port/path?q=a&b=c#frag",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := config.NewUri(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid uri")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantScheme, uri.Scheme())
			assert.Equal(t, tt.wantValue, uri.Value())
		})
	}
}

func TestMustNewUri(t *testing.T) {
	t.Run("valid URI should not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			uri := config.MustNewUri("http://example.com")
			assert.Equal(t, "http", uri.Scheme())
			assert.Equal(t, "//example.com", uri.Value())
		})
	})

	t.Run("invalid URI should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			config.MustNewUri("invalid-uri")
		})
	})

	t.Run("empty URI should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			config.MustNewUri("")
		})
	})
}

func TestUri_Scheme(t *testing.T) {
	uri := config.MustNewUri("https://example.com/path")
	assert.Equal(t, "https", uri.Scheme())
}

func TestUri_Value(t *testing.T) {
	uri := config.MustNewUri("https://example.com/path")
	assert.Equal(t, "//example.com/path", uri.Value())
}

func TestUriRegex(t *testing.T) {
	// Test the regex directly to ensure it matches expected patterns
	tests := []struct {
		name    string
		input   string
		matches bool
	}{
		{"valid scheme", "http://example.com", true},
		{"scheme with numbers", "h2://example.com", true},
		{"scheme with plus", "git+ssh://example.com", true},
		{"scheme with dot", "x.y://example.com", true},
		{"scheme with hyphen", "x-y://example.com", true},
		{"no scheme", "//example.com", false},
		{"scheme starting with number", "2http://example.com", false},
		{"scheme with underscore", "ht_tp://example.com", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := config.UriRegex.MatchString(tt.input)
			assert.Equal(t, tt.matches, matches)
		})
	}
}
