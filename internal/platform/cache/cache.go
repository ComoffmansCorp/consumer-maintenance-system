// Package cache is a thin wrapper around a Redis client, used for two
// independent purposes: access-token revocation (see internal/auth's
// Logout flow and internal/platform/middleware's JWTAuth) and cache-aside
// reads for the catalog domain. Kept minimal on purpose -- only the
// Get/Set/Del/Exists shape those two use cases actually need.
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func Connect(ctx context.Context, redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	rdb := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := rdb.Ping(pingCtx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Client{Client: rdb}, nil
}

// Get returns the cached value and whether it was present. A missing key
// (redis.Nil) is reported as found=false with a nil error; any other error
// is a real cache-layer failure the caller should log and fall through to
// the source of truth for, rather than a 500.
func (c *Client) Get(ctx context.Context, key string) (value string, found bool, err error) {
	value, err = c.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("cache get %s: %w", key, err)
	}
	return value, true, nil
}

// Set writes value under key with the given TTL. A zero or negative ttl
// means "no expiry", matching redis.Client.Set's own semantics.
func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if err := c.Client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("cache set %s: %w", key, err)
	}
	return nil
}

// Del removes zero or more keys. A no-op on an empty key list, since
// redis.Client.Del would otherwise error on a zero-argument call.
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := c.Client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("cache del %v: %w", keys, err)
	}
	return nil
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache exists %s: %w", key, err)
	}
	return n > 0, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
