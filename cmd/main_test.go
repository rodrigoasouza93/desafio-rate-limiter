package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rodrigoasouza93/desafio-rate-limiter/configs"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra/database/redis"
	"github.com/rodrigoasouza93/desafio-rate-limiter/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	configs, err := configs.GetConfig(".")
	if err != nil {
		panic(err)
	}
	repo := redis.NewRedisLimiterRepo("localhost:6379", "", 0)
	limiter := middleware.NewRateLimiter(repo, configs.GetLimitConfig())
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	ts := httptest.NewServer(limiter.Limit(router))
	defer ts.Close()

	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
			return
		}

		req.Header.Add("API_KEY", "apitoken")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error parsing response body: %v", err)
		}

		if i <= 4 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
			assert.Equal(t, "Hello World!", string(body), "Response body should be 'Hello World!'")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "Status should be 429 Too Many Requests")
			assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", string(body), "Response body should be 'you have reached the maximum number of requests or actions allowed within a certain time frame'")
		}
	}

	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
			return
		}

		req.Header.Add("API_KEY", "xyz987654")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error parsing response body: %v", err)
		}

		if i <= 2 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
			assert.Equal(t, "Hello World!", string(body), "Response body should be 'Hello World!'")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "Status should be 429 Too Many Requests")
			assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", string(body), "Response body should be 'you have reached the maximum number of requests or actions allowed within a certain time frame'")
		}
	}

	t.Log("Standing by for 6 seconds in order to test further")
	time.Sleep(6 * time.Second)

	for i := 0; i < 10; i++ {
		req, err := http.NewRequest("GET", ts.URL, nil)
		if err != nil {
			t.Fatalf("Error creating request: %v", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Error creating requestT: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error parsing body: %v", err)
		}

		if i <= 2 {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200 OK")
			assert.Equal(t, "Hello World!", string(body), "Response body should be 'Hello World!'")
		} else {
			assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "Status should be 429 Too Many Requests")
			assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame\n", string(body), "Response body should be 'you have reached the maximum number of requests or actions allowed within a certain time frame'")
		}
	}
}
