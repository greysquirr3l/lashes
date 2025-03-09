package lashes

import (
	"errors"
	"strings"
	"testing"
)

func TestValidationError(t *testing.T) {
	testCases := []struct {
		name        string
		proxyID     string
		proxyURL    string
		reason      string
		statusCode  int
		wantContain []string
	}{
		{
			name:       "Full error information",
			proxyID:    "proxy-123",
			proxyURL:   "http://example.com:8080",
			reason:     "connection timeout",
			statusCode: 504,
			wantContain: []string{
				"proxy validation failed for http://example.com:8080",
				"connection timeout",
				"status code: 504",
			},
		},
		{
			name:       "No status code",
			proxyID:    "proxy-123",
			proxyURL:   "http://example.com:8080",
			reason:     "connection refused",
			statusCode: 0,
			wantContain: []string{
				"proxy validation failed for http://example.com:8080",
				"connection refused",
			},
		},
		{
			name:       "No reason",
			proxyID:    "proxy-123",
			proxyURL:   "http://example.com:8080",
			reason:     "",
			statusCode: 403,
			wantContain: []string{
				"proxy validation failed for http://example.com:8080",
				"status code: 403",
			},
		},
		{
			name:       "Only URL information",
			proxyID:    "proxy-123",
			proxyURL:   "http://example.com:8080",
			reason:     "",
			statusCode: 0,
			wantContain: []string{
				"proxy validation failed for http://example.com:8080",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewValidationError(tc.proxyID, tc.proxyURL, tc.reason, tc.statusCode)

			// Check error message contains expected parts
			errMsg := err.Error()
			for _, want := range tc.wantContain {
				if !strings.Contains(errMsg, want) {
					t.Errorf("Error message does not contain %q, got: %q", want, errMsg)
				}
			}

			// Test errors.Is implementation
			if !errors.Is(err, ErrValidationFailed) {
				t.Error("ValidationError should match ErrValidationFailed with errors.Is")
			}

			// Test that the proxy ID is stored correctly
			if err.ProxyID != tc.proxyID {
				t.Errorf("ProxyID = %q, want %q", err.ProxyID, tc.proxyID)
			}
		})
	}
}
