package examples

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// MetricsExample demonstrates how to access proxy performance metrics
func MetricsExample() {
	// Initialize with default options
	rotator, err := lashes.New(lashes.DefaultOptions())
	if err != nil {
		log.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()
	
	// Add some proxies
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

	// Make some sample requests to generate metrics
	client, err := rotator.Client(ctx)
	if err != nil {
		log.Fatalf("Failed to get HTTP client: %v", err)
	}

	// Generate some sample traffic
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := client.Get("https://example.com")
		duration := time.Since(start)
		
		if err != nil {
			fmt.Printf("Request %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Request %d succeeded in %v\n", i+1, duration)
		}
		
		// Get a new client for the next request to use a different proxy
		client, err = rotator.Client(ctx)
		if err != nil {
			log.Fatalf("Failed to get HTTP client: %v", err)
		}
	}
	
	// Wait a moment for metrics to be processed
	time.Sleep(100 * time.Millisecond)
	
	// Get metrics for all proxies
	metrics, err := rotator.GetAllMetrics(ctx)
	if err != nil {
		log.Fatalf("Failed to get metrics: %v", err)
	}
	
	fmt.Printf("\nProxy performance metrics:\n")
	fmt.Printf("%-30s %-10s %-10s %-10s %-10s\n", 
		"URL", "Success %", "Calls", "Avg ms", "Errors")
	fmt.Println(strings.Repeat("-", 70))
	
	for _, m := range metrics {
		fmt.Printf("%-30s %-10.1f %-10d %-10.2f %-10d\n",
			m.URL,
			m.SuccessRate * 100,
			m.TotalCalls,
			float64(m.AvgLatency) / float64(time.Millisecond),
			m.ErrorCount)
	}
}
