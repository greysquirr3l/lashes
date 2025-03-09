package lashes

import (
	"context"
	"net/http"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// DiscoveryProvider defines an interface for services that can discover proxies
type DiscoveryProvider interface {
	// GetProxies returns a list of proxies from the provider
	GetProxies(ctx context.Context) ([]*Proxy, error)
}

// DiscoveryOptions configures the behavior of proxy discovery
type DiscoveryOptions struct {
	// Timeout for discovery operations
	Timeout time.Duration

	// UserAgent used for making requests to proxy providers
	UserAgent string

	// Validate determines whether discovered proxies should be validated
	Validate bool

	// ValidationURL is the URL used for validating discovered proxies
	ValidationURL string

	// Client is the HTTP client to use for discovery operations
	Client *http.Client
}

// DefaultDiscoveryOptions returns sensible default options for proxy discovery
func DefaultDiscoveryOptions() DiscoveryOptions {
	return DiscoveryOptions{
		Timeout:       30 * time.Second,
		UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		Validate:      true,
		ValidationURL: "https://api.ipify.org?format=json",
		Client:        http.DefaultClient,
	}
}

// ParseProxiesFromText parses a text list of proxies in IP:Port format
func ParseProxiesFromText(text string, proxyType ProxyType) ([]*Proxy, error) {
	// Create parser based on existing internal implementation
	proxies := []*domain.Proxy{} // Initialize as empty slice instead of nil
	// ... parsing logic ...

	// Convert to public types
	result := make([]*Proxy, len(proxies))
	for i, p := range proxies {
		result[i] = p
	}

	return result, nil
}

// ImportProxies adds multiple proxies to the rotator
func (r *rotator) ImportProxies(ctx context.Context, proxies []*Proxy) (int, error) {
	var imported int

	for _, proxy := range proxies {
		err := r.AddProxy(ctx, proxy.URL, proxy.Type)
		if err != nil {
			// Continue with other proxies even if some fail
			continue
		}
		imported++
	}

	return imported, nil
}
