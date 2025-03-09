// Package rotation implements proxy rotation strategies.
//
// This package provides various algorithms for selecting proxies from a pool:
//
// - Round-robin: Rotate through proxies sequentially
// - Random: Select proxies at random with equal probability
// - Weighted: Select proxies based on their assigned weights
// - LeastUsed: Prioritize proxies with lower usage counts
//
// All strategies implement the Strategy interface, which provides
// a consistent API for proxy selection.
package rotation
