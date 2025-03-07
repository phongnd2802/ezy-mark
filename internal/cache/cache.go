package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache key not found")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetObj(ctx context.Context, key string, obj interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	SetEx(ctx context.Context, key string, value interface{}, expirationTimeSecond int64) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (int64, error)
}

type redisCache struct {
	client *redis.Client
}

// GetObj implements Cache.
func (r *redisCache) GetObj(ctx context.Context, key string, obj interface{}) error {
	result, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%w: %s", ErrCacheMiss, key)
		}
		return err
	}
	// convert result to obj
	if err := sonic.Unmarshal(result, obj); err != nil {
		return err
	}
	return nil
}

// Exists implements Cache.
func (r *redisCache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Del implements Cache.
func (r *redisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// TTL implements Cache.
func (r *redisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// SetEx implements Cache.
func (r *redisCache) SetEx(ctx context.Context, key string, value interface{}, expirationTimeSecond int64) error {
	return r.client.Set(ctx, key, value, time.Duration(expirationTimeSecond)*time.Second).Err()
}

// Get implements Cache.
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set implements Cache.
func (r *redisCache) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func NewRedisClient(client *redis.Client) Cache {
	return &redisCache{
		client: client,
	}
}
