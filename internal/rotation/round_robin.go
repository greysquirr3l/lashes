package rotation

import (
	"context"
	"sort"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// RoundRobinStrategyImpl implements a round robin proxy rotation strategy
type roundRobinStrategyImpl struct {
	index int // tracks the current position in the rotation
}

// NewRoundRobinStrategy creates a new instance of the round robin rotation strategy
func NewRoundRobinStrategy() *roundRobinStrategyImpl {
	return &roundRobinStrategyImpl{index: 0}
}

// Implementation of the round robin rotation strategy
func (s *roundRobinStrategyImpl) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}
	// Sort proxies by URL for deterministic rotation
	sort.SliceStable(proxies, func(i, j int) bool {
		return proxies[i].URL < proxies[j].URL
	})
	proxy := proxies[s.index%len(proxies)]
	s.index = (s.index + 1) % len(proxies)
	return proxy, nil
}