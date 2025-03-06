package migrations

import (
	"database/sql"
	"fmt"

	"github.com/greysquirr3l/lashes/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteMigrator struct {
	db *sql.DB
}

func NewSQLiteMigrator(db *sql.DB) storage.Migrator {
	return &sqliteMigrator{db: db}
}

func (m *sqliteMigrator) Migrate(opts storage.Options) error {
	queries := []string{
		`PRAGMA foreign_keys = ON;`,
		`CREATE TABLE IF NOT EXISTS proxies (
            id TEXT PRIMARY KEY,
            url TEXT NOT NULL,
            type TEXT NOT NULL,
            last_used TIMESTAMP,
            last_check TIMESTAMP,
            latency INTEGER,
            is_active BOOLEAN DEFAULT TRUE,
            weight INTEGER DEFAULT 1,
            max_retries INTEGER DEFAULT 3,
            timeout_ms INTEGER DEFAULT 30000,
            success_count INTEGER DEFAULT 0,
            failure_count INTEGER DEFAULT 0,
            total_requests INTEGER DEFAULT 0,
            avg_latency INTEGER DEFAULT 0,
            last_status_code INTEGER DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`,
		`CREATE INDEX IF NOT EXISTS idx_proxies_type ON proxies(type);`,
		`CREATE INDEX IF NOT EXISTS idx_proxies_is_active ON proxies(is_active);`,
	}

	for _, query := range queries {
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}
	}
	return nil
}

func (m *sqliteMigrator) Drop() error {
	_, err := m.db.Exec(`DROP TABLE IF EXISTS proxies;`)
	return err
}
