package infra

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisPersistence struct {
	client *redis.Client
}

func NewRedisPersistence(client *redis.Client) *RedisPersistence {
	return &RedisPersistence{client: client}
}

func (r *RedisPersistence) Set(key string, value string, ttl int) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisPersistence) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}
