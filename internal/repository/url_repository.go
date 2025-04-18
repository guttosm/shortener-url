package repository

import (
	"context"

	"github.com/guttosm/url-shortener/internal/entity"
)

type URLRepository interface {
	Save(ctx context.Context, url *entity.URL) error
	FindByShortID(ctx context.Context, shortID string) (*entity.URL, error)
	FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error)
}
