package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/storage"
)

// BasicUsageExample shows how to use the lashes proxy rotator
func BasicUsageExample() {
	// Initialize with SQLite
	opts := lashes.Options{
		Storage: &storage.Options{
			Type: storage.SQLite,
			DSN:  "file:proxies.db?cache=shared&mode=rwc",
		},
	}

	// Create proxy rotator
	rotator, err := lashes.New(opts)
	if err != nil {
		log.Fatal(err)
	}

	// Add some proxies
	ctx := context.Background()
	err = rotator.AddProxy(ctx, "http://proxy1.example.com:8080", lashes.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	// Use the proxy rotator
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using proxy: %s\n", proxy.URL)
}
