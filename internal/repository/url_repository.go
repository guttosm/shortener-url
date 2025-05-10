package repository

import (
	"context"

	"github.com/guttosm/url-shortener/internal/entity"
)

// URLRepository defines the interface for interacting with the persistent storage of URL entities.
//
// Methods:
// - Save: Stores a new URL entity in the database.
// - FindByShortID: Retrieves a URL entity by its shortened ID.
// - FindByOriginalURL: Retrieves a URL entity by its original URL.
type URLRepository interface {
	// Save stores a new URL entity in the database.
	//
	// Parameters:
	// - ctx (context.Context): The context for the operation.
	// - url (*entity.URL): The URL entity to be saved.
	//
	// Returns:
	// - error: An error if the save operation fails.
	Save(ctx context.Context, url *entity.URL) error

	// FindByShortID retrieves a URL entity by its shortened ID.
	//
	// Parameters:
	// - ctx (context.Context): The context for the operation.
	// - shortID (string): The shortened ID to search for.
	//
	// Returns:
	// - *entity.URL: The URL entity if found, or nil if no matching document exists.
	// - error: An error if the query fails.
	FindByShortID(ctx context.Context, shortID string) (*entity.URL, error)

	// FindByOriginalURL retrieves a URL entity by its original URL.
	//
	// Parameters:
	// - ctx (context.Context): The context for the operation.
	// - originalURL (string): The original URL to search for.
	//
	// Returns:
	// - *entity.URL: The URL entity if found, or nil if no matching document exists.
	// - error: An error if the query fails.
	FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error)
}
