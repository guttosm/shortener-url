package app

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/guttosm/url-shortener/config"
	"github.com/guttosm/url-shortener/internal/http"
	"github.com/redis/go-redis/v9"
)

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

	cleanup := func() {
		_ = client.Disconnect(context.Background())
		_ = redisClient.Close()
	}

	return router, cleanup, nil
}
