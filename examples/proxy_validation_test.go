package examples

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/client/mock"
	"github.com/greysquirr3l/lashes/internal/domain"
)

func TestProxyValidation(t *testing.T) {
	// Create a mock response with real content
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"ip": "127.0.0.1"}`)),
		Header:     make(http.Header),
	}

	// Set up URL parser to ensure all URLs parse successfully
	resetURLParser := mock.SetURLParser(func(rawURL string) (*url.URL, error) {
		parsedURL, _ := url.Parse("http://example.com:8080")
		return parsedURL, nil
	})
	defer resetURLParser()

	// Set up a mock HTTP client with reliable responses
	resetClient := client.SetClientCreator(func(proxy *domain.Proxy, options client.Options) (*http.Client, error) {
		return &http.Client{
			Transport: &mock.MockTransport{
				Response: mockResp,
				Delay:    5 * time.Millisecond,
			},
		}, nil
	})
	defer resetClient()

	// Create simplified options for testing
	opts := lashes.DefaultOptions()
	opts.ValidateOnStart = true
	opts.ValidationTimeout = 100 * time.Millisecond
	opts.TestURL = "http://test-url.local" // This URL doesn't matter since we're mocking

	rotator, err := lashes.New(opts)
	if err != nil {
		t.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()
	
	// Add a test proxy
	err = rotator.AddProxy(ctx, "http://example.com:8080", domain.HTTPProxy)
	if err != nil {
		t.Fatalf("Failed to add proxy: %v", err)
	}

	// Validate should now succeed with our mocks
	if err := rotator.ValidateAll(ctx); err != nil {
		t.Fatalf("Failed to validate proxies: %v", err)
	}

	// Verify proxy is still available after validation
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		t.Fatalf("Failed to get proxy after validation: %v", err)
	}
	
	if proxy.URL != "http://example.com:8080" {
		t.Errorf("Expected proxy URL 'http://example.com:8080', got '%s'", proxy.URL)
	}
}
