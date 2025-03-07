package domain

import (
	"sync"
	"time"
)

// Metrics tracks usage and performance data for proxies
type Metrics struct {
    mu              sync.RWMutex
    TotalRequests   int64
    SuccessCount    int64
    FailureCount    int64
    LastStatusCode  int
    LastUsedAt      time.Time
    AverageLatency  time.Duration
    LatencyCount    int64
    TotalLatency    time.Duration
    ConsecutiveFails int
}

// IncrementLatency adds a new latency measurement and updates the average
func (m *Metrics) IncrementLatency(latency time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.LatencyCount++
    m.TotalLatency += latency
    m.AverageLatency = time.Duration(int64(m.TotalLatency) / m.LatencyCount)
    m.LastUsedAt = time.Now()
}

// RecordSuccess records a successful request
func (m *Metrics) RecordSuccess(statusCode int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.TotalRequests++
    m.SuccessCount++
    m.LastStatusCode = statusCode
    m.ConsecutiveFails = 0
    m.LastUsedAt = time.Now()
}

// RecordFailure records a failed request
func (m *Metrics) RecordFailure(statusCode int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.TotalRequests++
    m.FailureCount++
    m.LastStatusCode = statusCode
    m.ConsecutiveFails++
    m.LastUsedAt = time.Now()
}

// GetSuccessRate returns the percentage of successful requests
func (m *Metrics) GetSuccessRate() float64 {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    if m.TotalRequests == 0 {
        return 0.0
    }
    return float64(m.SuccessCount) / float64(m.TotalRequests) * 100.0
}

// Reset resets all metrics to their zero values
func (m *Metrics) Reset() {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.TotalRequests = 0
    m.SuccessCount = 0
    m.FailureCount = 0
    m.LastStatusCode = 0
    m.AverageLatency = 0
    m.LatencyCount = 0
    m.TotalLatency = 0
    m.ConsecutiveFails = 0
    // Intentionally not resetting LastUsedAt to keep track of idle proxies
}

// NewMetrics creates a new Metrics instance with initialized values
func NewMetrics() *Metrics {
    return &Metrics{
        LastUsedAt: time.Now(),
    }
}