package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository/custom"
)

// CustomRepositoryExample demonstrates how to use a custom repository implementation
func CustomRepositoryExample() {
	repo := custom.NewProxyRepository()
	
	ctx := context.Background()
	
	// Add some example proxies
	now := time.Now()
	// Use pointers for time fields that require them
	nowPtr := &now
	
	proxy1 := &domain.Proxy{
		ID:        "proxy1",
		URL:       "http://example1.com:8080",
		Type:      domain.HTTPProxy,
		LastUsed:  nowPtr,
		IsActive:  true,
		CreatedAt: now,
	}
	
	proxy2 := &domain.Proxy{
		ID:        "proxy2",
		URL:       "http://example2.com:8080",
		Type:      domain.SOCKS5Proxy, 
		LastUsed:  nowPtr,
		IsActive:  true,
		CreatedAt: now,
	}
	
	// Create proxies
	if err := repo.Create(ctx, proxy1); err != nil {
		fmt.Printf("Failed to create proxy: %v\n", err)
		return
	}
	
	if err := repo.Create(ctx, proxy2); err != nil {
		fmt.Printf("Failed to create proxy: %v\n", err)
		return
	}
	
	// List all proxies
	proxies, err := repo.List(ctx)
	if err != nil {
		fmt.Printf("Failed to list proxies: %v\n", err)
		return
	}
	
	fmt.Printf("Custom repository has %d proxies\n", len(proxies))
	
	for _, p := range proxies {
		fmt.Printf("- %s (%s)\n", p.URL, p.Type)
	}
}
