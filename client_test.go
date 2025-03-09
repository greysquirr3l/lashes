package lashes

import (
	"net/http"
	"testing"
	"time"
)

func TestDefaultClientOptions(t *testing.T) {
	opts := DefaultClientOptions()

	// Check default values
	if opts.Timeout != 30*time.Second {
		t.Errorf("Default Timeout = %v, want %v", opts.Timeout, 30*time.Second)
	}

	if opts.MaxRetries != 3 {
		t.Errorf("Default MaxRetries = %d, want %d", opts.MaxRetries, 3)
	}

	if !opts.FollowRedirects {
		t.Error("Default FollowRedirects should be true")
	}

	if !opts.VerifyCerts {
		t.Error("Default VerifyCerts should be true")
	}

	if opts.Headers != nil {
		t.Error("Default Headers should be nil")
	}

	if opts.UserAgent != "" {
		t.Errorf("Default UserAgent should be empty, got %q", opts.UserAgent)
	}
}

func TestNewClient(t *testing.T) {
	// Create a test proxy
	proxy := &Proxy{
		ID:   "test-proxy",
		URL:  "http://example.com:8080",
		Type: HTTP,
	}

	// Test with default options
	client, err := NewClient(proxy, DefaultClientOptions())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewClient returned nil client")
	}

	// Test with custom options
	customOpts := ClientOptions{
		Timeout:         5 * time.Second,
		MaxRetries:      2,
		FollowRedirects: false,
		VerifyCerts:     false,
		Headers: http.Header{
			"X-Test-Header": []string{"test-value"},
		},
		UserAgent: "TestUserAgent/1.0",
	}

	customClient, err := NewClient(proxy, customOpts)
	if err != nil {
		t.Fatalf("NewClient with custom options failed: %v", err)
	}

	if customClient == nil {
		t.Fatal("NewClient with custom options returned nil client")
	}

	// Can't easily test more details without mocking, but at least we confirmed
	// the client creation succeeds with different options
}
