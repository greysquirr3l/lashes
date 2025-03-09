package pool_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/greysquirr3l/lashes/internal/pool"
)

func TestRequestPool(t *testing.T) {
	p := pool.NewRequestPool()

	t.Run("Get returns new request", func(t *testing.T) {
		req := p.Get()
		if req == nil {
			t.Fatal("Get() returned nil")
		}
	})

	t.Run("Put resets request", func(t *testing.T) {
		req := p.Get()

		// Modify the request
		req.Method = "POST"
		req.URL, _ = url.Parse("https://example.com")
		req.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
		req.Body = http.NoBody
		req.ContentLength = 10
		req.Host = "example.com"

		// Return it to the pool
		p.Put(req)

		// Get another request - may be the same instance or a new one
		req2 := p.Get()

		// Verify it's been reset
		if req2.Method != "" {
			t.Errorf("Method not reset, got %q", req2.Method)
		}
		if req2.URL != nil {
			t.Errorf("URL not reset, got %v", req2.URL)
		}
		if len(req2.Header) != 0 {
			t.Errorf("Header not reset, got %v", req2.Header)
		}
		if req2.Body != nil {
			t.Error("Body not reset")
		}
		if req2.ContentLength != 0 {
			t.Errorf("ContentLength not reset, got %d", req2.ContentLength)
		}
		if req2.Host != "" {
			t.Errorf("Host not reset, got %q", req2.Host)
		}
	})
}

func TestResponsePool(t *testing.T) {
	p := pool.NewResponsePool()

	t.Run("Get returns new response", func(t *testing.T) {
		resp := p.Get()
		if resp == nil {
			t.Fatal("Get() returned nil")
		}
		// No need to close body as it's nil for a fresh response
		if resp.Body != nil {
			resp.Body.Close()
		}
	})

	t.Run("Put resets response", func(t *testing.T) {
		resp := p.Get()
		// Make sure response body is nil before we create a new one
		if resp.Body != nil {
			resp.Body.Close()
		}

		// Create a body that needs closing
		bodyContent := "test body"
		resp.Body = io.NopCloser(strings.NewReader(bodyContent))

		// Modify the response
		resp.Status = "200 OK"
		resp.StatusCode = 200
		resp.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
		resp.ContentLength = int64(len(bodyContent))
		resp.Request = &http.Request{}

		// Return it to the pool (this will close the body)
		p.Put(resp)

		// Get another response - may be the same instance or a new one
		resp2 := p.Get()
		// Make sure we close any body if it exists (should be nil after reset)
		if resp2.Body != nil {
			resp2.Body.Close()
		}

		// Verify it's been reset
		if resp2.Status != "" {
			t.Errorf("Status not reset, got %q", resp2.Status)
		}
		if resp2.StatusCode != 0 {
			t.Errorf("StatusCode not reset, got %d", resp2.StatusCode)
		}
		if len(resp2.Header) != 0 {
			t.Errorf("Header not reset, got %v", resp2.Header)
		}
		if resp2.Body != nil {
			t.Error("Body not reset")
		}
		if resp2.ContentLength != 0 {
			t.Errorf("ContentLength not reset, got %d", resp2.ContentLength)
		}
		if resp2.Request != nil {
			t.Error("Request not reset")
		}
	})
}
