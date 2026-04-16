package limiter

import (
	"context"
	"time"

	"login-rate-limiter/store"

	"github.com/redis/go-redis/v9"
)

// Lua script for atomic sliding window operation in Redis
var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local window_ms = tonumber(ARGV[1])
local max_req = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local cutoff = now - window_ms

-- Remove expired entries
redis.call('ZREMRANGEBYSCORE', key, '-inf', cutoff)

-- Count current requests in window
local count = redis.call('ZCARD', key)

if count >= max_req then
    return 0
end

-- Add current request
redis.call('ZADD', key, now, now .. '-' .. math.random(1000000))
redis.call('PEXPIRE', key, window_ms)
return 1
`)

type SlidingWindow struct {
	key        string
	windowSize time.Duration
	maxReq     int
}

func NewSlidingWindow(key string, windowSize time.Duration, maxReq int) *SlidingWindow {
	return &SlidingWindow{
		key:        key,
		windowSize: windowSize,
		maxReq:     maxReq,
	}
}

func (sw *SlidingWindow) Allow() bool {
	ctx := context.Background()
	now := time.Now().UnixMilli()

	result, err := slidingWindowScript.Run(ctx, store.Client,
		[]string{sw.key},
		sw.windowSize.Milliseconds(), sw.maxReq, now,
	).Int()

	if err != nil {
		return false
	}

	return result == 1
}
