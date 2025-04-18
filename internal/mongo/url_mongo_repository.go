package mongo

import (
	"context"
	"time"

	"encoding/json"
	"github.com/guttosm/url-shortener/internal/entity"
	"github.com/guttosm/url-shortener/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type urlMongoRepository struct {
	collection *mongo.Collection
	redis      *redis.Client
}

func NewURLMongoRepository(col *mongo.Collection, redisClient *redis.Client) repository.URLRepository {
	return &urlMongoRepository{
		collection: col,
		redis:      redisClient,
	}
}

func (r *urlMongoRepository) Save(ctx context.Context, url *entity.URL) error {
	url.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, url)
	return err
}

func (r *urlMongoRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URL, error) {
	// Tentar buscar no Redis primeiro
	cacheKey := "url:original:" + originalURL
	cachedData, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var url entity.URL
		if err := json.Unmarshal([]byte(cachedData), &url); err == nil {
			return &url, nil
		}
	}

	var url entity.URL
	err = r.collection.FindOne(ctx, bson.M{"original": originalURL}).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	data, _ := json.Marshal(url)
	r.redis.Set(ctx, cacheKey, data, time.Hour)

	return &url, nil
}

func (r *urlMongoRepository) FindByShortID(ctx context.Context, shortID string) (*entity.URL, error) {
	// Tentar buscar no Redis primeiro
	cacheKey := "url:short_id:" + shortID
	cachedData, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Se encontrado no Redis, deserializar e retornar
		var url entity.URL
		if err := json.Unmarshal([]byte(cachedData), &url); err == nil {
			return &url, nil
		}
	}

	// Se não encontrado no Redis, buscar no MongoDB
	var url entity.URL
	err = r.collection.FindOne(ctx, bson.M{"short_id": shortID}).Decode(&url)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// Salvar no Redis para futuras consultas
	data, _ := json.Marshal(url)
	r.redis.Set(ctx, cacheKey, data, time.Hour) // Cache por 1 hora

	return &url, nil
}
