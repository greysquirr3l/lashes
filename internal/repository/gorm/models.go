package gorm

import (
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// ProxyModel is the GORM model for proxy data
type ProxyModel struct {
	ID             string `gorm:"primaryKey"`
	URL            string
	Type           string
	Username       string
	Password       string
	CountryCode    string
	Weight         int       `gorm:"default:1"`
	LastUsed       time.Time // Store as time.Time in the database
	LastCheck      time.Time
	Latency        int64 // Store as int64 in milliseconds
	IsActive       bool  `gorm:"default:true"`
	SuccessRate    float64
	UsageCount     int64
	ErrorCount     int64
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	MaxRetries     int
	Timeout        time.Duration
	SuccessCount   int64
	FailureCount   int64
	TotalRequests  int64
	AvgLatency     time.Duration
	LastStatusCode int
}

// ToDomain converts a GORM model to a domain model
func (m *ProxyModel) ToDomain() (*domain.Proxy, error) {
	// Create a pointer to LastUsed and LastCheck for the domain model
	var lastUsed *time.Time
	if !m.LastUsed.IsZero() {
		t := m.LastUsed // Make a copy
		lastUsed = &t
	}
	
	var lastCheck *time.Time
	if !m.LastCheck.IsZero() {
		t := m.LastCheck // Make a copy
		lastCheck = &t
	}
	
	return &domain.Proxy{
		ID:          m.ID,
		URL:         m.URL, // URL is already a string
		Type:        domain.ProxyType(m.Type),
		Username:    m.Username,
		Password:    m.Password,
		CountryCode: m.CountryCode,
		Weight:      m.Weight,
		LastUsed:    lastUsed, // Convert to pointer
		LastCheck:   lastCheck, // Convert to pointer
		Latency:     m.Latency, // Already int64
		Enabled:     m.IsActive,
		IsActive:    m.IsActive, // Compatibility field
		SuccessRate: m.SuccessRate,
		UsageCount:  m.UsageCount,
		ErrorCount:  m.ErrorCount,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		MaxRetries:  m.MaxRetries,
		Timeout:     m.Timeout,
		Metrics: domain.ProxyMetrics{
			SuccessCount:   m.SuccessCount,
			FailureCount:   m.FailureCount,
			TotalRequests:  m.TotalRequests,
			AvgLatency:     m.AvgLatency,
			LastStatusCode: m.LastStatusCode,
		},
	}, nil
}

// FromDomain converts a domain model to a GORM model
func FromDomain(proxy *domain.Proxy) *ProxyModel {
	// Convert pointers to values, using zero values if nil
	var lastUsed time.Time
	if proxy.LastUsed != nil {
		lastUsed = *proxy.LastUsed
	}
	
	var lastCheck time.Time
	if proxy.LastCheck != nil {
		lastCheck = *proxy.LastCheck
	}
	
	return &ProxyModel{
		ID:          proxy.ID,
		URL:         proxy.URL, // URL is already a string
		Type:        string(proxy.Type),
		Username:    proxy.Username,
		Password:    proxy.Password,
		CountryCode: proxy.CountryCode,
		Weight:      proxy.Weight,
		LastUsed:    lastUsed,    // Convert from pointer
		LastCheck:   lastCheck,   // Convert from pointer
		Latency:     proxy.Latency, // Already int64
		IsActive:    proxy.Enabled,
		SuccessRate: proxy.SuccessRate,
		UsageCount:  proxy.UsageCount,
		ErrorCount:  proxy.ErrorCount,
		MaxRetries:  proxy.MaxRetries,
		Timeout:     proxy.Timeout,
		SuccessCount: proxy.Metrics.SuccessCount,
		FailureCount: proxy.Metrics.FailureCount,
		TotalRequests: proxy.Metrics.TotalRequests,
		AvgLatency:   proxy.Metrics.AvgLatency,
		LastStatusCode: proxy.Metrics.LastStatusCode,
	}
}
