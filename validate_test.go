package lashes

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

func TestValidateProxy(t *testing.T) {
	// Create a mock rotator for testing
	r := &rotator{
		opts: Options{
			ValidationTimeout: time.Second * 5,
			TestURL:           "https://example.com",
			MaxRetries:        2,
		},
	}

	// Create a test proxy
	proxy := &domain.Proxy{
		ID:      "test-proxy",
		URL:     "http://example.com:8080",
		Type:    domain.HTTPProxy,
		Enabled: true,
	}

	// We can't really test the actual validation here without mocking HTTP
	// responses, but we can test that context handling works correctly

	t.Run("Context cancellation", func(t *testing.T) {
		// Create a canceled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Validate with canceled context
		_, _, err := r.ValidateProxy(ctx, proxy, "https://example.com")

		// Should fail with context canceled error
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("Validate all with no proxies", func(t *testing.T) {
		// Setup the rotator with a mock repository that returns no proxies
		r.repo = &mockRepository{proxies: map[string]*domain.Proxy{}}

		// Call ValidateAll
		err := r.ValidateAll(context.Background())

		// Should fail with ErrNoProxiesAvailable
		if !errors.Is(err, ErrNoProxiesAvailable) {
			t.Errorf("Expected ErrNoProxiesAvailable, got %v", err)
		}
	})

	t.Run("Validate all with repository error", func(t *testing.T) {
		// Create a mock repository that returns an error
		mockRepo := &errorMockRepository{err: errors.New("repository error")}
		r.repo = mockRepo

		// Call ValidateAll
		err := r.ValidateAll(context.Background())

		// Should fail with the repository error
		if err == nil || !errors.Is(err, mockRepo.err) {
			t.Errorf("Expected repository error, got %v", err)
		}
	})
}

// errorMockRepository returns errors for all operations
type errorMockRepository struct {
	err error
}

func (r *errorMockRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	return r.err
}

func (r *errorMockRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	return nil, r.err
}

func (r *errorMockRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	return r.err
}

func (r *errorMockRepository) Delete(ctx context.Context, id string) error {
	return r.err
}

func (r *errorMockRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	return nil, r.err
}

func (r *errorMockRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	return nil, r.err
}
