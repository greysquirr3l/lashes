package lashes

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter provides rate limiting functionality for proxies
type RateLimiter interface {
	// Wait blocks until the rate limit allows an event to happen
	Wait(ctx context.Context) error

	// Allow reports whether an event may happen now
	Allow() bool

	// SetLimit updates the rate limit
	SetLimit(limit rate.Limit)

	// SetBurst updates the burst limit
	SetBurst(burst int)
}

// ProxyRateLimiter manages rate limiting for individual proxies
type ProxyRateLimiter struct {
	limiters sync.Map // map[string]*rate.Limiter
	defaults struct {
		limit rate.Limit
		burst int
	}
}

// NewProxyRateLimiter creates a new rate limiter for proxies
func NewProxyRateLimiter(requestsPerSecond float64, burst int) *ProxyRateLimiter {
	prl := &ProxyRateLimiter{}
	prl.defaults.limit = rate.Limit(requestsPerSecond)
	prl.defaults.burst = burst
	return prl
}

// GetLimiter returns the rate limiter for a specific proxy
func (prl *ProxyRateLimiter) GetLimiter(proxyID string) *rate.Limiter {
	// Get or create limiter for this proxy
	if limiter, ok := prl.limiters.Load(proxyID); ok {
		if rateLimiter, ok := limiter.(*rate.Limiter); ok {
			return rateLimiter
		}
		// If type assertion fails, create a new limiter
	}

	// Create new limiter with default settings
	limiter := rate.NewLimiter(prl.defaults.limit, prl.defaults.burst)
	prl.limiters.Store(proxyID, limiter)
	return limiter
}

// Wait blocks until the rate limit for a proxy allows an event to happen
func (prl *ProxyRateLimiter) Wait(ctx context.Context, proxyID string) error {
	return prl.GetLimiter(proxyID).Wait(ctx)
}

// Allow reports whether an event may happen for a proxy now
func (prl *ProxyRateLimiter) Allow(proxyID string) bool {
	return prl.GetLimiter(proxyID).Allow()
}

// SetProxyLimit updates the rate limit for a specific proxy
func (prl *ProxyRateLimiter) SetProxyLimit(proxyID string, limit rate.Limit, burst int) {
	limiter := prl.GetLimiter(proxyID)
	limiter.SetLimit(limit)
	limiter.SetBurst(burst)
}

// UseRateLimit applies rate limiting to a rotator
func (r *rotator) UseRateLimit(requestsPerSecond float64, burst int) *ProxyRateLimiter {
	rateLimiter := NewProxyRateLimiter(requestsPerSecond, burst)
	return rateLimiter
}
