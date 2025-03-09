package benchmarks

import (
	"context"
	"fmt"
	"testing"

	"github.com/greysquirr3l/lashes"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/rotation"
)

func setupRotator(b *testing.B, strategy rotation.StrategyType) lashes.ProxyRotator {
	opts := lashes.DefaultOptions()
	opts.Strategy = strategy

	rotator, err := lashes.New(opts)
	if err != nil {
		b.Fatalf("Failed to create rotator: %v", err)
	}

	ctx := context.Background()

	// Add a bunch of test proxies
	for i := 0; i < 100; i++ {
		proxyURL := fmt.Sprintf("http://proxy%d.example.com:8080", i)
		if err := rotator.AddProxy(ctx, proxyURL, domain.HTTPProxy); err != nil {
			b.Fatalf("Failed to add proxy: %v", err)
		}
	}

	return rotator
}

func BenchmarkRoundRobinStrategy(b *testing.B) {
	rotator := setupRotator(b, rotation.RoundRobinStrategy)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rotator.GetProxy(ctx)
		if err != nil {
			b.Fatalf("Failed to get proxy: %v", err)
		}
	}
}

func BenchmarkRandomStrategy(b *testing.B) {
	rotator := setupRotator(b, rotation.RandomStrategy)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rotator.GetProxy(ctx)
		if err != nil {
			b.Fatalf("Failed to get proxy: %v", err)
		}
	}
}

func BenchmarkLeastUsedStrategy(b *testing.B) {
	rotator := setupRotator(b, rotation.LeastUsedStrategy)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rotator.GetProxy(ctx)
		if err != nil {
			b.Fatalf("Failed to get proxy: %v", err)
		}
	}
}
