package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// CustomStorage is a simple in-memory storage for demonstration purposes.
type CustomStorage struct {
	proxies []string
}

// NewCustomStorage initializes a new CustomStorage with some example proxies.
func NewCustomStorage() *CustomStorage {
	return &CustomStorage{
		proxies: []string{
			"http://example-proxy1.com:8080",
			"http://example-proxy2.com:8080",
		},
	}
}

// CustomStorageExample demonstrates how to use lashes with a custom storage backend
func CustomStorageExample() {
	// Initialize with custom storage options
	opts := lashes.DefaultOptions()
	// In a real example, we would configure storage options here
	
	rotator, err := lashes.New(opts)
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

	fmt.Printf("Using proxy with custom storage: %s\n", proxy.URL)
}
