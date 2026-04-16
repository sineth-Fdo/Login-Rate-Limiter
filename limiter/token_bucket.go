package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int
	lastRefilled time.Time
	mu           sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillRate:   refillRate,
		lastRefilled: time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefilled).Seconds()

	refill := int(elapsed * float64(tb.refillRate))
	if refill > 0 {
		tb.tokens += refill
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lastRefilled = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}