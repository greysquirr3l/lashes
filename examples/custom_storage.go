package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/lashes"
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

func CustomStorageExample() {
	ctx := context.Background()

	// Initialize with default options (in-memory)
	rotator, err := lashes.New(lashes.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}

	// Add example proxies
	for _, proxyURL := range []string{
		"http://example-proxy1.com:8080",
		"http://example-proxy2.com:8080",
	} {
		err = rotator.AddProxy(ctx, proxyURL, lashes.HTTP)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Use the rotator
	proxy, err := rotator.GetProxy(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using proxy: %s\n", proxy.URL)
}
