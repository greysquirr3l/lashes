// filepath: /Users/nickcampbell/Projects/go/lashes/examples/proxy_validation_test.go
package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// ProxyValidationExample demonstrates how to validate proxies
func ProxyValidationExample() {
	// Initialize with validation options
	opts := lashes.DefaultOptions()
	opts.ValidateOnStart = true
	opts.ValidationTimeout = 5 * time.Second
	opts.TestURL = "https://api.ipify.org?format=json"

	rotator, err := lashes.New(opts)
	if err != nil {
		log.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()
	
	// Add a proxy that will be validated on insertion
	err = rotator.AddProxy(ctx, "http://example.com:8080", domain.HTTPProxy)
	if err != nil {
		log.Fatalf("Failed to add proxy: %v", err)
	}

	// Manually trigger validation for all proxies
	if err := rotator.ValidateAll(ctx); err != nil {
		log.Printf("Warning: Some proxies failed validation: %v\n", err)
	}

	// Get a validated proxy
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		log.Fatalf("Failed to get proxy: %v", err)
	}

	fmt.Printf("Using validated proxy: %s\n", proxy.URL)
}
