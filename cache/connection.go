package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {
	url := os.Getenv("REDIS_URL")

	if url == "" {
		log.Fatalf("missing or empty env variable REDIS_URL\n")
	}

	opts, err := redis.ParseURL(url)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	_, err = client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
