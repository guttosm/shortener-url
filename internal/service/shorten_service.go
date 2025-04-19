package service

import (
	"context"
	"time"

	"github.com/guttosm/url-shortener/internal/entity"
	"github.com/guttosm/url-shortener/internal/repository"
	"github.com/rs/xid"
)

// URLService defines the interface for URL shortening and retrieval services.
//
// Methods:
// - Shorten: Shortens a given original URL and stores it in the database and cache.
// - FindByShortID: Retrieves the original URL associated with a given shortened ID.
type URLService interface {
	// Shorten shortens a given original URL and stores it in the database and cache.
	//
	// Parameters:
	// - ctx (context.Context): The context for the operation.
	// - originalURL (string): The original URL to be shortened.
	//
	// Returns:
	// - *entity.URL: The shortened URL entity.
	// - error: An error if the operation fails.
	Shorten(ctx context.Context, originalURL string) (*entity.URL, error)

	// FindByShortID retrieves the original URL associated with a given shortened ID.
	//
	// Parameters:
	// - ctx (context.Context): The context for the operation.
	// - shortID (string): The shortened ID to search for.
	//
	// Returns:
	// - *entity.URL: The URL entity if found, or nil if no matching record exists.
	// - error: An error if the operation fails.
	FindByShortID(ctx context.Context, shortID string) (*entity.URL, error)
}

type urlService struct {
	repo      repository.URLRepository
	cacheRepo repository.URLCacheRepository
}

// NewURLService creates a new instance of URLService.
//
// Parameters:
// - repo (repository.URLRepository): The repository for persistent URL storage.
// - cache (repository.URLCacheRepository): The repository for caching URL entities.
//
// Returns:
// - URLService: An instance of the URLService interface.
func NewURLService(repo repository.URLRepository, cache repository.URLCacheRepository) URLService {
	return &urlService{
		repo:      repo,
		cacheRepo: cache,
	}
}

// Shorten shortens a given original URL and stores it in the database and cache.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - originalURL (string): The original URL to be shortened.
//
// Behavior:
// - Checks the cache for the original URL. If found, returns it.
// - Checks the database for the original URL. If found, caches it and returns it.
// - If not found, generates a new shortened ID, stores it in the database, and caches it.
//
// Returns:
// - *entity.URL: The shortened URL entity.
// - error: An error if the operation fails.
func (s *urlService) Shorten(ctx context.Context, originalURL string) (*entity.URL, error) {
	url, err := s.cacheRepo.GetByOriginalURL(ctx, originalURL)
	if err == nil && url != nil {
		return url, nil
	}

	url, err = s.repo.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return nil, err
	}
	if url != nil {
		_ = s.cacheRepo.SetByOriginalURL(ctx, url)
		_ = s.cacheRepo.SetByShortID(ctx, url)
		return url, nil
	}

	url = &entity.URL{
		ShortID:   xid.New().String()[:6],
		Original:  originalURL,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(ctx, url); err != nil {
		return nil, err
	}
	_ = s.cacheRepo.SetByOriginalURL(ctx, url)
	_ = s.cacheRepo.SetByShortID(ctx, url)

	return url, nil
}

// FindByShortID retrieves the original URL associated with a given shortened ID.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - shortID (string): The shortened ID to search for.
//
// Behavior:
// - Checks the cache for the shortened ID. If found, returns it.
// - Checks the database for the shortened ID. If found, caches it and returns it.
//
// Returns:
// - *entity.URL: The URL entity if found, or nil if no matching record exists.
// - error: An error if the operation fails.
func (s *urlService) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	url, err := s.cacheRepo.GetByShortID(ctx, shortID)
	if err == nil && url != nil {
		return url, nil
	}

	url, err = s.repo.FindByShortID(ctx, shortID)
	if err != nil {
		return nil, err
	}
	if url != nil {
		_ = s.cacheRepo.SetByOriginalURL(ctx, url)
		_ = s.cacheRepo.SetByShortID(ctx, url)
	}
	return url, nil
}
