package rotation

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// weightedStrategyImpl implements weighted random proxy selection.
type weightedStrategyImpl struct{}

// NewWeightedStrategy creates a new instance of weighted strategy.
func NewWeightedStrategy() *weightedStrategyImpl {
	return &weightedStrategyImpl{}
}

// cryptoRandFloat64 returns a random float64 in [0,total) using crypto/rand.
func cryptoRandFloat64(total float64) (float64, error) {
	max := big.NewInt(math.MaxInt64)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}
	return (float64(n.Int64()) / float64(math.MaxInt64)) * total, nil
}

// cryptoRandInt returns a random int in [0,max) using crypto/rand.
func cryptoRandInt(max int) (int, error) {
	if max <= 0 {
		return 0, nil
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// selectRandomProxy selects a proxy randomly without considering weights.
func (s *weightedStrategyImpl) selectRandomProxy(proxies []*domain.Proxy) (*domain.Proxy, error) {
	randomIndex, err := cryptoRandInt(len(proxies))
	if err != nil {
		return nil, fmt.Errorf("failed to generate random index: %w", err)
	}
	return proxies[randomIndex], nil
}

// Next selects the next proxy based on weight.
// All proxies are always considered; proxies with non-positive weight are given an effective weight epsilon.
func (s *weightedStrategyImpl) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// Lower epsilon to 0.0001 to make zero-weight proxies far less likely.
	const epsilon = 0.0001
	totalWeight := 0.0

	// Calculate the total effective weight for every proxy.
	for _, p := range proxies {
		weight := float64(p.Weight)
		if weight <= 0 {
			weight = epsilon
		}
		totalWeight += weight
	}

	// If totalWeight is zero, fallback to a random selection.
	if totalWeight == 0 {
		return s.selectRandomProxy(proxies)
	}

	r, err := cryptoRandFloat64(totalWeight)
	if err != nil {
		return nil, err
	}

	sum := 0.0
	// Select a proxy based on the weighted random number.
	for _, p := range proxies {
		weight := float64(p.Weight)
		if weight <= 0 {
			weight = epsilon
		}
		sum += weight
		if r < sum {
			return p, nil
		}
	}

	return proxies[len(proxies)-1], nil
}