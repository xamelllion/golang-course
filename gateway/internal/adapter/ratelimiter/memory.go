package ratelimiter

import (
	"context"
	"time"
)

type clientBucket struct {
	tokens     float64
	lastRefill time.Time
}

type MemoryBucketRateLimiter struct {
	capacity  int
	reqPerSec float64
	clients   map[string]*clientBucket
}

func NewMemoryBucketRateLimiter(capacity int, reqPerSec float64) MemoryBucketRateLimiter {
	return MemoryBucketRateLimiter{
		capacity:  capacity,
		reqPerSec: reqPerSec,
		clients:   map[string]*clientBucket{},
	}
}

func (rl MemoryBucketRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	key = "rl:ip:" + key

	client, ok := rl.clients[key]
	if !ok {
		client = &clientBucket{
			tokens:     float64(rl.capacity),
			lastRefill: time.Now(),
		}

		rl.clients[key] = client
	}

	newTokens := time.Since(client.lastRefill).Seconds() * rl.reqPerSec
	client.tokens = min(float64(rl.capacity), client.tokens+newTokens)
	client.lastRefill = time.Now()

	if client.tokens >= 1.0 {
		client.tokens -= 1.0
		return true, nil
	} else {
		return false, nil
	}
}
