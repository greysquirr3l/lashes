package validation_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/client/mock"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/validation"
)

func TestValidateProxy(t *testing.T) {
	// Create test proxy
	proxy := &domain.Proxy{
		ID:  "test-proxy",
		URL: "http://example.com:8080",
		Type: domain.HTTPProxy,
	}

	tests := []struct {
		name           string
		mockStatus     int
		mockDelay      time.Duration
		maxLatency     time.Duration
		expectValid    bool
		expectError    bool
	}{
		{
			name:        "Successful validation",
			mockStatus:  http.StatusOK,
			mockDelay:   10 * time.Millisecond,
			maxLatency:  100 * time.Millisecond,
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Server error",
			mockStatus:  http.StatusInternalServerError,
			mockDelay:   10 * time.Millisecond,
			maxLatency:  100 * time.Millisecond,
			expectValid: false,
			expectError: true,
		},
		{
			name:        "High latency",
			mockStatus:  http.StatusOK,
			mockDelay:   50 * time.Millisecond,
			maxLatency:  20 * time.Millisecond,
			expectValid: false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response based on the test case
			mockResp := &http.Response{
				StatusCode: tt.mockStatus,
				Body:       io.NopCloser(strings.NewReader(`{"ip":"127.0.0.1"}`)),
				Header:     make(http.Header),
			}

				// Setup mock client creator for this test
			resetClient := client.SetClientCreator(func(proxy *domain.Proxy, options client.Options) (*http.Client, error) {
				return &http.Client{
					Transport: &mock.MockTransport{
						Response: mockResp,
						Delay:    tt.mockDelay,
					},
				}, nil
			})
			defer resetClient() // Restore original client creator when done

			// Create validator with test configuration
			validator := validation.NewValidator(validation.Config{
				Timeout:    100 * time.Millisecond,
				RetryCount: 1,
				TestURL:    "http://test-url.com",
				MaxLatency: tt.maxLatency,
			})
			
			// Test validation
			ctx := context.Background()
			valid, latency, err := validator.Validate(ctx, proxy)

			// Check the validation result
			if valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got %v", tt.expectValid, valid)
			}

			// Check the error result
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error=%v, got error=%v: %v", tt.expectError, err != nil, err)
			}

			// Additional check for the latency
			if valid && latency < tt.mockDelay {
				t.Errorf("Expected latency >= %v, got %v", tt.mockDelay, latency)
			}
		})
	}
}
