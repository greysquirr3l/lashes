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

### Core Functionality

- 🔄 Multiple proxy protocols (HTTP, SOCKS4, SOCKS5)
- 🎯 Smart rotation strategies (Round-robin, Random, Weighted, Least-used)
- 💾 Flexible storage options (Memory, SQLite, MySQL, PostgreSQL)
- 🔍 Health checking and validation
- 📊 Performance metrics tracking
- 🛡️ Rate limiting and retry support

### Key Benefits

- Zero dependencies for core functionality
- Thread-safe operations
- Context-aware API
- Comprehensive test coverage
- Production-ready defaults

## Installation

```bash
go get github.com/greysquirr3l/lashes@latest
```

Minimum Go version: 1.24.0

## Quick Start

### Basic Usage (In-Memory)

```go
import "github.com/greysquirr3l/lashes"

// Create with default options
rotator, err := lashes.New(lashes.DefaultOptions())

// Add proxies
ctx := context.Background()
err = rotator.AddProxy(ctx, "http://proxy1.example.com:8080", lashes.HTTP)
err = rotator.AddProxy(ctx, "socks5://proxy2.example.com:1080", lashes.SOCKS5)

// Get configured HTTP client
client, err := rotator.Client(ctx)
resp, err := client.Get("https://api.example.com")
```

### With Database Storage

```go
opts := lashes.Options{
    Storage: &storage.Options{
        Type: storage.SQLite,
        DSN:  "file:proxies.db",
        MaxConnections: 10,
        QueryTimeout: time.Second * 30,
    },
    Strategy: rotation.Weighted,
    ValidateOnStart: true,
    ValidationTimeout: time.Second * 5,
    MaxRetries: 3,
}

rotator, err := lashes.New(opts)
```

## Advanced Usage

### Rotation Strategies

We support four rotation strategies optimized for different use cases:

#### Round-Robin (Default)

```go
opts := lashes.DefaultOptions() // Uses round-robin by default
```

- Consistent and predictable rotation
- O(1) selection time
- Even distribution guaranteed

#### Random Selection

```go
opts := lashes.Options{
    Strategy: rotation.Random,
}
```

- Unpredictable rotation pattern
- O(1) selection time
- Good for avoiding detection

#### Weighted Distribution

```go
opts := lashes.Options{
    Strategy: rotation.WeightedStrategy,
}

// Set proxy weights
proxy.Weight = 100 // Higher weight = more frequent selection
```

- Success rate-based selection
- Cryptographically secure randomization
- 95% selection chance for positive-weighted proxies
- 5% selection chance for zero/negative-weighted proxies

#### Least Used

```go
opts := lashes.Options{
    Strategy: rotation.LeastUsed,
}
```

- Even load distribution
- O(log n) selection time
- Prevents proxy overuse

### Rate Limiting

Configure global and per-proxy rate limits:

```go
// Global rate limit
opts := lashes.Options{
    RateLimit: rate.NewLimiter(rate.Limit(10), 1), // 10 requests/second
}

// Per-proxy rate limit
proxy.Settings.RateLimit = &lashes.RateLimit{
    RequestsPerSecond: 5,
    Burst: 2,
}
```

### Error Handling

The library provides structured error handling:

```go
proxy, err := rotator.GetProxy(ctx)
if err != nil {
    switch {
    case errors.Is(err, lashes.ErrNoProxiesAvailable):
        // Handle missing proxy
    case errors.Is(err, lashes.ErrProxyNotFound):
        // Handle proxy not found
    case errors.Is(err, lashes.ErrInvalidProxy):
        // Handle invalid proxy
    default:
        // Handle unknown error
    }
}
```

With automatic retries:

```go
client := rotator.Client(ctx, lashes.ClientOptions{
    MaxRetries: 3,
    RetryBackoff: lashes.ExponentialBackoff{
        Initial: time.Second,
        Factor:  2,
        Max:     time.Minute,
    },
})
```

### Custom Rotation Strategy

```go
opts := lashes.DefaultOptions()
opts.Strategy = rotation.Weighted  // or Random, LeastUsed

rotator, err := lashes.New(opts)
```

### Health Checking

```go
opts := lashes.DefaultOptions()
opts.ValidateOnStart = true
opts.ValidationTimeout = time.Second * 5
opts.TestURL = "https://api.ipify.org?format=json"
```

### Proxy Settings

```go
proxy, err := rotator.GetProxy(ctx)
proxy.Settings.FollowRedirects = false
proxy.Settings.VerifyCerts = false
proxy.Settings.Headers = map[string][]string{
    "User-Agent": {"custom-agent"},
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

## Storage Backends

### Supported Databases

- 💽 SQLite (zero external network dependencies)
- 🗄️ MySQL (production-ready)
- 🐘 PostgreSQL (enterprise-grade)

### Performance Considerations

- Lazy database initialization
- Connection pooling
- Prepared statements
- Indexed queries

## Security Features

- 🔒 Cryptographically secure randomization (crypto/rand)
- 🔐 TLS configuration with minimum TLS 1.2
- ⚡ Rate limiting with standard library rate.Limiter
- 🛡️ Robust input validation
- 📝 Comprehensive error handling

## Metrics & Monitoring

- Success/failure rates
- Latency tracking
- Request volumes
- Health status
- Last-used timestamps

## Project Status

Current Version: v0.1.0

Status:

- ✅ Core functionality complete
- ✅ Test coverage >80%
- ✅ API documentation
- ✅ Security policy
- ✅ Production ready

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone repository
git clone https://github.com/greysquirr3l/lashes.git
cd lashes

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run linter
golangci-lint run
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

[MIT](LICENSE) © Nick Campbell and Contributors

## Support

- 📚 [Documentation](https://pkg.go.dev/github.com/greysquirr3l/lashes)
- 🐛 [Issue Tracker](https://github.com/greysquirr3l/lashes/issues)
- 💬 [Discussions](https://github.com/greysquirr3l/lashes/discussions)
- 🔒 [Security](SECURITY.md)

## Acknowledgments

Thanks to all contributors and the Go community for inspiration and feedback.
