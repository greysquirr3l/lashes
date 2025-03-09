package examples

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/client/mock"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/rotation"
)

func setupMocksForTests() (func(), func()) {
	// Setup mock URL parser that preserves the exact URL string
	resetURLParser := mock.SetURLParser(func(rawURL string) (*url.URL, error) {
		return url.Parse(rawURL)
	})

	// Setup mock HTTP client
	resetClient := client.SetClientCreator(func(proxy *domain.Proxy, options client.Options) (*http.Client, error) {
		return &http.Client{
			Transport: &mock.MockTransport{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       http.NoBody,
				},
			},
		}, nil
	})

	return resetURLParser, resetClient
}

// TestRoundRobinStrategy tests that proxies can be rotated
func TestRoundRobinStrategy(t *testing.T) {
	// Actually use the parameter
	if testing.Short() {
		t.Skip("Skipping extended test in short mode")
	}

	resetURLParser, resetClient := setupMocksForTests()
	defer resetURLParser()
	defer resetClient()

	// Create a rotator with round-robin strategy
	opts := lashes.Options{
		Strategy:        rotation.RoundRobinStrategy,
		ValidateOnStart: false, // Skip validation to focus on rotation
	}

	rotator, err := lashes.New(opts)
	if err != nil {
		t.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()

	// Add test proxies with very distinct URLs
	proxyURLs := []string{
		"http://unique-proxy1.example.com:8080",
		"http://unique-proxy2.example.com:8080",
	}

	for _, url := range proxyURLs {
		if err := rotator.AddProxy(ctx, url, domain.HTTPProxy); err != nil {
			t.Fatalf("Failed to add proxy %s: %v", url, err)
		}
	}

	// Get proxies and verify rotation works
	seen := make(map[string]int)
	for i := 0; i < 4; i++ { // Get enough proxies to see rotation
		proxy, err := rotator.GetProxy(ctx)
		if err != nil {
			t.Fatalf("Failed to get proxy at index %d: %v", i, err)
		}

		// Record this proxy URL
		seen[proxy.URL]++
	}

	// Verify that we saw each proxy at least once
	for _, url := range proxyURLs {
		count, found := seen[url]
		if !found {
			t.Errorf("Expected to see proxy %s, but it was never returned", url)
		} else if count < 1 {
			t.Errorf("Expected to see proxy %s at least once, saw it %d times", url, count)
		}
	}

	// Verify we didn't see any proxies we didn't add
	if len(seen) != len(proxyURLs) {
		t.Errorf("Expected to see exactly %d unique proxies, saw %d", len(proxyURLs), len(seen))
	}
}

// TestRandomStrategy can remain unchanged
func TestRandomStrategy(t *testing.T) {
	resetURLParser, resetClient := setupMocksForTests()
	defer resetURLParser()
	defer resetClient()

	opts := lashes.DefaultOptions()
	opts.Strategy = rotation.RandomStrategy
	opts.ValidateOnStart = false // Skip validation in tests

	rotator, err := lashes.New(opts)
	if err != nil {
		t.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()
	// Add multiple proxies
	proxies := []string{
		"http://proxy1.example.com:8080",
		"http://proxy2.example.com:8080",
		"http://proxy3.example.com:8080",
	}

	for _, proxyURL := range proxies {
		if err := rotator.AddProxy(ctx, proxyURL, domain.HTTPProxy); err != nil {
			t.Fatalf("Failed to add proxy: %v", err)
		}
	}

	// Just verify we can get proxies without errors
	for i := 0; i < 10; i++ {
		_, err := rotator.GetProxy(ctx)
		if err != nil {
			t.Fatalf("Failed to get proxy: %v", err)
		}
	}
}
