# Rate Limiting Documentation

## Configuration Options

### Global Rate Limiting

```go
import "golang.org/x/time/rate"

opts := lashes.Options{
    RateLimit: rate.NewLimiter(rate.Limit(10), 1), // 10 requests per second
}
```

### Per-Proxy Rate Limiting

```go
proxy.Settings.RateLimit = &lashes.RateLimit{
    RequestsPerSecond: 5,
    Burst: 2,
}
```

## Backoff Strategies

### 1. Linear Backoff

```go
opts.RetryBackoff = lashes.LinearBackoff{
    Initial: time.Second,
    Step:    time.Second,
    Max:     time.Second * 10,
}
```

### 2. Exponential Backoff

```go
opts.RetryBackoff = lashes.ExponentialBackoff{
    Initial: time.Second,
    Factor:  2,
    Max:     time.Minute,
}
```
