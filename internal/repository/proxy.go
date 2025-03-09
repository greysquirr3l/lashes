package repository

import (
	"context"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// ProxyRepository defines the interface for proxy storage operations
type ProxyRepository interface {
	// Create stores a new proxy in the repository
	Create(ctx context.Context, proxy *domain.Proxy) error

	// GetByID retrieves a proxy by its ID
	GetByID(ctx context.Context, id string) (*domain.Proxy, error)

	// Update modifies an existing proxy
	Update(ctx context.Context, proxy *domain.Proxy) error

	// Delete removes a proxy by ID
	Delete(ctx context.Context, id string) error

	// List returns all proxies in the repository
	List(ctx context.Context) ([]*domain.Proxy, error)

	// GetNext returns the next proxy according to the defined strategy
	GetNext(ctx context.Context) (*domain.Proxy, error)
}
