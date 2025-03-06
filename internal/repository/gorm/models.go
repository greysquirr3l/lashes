package gorm

import (
	"net/url"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

type ProxyModel struct {
	ID             string `gorm:"primaryKey"`
	URL            string `gorm:"not null"`
	Type           string `gorm:"not null"`
	LastUsed       time.Time
	LastCheck      time.Time
	Latency        time.Duration
	IsActive       bool
	Weight         int
	MaxRetries     int
	Timeout        time.Duration
	SuccessCount   int64
	FailureCount   int64
	TotalRequests  int64
	AvgLatency     time.Duration
	LastStatusCode int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *ProxyModel) ToDomain() (*domain.Proxy, error) {
	parsedURL, err := url.Parse(m.URL)
	if err != nil {
		return nil, err
	}

	return &domain.Proxy{
		ID:         m.ID,
		URL:        parsedURL,
		Type:       domain.ProxyType(m.Type),
		LastUsed:   m.LastUsed,
		LastCheck:  m.LastCheck,
		Latency:    m.Latency,
		IsActive:   m.IsActive,
		Weight:     m.Weight,
		MaxRetries: m.MaxRetries,
		Timeout:    m.Timeout,
		Metrics: domain.ProxyMetrics{
			SuccessCount:   m.SuccessCount,
			FailureCount:   m.FailureCount,
			TotalRequests:  m.TotalRequests,
			AvgLatency:     m.AvgLatency,
			LastStatusCode: m.LastStatusCode,
			Created:        m.CreatedAt,
			Updated:        m.UpdatedAt,
		},
	}, nil
}
