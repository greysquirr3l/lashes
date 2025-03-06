package migrations

import (
	"database/sql"
	"fmt"

	"github.com/greysquirr3l/lashes/internal/storage"
	_ "github.com/lib/pq"
)

type postgresMigrator struct {
	db *sql.DB
}

func NewPostgresMigrator(db *sql.DB) storage.Migrator {
	return &postgresMigrator{db: db}
}

func (m *postgresMigrator) Migrate(opts storage.Options) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS proxies (
            id TEXT PRIMARY KEY,
            url TEXT NOT NULL,
            type TEXT NOT NULL,
            last_used TIMESTAMP WITH TIME ZONE,
            last_check TIMESTAMP WITH TIME ZONE,
            latency BIGINT,
            is_active BOOLEAN DEFAULT TRUE,
            weight INTEGER DEFAULT 1,
            max_retries INTEGER DEFAULT 3,
            timeout_ms BIGINT DEFAULT 30000,
            success_count BIGINT DEFAULT 0,
            failure_count BIGINT DEFAULT 0,
            total_requests BIGINT DEFAULT 0,
            avg_latency BIGINT DEFAULT 0,
            last_status_code INTEGER DEFAULT 0,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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

func (m *postgresMigrator) Drop() error {
	_, err := m.db.Exec(`DROP TABLE IF EXISTS proxies CASCADE;`)
	return err
}
