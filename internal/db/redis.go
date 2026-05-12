package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var redis_url string

func InitRedis() {
	redis_url = os.Getenv("REDIS_URL")
	if redis_url == "" {
		redis_url = "redis://some-url"
	}
	
	var err error
	Redis, err = newRedisClient()
	if err != nil {
		log.Fatalf("Unable to connect to redis - %s", redis_url)
	}
}

// newRedisClient establishes a SOTA connection pool to Redis.
func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_url,
		Password: "",  // No password set in our docker-compose
		DB:       0,   // Use default DB
		PoolSize: 100, // SOTA: Configure connection pooling for high concurrency
	})

	// Ping to ensure connection is alive
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
