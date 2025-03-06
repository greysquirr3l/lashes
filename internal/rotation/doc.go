// Package rotation implements proxy rotation strategies.
//
// Available strategies:
//   - Round-robin: Rotates through proxies sequentially
//   - Random: Selects proxies randomly
//   - Weighted: Uses proxy weights to influence selection probability
//   - Least-used: Prioritizes proxies with fewer requests
//
// All strategies are thread-safe and can be used concurrently.
package rotation
