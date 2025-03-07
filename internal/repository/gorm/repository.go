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

func (r *proxyRepository) Create(ctx context.Context, proxy *domain.Proxy) error {
	if err := validateProxy(proxy); err != nil {
		return fmt.Errorf("%w: %s", repository.ErrInvalidProxy, err.Error())
	}

	model := toModel(proxy)
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

func (r *proxyRepository) GetByID(ctx context.Context, id string) (*domain.Proxy, error) {
	var model ProxyModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrProxyNotFound
		}
		return nil, result.Error
	}
	return model.ToDomain()
}

func (r *proxyRepository) Update(ctx context.Context, proxy *domain.Proxy) error {
	model := toModel(proxy)

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repository.ErrProxyNotFound
	}
	return nil
}

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
	return model.ToDomain()
}

func validateProxy(proxy *domain.Proxy) error {
	if proxy == nil {
		return errors.New("proxy cannot be nil")
	}
	if proxy.URL == nil {
		return errors.New("proxy URL cannot be nil")
	}
	if proxy.Type == "" {
		return errors.New("proxy type cannot be empty")
	}
	return nil
}

func toModel(proxy *domain.Proxy) *ProxyModel {
	return &ProxyModel{
		ID:             proxy.ID,
		URL:            proxy.URL.String(),
		Type:           string(proxy.Type),
		LastUsed:       proxy.LastUsed,
		LastCheck:      proxy.LastCheck,
		Latency:        proxy.Latency,
		IsActive:       proxy.IsActive,
		Weight:         proxy.Weight,
		MaxRetries:     proxy.MaxRetries,
		Timeout:        proxy.Timeout,
		SuccessCount:   proxy.Metrics.SuccessCount,
		FailureCount:   proxy.Metrics.FailureCount,
		TotalRequests:  proxy.Metrics.TotalRequests,
		AvgLatency:     proxy.Metrics.AvgLatency,
		LastStatusCode: proxy.Metrics.LastStatusCode,
	}
}
