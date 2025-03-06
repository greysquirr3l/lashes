// Package repository provides proxy storage implementations for different backends.
//
// The repository package implements the domain.ProxyRepository interface using
// different storage backends:
//   - Memory: Thread-safe in-memory storage (default)
//   - GORM: SQL storage using GORM (SQLite, MySQL, PostgreSQL)
//
// All implementations are safe for concurrent use and follow standard CRUD patterns.
package repository
