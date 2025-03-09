package domain

import (
	"context"
	"time"
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

// ProxyProvider defines the minimal interface for getting proxies
type ProxyProvider interface {
	// GetProxy returns the next proxy according to the configured rotation strategy.
	GetProxy(ctx context.Context) (*Proxy, error)

	// List returns all available proxies in the pool.
	List(ctx context.Context) ([]*Proxy, error)
}

// ProxyManager defines the interface for managing proxy entries
type ProxyManager interface {
	// AddProxy adds a new proxy to the rotation pool.
	AddProxy(ctx context.Context, proxyURL string, proxyType ProxyType) error

	// RemoveProxy removes a proxy from the rotation pool.
	RemoveProxy(ctx context.Context, proxyURL string) error
}

// ProxyValidator defines the interface for proxy validation
type ProxyValidator interface {
	// ValidateProxy validates a single proxy against a target URL.
	ValidateProxy(ctx context.Context, proxy *Proxy, targetURL string) (bool, time.Duration, error)

	// ValidateAll validates all proxies in the pool using the configured test URL.
	ValidateAll(ctx context.Context) error
}

// MetricsProvider defines the interface for accessing proxy metrics
type MetricsProvider interface {
	// GetProxyMetrics returns performance metrics for a specific proxy
	GetProxyMetrics(ctx context.Context, proxyID string) (*ProxyMetrics, error)

	// GetAllMetrics returns performance metrics for all proxies
	GetAllMetrics(ctx context.Context) ([]*ProxyMetrics, error)
}
