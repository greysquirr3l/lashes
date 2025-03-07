package validation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/domain"
)

type Validator interface {
	Validate(ctx context.Context, proxy *domain.Proxy) (bool, time.Duration, error)
	ValidateWithTarget(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error)
}

type Config struct {
	Timeout    time.Duration
	RetryCount int
	TestURL    string
	MaxLatency time.Duration
	Concurrent int
}

type validator struct {
	config Config
}

func NewValidator(config Config) Validator {
	// Apply defaults if not specified
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.TestURL == "" {
		config.TestURL = "https://httpbin.org/ip" // Default test URL
	}
	if config.MaxLatency == 0 {
		config.MaxLatency = 5 * time.Second // Default max acceptable latency
	}
	
	return &validator{
		config: config,
	}
}

// Validate checks if the proxy is working properly
func (v *validator) Validate(ctx context.Context, proxy *domain.Proxy) (bool, time.Duration, error) {
	return v.ValidateWithTarget(ctx, proxy, v.config.TestURL)
}

// ValidateWithTarget validates a proxy against a specific target URL
func (v *validator) ValidateWithTarget(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, v.config.Timeout)
	defer cancel()
	
	// Create an HTTP client using the proxy
	httpClient, err := client.NewClient(proxy, client.Options{
		Timeout:         v.config.Timeout,
		MaxRetries:      v.config.RetryCount,
		VerifyCerts:     true,
		FollowRedirects: false,
	})
	if err != nil {
		return false, 0, fmt.Errorf("failed to create HTTP client: %w", err)
	}
	
	// Prepare the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return false, 0, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Execute the request and measure response time
	startTime := time.Now()
	resp, err := httpClient.Do(req)
	latency := time.Since(startTime)
	
	// If request failed, return error
	if err != nil {
		return false, latency, fmt.Errorf("request failed: %w", err)
	}
	
	// Ensure response body is closed properly and capture any close errors
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// If we already have an error, wrap the close error with it
			if err != nil {
				err = errors.Join(err, fmt.Errorf("error closing response body: %w", closeErr))
			} else {
				err = fmt.Errorf("error closing response body: %w", closeErr)
			}
		}
	}()
	
	// Verify the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, latency, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	
	// Check if latency is acceptable
	if latency > v.config.MaxLatency {
		return false, latency, fmt.Errorf("latency too high: %s", latency)
	}
	
	// Proxy validation successful
	return true, latency, nil
}
