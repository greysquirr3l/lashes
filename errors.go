package lashes

import "errors"

// Error definitions
var (
	// ErrInvalidOptions is returned when the provided options are invalid
	ErrInvalidOptions = errors.New("invalid options provided")

	// ErrNoProxiesAvailable is returned when no proxies are available
	ErrNoProxiesAvailable = errors.New("no proxies available")

	// ErrInvalidProxy is returned when the proxy configuration is invalid
	ErrInvalidProxy = errors.New("invalid proxy configuration")

	// ErrProxyNotFound is returned when a proxy cannot be found
	ErrProxyNotFound = errors.New("proxy not found")

	// ErrMetricsNotEnabled is returned when metrics functionality is requested but not enabled
	ErrMetricsNotEnabled = errors.New("metrics collection not enabled")
)
