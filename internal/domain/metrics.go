package domain

import (
	"sync"
	"time"
)

// Metrics tracks proxy performance measurements
type Metrics struct {
	SuccessCount   int64         `json:"success_count"`
	FailureCount   int64         `json:"failure_count"`
	TotalRequests  int64         `json:"total_requests"`
	AvgLatency     time.Duration `json:"avg_latency_ms"`
	LastStatusCode int           `json:"last_status_code"`
	mu             sync.RWMutex
}

// NewMetrics creates a new metrics tracker
func NewMetrics() *Metrics {
	return &Metrics{}
}

// RecordSuccess records a successful request with status code
func (m *Metrics) RecordSuccess(statusCode int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.SuccessCount++
	m.TotalRequests++
	m.LastStatusCode = statusCode
}

// RecordFailure records a failed request with status code
func (m *Metrics) RecordFailure(statusCode int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.FailureCount++
	m.TotalRequests++
	m.LastStatusCode = statusCode
}

// IncrementLatency adds a latency measurement and recalculates average
func (m *Metrics) IncrementLatency(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.TotalRequests == 0 {
		m.AvgLatency = latency
		return
	}
	
	// Calculate new average: ((avg * count) + new) / (count + 1)
	total := m.AvgLatency.Nanoseconds() * m.TotalRequests
	m.AvgLatency = time.Duration((total + latency.Nanoseconds()) / (m.TotalRequests + 1))
}

// GetSuccessRate returns the success rate as a float between 0 and 1
func (m *Metrics) GetSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.TotalRequests == 0 {
		return 0
	}
	
	return float64(m.SuccessCount) / float64(m.TotalRequests)
}