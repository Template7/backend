package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *client) IncrLoginFailures(ctx context.Context, key string, window time.Duration) (int64, error) {
	k := loginFailuresKey(key)
	n, err := c.core.Incr(ctx, k).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		// first failure starts the window
		if err := c.core.Expire(ctx, k, window).Err(); err != nil {
			return n, err
		}
	}
	return n, nil
}

func (c *client) GetLoginFailures(ctx context.Context, key string) (int64, error) {
	n, err := c.core.Get(ctx, loginFailuresKey(key)).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	return n, err
}

func (c *client) ResetLoginFailures(ctx context.Context, key string) error {
	return c.core.Del(ctx, loginFailuresKey(key)).Err()
}

func loginFailuresKey(key string) string {
	return cacheKeyLoginFailures + ":" + key
}
