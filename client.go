package lashes

import (
	"context"
	"net/http"
	"time"

	"github.com/greysquirr3l/lashes/internal/client"
)

// ClientOptions configures the behavior of HTTP clients created with proxy support.
type ClientOptions struct {
	// Timeout for requests
	Timeout time.Duration

	// MaxRetries sets the number of retry attempts for failed requests
	MaxRetries int

	// FollowRedirects determines whether HTTP redirects are followed
	FollowRedirects bool

	// VerifyCerts determines whether SSL certificates are verified
	VerifyCerts bool

	// Headers are default headers to include in all requests
	Headers http.Header

	// UserAgent overrides the default rotating User-Agent
	UserAgent string
}

// DefaultClientOptions returns sensible default options for proxy clients.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		FollowRedirects: true,
		VerifyCerts:     true,
	}
}

// NewClient creates an http.Client configured with the given proxy.
func NewClient(proxy *Proxy, opts ClientOptions) (*http.Client, error) {
	return client.NewClient(proxy, client.Options{
		Timeout:         opts.Timeout,
		MaxRetries:      opts.MaxRetries,
		FollowRedirects: opts.FollowRedirects,
		VerifyCerts:     opts.VerifyCerts,
		Headers:         opts.Headers,
	})
}

// GetNextClient returns an http.Client configured with the next proxy
// from the rotation according to the configured strategy.
func (r *rotator) GetNextClient(ctx context.Context) (*http.Client, error) {
	proxy, err := r.GetProxy(ctx)
	if err != nil {
		return nil, err
	}

	return client.NewClient(proxy, client.Options{
		Timeout:         r.opts.RequestTimeout,
		MaxRetries:      r.opts.MaxRetries,
		VerifyCerts:     true,
		FollowRedirects: true,
	})
}
