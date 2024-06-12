package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra"
	"golang.org/x/time/rate"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	r := mux.NewRouter()
	r.Use(infra.RateLimitMiddleware(rl))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
