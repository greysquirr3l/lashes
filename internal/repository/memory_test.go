package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
)

func TestMemoryRepository(t *testing.T) {
	repo := repository.NewMemoryRepository()
	ctx := context.Background()

	// Create a test proxy
	proxy := &domain.Proxy{
		ID:      "test-123",
		URL:     "http://example.com:8080",
		Type:    domain.HTTP,
		Enabled: true,
	}

	t.Run("Create new proxy", func(t *testing.T) {
		err := repo.Create(ctx, proxy)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	})

	t.Run("Try to create duplicate proxy", func(t *testing.T) {
		dupProxy := &domain.Proxy{
			ID:      "test-123", // Same ID as the first proxy
			URL:     "http://different.com:8080",
			Type:    domain.HTTP,
			Enabled: true,
		}
		err := repo.Create(ctx, dupProxy)
		if !errors.Is(err, repository.ErrDuplicateID) {
			t.Errorf("Create() error = %v, want %v", err, repository.ErrDuplicateID)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		got, err := repo.GetByID(ctx, proxy.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if got.URL != proxy.URL {
			t.Errorf("GetByID() URL = %v, want %v", got.URL, proxy.URL)
		}
	})

	t.Run("GetByID not found", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "non-existent-id")
		if !errors.Is(err, repository.ErrProxyNotFound) {
			t.Errorf("GetByID() error = %v, want %v", err, repository.ErrProxyNotFound)
		}
	})

	t.Run("List", func(t *testing.T) {
		proxies, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if len(proxies) != 1 {
			t.Errorf("List() returned %d proxies, want 1", len(proxies))
		}
	})

	t.Run("Update", func(t *testing.T) {
		proxy.URL = "http://updated.com:8080"
		err := repo.Update(ctx, proxy)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, err := repo.GetByID(ctx, proxy.ID)
		if err != nil {
			t.Fatalf("GetByID() after update error = %v", err)
		}
		if got.URL != "http://updated.com:8080" {
			t.Errorf("Update() didn't change URL, got %v", got.URL)
		}
	})

	t.Run("Update not found", func(t *testing.T) {
		notFound := &domain.Proxy{
			ID:  "non-existent-id",
			URL: "http://example.com",
		}
		err := repo.Update(ctx, notFound)
		if !errors.Is(err, repository.ErrProxyNotFound) {
			t.Errorf("Update() error = %v, want %v", err, repository.ErrProxyNotFound)
		}
	})

	t.Run("Enabled state updates", func(t *testing.T) {
		// Use the SetEnabled method to update both fields consistently
		proxy.SetEnabled(false)

		err := repo.Update(ctx, proxy)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		got, err := repo.GetByID(ctx, proxy.ID)
		if err != nil {
			t.Fatalf("GetByID() after update error = %v", err)
		}
		if got.Enabled != false {
			t.Errorf("Enabled = %v, want %v", got.Enabled, false)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, proxy.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		_, err = repo.GetByID(ctx, proxy.ID)
		if !errors.Is(err, repository.ErrProxyNotFound) {
			t.Errorf("GetByID() after delete error = %v, want %v", err, repository.ErrProxyNotFound)
		}

		// List should be empty
		proxies, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("List() after delete error = %v", err)
		}
		if len(proxies) != 0 {
			t.Errorf("List() after delete returned %d proxies, want 0", len(proxies))
		}
	})

	t.Run("Delete not found", func(t *testing.T) {
		err := repo.Delete(ctx, "non-existent-id")
		if !errors.Is(err, repository.ErrProxyNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, repository.ErrProxyNotFound)
		}
	})
}
