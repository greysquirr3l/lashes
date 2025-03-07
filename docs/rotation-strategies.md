# Rotation Strategy Documentation

## Available Strategies

### 1. Round-Robin

The round-robin strategy cycles through proxies in a deterministic order. This strategy is ideal when you want to distribute load evenly across all proxies and need predictable, sequential proxy selection.

```go
// Initialize with round-robin strategy (default)
opts := lashes.DefaultOptions() // Round-robin is the default
// or explicitly:
opts.Strategy = rotation.RoundRobinStrategy

rotator, err := lashes.New(opts)
```

Features:
- Deterministic rotation order (proxies are sorted by URL)
- Each proxy is used exactly once per rotation cycle
- O(1) time complexity for selection

### 2. Random

The random strategy selects proxies randomly using cryptographically secure randomization. This strategy is useful when you want to make traffic patterns unpredictable to avoid detection.

```go
// Initialize with random strategy
opts := lashes.DefaultOptions()
opts.Strategy = rotation.RandomStrategy

rotator, err := lashes.New(opts)
```

Features:
- Unpredictable selection using crypto/rand for secure randomization
- Equal probability for all proxies
- O(1) time complexity for selection

### 3. Weighted

The weighted strategy selects proxies based on their assigned weights. Proxies with higher weights are selected more frequently. Proxies with zero or negative weights are selected significantly less often (5% probability).

```go
// Initialize with weighted strategy
opts := lashes.DefaultOptions()
opts.Strategy = rotation.WeightedStrategy

rotator, err := lashes.New(opts)

// To adjust weights:
proxy.Weight = 100 // Higher weight = higher selection probability
```

Features:
- Cryptographically secure randomization using crypto/rand
- Strongly favors proxies with positive weights (95% probability)
- Rarely selects proxies with zero/negative weights (5% probability)
- Allows fine-grained control over proxy selection frequency

### 4. Least-Used

The least-used strategy prioritizes proxies with the lowest usage count and oldest last-used time. This strategy ensures even distribution over time and prevents overuse of specific proxies.

```go
// Initialize with least-used strategy
opts := lashes.DefaultOptions()
opts.Strategy = rotation.LeastUsedStrategy

rotator, err := lashes.New(opts)
```

Features:
- Prioritizes proxies with lowest usage count
- For equal usage counts, selects the proxy used least recently
- Automatically balances load across all proxies over time
- Useful for avoiding rate limits and IP blocking
