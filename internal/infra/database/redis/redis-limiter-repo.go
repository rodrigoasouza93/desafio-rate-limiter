package redis

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisLimiterRepo struct {
	client *redis.Client
}

func NewRedisLimiterRepo(address, password string, db int) *RedisLimiterRepo {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return &RedisLimiterRepo{
		client: client,
	}
}

func (r *RedisLimiterRepo) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisLimiterRepo) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisLimiterRepo) Increment(key string) error {
	return r.client.Incr(key).Err()
}
