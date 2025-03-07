package examples

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/client/mock"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// TestErrorHandling tests error handling functionality
func TestErrorHandling(t *testing.T) {
	// Create a mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"ip": "127.0.0.1"}`)),
		Header:     make(http.Header),
	}

	// Setup proper URL mocking
	resetURLParser := mock.SetURLParser(func(rawURL string) (*url.URL, error) {
		if rawURL == "invalid-url" {
			// Return a clear error for invalid URLs
			return nil, errors.New("invalid URL format")
		}
		
		// Pass through valid URLs to the real parser
		parsedURL, _ := url.Parse("http://example.com:8080")
		return parsedURL, nil
	})
	defer resetURLParser()

	// Create test HTTP client
	resetClient := client.SetClientCreator(func(proxy *domain.Proxy, options client.Options) (*http.Client, error) {
		return &http.Client{
			Transport: &mock.MockTransport{
				Response: mockResp,
			},
		}, nil
	})
	defer resetClient()

	// Setup rotator with test options
	opts := lashes.DefaultOptions()
	opts.ValidateOnStart = false
	
	rotator, err := lashes.New(opts)
	if err != nil {
		t.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()

	// Test invalid URL error handling
	err = rotator.AddProxy(ctx, "invalid-url", domain.HTTPProxy)
	if err == nil {
		t.Fatalf("ERROR: Expected error when adding invalid proxy URL, got nil")
	} else {
		t.Logf("Correctly got error for invalid URL: %v", err)
	}

	// Add and then remove a proxy to test no proxies available
	err = rotator.AddProxy(ctx, "http://example.com:8080", domain.HTTPProxy)
	if err != nil {
		t.Fatalf("Failed to add proxy: %v", err)
	}

	// Get and remove the proxy
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		t.Fatalf("Failed to get proxy: %v", err)
	}

	err = rotator.RemoveProxy(ctx, proxy.URL)
	if err != nil {
		t.Fatalf("Failed to remove proxy: %v", err)
	}

	// There should be no proxies available now
	_, err = rotator.GetProxy(ctx)
	if err == nil {
		t.Fatalf("ERROR: Expected error when getting proxy with none available, got nil")
	} else if errors.Is(err, lashes.ErrNoProxiesAvailable) {
		t.Logf("Correctly got ErrNoProxiesAvailable: %v", err)
	} else {
		t.Errorf("Expected ErrNoProxiesAvailable, got different error: %v", err)
	}
}
