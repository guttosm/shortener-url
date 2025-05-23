package app

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo establishes a connection to a MongoDB database and returns the client and database instances.
//
// Parameters:
// - uri (string): The connection URI for the MongoDB server.
// - dbName (string): The name of the database to connect to.
//
// Returns:
// - *mongo.Client: The MongoDB client instance.
// - *mongo.Database: The MongoDB database instance.
// - error: An error if the connection or ping to the database fails.
func ConnectMongo(uri, dbName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	log.Println("Connected to MongoDB")

	db := client.Database(dbName)
	return client, db, nil
}
