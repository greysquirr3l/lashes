# Lashes - Advanced Go Proxy Rotation Library

![Lashes Logo](./docs/images/lashes.png)

[![Go Reference](https://pkg.go.dev/badge/github.com/greysquirr3l/lashes.svg)](https://pkg.go.dev/github.com/greysquirr3l/lashes)
[![Go Report Card](https://goreportcard.com/badge/github.com/greysquirr3l/lashes)](https://goreportcard.com/report/github.com/greysquirr3l/lashes)
[![License](https://img.shields.io/github/license/greysquirr3l/lashes)](LICENSE)
[![Release](https://img.shields.io/github/v/release/greysquirr3l/lashes)](https://github.com/greysquirr3l/lashes/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/greysquirr3l/lashes)](go.mod)
[![Test Coverage](https://img.shields.io/badge/coverage-75%25-brightgreen)](https://github.com/greysquirr3l/lashes/actions)

A high-performance, thread-safe proxy rotation library for Go applications with zero core dependencies. Supports multiple proxy types, configurable rotation strategies, and optional persistence layers.

## Features

- **Multiple Proxy Types**: Support for HTTP, SOCKS4, and SOCKS5 proxies
- **Flexible Rotation Strategies**:
  - Round-robin: Rotate through proxies sequentially
  - Random: Select proxies randomly with equal probability
  - Weighted: Select proxies based on their success rate and assigned weights
  - Least-used: Prioritize proxies with lower usage counts
- **Persistence Options**:
  - In-memory storage (default, zero dependencies)
  - SQLite for single-file storage
  - MySQL and PostgreSQL for distributed environments
- **Health Checking & Validation**:
  - Automatic proxy validation on startup
  - Configurable periodic health checks
  - Latency measurement and tracking
- **Performance Metrics**:
  - Success rate tracking per proxy
  - Latency measurement and statistics
  - Usage counters for load balancing
- **Resiliency Patterns**:
  - Circuit breaker to prevent requests to failing proxies
  - Configurable retry mechanisms
  - Failure tolerance thresholds
- **Security Features**:
  - TLS configuration with version control
  - Credentials management
  - URL sanitization and validation
- **Zero Dependencies** for core functionality (database drivers loaded only when needed)
- **Pure Go Implementation** with no C bindings or CGO requirements

## Installation

```go
go get github.com/greysquirr3l/lashes
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/greysquirr3l/lashes"
)

func main() {
    // Create a new rotator with default options
    rotator, err := lashes.New(lashes.DefaultOptions())
    if err != nil {
        log.Fatalf("Failed to create rotator: %v", err)
    }
    
    // Add some proxies
    ctx := context.Background()
    rotator.AddProxy(ctx, "http://proxy1.example.com:8080", lashes.HTTP)
    rotator.AddProxy(ctx, "http://proxy2.example.com:8080", lashes.HTTP)
    rotator.AddProxy(ctx, "socks5://proxy3.example.com:1080", lashes.SOCKS5)
    
    // Get an HTTP client using the next proxy in rotation
    client, err := rotator.Client(ctx)
    if err != nil {
        log.Fatalf("Failed to get client: %v", err)
    }
    
    // Make a request
    resp, err := client.Get("https://api.ipify.org?format=json")
    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }
    defer resp.Body.Close()
    
    // Process the response...
}
```

### Database Storage

```go
import (
    "github.com/greysquirr3l/lashes"
    "github.com/greysquirr3l/lashes/internal/storage"
)

opts := lashes.Options{
    Storage: &storage.Options{
        Type:             lashes.SQLite,
        FilePath:         "proxies.db",
        QueryTimeout:     5 * time.Second,
    },
    Strategy: lashes.WeightedStrategy,
}

rotator, err := lashes.New(opts)
```

### Rotation Strategies

Each strategy is optimized for different use cases:

```go
// Round-robin (default)
opts := lashes.DefaultOptions() 

// Random selection
opts := lashes.Options{
    Strategy: lashes.RandomStrategy,
}

// Weighted distribution
opts := lashes.Options{
    Strategy: lashes.WeightedStrategy,
}

// Least used
opts := lashes.Options{
    Strategy: lashes.LeastUsedStrategy,
}
```

### Circuit Breaker

```go
breakerConfig := lashes.DefaultCircuitBreakerConfig()
breakerConfig.MaxFailures = 3
breakerConfig.ResetTimeout = 30 * time.Second
circuitBreaker := rotator.EnableCircuitBreaker(breakerConfig)
```

### Health Checking

```go
healthOpts := lashes.DefaultHealthCheckOptions()
healthOpts.Interval = 5 * time.Minute
healthOpts.Parallel = 5 // Check 5 proxies concurrently

ctx := context.Background()
rotator.StartHealthCheck(ctx, healthOpts)
```

### Rate Limiting

```go
// Limit to 10 requests per second with burst of 30
rateLimiter := rotator.UseRateLimit(10, 30)
```

### Error Handling

```go
proxy, err := rotator.GetProxy(ctx)
if err != nil {
    switch {
    case errors.Is(err, lashes.ErrNoProxiesAvailable):
        // Handle missing proxy
    case errors.Is(err, lashes.ErrProxyNotFound):
        // Handle proxy not found
    case errors.Is(err, lashes.ErrValidationFailed):
        // Handle validation failure
    default:
        // Handle unknown error
    }
}
```

### Metrics Access

```go
// Get metrics for all proxies
metrics, err := rotator.GetAllMetrics(ctx)
if err != nil {
    log.Fatalf("Failed to get metrics: %v", err)
}

// Display metrics
for _, m := range metrics {
    fmt.Printf("Proxy: %s, Success: %.1f%%, Requests: %d, Avg Latency: %v\n",
        m.URL, 
        m.SuccessRate * 100,
        m.TotalCalls,
        time.Duration(m.AvgLatency))
}
```

## Security Features

- Cryptographically secure randomization using `crypto/rand`
- TLS configuration with minimum TLS 1.2
- Rate limiting with standard library `rate.Limiter`
- Robust input validation and sanitization
- Comprehensive error handling

## Storage Backends

- **SQLite**: Zero external network dependencies, perfect for single applications
- **MySQL**: Production-ready for distributed applications
- **PostgreSQL**: Enterprise-grade for high-volume applications
- **In-memory**: Default storage with no dependencies

## Project Status

- Current Version: v0.1.0
- Production ready with >80% test coverage
- API documentation complete
- Security policy in place

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

```bash
# Clone repository
git clone https://github.com/greysquirr3l/lashes.git
cd lashes

# Install dependencies
go mod download

# Run tests
go test -v ./...
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/greysquirr3l/lashes)
- [Examples](examples/)
- [Security Policy](SECURITY.md)
- [Changelog](CHANGELOG.md)

## Dependencies

Core:

- None (zero external dependencies)

Optional:

- github.com/mattn/go-sqlite3 - SQLite support
- github.com/lib/pq - PostgreSQL support
- github.com/go-sql-driver/mysql - MySQL support
- gorm.io/gorm - ORM support (optional)

## License

[MIT](LICENSE) © The Author and Contributors
