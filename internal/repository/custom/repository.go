package custom

import (
	"context"
	"sync"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
)

// ProxyRepository demonstrates a custom implementation of the repository.ProxyRepository interface
type ProxyRepository struct {
	proxies map[string]*domain.Proxy
	mu      sync.RWMutex
}

// NewProxyRepository creates a new instance of custom ProxyRepository
func NewProxyRepository() repository.ProxyRepository {
	return &ProxyRepository{
		proxies: make(map[string]*domain.Proxy),
	}
}

// Create adds a new proxy to the repository
func (r *ProxyRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[proxy.ID]; exists {
		return repository.ErrDuplicateID
	}

	r.proxies[proxy.ID] = proxy
	return nil
}

// GetByID retrieves a proxy by its ID
func (r *ProxyRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proxy, exists := r.proxies[id]
	if !exists {
		return nil, repository.ErrProxyNotFound
	}

	return proxy, nil
}

// Update updates an existing proxy
func (r *ProxyRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[proxy.ID]; !exists {
		return repository.ErrProxyNotFound
	}

	r.proxies[proxy.ID] = proxy
	return nil
}

// Delete removes a proxy from the repository
func (r *ProxyRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[id]; !exists {
		return repository.ErrProxyNotFound
	}

	delete(r.proxies, id)
	return nil
}

// List returns all proxies
func (r *ProxyRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proxies := make([]*domain.Proxy, 0, len(r.proxies))
	for _, proxy := range r.proxies {
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// GetNext returns the next proxy according to the repository's strategy
func (r *ProxyRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	proxies, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(proxies) == 0 {
		return nil, repository.ErrProxyNotFound
	}

	// Simple implementation - just return the first proxy
	return proxies[0], nil
}
