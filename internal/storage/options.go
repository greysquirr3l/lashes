package storage

import "time"

// DatabaseType defines the type of database to use
type DatabaseType string

const (
	// Memory represents in-memory storage (not a database)
	Memory DatabaseType = "memory"
	// SQLite is a file-based database
	SQLite DatabaseType = "sqlite"
	// MySQL is a popular open source database
	MySQL DatabaseType = "mysql"
	// Postgres is a powerful object-relational database
	Postgres DatabaseType = "postgres"
)

// Options configures database storage
type Options struct {
	// Type of database to use
	Type DatabaseType
	// FilePath is the path to the database file (for SQLite)
	FilePath string
	// ConnectionString is the DSN (for MySQL/Postgres)
	ConnectionString string
	// QueryTimeout is the maximum time for operations
	QueryTimeout time.Duration
}

type Migrator interface {
	Migrate(opts Options) error
	Drop() error
}
