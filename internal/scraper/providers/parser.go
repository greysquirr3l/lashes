package providers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/greysquirr3l/lashes/internal/domain"
)

func (p *freeProxyProvider) parseTextList(resp *http.Response) ([]*domain.Proxy, error) {
	var proxies []*domain.Proxy
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Handle IP:Port format
		parsedURL, err := url.Parse("http://" + line)
		if err != nil {
			continue
		}

		proxy := &domain.Proxy{
			ID:         uuid.New().String(),
			URL:        parsedURL.String(), // Convert URL to string
			Type:       domain.HTTPProxy,   // Use proper constant
			Enabled:    true,
			MaxRetries: 3,
		}
		proxies = append(proxies, proxy)
	}

	return proxies, scanner.Err()
}

func (p *freeProxyProvider) parseProxyScrape(resp *http.Response) ([]*domain.Proxy, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(body), "\n")
	var proxies []*domain.Proxy

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parsedURL, err := url.Parse("http://" + line)
		if err != nil {
			continue
		}

		proxy := &domain.Proxy{
			ID:         uuid.New().String(),
			URL:        parsedURL.String(), // Convert URL to string
			Type:       domain.HTTPProxy,   // Use proper constant
			Enabled:    true,
			MaxRetries: 3,
		}
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// ParseProxyLine parses a line of text into a proxy object.
// This exported function can be used by external packages needing to convert
// IP:Port formatted strings to proxy objects.
func ParseProxyLine(line string) (*domain.Proxy, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	parsedURL, err := url.Parse("http://" + line)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
	}

	proxy := &domain.Proxy{
		ID:         uuid.New().String(),
		URL:        parsedURL.String(), // Convert URL to string
		Type:       domain.HTTPProxy,   // Use proper constant
		Enabled:    true,
		MaxRetries: 3,
	}

	return proxy, nil
}
