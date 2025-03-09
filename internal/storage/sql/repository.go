package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
)

const (
	// SQL statement to create the proxies table - updated schema
	createTableSQL = `
	CREATE TABLE IF NOT EXISTS proxies (
		id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		type TEXT NOT NULL,
		username TEXT,
		password TEXT,
		country_code TEXT,
		weight INTEGER DEFAULT 1,
		last_used TIMESTAMP,
		enabled BOOLEAN DEFAULT true,
		latency INTEGER DEFAULT 0,
		success_rate REAL DEFAULT 0,
		usage_count INTEGER DEFAULT 0,
		error_count INTEGER DEFAULT 0,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	)
	`
)

type sqlRepository struct {
	db      *sql.DB
	timeout time.Duration
}

// NewSQLRepository creates a new SQL-based repository
func NewSQLRepository(db *sql.DB, timeout time.Duration) repository.ProxyRepository {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	repo := &sqlRepository{
		db:      db,
		timeout: timeout,
	}

	// Initialize tables and check for errors
	if err := repo.init(); err != nil {
		// Log the error but don't fail - tables might already exist
		// In a production system, consider adding proper logging here
		fmt.Printf("Warning: repository initialization error: %v\n", err)
	}

	return repo
}

// init creates the necessary database tables if they don't exist
func (r *sqlRepository) init() error {
	_, err := r.db.Exec(createTableSQL)
	return err
}

func (r *sqlRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	// Create a timeout context that inherits from the provided context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
        INSERT INTO proxies (
            id, url, type, username, password, country_code, weight, 
            last_used, enabled, latency, success_rate, 
            usage_count, error_count, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	now := time.Now()
	if proxy.CreatedAt.IsZero() {
		proxy.CreatedAt = now
	}
	proxy.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		proxy.ID,
		proxy.URL,
		string(proxy.Type),
		proxy.Username,
		proxy.Password,
		proxy.CountryCode,
		proxy.Weight,
		proxy.LastUsed,
		proxy.Enabled,
		proxy.Latency,
		proxy.SuccessRate,
		proxy.UsageCount,
		proxy.ErrorCount,
		proxy.CreatedAt,
		proxy.UpdatedAt,
	)
	if err != nil {
		if isSQLiteConstraintError(err) || isPostgresConstraintError(err) {
			return repository.ErrDuplicateID
		}
		return fmt.Errorf("failed to create proxy: %w", err)
	}
	return nil
}

func (r *sqlRepository) Delete(ctx context.Context, id string) error {
	// Create a timeout context that inherits from the provided context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `DELETE FROM proxies WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete proxy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrProxyNotFound
	}

	return nil
}

func (r *sqlRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	// Create a timeout context that inherits from the provided context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
	SELECT 
		id, url, type, username, password, country_code, weight, 
		last_used, enabled, latency, success_rate, 
		usage_count, error_count, created_at, updated_at
	FROM proxies 
	WHERE id = ?
	`

	proxy := &domain.Proxy{}
	var lastUsed, createdAt, updatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&proxy.ID,
		&proxy.URL,
		&proxy.Type,
		&proxy.Username,
		&proxy.Password,
		&proxy.CountryCode,
		&proxy.Weight,
		&lastUsed,
		&proxy.Enabled,
		&proxy.Latency,
		&proxy.SuccessRate,
		&proxy.UsageCount,
		&proxy.ErrorCount,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProxyNotFound
		}
		return nil, fmt.Errorf("failed to get proxy: %w", err)
	}

	// Convert nullable fields
	if lastUsed.Valid {
		t := lastUsed.Time
		proxy.LastUsed = &t
	}

	if createdAt.Valid {
		proxy.CreatedAt = createdAt.Time
	}

	if updatedAt.Valid {
		proxy.UpdatedAt = updatedAt.Time
	}

	return proxy, nil
}

func (r *sqlRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	// Create a timeout context that inherits from the provided context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Using round-robin strategy by default
	query := `SELECT id, url, type, last_used, enabled, latency, 
            weight, max_retries, timeout_ms FROM proxies 
            WHERE enabled = true 
            ORDER BY last_used ASC LIMIT 1`

	row := r.db.QueryRowContext(ctx, query)

	var proxy domain.Proxy
	var urlStr string

	err := row.Scan(
		&proxy.ID,
		&urlStr,
		&proxy.Type,
		&proxy.LastUsed,
		&proxy.Enabled,
		&proxy.Latency,
		&proxy.Weight,
		&proxy.MaxRetries,
		&proxy.Timeout,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProxyNotFound
		}
		return nil, fmt.Errorf("failed to get next proxy: %w", err)
	}

	// Set URL directly
	proxy.URL = urlStr

	// Update last used time
	updateQuery := `UPDATE proxies SET last_used = ? WHERE id = ?`
	_, err = r.db.ExecContext(ctx, updateQuery, time.Now(), proxy.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last_used time: %w", err)
	}

	return &proxy, nil
}

func (r *sqlRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	// Create a timeout context that inherits from the provided context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, url, type, last_used, enabled, latency, 
              weight, max_retries, timeout_ms FROM proxies`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list proxies: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Consider logging this error or handling it appropriately
		}
	}()

	var proxies []*domain.Proxy

	for rows.Next() {
		var proxy domain.Proxy
		var urlStr string

		err := rows.Scan(
			&proxy.ID,
			&urlStr,
			&proxy.Type,
			&proxy.LastUsed,
			&proxy.Enabled,
			&proxy.Latency,
			&proxy.Weight,
			&proxy.MaxRetries,
			&proxy.Timeout,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}

		// Set URL directly
		proxy.URL = urlStr
		proxies = append(proxies, &proxy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating proxies: %w", err)
	}

	return proxies, nil
}

func (r *sqlRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	// Create a timeout context that inherits from the provided context
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
        UPDATE proxies SET
            url = ?, type = ?, username = ?, password = ?, country_code = ?, 
            weight = ?, last_used = ?, enabled = ?, latency = ?, 
            success_rate = ?, usage_count = ?, error_count = ?, updated_at = ?
        WHERE id = ?
    `

	proxy.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		proxy.URL,
		string(proxy.Type),
		proxy.Username,
		proxy.Password,
		proxy.CountryCode,
		proxy.Weight,
		proxy.LastUsed,
		proxy.Enabled,
		proxy.Latency,
		proxy.SuccessRate,
		proxy.UsageCount,
		proxy.ErrorCount,
		proxy.UpdatedAt,
		proxy.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update proxy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrProxyNotFound
	}

	return nil
}
