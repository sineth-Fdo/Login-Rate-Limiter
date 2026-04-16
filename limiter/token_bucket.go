package limiter

import (
	"context"
	"time"

	"login-rate-limiter/store"

	"github.com/redis/go-redis/v9"
)

// Lua script for atomic token bucket operation in Redis
var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

if tokens == nil then
    tokens = capacity
    last_refill = now
end

local elapsed = now - last_refill
local refill = math.floor(elapsed * refill_rate)
if refill > 0 then
    tokens = math.min(capacity, tokens + refill)
    last_refill = now
end

if tokens > 0 then
    tokens = tokens - 1
    redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)
    redis.call('EXPIRE', key, capacity / refill_rate + 10)
    return 1
end

redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)
redis.call('EXPIRE', key, capacity / refill_rate + 10)
return 0
`)

type TokenBucket struct {
	key        string
	capacity   int
	refillRate int
}

func NewTokenBucket(key string, capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		key:        key,
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func (tb *TokenBucket) Allow() bool {
	ctx := context.Background()
	now := float64(time.Now().UnixMilli()) / 1000.0

	result, err := tokenBucketScript.Run(ctx, store.Client,
		[]string{tb.key},
		tb.capacity, tb.refillRate, now,
	).Int()

	if err != nil {
		return false
	}

	return result == 1
}
