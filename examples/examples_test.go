package examples

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/greysquirr3l/lashes/internal/client"
	"github.com/greysquirr3l/lashes/internal/client/mock"
	"github.com/greysquirr3l/lashes/internal/domain"
)

// TestExampleFunctions ensures all example functions run without panicking
func TestExampleFunctions(t *testing.T) {
	// Setup a mock client for all examples
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}

	// Create a URL parser that always succeeds in the examples
	originalParse := url.Parse
	resetURLParser := mock.SetURLParser(func(rawURL string) (*url.URL, error) {
		// Always return a valid URL during tests, regardless of input
		return originalParse("http://example.com:8080")
	})
	defer resetURLParser()

	// Setup mock client creator for all tests
	resetClient := client.SetClientCreator(func(proxy *domain.Proxy, options client.Options) (*http.Client, error) {
		return mock.CreateMockClient(proxy, mockResp)
	})
	defer resetClient()

	tests := []struct {
		name     string
		function func()
	}{
		{"BasicUsageExample", BasicUsageExample},
		{"CustomStorageExample", CustomStorageExample},
		{"CustomRepositoryExample", CustomRepositoryExample},
		{"RoundRobinStrategyExample", RoundRobinStrategyExample},
		{"RandomStrategyExample", RandomStrategyExample},
		{"ProxyValidationExample", ProxyValidationExample},
		{"ErrorHandlingExample", ErrorHandlingExample},
		{"MetricsExample", MetricsExample},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This just verifies the example doesn't panic
			tt.function()
		})
	}
}
