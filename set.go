package gcache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type SetCache struct {
	cli *redis.Client
	key string
}

func NewSetCache(name string, cli *redis.Client) *SetCache {
	return &SetCache{
		cli: cli,
		key: name,
	}
}

func (c *SetCache) Get(ctx context.Context, _ string) (string, error) {
	return c.cli.SPop(ctx, c.key).Result()
}

func (c *SetCache) Set(ctx context.Context, member string, _ interface{}) error {
	return c.cli.SAdd(ctx, c.key, member).Err()
}

func (c *SetCache) Del(ctx context.Context, member string) error {
	return c.cli.SRem(ctx, c.key, member).Err()
}
func (c *SetCache) Exist(ctx context.Context, member string) (bool, error) {
	return c.cli.SIsMember(ctx, c.key, member).Result()
}
