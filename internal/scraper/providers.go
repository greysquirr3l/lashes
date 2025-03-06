package scraper

import (
	"context"
	"errors"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

type Provider interface {
	Name() string
	Fetch(ctx context.Context) ([]*domain.Proxy, error)
	LastUpdated() time.Time
}

// Common free proxy list providers
const (
	ProxyListDownload = "https://www.proxy-list.download/api/v1/get"
	FreeProxyList     = "https://free-proxy-list.net/"
	SSLProxies        = "https://www.sslproxies.org/"
	ProxyNova         = "https://www.proxynova.com/proxy-server-list/"
	Proxyscrape       = "https://api.proxyscrape.com/v2/"
)

var (
	ErrFetchFailed   = errors.New("failed to fetch proxy list")
	ErrInvalidFormat = errors.New("invalid proxy list format")
	ErrRateLimited   = errors.New("rate limited by provider")
)
