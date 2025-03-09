// Package repository provides data access interfaces and implementations for the lashes proxy rotator.
//
// The repository package follows the Repository pattern for data access, with implementations for
// various storage backends including in-memory storage and SQL databases through GORM.
//
// Key interfaces:
//   - ProxyRepository: Core interface for accessing proxy data
//
// Key implementations:
//   - memoryRepository: A simple in-memory repository using a map
//   - gorm.proxyRepository: A database-backed repository using GORM
//
// All repository implementations must be thread-safe to support concurrent access.
package repository
