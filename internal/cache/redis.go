package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// Init initializes the Redis client
func Init(url string) error {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	Client = redis.NewClient(opts)
	return Client.Ping(ctx).Err()
}

// Get retrieves a value from the cache
func Get(key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Set stores a value in the cache with a TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	return Client.Set(ctx, key, value, ttl).Err()
}

// Invalidate removes a key from the cache
func Invalidate(key string) error {
	return Client.Del(ctx, key).Err()
}

// InvalidatePattern removes keys matching a pattern
func InvalidatePattern(pattern string) error {
	iter := Client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := Client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
