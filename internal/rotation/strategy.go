package rotation

import (
	"context"
	"errors"

	"github.com/greysquirr3l/lashes/internal/domain"
)

type Strategy interface {
	Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error)
}

type StrategyType string

const (
	RoundRobin StrategyType = "round-robin"
	Random     StrategyType = "random"
	Weighted   StrategyType = "weighted"
	LeastUsed  StrategyType = "least-used"
)

var (
	ErrNoProxiesAvailable = errors.New("no proxies available")
	ErrInvalidStrategy    = errors.New("invalid rotation strategy")
)
