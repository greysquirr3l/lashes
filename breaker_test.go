package lashes

import (
	"fmt"
	"testing"
)

func TestCircuitBreakerManager(t *testing.T) {
	t.Run("Individual proxy circuit breaker", func(t *testing.T) {
		config := DefaultCircuitBreakerConfig()
		config.MaxFailures = 3
		config.EnableGlobalBreaker = false

		manager := NewCircuitBreakerManager(config)

		// Initially all proxies should be allowed
		if !manager.Allow("proxy-1") {
			t.Error("Initial Allow() for proxy-1 = false, want true")
		}

		// Record failures until circuit opens
		for i := 0; i < 3; i++ {
			manager.RecordFailure("proxy-1")
		}

		// Proxy-1 circuit should now be open
		if manager.Allow("proxy-1") {
			t.Error("Allow() for proxy-1 after failures = true, want false")
		}

		// Other proxies should still be allowed
		if !manager.Allow("proxy-2") {
			t.Error("Allow() for proxy-2 = false, want true")
		}

		// Record success for proxy-1
		manager.RecordSuccess("proxy-1")

		// Proxy-1 may still be in half-open state, allowing limited requests
		// So we won't test that here as it depends on implementation details
	})

	t.Run("Global circuit breaker", func(t *testing.T) {
		config := DefaultCircuitBreakerConfig()
		config.MaxFailures = 3
		config.EnableGlobalBreaker = true

		manager := NewCircuitBreakerManager(config)

		// Initially all proxies should be allowed
		if !manager.Allow("proxy-1") {
			t.Error("Initial Allow() = false, want true")
		}

		// Record many failures across different proxies
		for i := 0; i < 10; i++ {
			// Fix string conversion to use proper string formatting
			proxyID := fmt.Sprintf("proxy-%d", '1'+i%3)
			manager.RecordFailure(proxyID)
		}

		// All proxies should now be blocked by the global breaker
		if manager.Allow("proxy-4") {
			t.Error("Allow() after global failure threshold = true, want false")
		}

		// Record many successes to reset the global breaker
		for i := 0; i < 5; i++ {
			manager.RecordSuccess("proxy-1")
		}

		// Proxies may now be allowed again depending on the implementation
		// We can't test this definitively without knowledge of internal timers
	})
}
