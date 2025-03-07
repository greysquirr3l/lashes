package examples

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/rotation"
)

// RoundRobinStrategyExample demonstrates using the round-robin rotation strategy
func RoundRobinStrategyExample() {
	opts := lashes.DefaultOptions()
	opts.Strategy = rotation.RoundRobinStrategy

	rotator, err := lashes.New(opts)
	if err != nil {
		log.Fatalf("Failed to create rotator: %v", err)
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
			log.Fatalf("Failed to add proxy: %v", err)
		}
	}

	// Get proxies in round-robin order
	for i := 0; i < len(proxies); i++ {
		proxy, err := rotator.GetProxy(ctx)
		if err != nil {
			log.Fatalf("Failed to get proxy: %v", err)
		}
		fmt.Printf("Round %d: Using proxy %s\n", i+1, proxy.URL)
	}
}

// RandomStrategyExample demonstrates using the random rotation strategy
func RandomStrategyExample() {
	opts := lashes.DefaultOptions()
	opts.Strategy = rotation.RandomStrategy

	rotator, err := lashes.New(opts)
	if err != nil {
		log.Fatalf("Failed to create rotator: %v", err)
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
			log.Fatalf("Failed to add proxy: %v", err)
		}
	}

	// Get proxies in random order
	for i := 0; i < len(proxies); i++ {
		proxy, err := rotator.GetProxy(ctx)
		if err != nil {
			log.Fatalf("Failed to get proxy: %v", err)
		}
		fmt.Printf("Round %d: Using proxy %s\n", i+1, proxy.URL)
	}
}