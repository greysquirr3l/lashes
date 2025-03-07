package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// BasicUsageExample demonstrates how to initialize and use the lashes library
func BasicUsageExample() {
	// Initialize with default options
	rotator, err := lashes.New(lashes.DefaultOptions())
	if (err != nil) {
		log.Fatalf("Failed to create rotator: %v", err)
	}

	// Add a proxy
	ctx := context.Background()
	err = rotator.AddProxy(ctx, "http://example.com:8080", domain.HTTPProxy)
	if err != nil {
		log.Fatalf("Failed to add proxy: %v", err)
	}

	// Get a proxy
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		log.Fatalf("Failed to get proxy: %v", err)
	}

	fmt.Printf("Using proxy: %s\n", proxy.URL)
}
