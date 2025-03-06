package lashes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
	"github.com/greysquirr3l/lashes/internal/repository/gorm"
	"github.com/greysquirr3l/lashes/internal/rotation"
	"github.com/greysquirr3l/lashes/internal/validation"
)

type rotator struct {
	repo     domain.ProxyRepository
	strategy rotation.Strategy
	opts     Options
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

	return &rotator{
		repo:     repo,
		strategy: strategy,
		opts:     opts,
	}, nil
}

func (r *rotator) GetProxy(ctx context.Context) (*domain.Proxy, error) {
	proxies, err := r.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	proxy, err := r.strategy.Next(ctx, proxies)
	if err != nil {
		return nil, err
	}

	// Update last used timestamp
	proxy.LastUsed = time.Now()
	if err := r.repo.Update(ctx, proxy); err != nil {
		return nil, err
	}

	return proxy, nil
}

func (r *rotator) AddProxy(ctx context.Context, proxyURL string, proxyType domain.ProxyType) error {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	proxy := &domain.Proxy{
		ID:         uuid.New().String(),
		URL:        parsedURL,
		Type:       proxyType,
		IsActive:   true,
		LastCheck:  time.Now(),
		MaxRetries: r.opts.MaxRetries,
		Timeout:    r.opts.RequestTimeout,
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
		proxy.Latency = latency
	}

	return r.repo.Create(ctx, proxy)
}

func (r *rotator) RemoveProxy(ctx context.Context, proxyURL string) error {
	proxies, err := r.repo.List(ctx)
	if err != nil {
		return err
	}

	for _, proxy := range proxies {
		if proxy.URL.String() == proxyURL {
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
