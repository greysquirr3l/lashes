package examples

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/rotation"
)

func TestBasicUsageExample(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer mockServer.Close()

	// Initialize with proper rotation strategy
	opts := lashes.Options{
		Strategy:          rotation.RoundRobinStrategy,
		MaxRetries:        3,
		RequestTimeout:    1 * time.Second,
		ValidationTimeout: 500 * time.Millisecond,
		TestURL:           mockServer.URL, // Use our mock server
	}

	rotator, err := lashes.New(opts)
	if err != nil {
		t.Fatalf("Failed to create rotator: %v", err)
	}

	// Add a test proxy
	ctx := context.Background()
	err = rotator.AddProxy(ctx, "http://example.com:8080", domain.HTTPProxy)
	if err != nil {
		t.Fatalf("Failed to add proxy: %v", err)
	}

	// Test GetProxy
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		t.Fatalf("Failed to get proxy: %v", err)
	}

	if proxy.URL != "http://example.com:8080" {
		t.Errorf("Expected proxy URL 'http://example.com:8080', got '%s'", proxy.URL)
	}
}
