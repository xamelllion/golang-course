package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type RedisCacher struct {
	client *redis.Client
}

func NewRedisCacher(client *redis.Client) RedisCacher {
	return RedisCacher{
		client: client,
	}
}

func (rc RedisCacher) Get(ctx context.Context, key string) ([]byte, bool, error) {
	data, err := rc.client.Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return []byte(data), true, nil
}

func (rc RedisCacher) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	_, err := rc.client.Set(key, value, ttl).Result()
	if err != nil {
		return err
	}

	return nil
}
