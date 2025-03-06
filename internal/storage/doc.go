// Package storage provides database configuration and migrations.
//
// Supported storage backends:
//   - In-memory (default, zero dependencies)
//   - SQLite (local file storage)
//   - MySQL (network database)
//   - PostgreSQL (network database)
//
// Database support is optional and lazy-loaded only when explicitly configured.
package storage
