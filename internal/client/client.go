package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/greysquirr3l/lashes/internal/agent"
	"github.com/greysquirr3l/lashes/internal/domain"
)

type Client struct {
	http.Client // Embed http.Client
	maxRetries  int
	metrics     *domain.Metrics
}

type Options struct {
	Timeout         time.Duration
	MaxRetries      int
	FollowRedirects bool
	Headers         http.Header
	VerifyCerts     bool
}

func NewClient(proxy *domain.Proxy, opts Options) (*http.Client, error) {
	proxyURL, err := url.Parse(proxy.URL.String())
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	
	// Only modify TLS configuration if certificates should not be verified
	// This is explicitly set by the caller and is their responsibility
	if !opts.VerifyCerts {
		transport.TLSClientConfig = &tls.Config{
			// #nosec G402 -- InsecureSkipVerify is set only when needed for testing or non-production scenarios
			InsecureSkipVerify: true,
		}
	}

	// Set default headers if not provided
	if opts.Headers == nil {
		opts.Headers = make(http.Header)
	}
	if opts.Headers.Get("User-Agent") == "" {
		opts.Headers.Set("User-Agent", agent.GetRandomUserAgent())
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   opts.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !opts.FollowRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	return httpClient, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	
	startTime := time.Now()
	
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err = c.Client.Do(req)
		if err == nil {
			// Calculate request latency
			latency := time.Since(startTime)
			c.metrics.IncrementLatency(latency)
			
			// Update metrics based on status code
			if resp.StatusCode < 400 {
				c.metrics.RecordSuccess(resp.StatusCode)
			} else {
				c.metrics.RecordFailure(resp.StatusCode)
			}
			return resp, nil
		}
	}
	
	// Record failure if all attempts failed
	c.metrics.RecordFailure(0)
	return nil, err
}

// GetMetrics returns the client's metrics
func (c *Client) GetMetrics() *domain.Metrics {
	return c.metrics
}
