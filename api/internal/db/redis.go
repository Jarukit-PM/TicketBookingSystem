package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisPingTimeout = 5 * time.Second

// ConnectRedis parses redisURL and returns a connected client.
func ConnectRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis parse url: %w", err)
	}

	client := redis.NewClient(opts)
	return client, nil
}

// MustConnectRedis connects to Redis or panics.
func MustConnectRedis(redisURL string) *redis.Client {
	client, err := ConnectRedis(redisURL)
	if err != nil {
		panic(err)
	}
	return client
}

// PingRedis verifies the Redis connection.
func PingRedis(ctx context.Context, client *redis.Client) error {
	ctx, cancel := context.WithTimeout(ctx, redisPingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}
