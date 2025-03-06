package lashes

import "errors"

var (
	// ErrInvalidOptions is returned when the provided options are invalid
	ErrInvalidOptions = errors.New("invalid options provided")

	// ErrNoProxiesAvailable is returned when no proxies are available
	ErrNoProxiesAvailable = errors.New("no proxies available")

	// ErrInvalidProxy is returned when the proxy configuration is invalid
	ErrInvalidProxy = errors.New("invalid proxy configuration")
)
