package ratelimiter

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// https://github.com/redis/docs/blob/main/content/develop/use-cases/rate-limiter/go/token_bucket.go

const tokenBucketScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local refill_interval = tonumber(ARGV[3])
local now = tonumber(ARGV[4])

-- Get current state or initialize
local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
local tokens = tonumber(bucket[1])
local last_refill = tonumber(bucket[2])

-- Initialize if this is the first request
if tokens == nil then
    tokens = capacity
    last_refill = now
end

-- Calculate token refill
local time_passed = now - last_refill
local refills = math.floor(time_passed / refill_interval)

if refills > 0 then
    tokens = math.min(capacity, tokens + (refills * refill_rate))
    last_refill = last_refill + (refills * refill_interval)
end

-- Try to consume a token
local allowed = 0
if tokens >= 1 then
    tokens = tokens - 1
    allowed = 1
end

-- Update state
redis.call('HMSET', key, 'tokens', tokens, 'last_refill', last_refill)

-- Return result: allowed (1 or 0) and remaining tokens
return {allowed, tokens}
`

type RedisBucketRateLimiter struct {
	client         redis.Client
	capacity       int
	refillRate     float64
	refillInterval time.Duration
	scriptSHA      string
	scriptLoaded   bool
}

func NewRedisRateLimiter(client redis.Client, reqPerSecond float64) RedisBucketRateLimiter {
	h := sha1.New()
	h.Write([]byte(tokenBucketScript))
	sha := fmt.Sprintf("%x", h.Sum(nil))

	return RedisBucketRateLimiter{
		client:         client,
		capacity:       10,
		refillRate:     1.0,
		refillInterval: time.Duration(1000.0/reqPerSecond) * time.Millisecond,
		scriptSHA:      sha,
		scriptLoaded:   false,
	}
}

func (rl RedisBucketRateLimiter) ensureScriptLoaded(ctx context.Context) {
	if !rl.scriptLoaded {
		sha, err := rl.client.ScriptLoad(tokenBucketScript).Result()
		if err == nil {
			rl.scriptSHA = sha
			rl.scriptLoaded = true
		}
	}
}

func (rl RedisBucketRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	rl.ensureScriptLoaded(ctx)

	now := float64(time.Now().UnixMicro()) / 1e6

	args := []interface{}{
		rl.capacity,
		rl.refillRate,
		rl.refillInterval,
		now,
	}

	// Try EVALSHA first (faster if script is cached)
	result, err := rl.client.EvalSha(rl.scriptSHA, []string{key}, args...).Result()
	if err != nil {
		// Script not in cache, fall back to EVAL
		result, err = rl.client.Eval(tokenBucketScript, []string{key}, args...).Result()
		if err != nil {
			return false, fmt.Errorf("token bucket eval failed: %w", err)
		}
		rl.scriptLoaded = false
	}

	resultSlice, _ := result.([]int64)

	allowed := resultSlice[0] == 1

	return allowed, nil
}
