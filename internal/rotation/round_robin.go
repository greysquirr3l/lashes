package rotation

import (
	"context"
	"sort"
	"sync"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// roundRobinStrategy implements the Strategy interface
type roundRobinStrategy struct {
	current int // current index in the rotation
	mu      *sync.Mutex
}

// NewRoundRobinStrategy creates a new round robin strategy
func NewRoundRobinStrategy() Strategy {
	return &roundRobinStrategy{
		current: -1,
		mu:      &sync.Mutex{},
	}
}

// Next selects the next proxy in sequence
func (s *roundRobinStrategy) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// Sort proxies by URL for consistent ordering
	sort.SliceStable(proxies, func(i, j int) bool {
		return proxies[i].URL < proxies[j].URL
	})

	s.mu.Lock()
	defer s.mu.Unlock()

	s.current = (s.current + 1) % len(proxies)
	return proxies[s.current], nil
}
