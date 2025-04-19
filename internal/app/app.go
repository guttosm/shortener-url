package app

import (
	"context"

    // "github.com/gin-gonic/gin" // Removed unused import
	"github.com/guttosm/url-shortener/config"
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
func InitializeApp() (interface{}, func(), error) {
    // Example implementation
    mongoClient, _, err := ConnectMongo(config.AppConfig.MongoURI, config.AppConfig.MongoDB)
    if err != nil {
        return nil, nil, err
    }

    redisClient := redis.NewClient(&redis.Options{
        Addr: config.AppConfig.RedisURI,
    })
    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        return nil, nil, err
    }

    cleanup := func() {
        mongoClient.Disconnect(context.Background())
        redisClient.Close()
    }

    // Replace `nil` with your actual router initialization logic
    return nil, cleanup, nil
}
