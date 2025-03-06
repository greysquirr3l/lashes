package repository

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"testing"

	"github.com/greysquirr3l/lashes/internal/domain"
)

func TestMemoryRepository(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T, *memoryRepository)
	}{
		{"Create", testCreate},
		{"GetByID", testGetByID},
		{"Update", testUpdate},
		{"Delete", testDelete},
		{"List", testList},
		{"GetNext", testGetNext},
		{"ConcurrentAccess", testConcurrentAccess},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMemoryRepository()
			tt.fn(t, repo.(*memoryRepository))
		})
	}
}

func createTestProxy(id string) *domain.Proxy {
	u, _ := url.Parse("http://example.com:8080")
	return &domain.Proxy{
		ID:       id,
		URL:      u,
		Type:     domain.HTTP,
		IsActive: true,
	}
}

func testCreate(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	proxy := createTestProxy("test1")

	err := repo.Create(ctx, proxy)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	// Test duplicate creation
	err = repo.Create(ctx, proxy)
	if err != ErrDuplicateID {
		t.Errorf("Expected ErrDuplicateID, got %v", err)
	}
}

func testGetByID(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	proxy := createTestProxy("test2")

	err := repo.Create(ctx, proxy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.GetByID(ctx, proxy.ID)
	if err != nil {
		t.Errorf("GetByID failed: %v", err)
	}
	if got.ID != proxy.ID {
		t.Errorf("Expected ID %s, got %s", proxy.ID, got.ID)
	}
}

func testUpdate(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	proxy := createTestProxy("test3")

	err := repo.Create(ctx, proxy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	proxy.IsActive = false
	err = repo.Update(ctx, proxy)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	got, _ := repo.GetByID(ctx, proxy.ID)
	if got.IsActive != false {
		t.Error("Update didn't persist")
	}
}

func testDelete(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	proxy := createTestProxy("test4")

	err := repo.Create(ctx, proxy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err = repo.Delete(ctx, proxy.ID)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = repo.GetByID(ctx, proxy.ID)
	if err != ErrProxyNotFound {
		t.Errorf("Expected ErrProxyNotFound, got %v", err)
	}
}

func testList(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		proxy := createTestProxy(fmt.Sprintf("test5_%d", i))
		err := repo.Create(ctx, proxy)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	proxies, err := repo.List(ctx)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}
	if len(proxies) != 3 {
		t.Errorf("Expected 3 proxies, got %d", len(proxies))
	}
}

func testGetNext(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	proxy := createTestProxy("test6")

	err := repo.Create(ctx, proxy)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.GetNext(ctx)
	if err != nil {
		t.Errorf("GetNext failed: %v", err)
	}
	if got.ID != proxy.ID {
		t.Errorf("Expected ID %s, got %s", proxy.ID, got.ID)
	}
}

func testConcurrentAccess(t *testing.T, repo *memoryRepository) {
	ctx := context.Background()
	var wg sync.WaitGroup
	n := 100

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			proxy := createTestProxy(fmt.Sprintf("test7_%d", i))
			_ = repo.Create(ctx, proxy)
		}(i)
	}

	wg.Wait()
	proxies, _ := repo.List(ctx)
	if len(proxies) != n {
		t.Errorf("Expected %d proxies, got %d", n, len(proxies))
	}
}
