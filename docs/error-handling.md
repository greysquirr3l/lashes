# Error Handling Documentation

## Common Errors

### 1. Proxy Errors

```go
var (
    ErrNoProxiesAvailable = errors.New("no proxies available")
    ErrProxyNotFound = errors.New("proxy not found")
    ErrInvalidProxy = errors.New("invalid proxy configuration")
)
```

These errors indicate issues with proxy configuration or availability:
- `ErrNoProxiesAvailable` - Returned when trying to get a proxy but none are available
- `ErrProxyNotFound` - Returned when trying to access a specific proxy that doesn't exist
- `ErrInvalidProxy` - Returned when a proxy is improperly configured

### 2. Configuration Errors

```go
var (
    ErrInvalidOptions = errors.New("invalid options provided")
    ErrMetricsNotEnabled = errors.New("metrics collection not enabled")
)
```

These errors indicate issues with library configuration:
- `ErrInvalidOptions` - Returned when the options provided to `New()` are invalid
- `ErrMetricsNotEnabled` - Returned when trying to access metrics but they aren't enabled

## Error Handling Examples

### Basic Error Handling

```go
proxy, err := rotator.GetProxy(ctx)
if err != nil {
    switch {
    case errors.Is(err, lashes.ErrNoProxiesAvailable):
        // Handle no proxies available
        log.Println("No proxies available, adding some...")
        addDefaultProxies(rotator)
    case errors.Is(err, lashes.ErrProxyNotFound):
        // Handle proxy not found
        log.Println("Proxy not found")
    default:
        // Handle unknown error
        log.Printf("Unexpected error: %v", err)
    }
}
```

### Validation Error Handling

```go
isValid, latency, err := rotator.ValidateProxy(ctx, proxy, "https://example.com")
if err != nil {
    log.Printf("Validation error: %v", err)
    return
}

if !isValid {
    log.Printf("Proxy failed validation with latency %v", latency)
    // Handle invalid proxy
    err = rotator.RemoveProxy(ctx, proxy.URL)
    if err != nil {
        log.Printf("Failed to remove invalid proxy: %v", err)
    }
}
```

### With Retry Logic

Custom retry logic can be implemented for handling temporary failures:

```go
const maxRetries = 3
var lastError error

for attempt := 0; attempt < maxRetries; attempt++ {
    proxy, err := rotator.GetProxy(ctx)
    if err != nil {
        if errors.Is(err, lashes.ErrNoProxiesAvailable) {
            // No point in retrying if no proxies are available
            return nil, err
        }
        lastError = err
        continue
    }
    
    client, err := rotator.Client(ctx)
    if err != nil {
        lastError = err
        continue
    }
    
    resp, err := client.Get("https://api.example.com")
    if err != nil {
        lastError = err
        continue
    }
    
    return resp, nil
}

// If we've exhausted all retries
return nil, fmt.Errorf("max retries exceeded: %w", lastError)
```
