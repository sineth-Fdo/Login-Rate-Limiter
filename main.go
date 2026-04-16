package main

import (
	"log"
	"net/http"

	"login-rate-limiter/handler"
	"login-rate-limiter/middleware"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handler.LoginHandler)

	handlerWithMiddleware := middleware.RateLimiter(mux)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))
}