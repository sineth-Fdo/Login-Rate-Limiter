package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"login-rate-limiter/limiter"
)

// Per-IP limiter
var ipLimiters = make(map[string]*limiter.SlidingWindow)

// Per-user limiter
var userLimiters = make(map[string]*limiter.SlidingWindow)

var mu sync.Mutex

// Global limiter
var globalLimiter = limiter.NewTokenBucket(100, 50) // 100 burst, 50 req/sec

func getIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getUser(r *http.Request) string {
	// Example: from header (replace with JWT in real app)
	user := r.Header.Get("X-User-ID")
	if user == "" {
		return "anonymous"
	}
	return user
}

func getIPLimiter(ip string) *limiter.SlidingWindow {
	mu.Lock()
	defer mu.Unlock()

	if l, exists := ipLimiters[ip]; exists {
		return l
	}

	// 10 requests per 10 seconds per IP
	l := limiter.NewSlidingWindow(10*time.Second, 10)
	ipLimiters[ip] = l
	return l
}

func getUserLimiter(user string) *limiter.SlidingWindow {
	mu.Lock()
	defer mu.Unlock()

	if l, exists := userLimiters[user]; exists {
		return l
	}

	// 20 requests per 10 seconds per user
	l := limiter.NewSlidingWindow(10*time.Second, 20)
	userLimiters[user] = l
	return l
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1️⃣ Global protection
		if !globalLimiter.Allow() {
			http.Error(w, "Global rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		ip := getIP(r)
		user := getUser(r)

		// 2️⃣ Per-IP
		if !getIPLimiter(ip).Allow() {
			http.Error(w, "Too many requests from this IP", http.StatusTooManyRequests)
			return
		}

		// 3️⃣ Per-user
		if !getUserLimiter(user).Allow() {
			http.Error(w, "User rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}