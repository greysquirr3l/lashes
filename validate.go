package lashes

import (
	"context"
	"fmt"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/validation"
)

// ValidateProxy validates a single proxy against a target URL
func (r *rotator) ValidateProxy(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error) {
	validator := validation.NewValidator(validation.Config{
		Timeout:    r.opts.ValidationTimeout,
		RetryCount: r.opts.MaxRetries,
		TestURL:    targetURL,
	})

	return validator.Validate(ctx, proxy)
}

// ValidateAll validates all proxies in the pool
func (r *rotator) ValidateAll(ctx context.Context) error {
	proxies, err := r.repo.List(ctx)
	if err != nil {
		return err
	}

	validator := validation.NewValidator(validation.Config{
		Timeout:    r.opts.ValidationTimeout,
		RetryCount: r.opts.MaxRetries,
		TestURL:    r.opts.TestURL,
	})

	for _, proxy := range proxies {
		valid, latency, err := validator.Validate(ctx, proxy)
		if err != nil {
			// Log error but continue with next proxy
			continue
		}

		// Update proxy status
		proxy.IsActive = valid
		if valid {
			proxy.Latency = int64(latency.Milliseconds())
			now := time.Now()
			proxy.LastCheck = &now
		} else {
			// Increment failure count or mark as inactive
		}

		// Record metrics
		if r.metrics != nil {
			if metricErr := r.metrics.RecordRequest(ctx, proxy.ID, latency, valid); metricErr != nil {
				// Log the error but continue with the validation process
				fmt.Printf("Failed to record metrics: %v\n", metricErr)
				// Alternatively, could use a more structured approach with a logger
				// r.logger.Error("Failed to record metrics", "error", metricErr)
			}
		}

		// Update the proxy in the repository
		if err := r.repo.Update(ctx, proxy); err != nil {
			// Log error but continue with next proxy
			continue
		}
	}

	return nil
}

// validateProxy validates a proxy against a test URL.
// This is an internal helper method used by the validation system.
// nolint:unused // Intentionally kept for API completeness
func (r *rotator) validateProxy(proxy *domain.Proxy) error {
	if (!proxy.Enabled) {
		return fmt.Errorf("proxy %s is disabled", proxy.ID)
	}

	valid, latency, err := r.ValidateProxy(context.Background(), proxy, r.opts.TestURL)
	if err != nil {
		return err
	}

	proxy.IsActive = valid
	proxy.Latency = int64(latency.Milliseconds()) // Convert time.Duration to int64
	
	now := time.Now() // Create a new time.Time value
	proxy.LastCheck = &now // Assign its address to LastCheck

	return nil
}
