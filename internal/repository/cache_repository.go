package repository

import (
	"context"

	"github.com/guttosm/url-shortener/internal/entity"
)

// URLCacheRepository defines the interface for caching URL entities.
//
// Methods:
// - GetByOriginalURL: Retrieves a URL entity from the cache using its original URL as the key.
// - SetByOriginalURL: Caches a URL entity using its original URL as the key.
// - GetByShortID: Retrieves a URL entity from the cache using its shortened ID as the key.
// - SetByShortID: Caches a URL entity using its shortened ID as the key.
type URLCacheRepository interface {
	// GetByOriginalURL retrieves a URL entity from the cache using its original URL as the key.
    //
    // Parameters:
    // - ctx (context.Context): The context for the operation.
    // - originalURL (string): The original URL to search for.
    //
    // Returns:
    // - *entity.URL: The URL entity if found, or nil if no matching key exists.
    // - error: An error if the retrieval fails.
	GetByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error)
	
	// SetByOriginalURL caches a URL entity using its original URL as the key.
    //
    // Parameters:
    // - ctx (context.Context): The context for the operation.
    // - url (*entity.URL): The URL entity to be cached.
    //
    // Returns:
    // - error: An error if the caching operation fails.
	SetByOriginalURL(ctx context.Context, url *entity.URL) error


    // GetByShortID retrieves a URL entity from the cache using its shortened ID as the key.
    //
    // Parameters:
    // - ctx (context.Context): The context for the operation.
    // - shortID (string): The shortened ID to search for.
    //
    // Returns:
    // - *entity.URL: The URL entity if found, or nil if no matching key exists.
    // - error: An error if the retrieval fails.
	GetByShortID(ctx context.Context, shortID string) (*entity.URL, error)
	
	// SetByShortID caches a URL entity using its shortened ID as the key.
    //
    // Parameters:
    // - ctx (context.Context): The context for the operation.
    // - url (*entity.URL): The URL entity to be cached.
    //
    // Returns:
    // - error: An error if the caching operation fails.
	SetByShortID(ctx context.Context, url *entity.URL) error
	
}
