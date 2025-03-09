// Package lashes provides a robust proxy rotation system with multiple proxy types
// (HTTP, SOCKS4, SOCKS5), configurable rotation strategies, and optional persistence.
//
// Lashes follows a repository pattern with domain-driven design principles. It provides
// a clean API for managing proxies, validating their health, collecting performance metrics,
// and using proxies for HTTP requests.
//
// # Core Features
//
// - Multiple proxy types (HTTP, SOCKS4, SOCKS5)
// - Configurable rotation strategies (round-robin, random, weighted, least-used)
// - Storage backends (memory, SQLite, MySQL, PostgreSQL)
// - Proxy health checking and validation
// - Performance metrics tracking
//
// # Basic Usage
//
//	// Create a rotator with default options
//	rotator, err := lashes.New(lashes.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Add some proxies
//	ctx := context.Background()
//	rotator.AddProxy(ctx, "http://proxy1.example.com:8080", lashes.HTTP)
//	rotator.AddProxy(ctx, "socks5://proxy2.example.com:1080", lashes.SOCKS5)
//
//	// Get an http.Client that uses the next proxy
//	client, err := rotator.Client(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Make a request using the proxy
//	resp, err := client.Get("https://api.ipify.org?format=json")
//
// # Rotation Strategies
//
// Lashes supports four proxy rotation strategies:
//
//   - RoundRobin: Rotates through proxies in sequence
//   - Random: Selects proxies randomly
//   - Weighted: Selects proxies based on their weights and success rates
//   - LeastUsed: Prioritizes proxies with the lowest usage count
//
// # Database Storage
//
//	opts := lashes.Options{
//		Storage: &storage.Options{
//			Type:     lashes.SQLite,
//			FilePath: "proxies.db",
//		},
//		Strategy: lashes.RoundRobin,
//	}
//	rotator, err := lashes.New(opts)
//
// # Health Checking
//
//	// Configure health check options
//	healthOpts := lashes.DefaultHealthCheckOptions()
//	healthOpts.Interval = time.Minute * 5
//
//	// Start periodic health checking
//	ctx := context.Background()
//	rotator.StartHealthCheck(ctx, healthOpts)
//
// # Customization
//
// Lashes provides extensive customization options for proxy rotation, validation,
// HTTP client behavior, and persistence. See the Options struct for configuration details.
package lashes
