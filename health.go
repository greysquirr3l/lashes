package lashes

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthCheckOptions configures health check behavior
type HealthCheckOptions struct {
	// Interval between health checks
	Interval time.Duration

	// Timeout for each health check
	Timeout time.Duration

	// HealthURL is the URL used for health checks
	HealthURL string

	// MaxFailures is the number of failures before marking a proxy as inactive
	MaxFailures int

	// Parallel is the number of concurrent health checks
	Parallel int
}

// DefaultHealthCheckOptions returns sensible default options for health checking
func DefaultHealthCheckOptions() HealthCheckOptions {
	return HealthCheckOptions{
		Interval:    time.Minute * 10,
		Timeout:     time.Second * 5,
		HealthURL:   "https://api.ipify.org?format=json",
		MaxFailures: 3,
		Parallel:    10,
	}
}

// StartHealthCheck starts periodic health checking of all proxies
func (r *rotator) StartHealthCheck(ctx context.Context, opts HealthCheckOptions) error {
	go func() {
		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := r.performHealthCheck(ctx, opts); err != nil {
					// Log error or handle it if needed
					// But don't exit the goroutine
				}
			}
		}
	}()

	return nil
}

// performHealthCheck runs a health check on all proxies
func (r *rotator) performHealthCheck(ctx context.Context, opts HealthCheckOptions) error {
	proxies, err := r.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list proxies: %w", err)
	}

	if len(proxies) == 0 {
		return nil // No proxies to check, not an error
	}

	// Use a wait group to limit concurrency
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, opts.Parallel)

	// Track errors
	var errMu sync.Mutex
	var errors []error

	for _, proxy := range proxies {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(proxy *Proxy) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Create context with timeout
			checkCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
			defer cancel()

			// Check the proxy health
			valid, _, err := r.ValidateProxy(checkCtx, proxy, opts.HealthURL)
			if err != nil {
				// If validation fails, consider the proxy invalid
				valid = false
			}

			// Update proxy status if needed
			if proxy.Enabled != valid {
				proxy.SetEnabled(valid) // This updates both Enabled and IsActive

				// Update the proxy in the repository
				updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
				defer updateCancel()

				if updateErr := r.repo.Update(updateCtx, proxy); updateErr != nil {
					errMu.Lock()
					errors = append(errors, fmt.Errorf("failed to update proxy %s: %w", proxy.ID, updateErr))
					errMu.Unlock()
				}
			}
		}(proxy)
	}

	wg.Wait()

	// Handle any errors
	if len(errors) > 0 {
		return fmt.Errorf("health check completed with %d errors", len(errors))
	}

	return nil
}

// GetHealthStatus returns the current health status of all proxies
func (r *rotator) GetHealthStatus(ctx context.Context) (map[string]bool, error) {
	proxies, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	status := make(map[string]bool, len(proxies))
	for _, proxy := range proxies {
		status[proxy.ID] = proxy.Enabled
	}

	return status, nil
}
