package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/greysquirr3l/lashes/internal/domain"
	"github.com/greysquirr3l/lashes/internal/repository"
	"gorm.io/gorm"
)

type proxyRepository struct {
	db *gorm.DB
}

type Options struct {
	Debug         bool
	QueryTimeout  time.Duration
	RetryAttempts int
}

func NewProxyRepository(db *gorm.DB, opts Options) repository.ProxyRepository {
	if opts.QueryTimeout == 0 {
		opts.QueryTimeout = 30 * time.Second
	}

	gormDB := db
	if opts.Debug {
		gormDB = db.Session(&gorm.Session{})
	}

	return &proxyRepository{
		db: gormDB.Set("gorm:query_timeout", opts.QueryTimeout),
	}
}

// Create implements the repository method to create a proxy record
func (r *proxyRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	if err := validateProxy(proxy); err != nil {
		return fmt.Errorf("%w: %s", repository.ErrInvalidProxy, err.Error())
	}

	model := FromDomain(proxy)
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(model).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository.ErrDuplicateID
		}
		return fmt.Errorf("failed to create proxy: %w", err)
	}

	return tx.Commit().Error
}

// GetByID retrieves a proxy by its ID
func (r *proxyRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	var model ProxyModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrProxyNotFound
		}
		return nil, result.Error
	}

	// Convert model to domain object
	proxy, err := model.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to domain model: %w", err)
	}

	return proxy, nil
}

// Update updates an existing proxy
func (r *proxyRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	model := FromDomain(proxy)

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrProxyNotFound
	}
	return nil
}

// Delete removes a proxy by ID
func (r *proxyRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&ProxyModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrProxyNotFound
	}
	return nil
}

// List retrieves all proxies
func (r *proxyRepository) List(ctx context.Context) ([]*domain.Proxy, error) {
	var models []ProxyModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	proxies := make([]*domain.Proxy, 0, len(models))
	for _, model := range models {
		proxy, err := model.ToDomain()
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}

// GetNext retrieves the next proxy to use
func (r *proxyRepository) GetNext(ctx context.Context) (*domain.Proxy, error) {
	var model ProxyModel
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("last_used ASC").
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrProxyNotFound
		}
		return nil, err
	}

	proxy, err := model.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to domain model: %w", err)
	}

	return proxy, nil
}

// validateProxy checks if a proxy is valid
func validateProxy(proxy *domain.Proxy) error {
	if proxy == nil {
		return errors.New("proxy cannot be nil")
	}
	if proxy.ID == "" {
		return errors.New("proxy ID cannot be empty")
	}
	if proxy.URL == "" {
		return errors.New("proxy URL cannot be empty")
	}
	return nil
}
