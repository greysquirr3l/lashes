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
	Enabled        bool      `gorm:"default:true"` // Renamed from IsActive
	Latency        int64     // Store as int64 in milliseconds
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
	// Create a pointer to LastUsed for the domain model
	var lastUsed *time.Time
	if !m.LastUsed.IsZero() {
		t := m.LastUsed // Make a copy
		lastUsed = &t
	}

	proxy := &domain.Proxy{
		ID:          m.ID,
		URL:         m.URL,
		Type:        domain.ProxyType(m.Type),
		Username:    m.Username,
		Password:    m.Password,
		CountryCode: m.CountryCode,
		Weight:      m.Weight,
		LastUsed:    lastUsed,
		Latency:     m.Latency,
		Enabled:     m.Enabled,
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
	}

	return proxy, nil
}

// FromDomain converts a domain model to a GORM model
func FromDomain(proxy *domain.Proxy) *ProxyModel {
	// Convert pointers to values, using zero values if nil
	var lastUsed time.Time
	if proxy.LastUsed != nil {
		lastUsed = *proxy.LastUsed
	}

	return &ProxyModel{
		ID:             proxy.ID,
		URL:            proxy.URL,
		Type:           string(proxy.Type),
		Username:       proxy.Username,
		Password:       proxy.Password,
		CountryCode:    proxy.CountryCode,
		Weight:         proxy.Weight,
		LastUsed:       lastUsed,
		Enabled:        proxy.Enabled,
		Latency:        proxy.Latency,
		SuccessRate:    proxy.SuccessRate,
		UsageCount:     proxy.UsageCount,
		ErrorCount:     proxy.ErrorCount,
		MaxRetries:     proxy.MaxRetries,
		Timeout:        proxy.Timeout,
		SuccessCount:   proxy.Metrics.SuccessCount,
		FailureCount:   proxy.Metrics.FailureCount,
		TotalRequests:  proxy.Metrics.TotalRequests,
		AvgLatency:     proxy.Metrics.AvgLatency,
		LastStatusCode: proxy.Metrics.LastStatusCode,
	}
}
