package rotation_test

import (
	"context"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/rotation"
)

func TestValidRotationStrategies(t *testing.T) {
	testCases := []struct {
		name     string
		strategy rotation.StrategyType
		valid    bool
	}{
		{
			name:     "Round Robin Strategy",
			strategy: rotation.RoundRobinStrategy,
			valid:    true,
		},
		{
			name:     "Random Strategy",
			strategy: rotation.RandomStrategy,
			valid:    true,
		},
		{
			name:     "Weighted Strategy",
			strategy: rotation.WeightedStrategy,
			valid:    true,
		},
		{
			name:     "Least Used Strategy",
			strategy: rotation.LeastUsedStrategy,
			valid:    true,
		},
		{
			name:     "Invalid Strategy",
			strategy: rotation.StrategyType("invalid-strategy-name"),
			valid:    false,
		},
		{
			name:     "Empty Strategy",
			strategy: rotation.StrategyType(""),
			valid:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rotation.NewStrategy(tc.strategy)
			if tc.valid && err != nil {
				t.Errorf("Expected strategy %q to be valid, got error: %v", tc.strategy, err)
			}
			if !tc.valid && err == nil {
				t.Errorf("Expected strategy %q to be invalid, but no error was returned", tc.strategy)
			}
		})
	}
}

func TestStrategyImplementations(t *testing.T) {
	strategies := []rotation.StrategyType{
		rotation.RoundRobinStrategy,
		rotation.RandomStrategy,
		rotation.WeightedStrategy,
		rotation.LeastUsedStrategy,
	}

	for _, strategyType := range strategies {
		t.Run(string(strategyType), func(t *testing.T) {
			strategy, err := rotation.NewStrategy(strategyType)
			if err != nil {
				t.Fatalf("Failed to create strategy %q: %v", strategyType, err)
			}

			ctx := context.Background()
			
			// Test with empty list
			_, err = strategy.Next(ctx, nil)
			if err == nil {
				t.Error("Expected error for empty proxy list, got nil")
			}

			// Test with single proxy
			proxyID := "test-proxy-id"
			proxy := &domain.Proxy{ID: proxyID}
			result, err := strategy.Next(ctx, []*domain.Proxy{proxy})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result.ID != proxyID {
				t.Errorf("Expected proxy ID %q, got %q", proxyID, result.ID)
			}
		})
	}
}

func TestRoundRobinStrategy(t *testing.T) {
	strategy, _ := rotation.NewStrategy(rotation.RoundRobinStrategy)
	ctx := context.Background()
	
	proxies := []*domain.Proxy{
		{ID: "proxy1"},
		{ID: "proxy2"},
		{ID: "proxy3"},
	}
	
	// The strategy should iterate through all proxies in sequence
	for i := 0; i < len(proxies)*2; i++ {
		proxy, err := strategy.Next(ctx, proxies)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := proxies[i%len(proxies)].ID
		if proxy.ID != expected {
			t.Errorf("Iteration %d: Expected proxy ID %q, got %q", i, expected, proxy.ID)
		}
	}
}

func TestWeightedStrategy(t *testing.T) {
	strategy, _ := rotation.NewStrategy(rotation.WeightedStrategy)
	ctx := context.Background()
	
	proxies := []*domain.Proxy{
		{
			ID: "proxy1",
			Weight: 1,
			SuccessRate: 0.9,
			UsageCount: 100,
		},
		{
			ID: "proxy2",
			Weight: 5,
			SuccessRate: 0.5,
			UsageCount: 100,
		},
		{
			ID: "proxy3",
			Weight: 2,
			SuccessRate: 0.1,
			UsageCount: 100,
		},
	}
	
	// Run multiple selections to verify distribution
	counts := make(map[string]int)
	iterations := 1000
	
	for i := 0; i < iterations; i++ {
		proxy, err := strategy.Next(ctx, proxies)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		counts[proxy.ID]++
	}
	
	// Check if proxy2 (high weight) is selected more often than proxy3 (low weight)
	// This is statistical, so we only check general trends
	if counts["proxy2"] <= counts["proxy3"] {
		t.Errorf("Expected proxy2 (higher weight) to be selected more than proxy3 (lower weight): got proxy2=%d, proxy3=%d", 
			counts["proxy2"], counts["proxy3"])
	}
	
	// Test with one zero-weight proxy and one regular
	zeroWeightProxy := []*domain.Proxy{
		{
			ID: "zero",
			Weight: 0,
			SuccessRate: 0,
			UsageCount: 0,
		},
		{
			ID: "normal",
			Weight: 1,
			SuccessRate: 0.5,
			UsageCount: 10,
		},
	}
	
	zeroCount := 0
	normalCount := 0
	
	for i := 0; i < 100; i++ {
		proxy, _ := strategy.Next(ctx, zeroWeightProxy)
		if proxy.ID == "zero" {
			zeroCount++
		} else {
			normalCount++
		}
	}
	
	// Zero weight proxy should be selected less often
	if zeroCount >= normalCount {
		t.Errorf("Expected zero-weight proxy to be selected less often: got zero=%d, normal=%d", 
			zeroCount, normalCount)
	}
}

func TestLeastUsedStrategy(t *testing.T) {
	strategy, _ := rotation.NewStrategy(rotation.LeastUsedStrategy)
	ctx := context.Background()
	
	now := time.Now()
	earlier := now.Add(-time.Hour)
	
	// Create proxies with different usage counts
	proxies := []*domain.Proxy{
		{
			ID: "high-usage",
			UsageCount: 100,
			LastUsed: &now,
		},
		{
			ID: "medium-usage",
			UsageCount: 50,
			LastUsed: &now,
		},
		{
			ID: "low-usage",
			UsageCount: 10,
			LastUsed: &now,
		},
	}
	
	// The least used proxy should be selected
	proxy, err := strategy.Next(ctx, proxies)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if proxy.ID != "low-usage" {
		t.Errorf("Expected least used proxy (low-usage), got %s", proxy.ID)
	}
	
	// Test same usage count but different last used time
	sameCountProxies := []*domain.Proxy{
		{
			ID: "older",
			UsageCount: 10,
			LastUsed: &earlier,
		},
		{
			ID: "newer",
			UsageCount: 10,
			LastUsed: &now,
		},
	}
	
	proxy, err = strategy.Next(ctx, sameCountProxies)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if proxy.ID != "older" {
		t.Errorf("Expected older proxy to be selected when usage counts are equal, got %s", proxy.ID)
	}
	
	// Test with nil LastUsed
	nilLastUsedProxies := []*domain.Proxy{
		{
			ID: "with-last-used",
			UsageCount: 10,
			LastUsed: &now,
		},
		{
			ID: "nil-last-used",
			UsageCount: 10,
			LastUsed: nil,
		},
	}
	
	proxy, err = strategy.Next(ctx, nilLastUsedProxies)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if proxy.ID != "nil-last-used" {
		t.Errorf("Expected proxy with nil LastUsed to be selected, got %s", proxy.ID)
	}
}
