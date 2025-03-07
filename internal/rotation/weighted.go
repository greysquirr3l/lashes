package rotation

import (
	"context"
	"math/rand"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// weightedStrategyImpl implements weighted random proxy selection.
type weightedStrategyImpl struct{}

// NewWeightedStrategy creates a new instance of weighted strategy.
func NewWeightedStrategy() *weightedStrategyImpl {
	return &weightedStrategyImpl{}
}

// Next selects the next proxy based on weight.
// Proxies with zero or negative weight are given a minimal effective weight (epsilon)
// so they are less likely to be selected without being excluded entirely.
func (s *weightedStrategyImpl) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	epsilon := 0.001
	totalWeight := 0.0
	// Always include all proxies assigning an effective weight.
	effectiveWeights := make([]float64, len(proxies))
	for i, p := range proxies {
		if p.Weight <= 0 {
			effectiveWeights[i] = epsilon
		} else {
			effectiveWeights[i] = float64(p.Weight)
		}
		totalWeight += effectiveWeights[i]
	}
	// If totalWeight is zero, fallback to equal random selection.
	if totalWeight == 0 {
		return proxies[rand.Intn(len(proxies))], nil
	}
	
	r := rand.Float64() * totalWeight
	sum := 0.0
	for i, p := range proxies {
		sum += effectiveWeights[i]
		if r < sum {
			return p, nil
		}
	}
	
	return proxies[len(proxies)-1], nil
}