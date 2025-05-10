package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/guttosm/url-shortener/internal/entity"
	"github.com/guttosm/url-shortener/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// urlMongoRepository is a MongoDB implementation of the URLRepository interface.
//
// Fields:
// - collection (*mongo.Collection): The MongoDB collection used to store and retrieve URL entities.
type urlMongoRepository struct {
	collection *mongo.Collection
}

// NewURLMongoRepository creates a new instance of urlMongoRepository.
//
// Parameters:
// - col (*mongo.Collection): The MongoDB collection to be used for URL storage.
//
// Returns:
// - repository.URLRepository: An instance of the URLRepository interface backed by MongoDB.
func NewURLMongoRepository(col *mongo.Collection) repository.URLRepository {
	return &urlMongoRepository{
		collection: col,
	}
}

// Save stores a new URL entity in the MongoDB collection.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - url (*entity.URL): The URL entity to be saved.
//
// Behavior:
// - Sets the CreatedAt field of the URL entity to the current time.
// - Inserts the URL entity into the MongoDB collection.
//
// Returns:
// - error: An error if the insertion fails.
func (r *urlMongoRepository) Save(ctx context.Context, url *entity.URL) error {
	url.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, url)
	return err
}

// FindByOriginalURL retrieves a URL entity by its original URL.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - originalURL (string): The original URL to search for.
//
// Returns:
// - *entity.URL: The URL entity if found, or nil if no matching document exists.
// - error: An error if the query fails.
func (r *urlMongoRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	var url entity.URL
	err := r.collection.FindOne(ctx, bson.M{"original": originalURL}).Decode(&url)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}

// FindByShortID retrieves a URL entity by its shortened ID.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - shortID (string): The shortened ID to search for.
//
// Returns:
// - *entity.URL: The URL entity if found, or nil if no matching document exists.
// - error: An error if the query fails.
func (r *urlMongoRepository) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	var url entity.URL
	err := r.collection.FindOne(ctx, bson.M{"short_id": shortID}).Decode(&url)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}