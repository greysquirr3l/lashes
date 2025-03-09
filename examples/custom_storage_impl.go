package examples

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// CustomProxyRepository demonstrates a custom implementation of the ProxyRepository interface
type CustomProxyRepository struct {
	proxies map[string]*domain.Proxy
	mu      sync.RWMutex
}

// NewCustomProxyRepository creates a new instance of CustomProxyRepository
func NewCustomProxyRepository() *CustomProxyRepository {
	return &CustomProxyRepository{
		proxies: make(map[string]*domain.Proxy),
	}
}

// Create adds a new proxy to the repository
func (r *CustomProxyRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[proxy.URL]; exists {
		return fmt.Errorf("proxy with URL %s already exists", proxy.URL)
	}

	r.proxies[proxy.URL] = proxy
	return nil
}

// GetByURL retrieves a proxy by its URL
func (r *CustomProxyRepository) GetByURL(ctx context.Context, url string) (*domain.Proxy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proxy, exists := r.proxies[url]
	if !exists {
		return nil, fmt.Errorf("proxy with URL %s not found", url)
	}

	return proxy, nil
}

// Update updates an existing proxy
func (r *CustomProxyRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[proxy.URL]; !exists {
		return fmt.Errorf("proxy with URL %s not found", proxy.URL)
	}

	r.proxies[proxy.URL] = proxy
	return nil
}

// Delete removes a proxy from the repository
func (r *CustomProxyRepository) Delete(ctx context.Context, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.proxies[url]; !exists {
		return fmt.Errorf("proxy with URL %s not found", url)
	}

	delete(r.proxies, url)
	return nil
}

// List returns all proxies
func (r *CustomProxyRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proxies := make([]*domain.Proxy, 0, len(r.proxies))
	for _, proxy := range r.proxies {
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// CustomStorageImplementationExample demonstrates how to use a custom repository implementation
func CustomStorageImplementationExample() {
	repo := NewCustomProxyRepository()

	ctx := context.Background()

	// Add some example proxies - use minimal fields that are guaranteed to exist
	now := time.Now()
	nowPtr := &now

	proxy1 := &domain.Proxy{
		URL:       "http://example1.com:8080",
		Type:      domain.HTTPProxy,
		LastUsed:  nowPtr,
		CreatedAt: now,
	}

	proxy2 := &domain.Proxy{
		URL:       "http://example2.com:8080",
		Type:      domain.SOCKS5Proxy,
		LastUsed:  nowPtr,
		CreatedAt: now,
	}

	// Create proxies
	if err := repo.Create(ctx, proxy1); err != nil {
		fmt.Printf("Failed to create proxy: %v\n", err)
		return
	}

	if err := repo.Create(ctx, proxy2); err != nil {
		fmt.Printf("Failed to create proxy: %v\n", err)
		return
	}

	// List all proxies
	proxies, err := repo.List(ctx)
	if err != nil {
		fmt.Printf("Failed to list proxies: %v\n", err)
		return
	}

	fmt.Printf("Custom repository has %d proxies\n", len(proxies))

	for _, p := range proxies {
		fmt.Printf("- %s (%s)\n", p.URL, p.Type)
	}
}
