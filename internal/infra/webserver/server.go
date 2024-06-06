package webserver

import (
	"log"
	"net/http"
)

type (
	APIServer struct {
		addr    string
		limiter Limiter
	}

	Limiter interface {
		Limit(next http.Handler) http.HandlerFunc
	}
)

func NewAPIServer(addr string, limiter Limiter) *APIServer {
	return &APIServer{
		addr:    ":" + addr,
		limiter: limiter,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	server := http.Server{
		Addr:    s.addr,
		Handler: s.limiter.Limit(router),
	}

	log.Printf("Server has started on %s\n", s.addr)

	return server.ListenAndServe()
}
