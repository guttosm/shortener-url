package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnectMongoWithTestContainer(t *testing.T) {
	// Create a MongoDB container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo:5.0", // Use a specific MongoDB version
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(30 * time.Second),
	}
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err, "Failed to start MongoDB container")
	defer func(mongoC testcontainers.Container, ctx context.Context, opts ...testcontainers.TerminateOption) {
		err := mongoC.Terminate(ctx, opts...)
		if err != nil {

		}
	}(mongoC, ctx)

	// Get the container's host and port
	host, err := mongoC.Host(ctx)
	assert.NoError(t, err, "Failed to get container host")
	port, err := mongoC.MappedPort(ctx, "27017")
	assert.NoError(t, err, "Failed to get container port")

	// Build the MongoDB URI
	uri := "mongodb://" + host + ":" + port.Port()
	dbName := "testdb"

	// Call the ConnectMongo function
	client, db, err := ConnectMongo(uri, dbName)
	if err == nil {
		defer func() {
			if client != nil {
				_ = client.Disconnect(context.TODO())
			}
		}()
	}

	// Assert no error occurred
	assert.NoError(t, err, "Expected no error while connecting to MongoDB")

	// Assert client and database are not nil
	assert.NotNil(t, client, "Expected client to be non-nil")
	assert.NotNil(t, db, "Expected database to be non-nil")
}
