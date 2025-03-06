package storage

import "time"

type DatabaseType string

const (
	SQLite   DatabaseType = "sqlite"
	MySQL    DatabaseType = "mysql"
	Postgres DatabaseType = "postgres"
)

type Options struct {
	Type           DatabaseType
	DSN            string
	MaxConnections int
	RetentionDays  int
	MetricsEnabled bool
	ConnTimeout    time.Duration
	QueryTimeout   time.Duration
}

type Migrator interface {
	Migrate(opts Options) error
	Drop() error
}
