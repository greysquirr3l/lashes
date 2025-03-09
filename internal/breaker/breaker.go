// Package breaker implements the circuit breaker pattern for proxy resilience.
package breaker

import (
	"sync"
	"time"
)

// State represents the current state of a circuit breaker
type State int

const (
	// StateClosed is the normal operating state (requests flow normally)
	StateClosed State = iota
	// StateOpen is the failure state (all requests are blocked)
	StateOpen
	// StateHalfOpen is the testing state (allows a limited number of requests)
	StateHalfOpen
)

// Config defines circuit breaker behavior
type Config struct {
	// MaxFailures is the threshold of failures before opening the circuit
	MaxFailures int
	// ResetTimeout is the time to wait before transitioning from Open to HalfOpen
	ResetTimeout time.Duration
	// MaxHalfOpenRequests is the number of requests allowed in the HalfOpen state
	MaxHalfOpenRequests int
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() Config {
	return Config{
		MaxFailures:         5,
		ResetTimeout:        time.Second * 30,
		MaxHalfOpenRequests: 1,
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config Config
	state  State

	failures        int
	lastStateChange time.Time
	halfOpenCount   int

	mu sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker with the given config
func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	return &CircuitBreaker{
		config:          cfg,
		state:           StateClosed,
		lastStateChange: time.Now(),
	}
}

// Allow returns whether a request should be permitted
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		// Check if it's time to go to half-open state
		if time.Since(cb.lastStateChange) > cb.config.ResetTimeout {
			cb.state = StateHalfOpen
			cb.lastStateChange = time.Now()
			cb.halfOpenCount = 0 // Start with 0 so we can increment before returning
			cb.halfOpenCount++
			return true
		}
		return false
	case StateHalfOpen:
		// Only allow limited requests in half-open state
		if cb.halfOpenCount < cb.config.MaxHalfOpenRequests {
			cb.halfOpenCount++
			return true
		}
		return false
	default:
		return false
	}
}

// RecordSuccess records a successful request
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		// Reset failures counter
		cb.failures = 0
	case StateHalfOpen:
		// If successful in half-open, close the circuit immediately
		cb.state = StateClosed
		cb.failures = 0
		cb.lastStateChange = time.Now()
		cb.halfOpenCount = 0
	case StateOpen:
		// No action needed for StateOpen
	}
}

// RecordFailure records a failed request
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		// Increment failure counter
		cb.failures++
		if cb.failures >= cb.config.MaxFailures {
			cb.state = StateOpen
			cb.lastStateChange = time.Now()
		}
	case StateHalfOpen:
		// If failed in half-open, go back to open immediately
		cb.state = StateOpen
		cb.failures = cb.config.MaxFailures // Reset to max failures
		cb.lastStateChange = time.Now()
		cb.halfOpenCount = 0
	case StateOpen:
		// No action needed for StateOpen
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}
