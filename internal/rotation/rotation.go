package rotation

import (
	"context"
	"sync/atomic"

	"github.com/greysquirr3l/lashes/internal/domain"
)

func NewStrategy(strategyType StrategyType) (Strategy, error) {
	switch strategyType {
	case RoundRobin:
		return &roundRobinStrategy{counter: 0}, nil
	// ...add other strategies
	default:
		return nil, ErrInvalidStrategy
	}
}

type roundRobinStrategy struct {
	counter uint64
}

func (s *roundRobinStrategy) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	index := atomic.AddUint64(&s.counter, 1) % uint64(len(proxies))
	return proxies[index], nil
}
