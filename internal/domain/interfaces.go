package domain

import (
	"context"
)

// ProxyRepository defines the interface for proxy storage operations.
// Implementations must be safe for concurrent use.
type ProxyRepository interface {
	// Create stores a new proxy.
	// Returns ErrDuplicateID if a proxy with the same ID already exists.
	Create(ctx context.Context, proxy *Proxy) error

	// GetByID retrieves a proxy by its ID.
	// Returns ErrProxyNotFound if the proxy doesn't exist.
	GetByID(ctx context.Context, id string) (*Proxy, error)

	// Update modifies an existing proxy.
	// Returns ErrProxyNotFound if the proxy doesn't exist.
	Update(ctx context.Context, proxy *Proxy) error

	// Delete removes a proxy by its ID.
	// Returns ErrProxyNotFound if the proxy doesn't exist.
	Delete(ctx context.Context, id string) error

	// List returns all available proxies.
	List(ctx context.Context) ([]*Proxy, error)

	// GetNext returns the next proxy according to the repository's strategy.
	// Returns ErrNoProxiesAvailable if the repository is empty.
	GetNext(ctx context.Context) (*Proxy, error)
}
