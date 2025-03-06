package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu           sync.RWMutex
	proxyMetrics map[string]*proxyMetric
}

type proxyMetric struct {
	SuccessCount   int
	FailureCount   int
	TotalRequests  int
	AverageLatency int64
	LastStatusCode int
	LastUsed       time.Time
	// Add additional metrics:
	// - Protocol-specific metrics
	// - Geographic distribution
	// - Response time percentiles
	// - Error type tracking
}

func NewMetrics() *Metrics {
	return &Metrics{
		proxyMetrics: make(map[string]*proxyMetric),
	}
}

// RecordMetric records a metric for the given proxy
func (m *Metrics) RecordMetric(proxyID string, statusCode int, latency int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metric, exists := m.proxyMetrics[proxyID]
	if !exists {
		metric = &proxyMetric{}
		m.proxyMetrics[proxyID] = metric
	}

	metric.TotalRequests++
	metric.LastStatusCode = statusCode
	metric.LastUsed = time.Now()

	if statusCode < 400 {
		metric.SuccessCount++
	} else {
		metric.FailureCount++
	}

	// Update average latency
	oldTotal := metric.AverageLatency * int64(metric.TotalRequests-1)
	metric.AverageLatency = (oldTotal + latency) / int64(metric.TotalRequests)
}

// GetSuccessCount returns the success count for the given proxy
func (m *Metrics) GetSuccessCount(proxyID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if metric, exists := m.proxyMetrics[proxyID]; exists {
		return metric.SuccessCount
	}
	return 0
}

// GetFailureCount returns the failure count for the given proxy
func (m *Metrics) GetFailureCount(proxyID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if metric, exists := m.proxyMetrics[proxyID]; exists {
		return metric.FailureCount
	}
	return 0
}

// GetAverageLatency returns the average latency for the given proxy
func (m *Metrics) GetAverageLatency(proxyID string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if metric, exists := m.proxyMetrics[proxyID]; exists {
		return metric.AverageLatency
	}
	return 0
}
