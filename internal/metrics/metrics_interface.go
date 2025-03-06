package metrics

import "context"

// ProxyMetrics defines the interface for metrics tracking in the proxy rotation library.
type ProxyMetrics interface {
	RecordSuccess(ctx context.Context, proxyID string)
	RecordFailure(ctx context.Context, proxyID string, err error)
	GetMetrics(ctx context.Context, proxyID string) (MetricsData, error)
}

// MetricsData holds the metrics information for a proxy.
type MetricsData struct {
	SuccessCount   int
	FailureCount   int
	AverageLatency float64
	LastStatusCode int
	LastUsed       int64 // Unix timestamp
}
