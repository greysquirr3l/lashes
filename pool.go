package lashes

import (
	"context"
	"errors"
)

// PoolManager provides methods for managing groups of proxies
type PoolManager interface {
	// CreatePool creates a new proxy pool with the given name
	CreatePool(ctx context.Context, name string) error

	// DeletePool removes a proxy pool
	DeletePool(ctx context.Context, name string) error

	// AddToPool adds a proxy to a pool
	AddToPool(ctx context.Context, poolName, proxyID string) error

	// RemoveFromPool removes a proxy from a pool
	RemoveFromPool(ctx context.Context, poolName, proxyID string) error

	// GetPoolProxies returns all proxies in a pool
	GetPoolProxies(ctx context.Context, poolName string) ([]*Proxy, error)

	// GetNextFromPool returns the next proxy from a specific pool
	GetNextFromPool(ctx context.Context, poolName string) (*Proxy, error)
}

// Pool related errors
var (
	ErrPoolNotFound = errors.New("pool not found")
	ErrPoolExists   = errors.New("pool already exists")
)

// GetProxiesByCountry returns all proxies for a specific country
func (r *rotator) GetProxiesByCountry(ctx context.Context, countryCode string) ([]*Proxy, error) {
	allProxies, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*Proxy
	for _, proxy := range allProxies {
		if proxy.CountryCode == countryCode {
			filtered = append(filtered, proxy)
		}
	}

	return filtered, nil
}

// GetProxiesByType returns all proxies of the specified type
func (r *rotator) GetProxiesByType(ctx context.Context, proxyType ProxyType) ([]*Proxy, error) {
	allProxies, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*Proxy
	for _, proxy := range allProxies {
		if proxy.Type == proxyType {
			filtered = append(filtered, proxy)
		}
	}

	return filtered, nil
}
