package domain

import (
	"net/http"
	"net/url"
	"time"
)

type ProxyType string

const (
	HTTP   ProxyType = "http"
	SOCKS4 ProxyType = "socks4"
	SOCKS5 ProxyType = "socks5"
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

type Proxy struct {
	ID         string
	URL        *url.URL
	Type       ProxyType
	LastUsed   time.Time
	LastCheck  time.Time
	Latency    time.Duration
	IsActive   bool
	Metrics    ProxyMetrics
	Weight     int           // For weighted rotation
	MaxRetries int           // Maximum retry attempts
	Timeout    time.Duration // Proxy-specific timeout
	Settings   ProxySettings
}
