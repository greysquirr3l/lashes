package mock

import (
	"net/url"
	"sync"
)

// URLParserFunc is the type for URL parsing functions
type URLParserFunc func(string) (*url.URL, error)

var (
	// defaultURLParser holds the current URL parsing function
	defaultURLParser URLParserFunc = url.Parse

	// urlParserMu is a mutex to protect access to defaultURLParser
	urlParserMu sync.RWMutex
)

// ParseURL provides a mockable function for URL parsing
func ParseURL(rawURL string) (*url.URL, error) {
	urlParserMu.RLock()
	parser := defaultURLParser
	urlParserMu.RUnlock()

	return parser(rawURL)
}

// SetURLParser temporarily replaces the URL parsing function for testing
func SetURLParser(fn URLParserFunc) func() {
	urlParserMu.Lock()
	original := defaultURLParser
	defaultURLParser = fn
	urlParserMu.Unlock()

	// Return function to restore the original parser
	return func() {
		urlParserMu.Lock()
		defaultURLParser = original
		urlParserMu.Unlock()
	}
}

// Expose the mock URL parser to be used in rotator.go
func UseURLParserForTesting(input string) (*url.URL, error) {
	return ParseURL(input)
}
