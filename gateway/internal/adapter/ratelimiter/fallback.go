package ratelimiter

import "context"

type RateLimiter interface {
	Allow(context.Context, string) (bool, error)
}

type FallbackRateLimiter struct {
	primary   RateLimiter
	secondary RateLimiter
}

func NewFallbackRateLimiter(primary RateLimiter, secondary RateLimiter) FallbackRateLimiter {
	return FallbackRateLimiter{
		primary:   primary,
		secondary: secondary,
	}
}

func (rl FallbackRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	allowed, err := rl.primary.Allow(ctx, key)
	if err == nil {
		return allowed, nil
	}

	return rl.secondary.Allow(ctx, key)
}
