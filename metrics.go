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

// ProxyMetrics contains performance metrics for a single proxy
type ProxyMetrics struct {
	ProxyID     string        `json:"proxy_id"`
	URL         string        `json:"url"`
	Type        string        `json:"type"`
	SuccessRate float64       `json:"success_rate"`
	TotalCalls  int64         `json:"total_calls"`
	AvgLatency  time.Duration `json:"avg_latency_ms"`
	MinLatency  time.Duration `json:"min_latency_ms"`
	MaxLatency  time.Duration `json:"max_latency_ms"`
	LastUsed    time.Time     `json:"last_used"`
	ErrorCount  int64         `json:"error_count"`
	IsActive    bool          `json:"is_active"`
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
		IsActive:   proxy.IsActive,
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
