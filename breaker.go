package lashes

import (
	"sync"
	"time"

	"github.com/greysquirr3l/lashes/internal/breaker"
)

// CircuitBreakerConfig configures the circuit breaker behavior
type CircuitBreakerConfig struct {
	// MaxFailures is the threshold of failures before opening the circuit
	MaxFailures int

	// ResetTimeout is the time to wait before trying again
	ResetTimeout time.Duration

	// MaxHalfOpenRequests is the number of requests allowed in the testing state
	MaxHalfOpenRequests int

	// EnableGlobalBreaker enables a circuit breaker for the entire proxy pool
	EnableGlobalBreaker bool
}

// DefaultCircuitBreakerConfig returns sensible defaults for circuit breakers
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:         5,
		ResetTimeout:        time.Second * 30,
		MaxHalfOpenRequests: 1,
		EnableGlobalBreaker: true,
	}
}

// CircuitBreakerManager manages circuit breakers for proxies
type CircuitBreakerManager struct {
	breakers      map[string]*breaker.CircuitBreaker
	globalBreaker *breaker.CircuitBreaker
	config        CircuitBreakerConfig
	mu            sync.RWMutex
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(config CircuitBreakerConfig) *CircuitBreakerManager {
	mgr := &CircuitBreakerManager{
		breakers: make(map[string]*breaker.CircuitBreaker),
		config:   config,
	}

	if config.EnableGlobalBreaker {
		mgr.globalBreaker = breaker.NewCircuitBreaker(breaker.Config{
			MaxFailures:         config.MaxFailures * 3, // Higher threshold for global breaker
			ResetTimeout:        config.ResetTimeout,
			MaxHalfOpenRequests: config.MaxHalfOpenRequests,
		})
	}

	return mgr
}

// Allow checks if a request should be allowed through for a proxy
func (m *CircuitBreakerManager) Allow(proxyID string) bool {
	// Check global circuit breaker first
	if m.globalBreaker != nil && !m.globalBreaker.Allow() {
		return false
	}

	m.mu.RLock()
	cb, exists := m.breakers[proxyID]
	m.mu.RUnlock()

	if !exists {
		m.mu.Lock()
		// Double-check after acquiring write lock
		if cb, exists = m.breakers[proxyID]; !exists {
			// Create new circuit breaker if it doesn't exist
			cb = breaker.NewCircuitBreaker(breaker.Config{
				MaxFailures:         m.config.MaxFailures,
				ResetTimeout:        m.config.ResetTimeout,
				MaxHalfOpenRequests: m.config.MaxHalfOpenRequests,
			})
			m.breakers[proxyID] = cb
		}
		m.mu.Unlock()
	}

	return cb.Allow()
}

// RecordSuccess records a successful request for a proxy
func (m *CircuitBreakerManager) RecordSuccess(proxyID string) {
	if m.globalBreaker != nil {
		m.globalBreaker.RecordSuccess()
	}

	m.mu.RLock()
	cb, exists := m.breakers[proxyID]
	m.mu.RUnlock()

	if exists {
		cb.RecordSuccess()
	}
}

// RecordFailure records a failed request for a proxy
func (m *CircuitBreakerManager) RecordFailure(proxyID string) {
	if m.globalBreaker != nil {
		m.globalBreaker.RecordFailure()
	}

	m.mu.RLock()
	cb, exists := m.breakers[proxyID]
	m.mu.RUnlock()

	if exists {
		cb.RecordFailure()
	}
}

// EnableCircuitBreaker adds circuit breaker support to the rotator
func (r *rotator) EnableCircuitBreaker(config CircuitBreakerConfig) *CircuitBreakerManager {
	// Initialize the circuit breaker manager
	return NewCircuitBreakerManager(config)
}
