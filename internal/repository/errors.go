package repository

import "errors"

// Repository error definitions
var (
	// ErrProxyNotFound is returned when a proxy cannot be found in the repository
	ErrProxyNotFound = errors.New("proxy not found")

	// ErrDuplicateID is returned when attempting to create a proxy with an existing ID
	ErrDuplicateID = errors.New("proxy with this ID already exists")

	// ErrInvalidProxy is returned when proxy data validation fails
	ErrInvalidProxy = errors.New("invalid proxy data")
)
