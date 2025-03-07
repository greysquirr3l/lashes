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

func setupMocksForTests(t *testing.T) (func(), func()) {
	// Setup mock URL parser
	resetURLParser := mock.SetURLParser(func(rawURL string) (*url.URL, error) {
		// Parse as-is but ignore errors for test proxies
		u, _ := url.Parse(rawURL)
		return u, nil
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

// TestRoundRobinStrategy verifies the behavior of round-robin rotation
func TestRoundRobinStrategy(t *testing.T) {
	resetURLParser, resetClient := setupMocksForTests(t)
	defer resetURLParser()
	defer resetClient()

	// Create a new rotator with round-robin strategy
	opts := lashes.Options{
		Strategy:        rotation.RoundRobinStrategy,
		ValidateOnStart: false,
	}

	rotator, err := lashes.New(opts)
	if err != nil {
		t.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()

	// Add just two proxies for simpler testing
	proxyURLs := []string{
		"http://proxy1.example.com:8080",
		"http://proxy2.example.com:8080",
	}

	for _, url := range proxyURLs {
		if err := rotator.AddProxy(ctx, url, domain.HTTPProxy); err != nil {
			t.Fatalf("Failed to add proxy %s: %v", url, err)
		}
	}

	// Get all proxies to verify they were added
	allProxies, err := rotator.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list proxies: %v", err)
	}
	
	// Verify we have exactly the number of proxies we added
	if len(allProxies) != len(proxyURLs) {
		t.Fatalf("Expected %d proxies, got %d", len(proxyURLs), len(allProxies))
	}

	// Record first round order
	firstRound := make([]string, 0, len(proxyURLs))
	for i := 0; i < len(proxyURLs); i++ {
		proxy, err := rotator.GetProxy(ctx)
		if err != nil {
			t.Fatalf("Failed to get proxy at index %d: %v", i, err)
			}
		firstRound = append(firstRound, proxy.URL)
	}

	// Verify we saw all proxies
	seenProxies := make(map[string]bool)
	for _, url := range firstRound {
		seenProxies[url] = true
	}
	if len(seenProxies) != len(proxyURLs) {
		t.Errorf("Expected to see %d unique proxies, saw %d", len(proxyURLs), len(seenProxies))
	}

	// Verify second round matches first round (round-robin behavior)
	for i := 0; i < len(proxyURLs); i++ {
		proxy, err := rotator.GetProxy(ctx)
		if err != nil {
			t.Fatalf("Failed to get proxy at index %d in second round: %v", i, err)
		}

		if proxy.URL != firstRound[i] {
			t.Errorf("Proxy order changed at position %d: expected %s, got %s", 
				i, firstRound[i], proxy.URL)
		}
	}
}

func TestRandomStrategy(t *testing.T) {
	resetURLParser, resetClient := setupMocksForTests(t)
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
