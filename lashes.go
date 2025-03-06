package lashes

import (
	"context"
	"net/http"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/rotation"
	"github.com/greysquirr3l/lashes/internal/storage"
)

// Storage types and options
type StorageType string

const (
	Memory   StorageType = "memory"
	SQLite   StorageType = "sqlite"
	MySQL    StorageType = "mysql"
	Postgres StorageType = "postgres"
)

type StorageOptions struct {
	Type     StorageType
	Database storage.Options
}

// Public type aliases
type (
	Proxy        = domain.Proxy
	ProxyType    = domain.ProxyType
	ProxyMetrics = domain.ProxyMetrics
)

// Public constants
const (
	HTTP   = domain.HTTP
	SOCKS4 = domain.SOCKS4
	SOCKS5 = domain.SOCKS5

	// Rotation strategies
	RoundRobin = rotation.RoundRobin
	Random     = rotation.Random
	Weighted   = rotation.Weighted
	LeastUsed  = rotation.LeastUsed
)

// ProxyRotator represents the main interface for proxy rotation operations.
// It provides methods for managing and rotating through a pool of proxies
// with support for different proxy types and rotation strategies.
type ProxyRotator interface {
	// GetProxy returns the next proxy according to the configured rotation strategy.
	GetProxy(ctx context.Context) (*Proxy, error)

	// AddProxy adds a new proxy to the rotation pool.
	// The proxy URL should be in the format scheme://host:port
	// Supported schemes are http, socks4, and socks5.
	AddProxy(ctx context.Context, proxyURL string, proxyType ProxyType) error

	// RemoveProxy removes a proxy from the rotation pool.
	// Returns ErrProxyNotFound if the proxy doesn't exist.
	RemoveProxy(ctx context.Context, proxyURL string) error

	// Client returns an http.Client configured with the next proxy in the rotation.
	// The client is configured with the rotator's timeout and retry settings.
	Client(ctx context.Context) (*http.Client, error)

	// List returns all available proxies in the pool.
	List(ctx context.Context) ([]*Proxy, error)

	// ValidateProxy validates a single proxy against a target URL.
	// Returns the validation result, latency, and any error that occurred.
	ValidateProxy(ctx context.Context, proxy *Proxy, targetURL string) (bool, time.Duration, error)

	// ValidateAll validates all proxies in the pool using the configured test URL.
	// Updates proxy status and latency metrics for each proxy.
	ValidateAll(ctx context.Context) error
}

// Options configures the proxy rotator behavior.
type Options struct {
	// Storage configuration (optional, defaults to in-memory)
	Storage *storage.Options

	// Strategy defines how proxies are rotated (round-robin, random, weighted, least-used)
	Strategy rotation.StrategyType

	// ValidationTimeout sets the maximum time to wait for proxy validation
	ValidationTimeout time.Duration

	// ValidateOnStart enables proxy validation when adding new proxies
	ValidateOnStart bool

	// TestURL is the URL used for proxy validation
	TestURL string

	// MaxRetries sets the number of retry attempts for failed requests
	MaxRetries int

	// RetryDelay sets the delay between retry attempts
	RetryDelay time.Duration

	// RequestTimeout sets the maximum time to wait for proxy requests
	RequestTimeout time.Duration
}

// New creates a new proxy rotator with the given options.
// If no storage is configured, uses in-memory storage.
func New(opts Options) (ProxyRotator, error) {
	return newRotator(opts)
}

// DefaultOptions returns sensible default options for the proxy rotator.
// Uses in-memory storage and round-robin rotation strategy.
func DefaultOptions() Options {
	return Options{
		Storage:           nil, // Use in-memory storage by default
		Strategy:          rotation.RoundRobin,
		ValidationTimeout: time.Second * 10,
		ValidateOnStart:   true,
		TestURL:           "https://api.ipify.org?format=json",
		MaxRetries:        3,
		RetryDelay:        time.Second,
		RequestTimeout:    time.Second * 30,
	}
}
