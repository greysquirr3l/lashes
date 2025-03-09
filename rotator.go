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

// BackoffStrategy defines how retry delays should be calculated
type BackoffStrategy interface {
	// NextDelay returns the delay to wait before the next retry attempt
	NextDelay(attempt int) time.Duration
}

// Remove unused code to fix linter warnings
// The commented out code could be reimplemented if needed in the future

// rotator is the implementation of the ProxyRotator interface
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

// GetProxy returns the next proxy according to the strategy
func (r *rotator) GetProxy(ctx context.Context) (*domain.Proxy, error) {
	proxies, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// Filter for only enabled proxies
	var enabledProxies []*domain.Proxy
	for _, proxy := range proxies {
		if proxy.GetEnabled() {
			enabledProxies = append(enabledProxies, proxy)
		}
	}

	if len(enabledProxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	proxy, err := r.strategy.Next(ctx, enabledProxies)
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
	// Use the proper URL parser
	parsedURL, err := mock.ParseURL(proxyURL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	// Create a pointer to the time.Time value
	now := time.Now()

	proxy := &domain.Proxy{
		ID:         uuid.New().String(),
		URL:        parsedURL.String(), // Store URL as string
		Type:       proxyType,
		Enabled:    true,
		LastUsed:   nil,
		MaxRetries: r.opts.MaxRetries,
		Timeout:    r.opts.RequestTimeout,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

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
		proxy.Latency = int64(latency.Milliseconds())
	}

	return r.repo.Create(ctx, proxy)
}

func (r *rotator) RemoveProxy(ctx context.Context, proxyURL string) error {
	proxies, err := r.repo.List(ctx)
	if err != nil {
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

// Helper method to get a proxy URL without exposing the proxy object
func (r *rotator) GetProxyURL() (string, error) {
	proxy, err := r.GetProxy(context.Background())
	if err != nil {
		return "", err
	}
	return proxy.URL, nil
}

// Implement the MetricsProvider interface methods
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
