# Proxy Validation Documentation

## Overview

Proxy validation is essential to ensure that proxies are functioning correctly before using them in production. The lashes library provides built-in validation capabilities that:

1. Verify the proxy can successfully handle HTTP requests
2. Measure the proxy's latency
3. Update proxy status and metrics based on validation results

## Validation Methods

### ValidateProxy

Validates a single proxy against a target URL:

```go
isValid, latency, err := rotator.ValidateProxy(ctx, proxy, "https://example.com")
if err != nil {
    log.Printf("Validation error: %v", err)
    return
}

fmt.Printf("Proxy valid: %v, Latency: %v\n", isValid, latency)
```

### ValidateAll

Validates all proxies in the pool:

```go
if err := rotator.ValidateAll(ctx); err != nil {
    log.Printf("Some proxies failed validation: %v", err)
}
```

## Configuration Options

The following options control validation behavior:

```go
opts := lashes.DefaultOptions()
opts.ValidateOnStart = true         // Validate proxies when adding them
opts.ValidationTimeout = 5 * time.Second  // Max time for validation
opts.TestURL = "https://api.ipify.org?format=json"  // URL to test against
```

## Validation Process

During validation:

1. A test HTTP request is made through the proxy to the test URL
2. The response status code is checked (200-299 is considered valid)
3. The request latency is measured
4. The proxy's status is updated based on the result

## Automatic Validation

Proxies can be validated automatically:

1. **On Addition**: When `ValidateOnStart` is true
2. **Periodically**: Using your own background validation routine

Example of a background validation routine:

```go
func startBackgroundValidation(ctx context.Context, rotator lashes.ProxyRotator) {
    ticker := time.NewTicker(15 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if err := rotator.ValidateAll(ctx); err != nil {
                log.Printf("Validation error: %v", err)
            }
        case <-ctx.Done():
            return
        }
    }
}

// Usage:
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go startBackgroundValidation(ctx, rotator)
```

## Custom Validation Logic

You can implement custom validation logic:

```go
proxies, err := rotator.List(ctx)
if err != nil {
    log.Fatalf("Failed to list proxies: %v", err)
}

for _, proxy := range proxies {
    // Custom validation logic
    isValid, latency, err := rotator.ValidateProxy(
        ctx,
        proxy,
        "https://your-custom-validation-endpoint.com",
    )
    
    // Update proxy based on custom validation
    if isValid {
        proxy.Weight += 10 // Increase weight for good proxies
    } else {
        proxy.Weight -= 10 // Decrease weight for bad proxies
    }
    
    // Update the proxy in storage
    if err := rotator.UpdateProxy(ctx, proxy); err != nil {
        log.Printf("Failed to update proxy: %v", err)
    }
}
```
