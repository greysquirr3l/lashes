package lashes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/client/mock"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
	"github.com/greysquirr3l/lashes/internal/repository/gorm"
	"github.com/greysquirr3l/lashes/internal/rotation"
	"github.com/greysquirr3l/lashes/internal/validation"
)

// Add metrics to the rotator implementation
type rotator struct {
	repo     domain.ProxyRepository
	strategy rotation.Strategy
	opts     Options
	metrics  MetricsCollector
}

func newRotator(opts Options) (*rotator, error) {
	var repo domain.ProxyRepository
	var err error

	// Initialize storage
	if opts.Storage == nil {
		repo = repository.NewMemoryRepository()
	} else {
		// Only initialize database if explicitly requested
		db, err := gorm.NewDB(*opts.Storage)
		if err != nil {
			return nil, err
		}

		repo = gorm.NewProxyRepository(db, gorm.Options{
			QueryTimeout: opts.Storage.QueryTimeout,
		})
	}

	strategy, err := rotation.NewStrategy(opts.Strategy)
	if err != nil {
		return nil, err
	}

	r := &rotator{
		repo:     repo,
		strategy: strategy,
		opts:     opts,
		metrics:  NewMetricsCollector(repo),
	}

	return r, nil
}

// Update the GetProxy method to return ErrNoProxiesAvailable when no proxies are available
func (r *rotator) GetProxy(ctx context.Context) (*domain.Proxy, error) {
	proxies, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	// Check if there are any proxies before trying to get one
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// Get the next proxy according to the strategy
	proxy, err := r.strategy.Next(ctx, proxies)
	if err != nil {
		return nil, err
	}

	// Update last used timestamp
	now := time.Now()
	proxy.LastUsed = &now
	if err := r.repo.Update(ctx, proxy); err != nil {
		return nil, err
	}

	return proxy, nil
}

func (r *rotator) AddProxy(ctx context.Context, proxyURL string, proxyType domain.ProxyType) error {
	// Use the mock URL parser instead of directly using url.Parse
	parsedURL, err := mock.ParseURL(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	// Create a pointer to the time.Time value
	now := time.Now()

	proxy := &domain.Proxy{
		ID:        uuid.New().String(),
		URL:       parsedURL.String(), // Convert URL to string
		Type:      proxyType,
		Enabled:   true,               // Use Enabled instead of IsActive
		LastUsed:  nil,                // Initialize as nil
		LastCheck: &now,               // Use pointer to time value
		MaxRetries: r.opts.MaxRetries,
		Timeout:    r.opts.RequestTimeout,
	}

	// Also set the backwards compatibility fields
	proxy.IsActive = proxy.Enabled

	if r.opts.ValidateOnStart {
		validator := validation.NewValidator(validation.Config{
			Timeout:    r.opts.ValidationTimeout,
			RetryCount: r.opts.MaxRetries,
			TestURL:    r.opts.TestURL,
		})

		valid, latency, err := validator.Validate(ctx, proxy)
		if err != nil {
			return err
		}
		if !valid {
			return fmt.Errorf("proxy validation failed")
		}
		proxy.Latency = int64(latency.Milliseconds()) // Convert time.Duration to int64 milliseconds
	}

	return r.repo.Create(ctx, proxy)
}

func (r *rotator) RemoveProxy(ctx context.Context, proxyURL string) error {
	proxies, err := r.repo.List(ctx)
	if (err != nil) {
		return err
	}

	for _, proxy := range proxies {
		if proxy.URL == proxyURL {
			return r.repo.Delete(ctx, proxy.ID)
		}
	}

	return repository.ErrProxyNotFound
}

func (r *rotator) Client(ctx context.Context) (*http.Client, error) {
	proxy, err := r.GetProxy(ctx)
	if err != nil {
		return nil, err
	}

	return client.NewClient(proxy, client.Options{
		Timeout:         r.opts.RequestTimeout,
		MaxRetries:      r.opts.MaxRetries,
		VerifyCerts:     true,
		FollowRedirects: true,
	})
}

func (r *rotator) List(ctx context.Context) ([]*domain.Proxy, error) {
	return r.repo.List(ctx)
}

func (r *rotator) GetProxyURL() (string, error) {
	proxy, err := r.GetProxy(context.Background())
	if err != nil {
		return "", err
	}
	return proxy.URL, nil // Return URL directly as it's already a string
}

// updateLastUsed updates the last used timestamp of a proxy.
// This method is part of the internal API and used by higher-level functions.
// nolint:unused // Intentionally kept for API completeness
func (r *rotator) updateLastUsed(proxyID string) error {
	ctx := context.Background()
	proxy, err := r.repo.GetByID(ctx, proxyID) // Use GetByID instead of Get
	if err != nil {
		return err
	}

	now := time.Now()
	proxy.LastUsed = &now
	return r.repo.Update(ctx, proxy)
}

// recordLatency updates the latency measurement for a proxy.
// This method is part of the internal API and used by higher-level functions.
// nolint:unused // Intentionally kept for API completeness
func (r *rotator) recordLatency(proxyID string, latency time.Duration) error {
	ctx := context.Background()
	proxy, err := r.repo.GetByID(ctx, proxyID) // Use GetByID instead of Get
	if err != nil {
		return err
	}

	proxy.Latency = int64(latency.Milliseconds()) // Convert time.Duration to int64 milliseconds
	if err := r.repo.Update(ctx, proxy); err != nil {
		return err
	}

	// Add metrics recording
	if r.metrics != nil {
		if err := r.metrics.RecordRequest(ctx, proxyID, latency, true); err != nil {
			// Just log the error, don't fail the operation
			// Consider adding a logger interface to the project
		}
	}

	return nil
}

// Add new methods to expose metrics functionality
func (r *rotator) GetProxyMetrics(ctx context.Context, proxyID string) (*ProxyMetrics, error) {
	if r.metrics == nil {
		return nil, ErrMetricsNotEnabled
	}
	return r.metrics.GetProxyMetrics(ctx, proxyID)
}

func (r *rotator) GetAllMetrics(ctx context.Context) ([]*ProxyMetrics, error) {
	if r.metrics == nil {
		return nil, ErrMetricsNotEnabled
	}
	return r.metrics.GetAllMetrics(ctx)
}
