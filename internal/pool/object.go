package pool

import (
	"net/http"
	"sync"
)

// RequestPool provides a pool of reusable HTTP requests
type RequestPool struct {
	pool sync.Pool
}

// NewRequestPool creates a new pool for HTTP requests
func NewRequestPool() *RequestPool {
	return &RequestPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &http.Request{}
			},
		},
	}
}

// Get retrieves an HTTP request from the pool
func (p *RequestPool) Get() *http.Request {
	// We can safely ignore the ok return value since our pool.New function
	// always returns a *http.Request, so the type assertion will never fail
	req, ok := p.pool.Get().(*http.Request)
	if !ok {
		// This should never happen, but if it does, return a new request
		// to maintain API compatibility
		return &http.Request{}
	}
	return req
}

// Put returns an HTTP request to the pool after resetting it
func (p *RequestPool) Put(req *http.Request) {
	// Reset request to be reused
	req.Method = ""
	req.URL = nil
	req.Header = make(http.Header)
	req.Body = nil
	req.GetBody = nil
	req.ContentLength = 0
	req.TransferEncoding = nil
	req.Close = false
	req.Host = ""
	req.Form = nil
	req.PostForm = nil
	req.MultipartForm = nil
	req.Trailer = nil

	p.pool.Put(req)
}

// ResponsePool provides a pool of reusable HTTP responses
type ResponsePool struct {
	pool sync.Pool
}

// NewResponsePool creates a new pool for HTTP responses
func NewResponsePool() *ResponsePool {
	return &ResponsePool{
		pool: sync.Pool{
			New: func() interface{} {
				return &http.Response{}
			},
		},
	}
}

// Get retrieves an HTTP response from the pool
func (p *ResponsePool) Get() *http.Response {
	// We can safely ignore the ok return value since our pool.New function
	// always returns a *http.Response, so the type assertion will never fail
	resp, ok := p.pool.Get().(*http.Response)
	if !ok {
		// This should never happen, but if it does, return a new response
		// to maintain API compatibility
		return &http.Response{}
	}
	return resp
}

// Put returns an HTTP response to the pool after resetting it
func (p *ResponsePool) Put(resp *http.Response) {
	// Make sure body is closed
	if resp.Body != nil {
		// We intentionally ignore the error from Close() as we're just cleaning up
		// Adding explicit error ignoring to satisfy linter
		closeErr := resp.Body.Close()
		_ = closeErr // Explicitly ignoring error
	}

	// Reset response to be reused
	resp.Status = ""
	resp.StatusCode = 0
	resp.Header = make(http.Header)
	resp.Body = nil
	resp.ContentLength = 0
	resp.TransferEncoding = nil
	resp.Close = false
	resp.Uncompressed = false
	resp.Trailer = nil
	resp.Request = nil

	p.pool.Put(resp)
}
