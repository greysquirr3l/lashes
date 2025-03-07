package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
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
	Timeout         interface{}
	MaxRetries      int
	FollowRedirects bool
	Headers         http.Header
	VerifyCerts     bool
}

// ClientCreator is the function type for creating HTTP clients
type ClientCreator func(proxy *domain.Proxy, options Options) (*http.Client, error)

var (
	defaultClientCreator ClientCreator = createDefaultClient
	clientCreatorMu      sync.RWMutex
)

// NewClient creates a new HTTP client using the given proxy
func NewClient(proxy *domain.Proxy, options Options) (*http.Client, error) {
	clientCreatorMu.RLock()
	creator := defaultClientCreator
	clientCreatorMu.RUnlock()
	return creator(proxy, options)
}

// createDefaultClient is the default implementation for creating HTTP clients
func createDefaultClient(proxy *domain.Proxy, options Options) (*http.Client, error) {
	// Parse the proxy URL
	proxyURL, err := url.Parse(proxy.URL)
	if err != nil {
		return nil, err
	}

	// Create the transport with proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	
	// Configure TLS settings
	configureTLSForTransport(transport, options.VerifyCerts)
	
	// Configure headers for future requests
	headers := configureHeaders(options.Headers)
	
	// Configure timeout
	timeout := parseTimeout(options.Timeout)

	// Create the HTTP client
	httpClient := &http.Client{
		Transport: &headerTransport{
			rt:      transport,
			headers: headers,
		},
		Timeout:   timeout,
		CheckRedirect: createRedirectHandler(options.FollowRedirects),
	}

	return httpClient, nil
}

// configureTLSForTransport sets up the TLS configuration for a transport
func configureTLSForTransport(transport *http.Transport, verifyCerts bool) {
	// Only modify TLS configuration if certificates should not be verified
	if !verifyCerts {
		// SECURITY: InsecureSkipVerify is used only when explicitly requested by the user.
		transport.TLSClientConfig = &tls.Config{
			// #nosec G402 -- Explicitly documented security trade-off
			InsecureSkipVerify: true,
		}
	}
}

// configureHeaders sets up the request headers
func configureHeaders(headers http.Header) http.Header {
	if headers == nil {
		headers = make(http.Header)
	}
	
	if headers.Get("User-Agent") == "" {
		headers.Set("User-Agent", agent.GetRandomUserAgent())
	}
	
	return headers
}

// parseTimeout converts various timeout formats to time.Duration
func parseTimeout(timeoutValue interface{}) time.Duration {
	switch t := timeoutValue.(type) {
	case time.Duration:
		return t
	case int:
		return time.Duration(t) * time.Second
	case int64:
		return time.Duration(t) * time.Second
	case float64:
		return time.Duration(t * float64(time.Second))
	case string:
		// Try to parse the string as a duration
		if parsedDuration, err := time.ParseDuration(t); err == nil {
			return parsedDuration
		}
	}
	
	// Default timeout if type is unknown or parsing fails
	return 30 * time.Second
}
// headerTransport adds default headers to all requests
type headerTransport struct {
	rt      http.RoundTripper
	headers http.Header
}

// RoundTrip implements the http.RoundTripper interface
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqCopy := req.Clone(req.Context())
	
	// Add default headers if not already set
	for key, values := range t.headers {
		if len(reqCopy.Header.Values(key)) == 0 {
			for _, value := range values {
				reqCopy.Header.Add(key, value)
			}
		}
	}
	
	// Call the underlying transport
	return t.rt.RoundTrip(reqCopy)
}

// createRedirectHandler returns a CheckRedirect function based on whether redirects should be followed
func createRedirectHandler(followRedirects bool) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if (!followRedirects) {
			return http.ErrUseLastResponse
		}
		return nil
	}
}

// SetClientCreator sets a custom client creator function for testing
func SetClientCreator(creator ClientCreator) func() {
	clientCreatorMu.Lock()
	prev := defaultClientCreator
	defaultClientCreator = creator
	clientCreatorMu.Unlock()
	
	// Return a reset function
	return func() {
		clientCreatorMu.Lock()
		defaultClientCreator = prev
		clientCreatorMu.Unlock()
	}
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

// ConfigureTLS configures the TLS settings for the client
func (c *Client) ConfigureTLS(verifyTLS bool) {
	// Create a new transport based on current TLS settings
	transport := &http.Transport{}
	
	if !verifyTLS {
		// SECURITY: InsecureSkipVerify is used only when explicitly requested.
		// This setting bypasses certificate validation and should be used with caution.
		transport.TLSClientConfig = &tls.Config{
			// #nosec G402 -- Explicitly documented security trade-off, see .snyk file for policy
			InsecureSkipVerify: true, 
		}
	} else {
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12, // Enforce minimum TLS version
		}
	}
	
	// Set the transport on the client
	c.Transport = transport
}

// GetMetrics returns the client's metrics
func (c *Client) GetMetrics() *domain.Metrics {
	return c.metrics
}
