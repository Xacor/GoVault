package redis

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(connURL string) (*RedisClient, error) {
	opt, err := redis.ParseURL(connURL)
	if err != nil {
		return nil, errors.Wrap(err, "connection to redis failed")
	}

	return &RedisClient{redis.NewClient(opt)}, nil
}

func (c *RedisClient) Set(ctx context.Context, key string, data any) error {
	err := c.rdb.HSet(ctx, key, data).Err()
	if err != nil {
		return errors.Wrapf(err, "unable to set, key: %+v, data: %+v", key, data)
	}

	return nil
}

func (c *RedisClient) Get(ctx context.Context, key string, dst any) error {
	err := c.rdb.HGetAll(ctx, key).Scan(dst)
	if err != nil {
		return errors.Wrapf(err, "unable to get data by key: %v", key)
	}

	return nil
}
