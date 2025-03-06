package lashes

import (
	"context"
	"time"

	"github.com/greysquirr3l/lashes/internal/validation"
)

// ValidateProxy validates a single proxy against a target URL
func (r *rotator) ValidateProxy(ctx context.Context, proxy *Proxy, targetURL string) (bool, time.Duration, error) {
	validator := validation.NewValidator(validation.Config{
		Timeout:    r.opts.ValidationTimeout,
		RetryCount: r.opts.MaxRetries,
		TestURL:    targetURL,
	})

	return validator.ValidateWithTarget(ctx, proxy, targetURL)
}

// ValidateAll validates all proxies in the pool
func (r *rotator) ValidateAll(ctx context.Context) error {
	proxies, err := r.List(ctx)
	if err != nil {
		return err
	}

	for _, proxy := range proxies {
		valid, latency, err := r.ValidateProxy(ctx, proxy, r.opts.TestURL)
		if err != nil {
			proxy.IsActive = false
		} else {
			proxy.IsActive = valid
			proxy.Latency = latency
			proxy.LastCheck = time.Now()
		}

		if err := r.repo.Update(ctx, proxy); err != nil {
			return err
		}
	}

	return nil
}
