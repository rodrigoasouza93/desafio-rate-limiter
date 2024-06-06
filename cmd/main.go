package main

import (
	"github.com/rodrigoasouza93/desafio-rate-limiter/configs"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra/database/redis"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra/middleware"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra/webserver"
)

func main() {
	configs, err := configs.GetConfig(".")
	if err != nil {
		panic(err)
	}

	repo := redis.NewRedisLimiterRepo(configs.Host, "", 0)
	limiter := middleware.NewRateLimiter(repo, configs.GetLimitConfig())
	server := webserver.NewAPIServer(configs.Port, limiter)
	server.Run()
}
