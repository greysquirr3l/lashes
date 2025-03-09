package lashes

import (
	"errors"
	"fmt"
)

// Base error types
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

	// ErrValidationFailed is returned when proxy validation fails
	ErrValidationFailed = errors.New("proxy validation failed")
)

// ValidationError provides detailed information about proxy validation failures
type ValidationError struct {
	ProxyID    string
	ProxyURL   string
	Reason     string
	StatusCode int
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	msg := fmt.Sprintf("proxy validation failed for %s", e.ProxyURL)
	if e.Reason != "" {
		msg += ": " + e.Reason
	}
	if e.StatusCode > 0 {
		msg += fmt.Sprintf(" (status code: %d)", e.StatusCode)
	}
	return msg
}

// Is allows this error to be matched with errors.Is() against ErrValidationFailed
func (e *ValidationError) Is(target error) bool {
	return target == ErrValidationFailed
}

// NewValidationError creates a new validation error
func NewValidationError(proxyID, proxyURL, reason string, statusCode int) *ValidationError {
	return &ValidationError{
		ProxyID:    proxyID,
		ProxyURL:   proxyURL,
		Reason:     reason,
		StatusCode: statusCode,
	}
}
