package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/auth"
	apphttp "github.com/guttosm/url-shortener/internal/http"
	"github.com/guttosm/url-shortener/internal/repository/mongo"
	redisrepo "github.com/guttosm/url-shortener/internal/repository/redis"
	"github.com/guttosm/url-shortener/internal/service"
	"github.com/redis/go-redis/v9"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitializeApp sets up the entire application and its dependencies.
func InitializeApp() (*gin.Engine, func(), error) {
	// --- MongoDB setup
	mongoClient, err := mongodriver.Connect(context.Background(), options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		return nil, nil, err
	}
	db := mongoClient.Database(config.AppConfig.MongoDB)
	urlMongoRepo := mongo.NewURLMongoRepository(db.Collection("urls"))

	// --- Redis setup
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.AppConfig.RedisURI,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, nil, err
	}
	urlRedisRepo := redisrepo.NewURLRedisRepository(redisClient)

	// --- Services
	urlService := service.NewURLService(urlMongoRepo, urlRedisRepo)

	// --- HTTP Handler and Router
	authCfg := config.AppConfig.Auth
	handler := apphttp.NewHandler(urlService, authCfg)
	validator := &auth.JWTValidator{}
	router := apphttp.NewRouter(handler, validator)

	// --- Cleanup resources
	cleanup := func() {
		_ = mongoClient.Disconnect(context.Background())
		_ = redisClient.Close()
	}

	return router, cleanup, nil
}
