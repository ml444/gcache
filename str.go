package gcache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"time"
)

type StrCache struct {
	cli    *redis.Client
	expire time.Duration
}

func NewStrCache(cli *redis.Client, expire time.Duration) *StrCache {
	return &StrCache{
		cli:    cli,
		expire: expire,
	}
}

func (c *StrCache) getRandExpire() time.Duration {
	n := rand.Intn(100)
	e := c.expire / 10
	r := e * time.Duration(n/100)
	return c.expire + r
}

func (c *StrCache) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

func (c *StrCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.cli.Set(ctx, key, value, c.getRandExpire()).Err()
}

func (c *StrCache) Del(ctx context.Context, key string) error {
	return c.cli.Del(ctx, key).Err()
}
func (c *StrCache) Exist(ctx context.Context, key string) (bool, error) {
	res, err := c.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res > 0 {
		return true, nil
	}
	return false, nil
}
