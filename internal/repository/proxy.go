package repository

import (
	"context"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// ProxyRepository defines the interface for proxy storage operations
type ProxyRepository interface {
	// Create stores a new proxy
	Create(ctx context.Context, proxy *domain.Proxy) error

	// GetByID retrieves a proxy by its ID
	GetByID(ctx context.Context, id string) (*domain.Proxy, error)

	// Update modifies an existing proxy
	Update(ctx context.Context, proxy *domain.Proxy) error

	// Delete removes a proxy by its ID
	Delete(ctx context.Context, id string) error

	// List returns all available proxies
	List(ctx context.Context) ([]*domain.Proxy, error)

	// GetNext returns the next proxy according to the repository's strategy
	GetNext(ctx context.Context) (*domain.Proxy, error)
}
