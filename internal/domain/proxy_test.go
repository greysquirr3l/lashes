package domain_test

import (
	"testing"

	"github.com/greysquirr3l/lashes/internal/domain"
)

func TestProxyEnabledCompatibility(t *testing.T) {
	// Test cases for Enabled field
	testCases := []struct {
		name    string
		enabled bool
		want    bool
	}{
		{
			name:    "Enabled proxy",
			enabled: true,
			want:    true,
		},
		{
			name:    "Disabled proxy",
			enabled: false,
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a proxy with the initial values
			proxy := &domain.Proxy{
				Enabled: tc.enabled,
			}

			// Call SetEnabled to update value
			proxy.SetEnabled(tc.want)

			// Check if field is updated correctly
			if proxy.Enabled != tc.want {
				t.Errorf("Enabled = %v, want %v", proxy.Enabled, tc.want)
			}

			// Test GetEnabled method
			if got := proxy.GetEnabled(); got != tc.want {
				t.Errorf("GetEnabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestProxyParseURL(t *testing.T) {
	testCases := []struct {
		name    string
		urlStr  string
		wantErr bool
	}{
		{
			name:    "Valid HTTP proxy",
			urlStr:  "http://example.com:8080",
			wantErr: false,
		},
		{
			name:    "Valid SOCKS5 proxy",
			urlStr:  "socks5://example.com:1080",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			urlStr:  "://invalid-url",
			wantErr: true,
		},
		{
			name:    "Empty URL",
			urlStr:  "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proxy := &domain.Proxy{
				URL: tc.urlStr,
			}

			url, err := proxy.ParseURL()
			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseURL() error = nil, want error for invalid URL %q", tc.urlStr)
				}
			} else {
				if err != nil {
					t.Errorf("ParseURL() error = %v, want nil for valid URL %q", err, tc.urlStr)
				}
				if url.String() != tc.urlStr {
					t.Errorf("ParseURL() = %q, want %q", url.String(), tc.urlStr)
				}
			}
		})
	}
}

func TestProxyIsValid(t *testing.T) {
	testCases := []struct {
		name  string
		proxy domain.Proxy
		want  bool
	}{
		{
			name: "Valid proxy",
			proxy: domain.Proxy{
				ID:   "test-id",
				URL:  "http://example.com:8080",
				Type: domain.HTTPProxy,
			},
			want: true,
		},
		{
			name: "Missing ID",
			proxy: domain.Proxy{
				URL:  "http://example.com:8080",
				Type: domain.HTTPProxy,
			},
			want: false,
		},
		{
			name: "Missing Type",
			proxy: domain.Proxy{
				ID:  "test-id",
				URL: "http://example.com:8080",
			},
			want: false,
		},
		{
			name: "Invalid URL",
			proxy: domain.Proxy{
				ID:   "test-id",
				URL:  "://invalid-url",
				Type: domain.HTTPProxy,
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.proxy.IsValid()
			if got != tc.want {
				t.Errorf("IsValid() = %v, want %v", got, tc.want)
			}
		})
	}
}
