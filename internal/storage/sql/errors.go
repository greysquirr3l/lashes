package sql

import (
	"database/sql"
	"errors"
	"strings"
)

// isSQLiteConstraintError checks if error is a SQLite constraint violation
func isSQLiteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "FOREIGN KEY constraint failed")
}

// isPostgresConstraintError checks if error is a PostgreSQL constraint violation
func isPostgresConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
		strings.Contains(err.Error(), "violates foreign key constraint")
}

// isMySQLConstraintError checks if error is a MySQL constraint violation
func isMySQLConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "foreign key constraint fails")
}

// IsConstraintViolation returns true if the error is a database constraint violation
func IsConstraintViolation(err error) bool {
	return isSQLiteConstraintError(err) ||
		isPostgresConstraintError(err) ||
		isMySQLConstraintError(err)
}

// IsNoRowsError returns true if the error is a "no rows found" error
func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// Common errors
var (
	// ErrInvalidProxy is returned when a proxy validation fails
	ErrInvalidProxy = errors.New("invalid proxy configuration")

	// ErrDuplicateID is returned when attempting to create a proxy with an existing ID
	ErrDuplicateID = errors.New("proxy with this ID already exists")

	// ErrNotImplemented is returned for methods not yet implemented
	ErrNotImplemented = errors.New("method not implemented")
)
