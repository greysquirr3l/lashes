package domain

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ProxyType defines the type of proxy
type ProxyType string

// Available proxy types
const (
	HTTPProxy   ProxyType = "http"
	SOCKS4Proxy ProxyType = "socks4"
	SOCKS5Proxy ProxyType = "socks5"
)

// Alias constants for better readability
const (
	HTTP   = HTTPProxy
	SOCKS4 = SOCKS4Proxy
	SOCKS5 = SOCKS5Proxy
)

type ProxyMetrics struct {
	SuccessCount   int64
	FailureCount   int64
	TotalRequests  int64
	AvgLatency     time.Duration
	LastStatusCode int
	Created        time.Time
	Updated        time.Time
}

type ProxySettings struct {
	FollowRedirects bool
	VerifyCerts     bool
	Headers         map[string][]string
	Cookies         []*http.Cookie
	UserAgent       string
}

// Proxy represents a proxy server configuration
type Proxy struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	URL         string     `json:"url"` // Standardized to string representation
	Type        ProxyType  `json:"type"`
	Username    string     `json:"username,omitempty"`
	Password    string     `json:"password,omitempty"`
	CountryCode string     `json:"country_code,omitempty"`
	Weight      int        `json:"weight" gorm:"default:1"`
	LastUsed    *time.Time `json:"last_used,omitempty"`
	Enabled     bool       `json:"enabled" gorm:"default:true"`
	Latency     int64      `json:"latency_ms" gorm:"default:0"` // in milliseconds
	SuccessRate float64    `json:"success_rate" gorm:"default:0"`
	UsageCount  int64      `json:"usage_count" gorm:"default:0"`
	ErrorCount  int64      `json:"error_count" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Metrics     ProxyMetrics
	Settings    ProxySettings
	MaxRetries  int           // Maximum retry attempts
	Timeout     time.Duration // Proxy-specific timeout
}

// ParseURL parses the proxy URL string into a URL object
func (p *Proxy) ParseURL() (*url.URL, error) {
	if p.URL == "" {
		return nil, fmt.Errorf("empty URL")
	}
	return url.Parse(p.URL)
}

// IsValid checks if the proxy configuration is valid
func (p *Proxy) IsValid() bool {
	_, err := p.ParseURL()
	return err == nil && p.ID != "" && p.Type != ""
}

// String returns the URL as a string representation of the proxy
func (p *Proxy) String() string {
	return p.URL
}

// GetEnabled returns the proxy's enabled state
func (p *Proxy) GetEnabled() bool {
	return p.Enabled
}

// SetEnabled sets the enabled state
func (p *Proxy) SetEnabled(enabled bool) {
	p.Enabled = enabled
}
