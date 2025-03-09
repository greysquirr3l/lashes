package breaker_test

import (
	"testing"
	"time"

	"github.com/greysquirr3l/lashes/internal/breaker"
)

func TestCircuitBreaker(t *testing.T) {
	t.Run("Initial state is closed", func(t *testing.T) {
		cb := breaker.NewCircuitBreaker(breaker.DefaultConfig())
		if !cb.Allow() {
			t.Error("New circuit breaker should allow requests")
		}
		if state := cb.GetState(); state != breaker.StateClosed {
			t.Errorf("Initial state = %v, want %v", state, breaker.StateClosed)
		}
	})

	t.Run("Opens after max failures", func(t *testing.T) {
		cb := breaker.NewCircuitBreaker(breaker.Config{
			MaxFailures:         3,
			ResetTimeout:        time.Second * 30,
			MaxHalfOpenRequests: 1,
		})

		// Record failures until breaker opens
		for i := 0; i < 3; i++ {
			if !cb.Allow() {
				t.Errorf("Allow() on attempt %d = false, want true", i)
			}
			cb.RecordFailure()
		}

		// Breaker should be open now
		if cb.Allow() {
			t.Error("Allow() = true, want false after max failures")
		}
		if state := cb.GetState(); state != breaker.StateOpen {
			t.Errorf("State after failures = %v, want %v", state, breaker.StateOpen)
		}
	})

	t.Run("Success in closed state resets failure counter", func(t *testing.T) {
		cb := breaker.NewCircuitBreaker(breaker.Config{
			MaxFailures:         3,
			ResetTimeout:        time.Second * 30,
			MaxHalfOpenRequests: 1,
		})

		// Record some failures but not enough to open
		for i := 0; i < 2; i++ {
			cb.RecordFailure()
		}

		// Record a success which should reset counter
		cb.RecordSuccess()

		// Record more failures - should still be below threshold
		for i := 0; i < 2; i++ {
			cb.RecordFailure()
		}

		// Breaker should still be closed
		if !cb.Allow() {
			t.Error("Allow() = false, want true after failure count reset")
		}
	})

	t.Run("Transitions to half-open after timeout", func(t *testing.T) {
		cb := breaker.NewCircuitBreaker(breaker.Config{
			MaxFailures:         2,
			ResetTimeout:        time.Millisecond * 50, // Short timeout for testing
			MaxHalfOpenRequests: 1,
		})

		// Open the circuit
		cb.RecordFailure()
		cb.RecordFailure()

		if cb.Allow() {
			t.Error("Allow() = true, want false when circuit is open")
		}

		// Wait for timeout
		time.Sleep(time.Millisecond * 60)

		// Should be in half-open state now
		if !cb.Allow() {
			t.Error("Allow() = false, want true after timeout (half-open state)")
		}

		// Should only allow one request in half-open state
		if cb.Allow() {
			t.Error("Allow() = true, want false for second request in half-open state")
		}
	})

	t.Run("Success in half-open state closes circuit", func(t *testing.T) {
		cb := breaker.NewCircuitBreaker(breaker.Config{
			MaxFailures:         2,
			ResetTimeout:        time.Millisecond * 50,
			MaxHalfOpenRequests: 1,
		})

		// First open the circuit
		cb.RecordFailure()
		cb.RecordFailure()

		// Verify it's open
		if cb.GetState() != breaker.StateOpen {
			t.Errorf("Expected circuit to be open after failures")
		}

		// Wait for timeout to transition to half-open
		time.Sleep(time.Millisecond * 60)

		// Now we need to actually attempt a request to trigger the transition
		allowed := cb.Allow()
		if !allowed {
			t.Error("Allow() = false, want true after timeout period")
		}

		// Verify state changed to half-open
		if state := cb.GetState(); state != breaker.StateHalfOpen {
			t.Errorf("State before success = %v, want %v", state, breaker.StateHalfOpen)
			return // Skip the rest of the test
		}

		// Record success on the test request
		cb.RecordSuccess()

		// Circuit should be closed now
		if state := cb.GetState(); state != breaker.StateClosed {
			t.Errorf("State = %v, want %v after success in half-open", state, breaker.StateClosed)
		}

		// Should allow requests now
		if !cb.Allow() {
			t.Error("Allow() = false, want true after success in half-open state")
		}
	})

	t.Run("Failure in half-open state reopens circuit", func(t *testing.T) {
		cb := breaker.NewCircuitBreaker(breaker.Config{
			MaxFailures:         2,
			ResetTimeout:        time.Millisecond * 50,
			MaxHalfOpenRequests: 1,
		})

		// First open the circuit
		cb.RecordFailure()
		cb.RecordFailure()

		// Verify it's open
		if cb.GetState() != breaker.StateOpen {
			t.Errorf("Expected circuit to be open after failures")
		}

		// Wait for timeout to transition to half-open
		time.Sleep(time.Millisecond * 60)

		// Need to call Allow() to actually transition to half-open
		allowed := cb.Allow()
		if !allowed {
			t.Error("Allow() = false, want true after timeout period")
		}

		// Verify state changed to half-open
		if state := cb.GetState(); state != breaker.StateHalfOpen {
			t.Errorf("State before failure = %v, want %v", state, breaker.StateHalfOpen)
			return // Skip the rest of the test
		}

		// Record failure in half-open state
		cb.RecordFailure()

		// Check that we're back in open state
		if state := cb.GetState(); state != breaker.StateOpen {
			t.Errorf("State = %v, want %v after failure in half-open", state, breaker.StateOpen)
		}

		// Should not allow requests
		if cb.Allow() {
			t.Error("Allow() = true, want false after failure in half-open state")
		}
	})
}
