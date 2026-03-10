// this file contains the implementation of a function to create a new Redis client using the provided configuration. The function attempts to connect to the Redis server and returns the client if successful, or an error if the connection fails.

package database

import (
	"cinema-booking/backend/internal/config"
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
