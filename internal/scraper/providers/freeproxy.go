package providers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/scraper"
)

type freeProxyProvider struct {
	client     *http.Client
	lastUpdate time.Time
}

func NewFreeProxyProvider() scraper.Provider {
	return &freeProxyProvider{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (p *freeProxyProvider) Name() string {
	return "FreeProxy"
}

func (p *freeProxyProvider) Fetch(ctx context.Context) ([]*domain.Proxy, error) {
	sources := []string{
		"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
		"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/proxy.txt",
		"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	}

	var allProxies []*domain.Proxy

	for _, source := range sources {
		req, err := http.NewRequestWithContext(ctx, "GET", source, nil)
		if err != nil {
			continue
		}

		resp, err := p.client.Do(req)
		if err != nil {
			continue
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				// Consider logging this error
			}
		}()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		// Parse response based on format
		if strings.Contains(source, "api.proxyscrape.com") {
			proxies, err := p.parseProxyScrape(resp)
			if err == nil {
				allProxies = append(allProxies, proxies...)
			}
		} else {
			proxies, err := p.parseTextList(resp)
			if err == nil {
				allProxies = append(allProxies, proxies...)
			}
		}
	}

	if len(allProxies) == 0 {
		return nil, scraper.ErrFetchFailed
	}

	p.lastUpdate = time.Now()
	return allProxies, nil
}

func (p *freeProxyProvider) LastUpdated() time.Time {
	return p.lastUpdate
}
