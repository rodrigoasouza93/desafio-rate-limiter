package infra

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func RateLimitMiddleware(rl *RateLimiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			token := r.Header.Get("API_KEY")

			var key string
			var isToken bool

			if token != "" {
				key = token
				isToken = true
			} else {
				key = ip
				isToken = false
			}

			if rl.IsBlocked(key) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			if !rl.Allow(key, isToken) {
				rl.Block(key, isToken)
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
