package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/guttosm/url-shortener/internal/entity"
	"github.com/guttosm/url-shortener/internal/repository"
	"github.com/redis/go-redis/v9"
)

// urlRedisRepository is a Redis implementation of the URLCacheRepository interface.
//
// Fields:
// - client (*redis.Client): The Redis client used to interact with the Redis database.
type urlRedisRepository struct {
	client *redis.Client
}

// NewURLRedisRepository creates a new instance of urlRedisRepository.
//
// Parameters:
// - client (*redis.Client): The Redis client to be used for caching URL entities.
//
// Returns:
// - repository.URLCacheRepository: An instance of the URLCacheRepository interface backed by Redis.
func NewURLRedisRepository(client *redis.Client) repository.URLCacheRepository {
	return &urlRedisRepository{client: client}
}

// GetByOriginalURL retrieves a URL entity from Redis by its original URL.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - originalURL (string): The original URL to search for.
//
// Behavior:
// - Constructs a Redis key using the original URL.
// - Retrieves the cached URL entity from Redis and unmarshal it into an entity.URL object.
//
// Returns:
// - *entity.URL: The URL entity if found, or nil if no matching key exists.
// - error: An error if the retrieval or unmarshalling fails.
func (r *urlRedisRepository) GetByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	key := "url:original:" + originalURL
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var url entity.URL
	err = json.Unmarshal([]byte(val), &url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// SetByOriginalURL caches a URL entity in Redis using its original URL as the key.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - url (*entity.URL): The URL entity to be cached.
//
// Behavior:
// - Constructs a Redis key using the original URL.
// - Marshals the URL entity into JSON and stores it in Redis with a 1-hour expiration.
//
// Returns:
// - error: An error if the caching operation fails.
func (r *urlRedisRepository) SetByOriginalURL(ctx context.Context, url *entity.URL) error {
	key := "url:original:" + url.Original
	data, _ := json.Marshal(url)
	return r.client.Set(ctx, key, data, time.Hour).Err()
}

// GetByShortID retrieves a URL entity from Redis by its shortened ID.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - shortID (string): The shortened ID to search for.
//
// Behavior:
// - Constructs a Redis key using the shortened ID.
// - Retrieves the cached URL entity from Redis and unmarshal it into an entity.URL object.
//
// Returns:
// - *entity.URL: The URL entity if found, or nil if no matching key exists.
// - error: An error if the retrieval or unmarshalling fails.
func (r *urlRedisRepository) GetByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	key := "url:short_id:" + shortID
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var url entity.URL
	err = json.Unmarshal([]byte(val), &url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// SetByShortID caches a URL entity in Redis using its shortened ID as the key.
//
// Parameters:
// - ctx (context.Context): The context for the operation.
// - url (*entity.URL): The URL entity to be cached.
//
// Behavior:
// - Constructs a Redis key using the shortened ID.
// - Marshals the URL entity into JSON and stores it in Redis with a 1-hour expiration.
//
// Returns:
// - error: An error if the caching operation fails.
func (r *urlRedisRepository) SetByShortID(ctx context.Context, url *entity.URL) error {
	key := "url:short_id:" + url.ShortID
	data, _ := json.Marshal(url)
	return r.client.Set(ctx, key, data, time.Hour).Err()
}
