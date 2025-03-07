# Metrics Documentation

## Overview

The lashes library provides a comprehensive metrics collection system that tracks:

1. Request success and failure counts
2. Request latencies (min, max, average)
3. Response status codes
4. Proxy usage frequency

## Accessing Metrics

### Get Metrics for a Specific Proxy

```go
proxyID := "your-proxy-id"
metrics, err := rotator.GetProxyMetrics(ctx, proxyID)
if err != nil {
    log.Fatalf("Failed to get metrics: %v", err)
}

fmt.Printf("Success rate: %.2f%%\n", metrics.SuccessRate * 100)
fmt.Printf("Total calls: %d\n", metrics.TotalCalls)
fmt.Printf("Average latency: %v\n", metrics.AvgLatency)
fmt.Printf("Error count: %d\n", metrics.ErrorCount)
```

### Get Metrics for All Proxies

```go
allMetrics, err := rotator.GetAllMetrics(ctx)
if err != nil {
    log.Fatalf("Failed to get all metrics: %v", err)
}

// Print a metrics table
fmt.Printf("%-30s %-10s %-10s %-10s\n", "URL", "Success %", "Calls", "Avg ms")
fmt.Println(strings.Repeat("-", 65))

for _, m := range allMetrics {
    fmt.Printf("%-30s %-10.1f %-10d %-10.2f\n",
        m.URL,
        m.SuccessRate * 100,
        m.TotalCalls,
        float64(m.AvgLatency) / float64(time.Millisecond))
}
```

## Metrics Structure

The `ProxyMetrics` struct contains:

```go
type ProxyMetrics struct {
    ProxyID     string        // Unique proxy identifier
    URL         string        // Proxy URL
    Type        string        // Proxy type (HTTP, SOCKS4, SOCKS5)
    SuccessRate float64       // Success rate (0.0-1.0)
    TotalCalls  int64         // Total number of requests
    AvgLatency  time.Duration // Average latency
    MinLatency  time.Duration // Minimum latency
    MaxLatency  time.Duration // Maximum latency
    LastUsed    time.Time     // Last time the proxy was used
    ErrorCount  int64         // Number of errors
    IsActive    bool          // Whether the proxy is active
}
```

## Using Metrics for Proxy Selection

You can use metrics to implement smart proxy selection logic:

```go
// Example: Select proxy with best recent performance
func selectBestProxy(ctx context.Context, rotator lashes.ProxyRotator) (*domain.Proxy, error) {
    metrics, err := rotator.GetAllMetrics(ctx)
    if err != nil {
        return nil, err
    }
    
    if len(metrics) == 0 {
        return nil, lashes.ErrNoProxiesAvailable
    }
    
    var bestMetrics *lashes.ProxyMetrics
    bestScore := -1.0
    
    for _, m := range metrics {
        // Skip inactive proxies
        if !m.IsActive {
            continue
        }
        
        // Calculate score based on success rate and latency
        // Higher score is better
        latencyScore := 1000.0 / float64(m.AvgLatency/time.Millisecond)
        score := m.SuccessRate*0.7 + latencyScore*0.3
        
        if score > bestScore {
            bestScore = score
            bestMetrics = m
        }
    }
    
    if bestMetrics == nil {
        return nil, lashes.ErrNoProxiesAvailable
    }
    
    // Get the actual proxy from the repository
    proxies, err := rotator.List(ctx)
    if err != nil {
        return nil, err
    }
    
    for _, proxy := range proxies {
        if proxy.URL == bestMetrics.URL {
            return proxy, nil
        }
    }
    
    return nil, lashes.ErrProxyNotFound
}
```

## Metrics Collection

Metrics are automatically collected during:

1. Proxy rotation (via `GetProxy()`)
2. HTTP requests (via `Client()`) 
3. Proxy validation (via `ValidateProxy()` and `ValidateAll()`)

The metrics are stored in memory and you can access them at any time.
