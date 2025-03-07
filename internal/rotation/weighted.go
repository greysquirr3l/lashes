package rotation

import (
	"context"
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/greysquirr3l/lashes/internal/domain"
)

type weightedStrategy struct {
	mu sync.Mutex
}

func NewWeightedStrategy() Strategy {
	return &weightedStrategy{}
}

func (s *weightedStrategy) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	totalWeight := 0
	for _, proxy := range proxies {
		totalWeight += proxy.Weight
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(totalWeight)))
	if err != nil {
		return proxies[0], nil
	}
	randWeight := int(n.Int64())

	for _, proxy := range proxies {
		if randWeight < proxy.Weight {
			return proxy, nil
		}
		randWeight -= proxy.Weight
	}

	return proxies[0], nil // Fallback to first proxy
}
