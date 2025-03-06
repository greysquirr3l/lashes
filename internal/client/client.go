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
	proxy       *domain.Proxy
	client      *http.Client
	maxRetries  int
	timeout     time.Duration
	verifyCerts bool
}

type Options struct {
	Timeout         time.Duration
	MaxRetries      int
	VerifyCerts     bool
	FollowRedirects bool
	Headers         http.Header
}

func NewClient(proxy *domain.Proxy, opts Options) (*http.Client, error) {
	proxyURL, err := url.Parse(proxy.URL.String())
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !opts.VerifyCerts,
		},
	}

	// Set default headers if not provided
	if opts.Headers == nil {
		opts.Headers = make(http.Header)
	}
	if opts.Headers.Get("User-Agent") == "" {
		opts.Headers.Set("User-Agent", agent.GetRandomUserAgent())
	}

	return &http.Client{
		Transport: transport,
		Timeout:   opts.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !opts.FollowRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err = c.client.Do(req)
		if err == nil {
			// Update metrics
			c.proxy.Metrics.TotalRequests++
			c.proxy.Metrics.LastStatusCode = resp.StatusCode
			if resp.StatusCode < 400 {
				c.proxy.Metrics.SuccessCount++
			} else {
				c.proxy.Metrics.FailureCount++
			}
			return resp, nil
		}
	}

	return nil, err
}
