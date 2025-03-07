package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
)

type sqlRepository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewSQLRepository(db *sql.DB, timeout time.Duration) repository.ProxyRepository {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &sqlRepository{
		db:      db,
		timeout: timeout,
	}
}

func (r *sqlRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
        INSERT INTO proxies (
            id, url, type, last_used, last_check, latency, is_active,
            weight, max_retries, timeout_ms
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := r.db.ExecContext(ctx, query,
		proxy.ID,
		proxy.URL, // URL is already a string
		proxy.Type,
		proxy.LastUsed,
		proxy.LastCheck,
		proxy.Latency,
		proxy.IsActive,
		proxy.Weight,
		proxy.MaxRetries,
		proxy.Timeout,
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
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `DELETE FROM proxies WHERE id = $1`

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
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, url, type, last_used, last_check, latency, is_active, 
              weight, max_retries, timeout_ms FROM proxies WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var proxy domain.Proxy
	var urlStr string

	err := row.Scan(
		&proxy.ID,
		&urlStr,
		&proxy.Type,
		&proxy.LastUsed,
		&proxy.LastCheck,
		&proxy.Latency,
		&proxy.IsActive,
		&proxy.Weight,
		&proxy.MaxRetries,
		&proxy.Timeout,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrProxyNotFound
		}
		return nil, fmt.Errorf("failed to get proxy: %w", err)
	}

	// Fixed: assign URL string directly
	proxy.URL = urlStr

	return &proxy, nil
}

func (r *sqlRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Using round-robin strategy by default
	query := `SELECT id, url, type, last_used, last_check, latency, is_active, 
            weight, max_retries, timeout_ms FROM proxies 
            WHERE is_active = true 
            ORDER BY last_used ASC LIMIT 1`

	row := r.db.QueryRowContext(ctx, query)

	var proxy domain.Proxy
	var urlStr string

	err := row.Scan(
		&proxy.ID,
		&urlStr,
		&proxy.Type,
		&proxy.LastUsed,
		&proxy.LastCheck,
		&proxy.Latency,
		&proxy.IsActive,
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

	// Fixed: assign URL string directly
	proxy.URL = urlStr

	// Update last used time
	updateQuery := `UPDATE proxies SET last_used = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, updateQuery, time.Now(), proxy.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update last_used time: %w", err)
	}

	return &proxy, nil
}

func (r *sqlRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `SELECT id, url, type, last_used, last_check, latency, is_active, 
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
			&proxy.LastCheck,
			&proxy.Latency,
			&proxy.IsActive,
			&proxy.Weight,
			&proxy.MaxRetries,
			&proxy.Timeout,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}

		// Fixed: assign URL string directly
		proxy.URL = urlStr
		proxies = append(proxies, &proxy)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating proxies: %w", err)
	}

	return proxies, nil
}

func (r *sqlRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
        UPDATE proxies SET
            url = $1, type = $2, last_used = $3, last_check = $4,
            latency = $5, is_active = $6, weight = $7, max_retries = $8,
            timeout_ms = $9, updated_at = CURRENT_TIMESTAMP
        WHERE id = $10
    `

	result, err := r.db.ExecContext(ctx, query,
		proxy.URL, // URL is already a string
		proxy.Type,
		proxy.LastUsed,
		proxy.LastCheck,
		proxy.Latency,
		proxy.IsActive,
		proxy.Weight,
		proxy.MaxRetries,
		proxy.Timeout,
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

// nolint:unused
func (r *sqlRepository) insertProxy(ctx context.Context, tx *sql.Tx, proxy *domain.Proxy) error {
	query := `
        INSERT INTO proxies (
            id, url, type, last_used, last_check, latency, is_active,
            weight, max_retries, timeout_ms
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			// Log the error or append it to the returned error if not nil
			err = errors.Join(err, fmt.Errorf("failed to close statement: %w", closeErr))
		}
	}()

	_, err = stmt.ExecContext(ctx, proxy.ID, proxy.URL, string(proxy.Type),
		proxy.LastUsed, proxy.LastCheck, proxy.Latency, proxy.IsActive,
		proxy.Weight, proxy.MaxRetries, proxy.Timeout)
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}

	return nil
}

// nolint:unused
func (r *sqlRepository) validateProxy(proxy *domain.Proxy) error {
	if proxy.URL == "" {
		return errors.New("proxy URL cannot be empty")
	}

	parsedURL, err := url.Parse(proxy.URL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	// Fixed: update URL with the string representation
	proxy.URL = parsedURL.String()
	return nil
}

// nolint:unused
func (r *sqlRepository) updateProxy(ctx context.Context, tx *sql.Tx, proxy *domain.Proxy) error {
	query := `
        UPDATE proxies SET
            url = $1, type = $2, last_used = $3, last_check = $4,
            latency = $5, is_active = $6, weight = $7, max_retries = $8,
            timeout_ms = $9, updated_at = CURRENT_TIMESTAMP
        WHERE id = $10
    `

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			// Log the error or append it to the returned error if not nil
			err = errors.Join(err, fmt.Errorf("failed to close statement: %w", closeErr))
		}
	}()

	parsedURL, err := url.Parse(proxy.URL)
	if err != nil {
		return fmt.Errorf("invalid proxy URL: %w", err)
	}

	// Fixed: update URL with the string representation
	proxy.URL = parsedURL.String()

	_, err = stmt.ExecContext(ctx, proxy.URL, string(proxy.Type), proxy.LastUsed,
		proxy.LastCheck, proxy.Latency, proxy.IsActive, proxy.Weight,
		proxy.MaxRetries, proxy.Timeout, proxy.ID)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	return nil
}
