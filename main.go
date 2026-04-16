package main

import (
	"log"
	"net/http"
	"os"

	"login-rate-limiter/handler"
	"login-rate-limiter/middleware"
	"login-rate-limiter/store"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	store.InitRedis(redisAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handler.LoginHandler)

	handlerWithMiddleware := middleware.RateLimiter(mux)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
}
