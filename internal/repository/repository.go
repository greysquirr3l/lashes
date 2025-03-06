package repository

import (
	"context"

	"github.com/greysquirr3l/lashes/internal/models"
)

type PostRepository interface {
    Create(ctx context.Context, post *models.Post) error
    GetByID(ctx context.Context, id int64) (*models.Post, error)
    Update(ctx context.Context, post *models.Post) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context) ([]*models.Post, error)
}
