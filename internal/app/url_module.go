package app

import (
	"github.com/guttosm/url-shortener/internal/repository"
	mongoRepo "github.com/guttosm/url-shortener/internal/repository/mongo"
	redisRepo "github.com/guttosm/url-shortener/internal/repository/redis"
	"github.com/guttosm/url-shortener/internal/service"
	"github.com/redis/go-redis/v9"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type URLModule struct {
	Repository repository.URLRepository
	Service    service.URLService
}

func InitURLModule(db *mongoDriver.Database, redisClient *redis.Client) *URLModule {
	urlCollection := db.Collection("urls")
	urlRepo := mongoRepo.NewURLMongoRepository(urlCollection)
	urlCacheRepo := redisRepo.NewURLRedisRepository(redisClient)
	urlService := service.NewURLService(urlRepo, urlCacheRepo)

	return &URLModule{
		Repository: urlRepo,
		Service:    urlService,
	}
}
