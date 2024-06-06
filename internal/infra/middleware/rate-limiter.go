package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/rodrigoasouza93/desafio-rate-limiter/configs"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/application/protocols"
)

type RateLimiter struct {
	configs configs.LimitConfig
	repo    protocols.LimiterRepo
}

const (
	token = "token"
	ip    = "ip"
)

func NewRateLimiter(repo protocols.LimiterRepo, configs configs.LimitConfig) *RateLimiter {
	return &RateLimiter{repo: repo, configs: configs}
}

func (rl *RateLimiter) Limit(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			apiKey              = r.Header.Get("API_KEY")
			ip                  = strings.Split(r.RemoteAddr, ":")[0]
			limiterKey, keyType = rl.getKeyType(apiKey, ip)
			err                 error
			blocked             bool
		)

		if blocked, err = rl.IsBlocked(limiterKey, keyType); blocked {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (rl *RateLimiter) IsBlocked(key, keyType string) (bool, error) {
	status, err := rl.repo.Get(key)

	if status == "blocked" {
		return true, nil
	}

	if errors.Is(err, redis.Nil) {
		rl.startCount(key, keyType)
		return false, nil
	}

	if err == nil {
		err = rl.increaseCount(status, key, keyType)
		return false, err
	}

	fmt.Println(err)
	return true, err
}

func (rl *RateLimiter) increaseCount(status, key, keyType string) error {
	rl.repo.Increment(key)
	reqs, err := strconv.Atoi(status)
	if err != nil {
		return err
	}
	rl.blockCheck(reqs, key, keyType)
	return nil
}

func (rl *RateLimiter) getKeyType(apiKey, ipAddress string) (limiterKey, keyType string) {
	if apiKey == rl.configs.AllowedToken {
		limiterKey = apiKey
		keyType = token
		return
	}
	limiterKey = ipAddress
	keyType = ip
	return
}

func (rl *RateLimiter) startCount(key, keyType string) {
	if keyType == token {
		rl.repo.Set(key, "1", time.Duration(rl.configs.MaxRequestsToken)*time.Second)
		return
	}
	rl.repo.Set(key, "1", time.Duration(rl.configs.MaxRequestsIp)*time.Second)
}

func (rl *RateLimiter) blockCheck(reqs int, key, keyType string) {
	if keyType == token {
		if reqs == rl.configs.MaxRequestsToken-1 {
			rl.repo.Set(key, "blocked", rl.configs.TokenBlockTime)
			return
		}
		return
	}
	if reqs == rl.configs.MaxRequestsIp-1 {
		rl.repo.Set(key, "blocked", rl.configs.IpBlockTime)
		return
	}
	return
}
