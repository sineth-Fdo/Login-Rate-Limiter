package middleware

import (
	"net"
	"net/http"
	"time"

	"login-rate-limiter/limiter"
)

// Global limiter: 100 burst, 50 req/sec
var globalLimiter = limiter.NewTokenBucket("ratelimit:global", 100, 50)

func getIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getUser(r *http.Request) string {
	user := r.Header.Get("X-User-ID")
	if user == "" {
		return "anonymous"
	}
	return user
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Global protection
		if !globalLimiter.Allow() {
			http.Error(w, "Global rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		ip := getIP(r)
		user := getUser(r)

		// 2. Per-IP: 10 requests per 10 seconds
		ipLimiter := limiter.NewSlidingWindow("ratelimit:ip:"+ip, 10*time.Second, 10)
		if !ipLimiter.Allow() {
			http.Error(w, "Too many requests from this IP", http.StatusTooManyRequests)
			return
		}

		// 3. Per-user: 20 requests per 10 seconds
		userLimiter := limiter.NewSlidingWindow("ratelimit:user:"+user, 10*time.Second, 20)
		if !userLimiter.Allow() {
			http.Error(w, "User rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
