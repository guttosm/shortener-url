package service

import (
	"context"
	"time"

	"github.com/guttosm/url-shortener/internal/entity"
	"github.com/guttosm/url-shortener/internal/repository"
	"github.com/rs/xid"
)

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (*entity.URL, error)
	FindByShortID(ctx context.Context, shortID string) (*entity.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(repo repository.URLRepository) URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(ctx context.Context, originalURL string) (*entity.URL, error) {
	existingURL, err := s.repo.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}

	// Se a URL já existir, retornar o registro existente
	if existingURL != nil {
		return existingURL, nil
	}

	url := &entity.URL{
		ShortID:   xid.New().String()[:6],
		Original:  originalURL,
		CreatedAt: time.Now(),
	}

	err = s.repo.Save(ctx, url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	return s.repo.FindByShortID(ctx, shortID)
}
