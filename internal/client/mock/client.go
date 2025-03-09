package mock

import (
	"net/http"
	"net/url"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// MockTransport implements http.RoundTripper for testing
type MockTransport struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
	Delay       time.Duration
}

// RoundTrip implements http.RoundTripper
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}
	if m.RoundTripFn != nil {
		return m.RoundTripFn(req)
	}
	return m.Response, nil
}

// NewMockClient creates a mock HTTP client for testing
func NewMockClient(response *http.Response, delay time.Duration) *http.Client {
	return &http.Client{
		Transport: &MockTransport{
			Response: response,
			Delay:    delay,
		},
	}
}

// CreateMockClient creates a mock client that can be used for testing
func CreateMockClient(proxy *domain.Proxy, opts ...interface{}) (*http.Client, error) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}

	// Allow customization through opts if needed
	for _, opt := range opts {
		if resp, ok := opt.(*http.Response); ok {
			mockResp = resp
		}
	}

	return NewMockClient(mockResp, 10*time.Millisecond), nil
}

// URLParseFunc is the type for URL parsing functions
type URLParseFunc func(string) (*url.URL, error)

// DefaultURLParse holds the current URL parsing function
var DefaultURLParse URLParseFunc = url.Parse

// SetupURLParsing temporarily replaces the URL parsing function for testing
func SetupURLParsing(fn URLParseFunc) func() {
	original := DefaultURLParse
	// Replace with mock implementation
	DefaultURLParse = fn
	// Return function to restore the original
	return func() {
		DefaultURLParse = original
	}
}
