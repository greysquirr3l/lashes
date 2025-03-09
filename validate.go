package lashes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/validation"
)

// ValidateProxy validates a single proxy against a target URL
func (r *rotator) ValidateProxy(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error) {
	// Create a proper timeout context if not already set
	ctx, cancel := context.WithTimeout(ctx, r.opts.ValidationTimeout)
	defer cancel()

	validator := validation.NewValidator(validation.Config{
		Timeout:    r.opts.ValidationTimeout,
		RetryCount: r.opts.MaxRetries,
		TestURL:    targetURL,
	})

	return validator.Validate(ctx, proxy)
}

// ValidateAll validates all proxies in the pool
func (r *rotator) ValidateAll(ctx context.Context) error {
	// Split function to reduce complexity
	proxies, err := r.getProxiesForValidation(ctx)
	if err != nil {
		return err
	}

	validator := validation.NewValidator(validation.Config{
		Timeout:    r.opts.ValidationTimeout,
		RetryCount: r.opts.MaxRetries,
		TestURL:    r.opts.TestURL,
	})

	var validationErrors []error
	for _, proxy := range proxies {
		if err := r.validateSingleProxy(ctx, proxy, validator, &validationErrors); err != nil {
			return err
		}
	}

	// If we had validation errors, return a combined error
	if len(validationErrors) > 0 {
		return fmt.Errorf("validation completed with %d errors: %w", len(validationErrors), errors.Join(validationErrors...))
	}

	return nil
}

// getProxiesForValidation gets the list of proxies to validate
func (r *rotator) getProxiesForValidation(ctx context.Context) ([]*domain.Proxy, error) {
	// Use the context directly instead of assigning it to a new variable
	// This fixes the "Non-inherited new context" warning

	// Create a timeout context derived from the parent context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Use the timeout context for the repository call
	proxies, err := r.repo.List(ctxWithTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to list proxies: %w", err)
	}

	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	return proxies, nil
}

// validateSingleProxy validates a single proxy
func (r *rotator) validateSingleProxy(
	ctx context.Context,
	proxy *domain.Proxy,
	validator validation.Validator,
	validationErrors *[]error,
) error {
	// Skip validation if context is done
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Continue validation
	}

	// Create a sub-context for this validation that inherits from ctx
	proxyCtx, cancel := context.WithTimeout(ctx, r.opts.ValidationTimeout)
	defer cancel()

	valid, latency, err := validator.Validate(proxyCtx, proxy)

	// Update proxy status
	proxy.SetEnabled(valid)

	if valid {
		proxy.Latency = int64(latency.Milliseconds())
		// Update timestamp in a proxy metadata field if needed
		// For now, we just record the success
	} else if err != nil {
		*validationErrors = append(*validationErrors, NewValidationError(
			proxy.ID,
			proxy.URL,
			err.Error(),
			0,
		))
	}

	// Record metrics and update repository
	return r.recordValidationResults(ctx, proxy, latency, valid, validationErrors)
}

// recordValidationResults records metrics and updates the repository
func (r *rotator) recordValidationResults(
	ctx context.Context,
	proxy *domain.Proxy,
	latency time.Duration,
	valid bool,
	validationErrors *[]error,
) error {
	// Record metrics
	if r.metrics != nil {
		metricCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		metricErr := r.metrics.RecordRequest(metricCtx, proxy.ID, latency, valid)
		cancel()

		if metricErr != nil {
			// Log error but continue with validation process
			*validationErrors = append(*validationErrors,
				fmt.Errorf("metrics recording for proxy %s: %w", proxy.ID, metricErr))
		}
	}

	// Update the proxy in the repository
	updateCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.repo.Update(updateCtx, proxy); err != nil {
		*validationErrors = append(*validationErrors,
			fmt.Errorf("failed to update proxy %s: %w", proxy.ID, err))
		return nil // Continue with other proxies
	}

	return nil
}
