package main

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra"
	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {

	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	rateLimitIP, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	rateLimitToken, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
	blockTime, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	persistence := infra.NewRedisPersistence(redisClient)
	rl := infra.NewRateLimiter(persistence, rate.Limit(rateLimitIP), rate.Limit(rateLimitToken), time.Duration(blockTime)*time.Second)

	key := "test_ip"
	isToken := false

	for i := 0; i < rateLimitIP; i++ {
		if !rl.Allow(key, isToken) {
			t.Errorf("Expected to allow request %d", i+1)
		}
	}

	if rl.Allow(key, isToken) {
		t.Errorf("Expected to block request")
	}

	rl.Block(key, isToken)

	if !rl.IsBlocked(key) {
		t.Errorf("Expected key to be blocked")
	}

	time.Sleep(time.Duration(blockTime) * time.Second)

	if rl.IsBlocked(key) {
		t.Errorf("Expected key to be unblocked")
	}

	key = "test_token"
	isToken = true

	for i := 0; i < rateLimitToken; i++ {
		if !rl.Allow(key, isToken) {
			t.Errorf("Expected to allow request %d", i+1)
		}
	}

	if rl.Allow(key, isToken) {
		t.Errorf("Expected to block request")
	}

	rl.Block(key, isToken)

	if !rl.IsBlocked(key) {
		t.Errorf("Expected key to be blocked")
	}

	time.Sleep(time.Duration(blockTime) * time.Second)

	if rl.IsBlocked(key) {
		t.Errorf("Expected key to be unblocked")
	}
}
