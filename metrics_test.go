package lashes

import (
	"context"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
)

// mockRepository implements a simplified domain.ProxyRepository for testing
type mockRepository struct {
	proxies map[string]*domain.Proxy
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		proxies: make(map[string]*domain.Proxy),
	}
}

func (r *mockRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	r.proxies[proxy.ID] = proxy
	return nil
}

func (r *mockRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	proxy, found := r.proxies[id]
	if !found {
		return nil, repository.ErrProxyNotFound
	}
	return proxy, nil
}

func (r *mockRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	r.proxies[proxy.ID] = proxy
	return nil
}

func (r *mockRepository) Delete(ctx context.Context, id string) error {
	delete(r.proxies, id)
	return nil
}

func (r *mockRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	var proxies []*domain.Proxy
	for _, p := range r.proxies {
		proxies = append(proxies, p)
	}
	return proxies, nil
}

func (r *mockRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	for _, p := range r.proxies {
		return p, nil
	}
	return nil, repository.ErrProxyNotFound
}

func TestMetricsCollector(t *testing.T) {
	// Create a mock repository and add a test proxy
	repo := newMockRepository()
	ctx := context.Background()

	testProxy := &domain.Proxy{
		ID:      "test-proxy",
		URL:     "http://example.com:8080",
		Type:    domain.HTTP,
		Enabled: true,
	}

	if err := repo.Create(ctx, testProxy); err != nil {
		t.Fatalf("Failed to create test proxy: %v", err)
	}

	t.Run("Basic metrics collection", func(t *testing.T) {
		collector := NewMetricsCollector(repo)

		// Record some successful requests
		err := collector.RecordRequest(ctx, testProxy.ID, 100*time.Millisecond, true)
		if err != nil {
			t.Fatalf("RecordRequest failed: %v", err)
		}

		err = collector.RecordRequest(ctx, testProxy.ID, 200*time.Millisecond, true)
		if err != nil {
			t.Fatalf("RecordRequest failed: %v", err)
		}

		// Record a failed request
		err = collector.RecordRequest(ctx, testProxy.ID, 300*time.Millisecond, false)
		if err != nil {
			t.Fatalf("RecordRequest failed: %v", err)
		}

		// Get metrics for this proxy
		metrics, err := collector.GetProxyMetrics(ctx, testProxy.ID)
		if err != nil {
			t.Fatalf("GetProxyMetrics failed: %v", err)
		}

		// Check metrics values
		if metrics.TotalCalls != 3 {
			t.Errorf("TotalCalls = %d, want 3", metrics.TotalCalls)
		}

		if metrics.ErrorCount != 1 {
			t.Errorf("ErrorCount = %d, want 1", metrics.ErrorCount)
		}

		// Check success rate calculation
		expectedSuccessRate := float64(2) / float64(3) // 2 successes out of 3 calls
		if metrics.SuccessRate != expectedSuccessRate {
			t.Errorf("SuccessRate = %f, want %f", metrics.SuccessRate, expectedSuccessRate)
		}

		// Check average latency
		expectedAvgLatency := (100*time.Millisecond + 200*time.Millisecond + 300*time.Millisecond) / 3
		if metrics.AvgLatency != expectedAvgLatency {
			t.Errorf("AvgLatency = %v, want %v", metrics.AvgLatency, expectedAvgLatency)
		}
	})

	t.Run("Cached metrics collection", func(t *testing.T) {
		// Create a cached collector with a short expiration time
		collector := NewCachedMetricsCollector(repo, 50*time.Millisecond)

		// Record an initial request
		err := collector.RecordRequest(ctx, testProxy.ID, 100*time.Millisecond, true)
		if err != nil {
			t.Fatalf("RecordRequest failed: %v", err)
		}

		// Get metrics - should be fresh data
		metrics1, err := collector.GetProxyMetrics(ctx, testProxy.ID)
		if err != nil {
			t.Fatalf("GetProxyMetrics failed: %v", err)
		}

		// Record another request - should invalidate cache
		err = collector.RecordRequest(ctx, testProxy.ID, 200*time.Millisecond, true)
		if err != nil {
			t.Fatalf("RecordRequest failed: %v", err)
		}

		// Get metrics again - should be fresh data with updated counts
		metrics2, err := collector.GetProxyMetrics(ctx, testProxy.ID)
		if err != nil {
			t.Fatalf("GetProxyMetrics failed: %v", err)
		}

		if metrics1.TotalCalls >= metrics2.TotalCalls {
			t.Errorf("Expected metrics2.TotalCalls (%d) > metrics1.TotalCalls (%d)",
				metrics2.TotalCalls, metrics1.TotalCalls)
		}

		// Get metrics third time quickly - should use cache
		before := time.Now()
		metrics3, err := collector.GetProxyMetrics(ctx, testProxy.ID)
		cacheFetchTime := time.Since(before)

		if err != nil {
			t.Fatalf("GetProxyMetrics failed: %v", err)
		}

		// Make sure the cache retrieval is very fast
		if cacheFetchTime > 5*time.Millisecond {
			t.Logf("Cache fetch took longer than expected: %v", cacheFetchTime)
		}

		// Cached metrics should be identical
		if metrics2.TotalCalls != metrics3.TotalCalls {
			t.Errorf("Cache not working: metrics2.TotalCalls = %d, metrics3.TotalCalls = %d",
				metrics2.TotalCalls, metrics3.TotalCalls)
		}

		// Wait for cache to expire
		time.Sleep(60 * time.Millisecond)

		// Record one more request
		err = collector.RecordRequest(ctx, testProxy.ID, 300*time.Millisecond, true)
		if err != nil {
			t.Fatalf("RecordRequest failed: %v", err)
		}

		// Get metrics again - should be fresh data after cache expired
		metrics4, err := collector.GetProxyMetrics(ctx, testProxy.ID)
		if err != nil {
			t.Fatalf("GetProxyMetrics failed: %v", err)
		}

		if metrics3.TotalCalls >= metrics4.TotalCalls {
			t.Errorf("Expected metrics4.TotalCalls (%d) > metrics3.TotalCalls (%d) after cache expired",
				metrics4.TotalCalls, metrics3.TotalCalls)
		}
	})
}
