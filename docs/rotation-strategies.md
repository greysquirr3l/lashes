# Rotation Strategy Documentation

## Available Strategies

### 1. Round-Robin

- Sequential rotation through proxy list
- Guarantees even distribution
- O(1) selection time
- Best for consistent load distribution

```go
opts := lashes.Options{
    Strategy: rotation.RoundRobin,
}
```

### 2. Random

- Randomly selects proxies
- Good for avoiding pattern detection
- O(1) selection time
- May have uneven distribution in short term

```go
opts := lashes.Options{
    Strategy: rotation.Random,
}
```

### 3. Weighted

- Uses proxy weights to influence selection
- Better proxies get higher weights
- O(log n) selection time
- Best for optimizing success rates

```go
// Set proxy weights
proxy.Weight = 5 // Higher weight = more likely to be chosen
opts := lashes.Options{
    Strategy: rotation.Weighted,
}
```

### 4. Least-Used

- Prioritizes less frequently used proxies
- Maintains usage counts
- O(log n) selection time
- Best for maximizing proxy lifespan

```go
opts := lashes.Options{
    Strategy: rotation.LeastUsed,
}
```
