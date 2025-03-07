// filepath: /Users/nickcampbell/Projects/go/lashes/examples/error_handling_test.go
package examples

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// ErrorHandlingExample demonstrates how to handle various error conditions with lashes
func ErrorHandlingExample() {
	opts := lashes.DefaultOptions()
	opts.MaxRetries = 2
	opts.RequestTimeout = 100 * time.Millisecond

	rotator, err := lashes.New(opts)
	if err != nil {
		log.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()

	// Example 1: Handle invalid proxy URL
	err = rotator.AddProxy(ctx, "invalid-url", domain.HTTPProxy)
	if err != nil {
		fmt.Printf("Expected error when adding invalid proxy URL: %v\n", err)
	}

	// Example 2: Add a valid proxy
	err = rotator.AddProxy(ctx, "http://example.com:8080", domain.HTTPProxy)
	if err != nil {
		log.Fatalf("Failed to add proxy: %v", err)
	}

	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		log.Fatalf("Failed to get proxy: %v", err)
	}
	
	fmt.Printf("Using proxy: %s\n", proxy.URL)
	
	// Example 3: Validate a proxy
	isValid, latency, err := rotator.ValidateProxy(ctx, proxy, "http://example.com/test")
	fmt.Printf("Proxy validation result: valid=%v, latency=%v, error=%v\n", 
		isValid, latency, err)
}
