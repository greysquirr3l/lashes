package queries

const (
	CreateProxyTable = `
        CREATE TABLE IF NOT EXISTS proxies (
            id TEXT PRIMARY KEY,
            url TEXT NOT NULL,
            type TEXT NOT NULL,
            last_used TIMESTAMP,
            last_check TIMESTAMP,
            latency BIGINT,
            is_active BOOLEAN,
            weight INTEGER DEFAULT 1,
            max_retries INTEGER DEFAULT 3,
            timeout_ms BIGINT DEFAULT 30000,
            success_count BIGINT DEFAULT 0,
            failure_count BIGINT DEFAULT 0,
            total_requests BIGINT DEFAULT 0,
            avg_latency BIGINT DEFAULT 0,
            last_status_code INTEGER DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`

	InsertProxy = `
        INSERT INTO proxies (
            id, url, type, last_used, last_check, latency, is_active,
            weight, max_retries, timeout_ms
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	UpdateProxy = `
        UPDATE proxies SET
            url = ?, type = ?, last_used = ?, last_check = ?,
            latency = ?, is_active = ?, weight = ?, max_retries = ?,
            timeout_ms = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?`

	SelectProxyByID = `
        SELECT * FROM proxies WHERE id = ?`

	SelectAllProxies = `
        SELECT * FROM proxies`

	DeleteProxy = `
        DELETE FROM proxies WHERE id = ?`

	UpdateMetrics = `
        UPDATE proxies SET
            success_count = ?, failure_count = ?, total_requests = ?,
            avg_latency = ?, last_status_code = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?`
)
