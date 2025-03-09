package rotation

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/greysquirr3l/lashes/internal/domain"
)

// StrategyType defines the type of rotation strategy to use
type StrategyType string

// Available rotation strategies with consistent naming
const (
	RoundRobinStrategy StrategyType = "round-robin"
	RandomStrategy     StrategyType = "random"
	WeightedStrategy   StrategyType = "weighted"
	LeastUsedStrategy  StrategyType = "least-used"
)

// Deprecated: Legacy constants for backward compatibility
// Use the newer versions with "Strategy" suffix instead.
const (
	RoundRobin = RoundRobinStrategy
	Random     = RandomStrategy
	Weighted   = WeightedStrategy
	LeastUsed  = LeastUsedStrategy
)

// Common errors
var (
	ErrNoProxiesAvailable = errors.New("no proxies available")
	ErrInvalidStrategy    = errors.New("invalid rotation strategy")
)

// Strategy defines the interface for proxy rotation strategies
type Strategy interface {
	Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error)
}

// NewStrategy creates a new rotation strategy based on the provided type
func NewStrategy(strategyType StrategyType) (Strategy, error) {
	switch strategyType {
	case RoundRobinStrategy, "RoundRobin":
		return NewRoundRobinStrategy(), nil
	case RandomStrategy, "Random":
		return &randomStrategy{}, nil
	case WeightedStrategy, "Weighted":
		return &weightedStrategy{}, nil
	case LeastUsedStrategy, "LeastUsed":
		return &leastUsedStrategy{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidStrategy, strategyType)
	}
}

// randomStrategy implements a random rotation strategy
type randomStrategy struct{}

func (s *randomStrategy) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable
	}

	// Use crypto/rand instead of math/rand for secure random number generation
	maxBig := big.NewInt(int64(len(proxies)))
	nBig, err := rand.Int(rand.Reader, maxBig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}

	index := nBig.Int64()
	return proxies[index], nil
}

// weightedStrategy implements a weighted rotation strategy based on success rate and explicit weights
type weightedStrategy struct {
	mu sync.Mutex
}

// calculateProxyWeight determines the weight for a single proxy
func calculateProxyWeight(proxy *domain.Proxy) float64 {
	// Base weight is the explicit weight with a minimum of 1
	baseWeight := float64(proxy.Weight)
	if baseWeight <= 0 {
		baseWeight = 1
	}

	// Adjust by success rate (default to 0.5 if no data)
	successRate := proxy.SuccessRate
	if successRate <= 0 && proxy.UsageCount == 0 {
		successRate = 0.5 // Default for new proxies
	} else if successRate <= 0 {
		successRate = 0.1 // Minimum for proxies with failures
	}

	// Combine factors for final weight
	return baseWeight * successRate
}

// calculateWeights computes weights for all proxies
func (s *weightedStrategy) calculateWeights(proxies []*domain.Proxy) ([]float64, float64) {
	weights := make([]float64, len(proxies))
	var totalWeight float64

	for i, proxy := range proxies {
		weights[i] = calculateProxyWeight(proxy)
		totalWeight += weights[i]
	}

	return weights, totalWeight
}

// selectProxyByWeight selects a proxy based on the calculated weights
func (s *weightedStrategy) selectProxyByWeight(proxies []*domain.Proxy, weights []float64, totalWeight float64) (*domain.Proxy, error) {
	// Generate random number in range [0, totalWeight)
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(totalWeight*1000)))
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}

	// Convert to float and scale to our weight range
	selector := float64(nBig.Int64()) / 1000.0

	// Select proxy based on weight
	var cumulativeWeight float64
	for i, weight := range weights {
		cumulativeWeight += weight
		if selector < cumulativeWeight {
			return proxies[i], nil
		}
	}

	// Fallback in case of rounding errors
	return proxies[len(proxies)-1], nil
}

func (s *weightedStrategy) handleEmptyOrSingleProxy(proxies []*domain.Proxy) (*domain.Proxy, error, bool) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable, true
	}

	if len(proxies) == 1 {
		return proxies[0], nil, true
	}

	return nil, nil, false
}

func (s *weightedStrategy) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	// Handle special cases
	if proxy, err, handled := s.handleEmptyOrSingleProxy(proxies); handled {
		return proxy, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Calculate weights for all proxies
	weights, totalWeight := s.calculateWeights(proxies)

	// Handle case where all proxies have zero weight
	if totalWeight <= 0 {
		// Fallback to first proxy
		return proxies[0], nil
	}

	// Select proxy based on weight
	return s.selectProxyByWeight(proxies, weights, totalWeight)
}

// leastUsedStrategy implements a strategy that selects the least used proxy
type leastUsedStrategy struct {
	mu sync.Mutex
}

func (s *leastUsedStrategy) handleEmptyOrSingleProxy(proxies []*domain.Proxy) (*domain.Proxy, error, bool) {
	if len(proxies) == 0 {
		return nil, ErrNoProxiesAvailable, true
	}

	if len(proxies) == 1 {
		return proxies[0], nil, true
	}

	return nil, nil, false
}

func (s *leastUsedStrategy) findMinimumUsageCandidates(proxies []*domain.Proxy) []*domain.Proxy {
	// Find minimum usage count
	minUsage := int64(-1)
	for _, proxy := range proxies {
		if minUsage == -1 || proxy.UsageCount < minUsage {
			minUsage = proxy.UsageCount
		}
	}

	// Collect all proxies with minimum usage
	var candidates []*domain.Proxy
	for _, proxy := range proxies {
		if proxy.UsageCount == minUsage {
			candidates = append(candidates, proxy)
		}
	}

	return candidates
}

func (s *leastUsedStrategy) sortCandidatesByLastUsed(candidates []*domain.Proxy) {
	if len(candidates) > 1 {
		sort.Slice(candidates, func(i, j int) bool {
			// Handle nil LastUsed (never used)
			if candidates[i].LastUsed == nil {
				return true
			}
			if candidates[j].LastUsed == nil {
				return false
			}
			// Compare by last used time
			return candidates[i].LastUsed.Before(*candidates[j].LastUsed)
		})
	}
}

func (s *leastUsedStrategy) Next(ctx context.Context, proxies []*domain.Proxy) (*domain.Proxy, error) {
	// Handle special cases
	if proxy, err, handled := s.handleEmptyOrSingleProxy(proxies); handled {
		return proxy, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	candidates := s.findMinimumUsageCandidates(proxies)
	s.sortCandidatesByLastUsed(candidates)

	return candidates[0], nil
}
