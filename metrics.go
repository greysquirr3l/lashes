package lashes

import (
	"context"
	"sync"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// MetricsCollector provides methods for collecting and accessing proxy metrics
type MetricsCollector interface {
	// RecordRequest records a successful request through a proxy
	RecordRequest(ctx context.Context, proxyID string, latency time.Duration, success bool) error

	// GetProxyMetrics returns metrics for a specific proxy
	GetProxyMetrics(ctx context.Context, proxyID string) (*ProxyMetrics, error)

	// GetAllMetrics returns metrics for all proxies
	GetAllMetrics(ctx context.Context) ([]*ProxyMetrics, error)
}

// defaultMetricsCollector implements MetricsCollector using in-memory storage
type defaultMetricsCollector struct {
	repo    domain.ProxyRepository
	metrics map[string]*proxyMetricsData
	mu      sync.RWMutex
}

type proxyMetricsData struct {
	totalCalls  int64
	totalErrors int64
	latencies   []time.Duration
	lastUsed    time.Time
	minLatency  time.Duration
	maxLatency  time.Duration
	sumLatency  time.Duration
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(repo domain.ProxyRepository) MetricsCollector {
	return &defaultMetricsCollector{
		repo:    repo,
		metrics: make(map[string]*proxyMetricsData),
	}
}

// RecordRequest implements MetricsCollector.RecordRequest
func (m *defaultMetricsCollector) RecordRequest(ctx context.Context, proxyID string, latency time.Duration, success bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, exists := m.metrics[proxyID]
	if !exists {
		data = &proxyMetricsData{
			minLatency: latency,
			maxLatency: latency,
		}
		m.metrics[proxyID] = data
	}

	// Update metrics
	data.totalCalls++
	if !success {
		data.totalErrors++
	}
	data.lastUsed = time.Now()
	data.latencies = append(data.latencies, latency)
	data.sumLatency += latency

	// Update min/max
	if latency < data.minLatency {
		data.minLatency = latency
	}
	if latency > data.maxLatency {
		data.maxLatency = latency
	}

	return nil
}

// GetProxyMetrics implements MetricsCollector.GetProxyMetrics
func (m *defaultMetricsCollector) GetProxyMetrics(ctx context.Context, proxyID string) (*ProxyMetrics, error) {
	m.mu.RLock()
	data, exists := m.metrics[proxyID]
	m.mu.RUnlock()

	if !exists {
		// Use the correct error from errors.go
		return nil, ErrProxyNotFound
	}

	proxy, err := m.repo.GetByID(ctx, proxyID)
	if err != nil {
		return nil, err
	}

	metrics := &ProxyMetrics{
		ProxyID:    proxyID,
		URL:        proxy.URL,
		Type:       string(proxy.Type),
		TotalCalls: data.totalCalls,
		LastUsed:   data.lastUsed,
		ErrorCount: data.totalErrors,
		IsActive:   proxy.Enabled, // Use Enabled for IsActive field in metrics
	}

	// Calculate derived metrics
	if data.totalCalls > 0 {
		metrics.SuccessRate = float64(data.totalCalls-data.totalErrors) / float64(data.totalCalls)
		metrics.AvgLatency = data.sumLatency / time.Duration(data.totalCalls)
		metrics.MinLatency = data.minLatency
		metrics.MaxLatency = data.maxLatency
	}

	return metrics, nil
}

// GetAllMetrics implements MetricsCollector.GetAllMetrics
func (m *defaultMetricsCollector) GetAllMetrics(ctx context.Context) ([]*ProxyMetrics, error) {
	proxies, err := m.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	metrics := make([]*ProxyMetrics, 0, len(proxies))

	for _, proxy := range proxies {
		proxyMetrics, err := m.GetProxyMetrics(ctx, proxy.ID)
		if err != nil {
			// Skip if metrics aren't available for this proxy
			continue
		}
		metrics = append(metrics, proxyMetrics)
	}

	return metrics, nil
}

// cachedMetricsCollector adds caching to the defaultMetricsCollector
type cachedMetricsCollector struct {
	defaultMetricsCollector
	cache      map[string]*ProxyMetrics
	cacheMu    sync.RWMutex
	expiration time.Duration
	lastUpdate time.Time
}

// NewCachedMetricsCollector creates a metrics collector with caching
func NewCachedMetricsCollector(repo domain.ProxyRepository, cacheExpiration time.Duration) MetricsCollector {
	return &cachedMetricsCollector{
		defaultMetricsCollector: defaultMetricsCollector{
			repo:    repo,
			metrics: make(map[string]*proxyMetricsData),
		},
		cache:      make(map[string]*ProxyMetrics),
		expiration: cacheExpiration,
		lastUpdate: time.Now(),
	}
}

// GetProxyMetrics implements MetricsCollector.GetProxyMetrics with caching
func (m *cachedMetricsCollector) GetProxyMetrics(ctx context.Context, proxyID string) (*ProxyMetrics, error) {
	// Check cache first
	m.cacheMu.RLock()
	metrics, ok := m.cache[proxyID]
	cacheValid := ok && time.Since(m.lastUpdate) < m.expiration
	m.cacheMu.RUnlock()

	if cacheValid {
		return metrics, nil
	}

	// Cache miss or expired, get fresh data
	metrics, err := m.defaultMetricsCollector.GetProxyMetrics(ctx, proxyID)
	if err != nil {
		return nil, err
	}

	// Update cache
	m.cacheMu.Lock()
	m.cache[proxyID] = metrics
	m.lastUpdate = time.Now()
	m.cacheMu.Unlock()

	return metrics, nil
}

// GetAllMetrics implements MetricsCollector.GetAllMetrics with caching
func (m *cachedMetricsCollector) GetAllMetrics(ctx context.Context) ([]*ProxyMetrics, error) {
	// Check if the cache is still valid
	m.cacheMu.RLock()
	cacheValid := time.Since(m.lastUpdate) < m.expiration && len(m.cache) > 0
	m.cacheMu.RUnlock()

	if cacheValid {
		// Return all cached metrics
		m.cacheMu.RLock()
		defer m.cacheMu.RUnlock()

		metrics := make([]*ProxyMetrics, 0, len(m.cache))
		for _, metric := range m.cache {
			metrics = append(metrics, metric)
		}
		return metrics, nil
	}

	// Cache expired or empty, get fresh data
	metrics, err := m.defaultMetricsCollector.GetAllMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	m.cacheMu.Lock()
	m.cache = make(map[string]*ProxyMetrics, len(metrics))
	for _, metric := range metrics {
		m.cache[metric.ProxyID] = metric
	}
	m.lastUpdate = time.Now()
	m.cacheMu.Unlock()

	return metrics, nil
}

// RecordRequest implements MetricsCollector.RecordRequest and invalidates cache
func (m *cachedMetricsCollector) RecordRequest(ctx context.Context, proxyID string, latency time.Duration, success bool) error {
	// Update metrics
	err := m.defaultMetricsCollector.RecordRequest(ctx, proxyID, latency, success)
	if err != nil {
		return err
	}

	// Invalidate cache for this proxy
	m.cacheMu.Lock()
	delete(m.cache, proxyID)
	m.cacheMu.Unlock()

	return nil
}
