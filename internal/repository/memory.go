package repository

import (
	"context"
	"sync"

	"github.com/greysquirr3l/lashes/internal/domain"
)

type memoryRepository struct {
	proxies map[string]*domain.Proxy
	mu      sync.RWMutex
}

func NewMemoryRepository() ProxyRepository {
	return &memoryRepository{
		proxies: make(map[string]*domain.Proxy),
	}
}

// Create implements ProxyRepository.Create
func (r *memoryRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	if err := validateProxy(proxy); err != nil {
		return ErrInvalidProxy
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[proxy.ID]; exists {
		return ErrDuplicateID
	}
	r.proxies[proxy.ID] = proxy
	return nil
}

func (r *memoryRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proxy, exists := r.proxies[id]
	if !exists {
		return nil, ErrProxyNotFound
	}
	return proxy, nil
}

func (r *memoryRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[proxy.ID]; !exists {
		return ErrProxyNotFound
	}
	r.proxies[proxy.ID] = proxy
	return nil
}

func (r *memoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[id]; !exists {
		return ErrProxyNotFound
	}
	delete(r.proxies, id)
	return nil
}

func (r *memoryRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proxies := make([]*domain.Proxy, 0, len(r.proxies))
	for _, proxy := range r.proxies {
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}

func (r *memoryRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	proxies, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(proxies) == 0 {
		return nil, ErrProxyNotFound
	}
	return proxies[0], nil
}

// validateProxy checks if a proxy is valid
func validateProxy(proxy *domain.Proxy) error {
	if proxy == nil {
		return ErrInvalidProxy
	}
	if proxy.URL == "" {
		return ErrInvalidProxy
	}
	if proxy.Type == "" {
		return ErrInvalidProxy
	}
	return nil
}
