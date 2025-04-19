package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/http"
	"github.com/guttosm/url-shortener/internal/middleware"
	"github.com/redis/go-redis/v9"
)

// InitializeApp initializes the application by setting up the database connections,
// Redis client, URL module, and HTTP router.
// It returns the configured Gin engine, a cleanup function to close resources, and an error if any occurs.
//
// Returns:
// - *gin.Engine: The configured Gin HTTP router.
// - func(): A cleanup function to close database and Redis connections.
// - error: An error if there is an issue during initialization.
func InitializeApp() (*gin.Engine, func(), error) {
	client, db, err := ConnectMongo(config.AppConfig.MongoURI, config.AppConfig.MongoDB)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to mongoDB: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.AppConfig.RedisURI,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	urlModule := InitURLModule(db, redisClient)

	handler := http.NewHandler(urlModule.Service)
	router := http.NewRouter(handler)

	// Register the error handler middleware
	router.Use(middleware.ErrorHandler)

	cleanup := func() {
		_ = client.Disconnect(context.Background())
		_ = redisClient.Close()
	}

	return router, cleanup, nil
}
