package mock

import (
	"context"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/validation"
)

// MockValidator provides a controllable validator implementation for testing
type MockValidator struct {
	ValidateFunc           func(ctx context.Context, proxy *domain.Proxy) (bool, time.Duration, error)
	ValidateWithTargetFunc func(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error)
}

// NewValidator creates a new mock validator
func NewValidator() validation.Validator {
	return &MockValidator{
		ValidateFunc: func(ctx context.Context, proxy *domain.Proxy) (bool, time.Duration, error) {
			// Default implementation: success with 100ms latency
			return true, 100 * time.Millisecond, nil
		},
		ValidateWithTargetFunc: func(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error) {
			// Default implementation: success with 100ms latency
			return true, 100 * time.Millisecond, nil
		},
	}
}

// Validate implements the Validator interface
func (m *MockValidator) Validate(ctx context.Context, proxy *domain.Proxy) (bool, time.Duration, error) {
	return m.ValidateFunc(ctx, proxy)
}

// ValidateWithTarget implements the Validator interface
func (m *MockValidator) ValidateWithTarget(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error) {
	return m.ValidateWithTargetFunc(ctx, proxy, targetURL)
}

// WithCustomResponse configures the mock with custom validation responses
func (m *MockValidator) WithCustomResponse(valid bool, latency time.Duration, err error) *MockValidator {
	m.ValidateFunc = func(ctx context.Context, proxy *domain.Proxy) (bool, time.Duration, error) {
		return valid, latency, err
	}
	m.ValidateWithTargetFunc = func(ctx context.Context, proxy *domain.Proxy, targetURL string) (bool, time.Duration, error) {
		return valid, latency, err
	}
	return m
}
