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
// This implementation heavily favors proxies with positive weights.
func (s *weightedStrategyImpl) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// Filter into positive-weight and zero-weight proxies
	var positiveWeightProxies []*domain.Proxy
	var zeroWeightProxies []*domain.Proxy

	totalPositiveWeight := 0.0

	for _, p := range proxies {
		if p.Weight > 0 {
			positiveWeightProxies = append(positiveWeightProxies, p)
			totalPositiveWeight += float64(p.Weight)
		} else {
			zeroWeightProxies = append(zeroWeightProxies, p)
		}
	}

	// If we have positive-weight proxies, select from those 95% of the time
	if len(positiveWeightProxies) > 0 {
		// Use crypto random to determine if we should pick from positive-weight proxies
		r, err := cryptoRandInt(100)
		if err != nil {
			return nil, err
		}

		// If r < 95, select from positive-weight proxies (95% chance)
		if r < 95 {
			return s.selectFromWeightedList(positiveWeightProxies, totalPositiveWeight)
		}

		// Otherwise, fall through to select from zero-weight proxies if available
	}

	// If we reach here, either:
	// 1. All proxies have zero weight, or
	// 2. We chose to select from zero-weight proxies (5% chance)

	if len(zeroWeightProxies) > 0 {
		// Select randomly from zero-weight proxies
		randomIndex, err := cryptoRandInt(len(zeroWeightProxies))
		if err != nil {
			return nil, err
		}
		return zeroWeightProxies[randomIndex], nil
	}

	// If no zero-weight proxies but we have positive-weight proxies
	if len(positiveWeightProxies) > 0 {
		return s.selectFromWeightedList(positiveWeightProxies, totalPositiveWeight)
	}

	// If we get here, something is wrong with our logic - pick any proxy
	return proxies[0], nil
}

// selectFromWeightedList selects a proxy from a list using weighted random selection
func (s *weightedStrategyImpl) selectFromWeightedList(proxies []*domain.Proxy, totalWeight float64) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// If all weights are zero, select randomly
	if totalWeight == 0 {
		return s.selectRandomProxy(proxies)
	}

	r, err := cryptoRandFloat64(totalWeight)
	if err != nil {
		return nil, err
	}

	sum := 0.0
	for _, p := range proxies {
		weight := float64(p.Weight)
		// Skip proxies with zero weight
		if weight <= 0 {
			continue
		}

		sum += weight
		if r < sum {
			return p, nil
		}
	}

	// Fallback to last proxy with positive weight
	return proxies[len(proxies)-1], nil
}
