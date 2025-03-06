# Error Handling Documentation

## Common Errors

### 1. Proxy Errors

```go
var (
    ErrProxyNotFound     = errors.New("proxy not found")
    ErrProxyUnavailable  = errors.New("proxy unavailable")
    ErrProxyTimeout      = errors.New("proxy timeout")
)
```

### 2. Configuration Errors

```go
var (
    ErrInvalidOptions    = errors.New("invalid options")
    ErrInvalidStrategy   = errors.New("invalid strategy")
)
```

## Error Handling Examples

### Basic Error Handling

```go
proxy, err := rotator.GetProxy(ctx)
if err != nil {
    switch {
    case errors.Is(err, lashes.ErrProxyNotFound):
        // Handle no proxies available
    case errors.Is(err, lashes.ErrProxyTimeout):
        // Handle timeout
    default:
        // Handle unknown error
    }
}
```

### With Retry Logic

```go
client := rotator.Client(ctx, lashes.ClientOptions{
    MaxRetries: 3,
    RetryBackoff: lashes.ExponentialBackoff{
        Initial: time.Second,
        Factor:  2,
        Max:     time.Minute,
    },
    RetryableErrors: []error{
        lashes.ErrProxyTimeout,
        lashes.ErrProxyUnavailable,
    },
})
```
