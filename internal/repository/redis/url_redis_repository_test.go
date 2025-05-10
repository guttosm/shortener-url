package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/guttosm/url-shortener/internal/entity"
	repo "github.com/guttosm/url-shortener/internal/repository/redis"
	redislib "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedis(t *testing.T) (*redislib.Client, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:6.2",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start Redis container: %v", err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get Redis host: %v", err)
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("failed to get Redis port: %v", err)
	}

	client := redislib.NewClient(&redislib.Options{
		Addr: host + ":" + port.Port(),
		DB:   0,
	})

	teardown := func() {
		client.Close()
		redisC.Terminate(ctx)
	}

	return client, teardown
}

func TestURLRedisRepository(t *testing.T) {
	client, teardown := setupRedis(t)
	defer teardown()

	repository := repo.NewURLRedisRepository(client)
	ctx := context.Background()

	url := &entity.URL{
		ShortID:   "abc123",
		Original:  "https://example.com",
		CreatedAt: time.Now(),
	}

	// Set and Get by Original URL
	err := repository.SetByOriginalURL(ctx, url)
	assert.NoError(t, err)

	result, err := repository.GetByOriginalURL(ctx, url.Original)
	assert.NoError(t, err)
	assert.Equal(t, url.ShortID, result.ShortID)

	// Set and Get by ShortID
	err = repository.SetByShortID(ctx, url)
	assert.NoError(t, err)

	result, err = repository.GetByShortID(ctx, url.ShortID)
	assert.NoError(t, err)
	assert.Equal(t, url.Original, result.Original)
}
