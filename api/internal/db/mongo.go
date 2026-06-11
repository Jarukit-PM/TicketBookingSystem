package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoConnectTimeout = 10 * time.Second

// ConnectMongo establishes a MongoDB client for the given URI.
func ConnectMongo(ctx context.Context, uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	return client, nil
}

// MustConnectMongo connects to MongoDB or panics.
func MustConnectMongo(ctx context.Context, uri string) *mongo.Client {
	client, err := ConnectMongo(ctx, uri)
	if err != nil {
		panic(err)
	}
	return client
}

// PingMongo verifies the MongoDB connection.
func PingMongo(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("mongo ping: %w", err)
	}
	return nil
}

// DisconnectMongo closes the MongoDB client.
func DisconnectMongo(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("mongo disconnect: %w", err)
	}
	return nil
}
