package infra

import (
	"fmt"
	"time"

	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/domain"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	persistence    domain.Persistence
	ipLimiter      map[string]*rate.Limiter
	tokenLimiter   map[string]*rate.Limiter
	rateLimitIP    rate.Limit
	rateLimitToken rate.Limit
	blockTime      time.Duration
}

func NewRateLimiter(persistence domain.Persistence, rateLimitIP, rateLimitToken rate.Limit, blockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		persistence:    persistence,
		ipLimiter:      make(map[string]*rate.Limiter),
		tokenLimiter:   make(map[string]*rate.Limiter),
		rateLimitIP:    rateLimitIP,
		rateLimitToken: rateLimitToken,
		blockTime:      blockTime,
	}
}

func (rl *RateLimiter) getLimiter(key string, isToken bool) *rate.Limiter {
	if isToken {
		if limiter, exists := rl.tokenLimiter[key]; exists {
			return limiter
		}
		limiter := rate.NewLimiter(rl.rateLimitToken, int(rl.rateLimitToken))
		rl.tokenLimiter[key] = limiter
		return limiter
	} else {
		if limiter, exists := rl.ipLimiter[key]; exists {
			return limiter
		}
		limiter := rate.NewLimiter(rl.rateLimitIP, int(rl.rateLimitIP))
		rl.ipLimiter[key] = limiter
		return limiter
	}
}

func (rl *RateLimiter) Allow(key string, isToken bool) bool {
	limiter := rl.getLimiter(key, isToken)
	return limiter.Allow()
}

func (rl *RateLimiter) Block(key string, isToken bool) {
	rl.persistence.Set(key, "blocked", int(rl.blockTime.Seconds()))
}

func (rl *RateLimiter) IsBlocked(key string) bool {
	val, err := rl.persistence.Get(key)
	if err != nil {
		if err.Error() == "redis: nil" {
			return false
		}
		fmt.Println("Error checking block status:", err)
		return false
	}
	return val == "blocked"
}
