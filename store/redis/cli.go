package redis

import (
	"context"
	"errors"
	"os"
	"strconv"

	log "github.com/ml444/glog"
	"github.com/redis/go-redis/v9"
)

const (
	RedisAddr       = "REDIS_ADDRESS"
	RedisUser       = "REDIS_USER"
	RedisPwd        = "REDIS_PWD"
	RedisDb         = "REDIS_DB"
	RedisClientName = "REDIS_CLIENT_NAME"
)

var ErrNotFoundRedisAddr = errors.New("not found redis address")

type RedisConfig = redis.Options

func GetConfig4Env() (*RedisConfig, error) {
	addr := os.Getenv(RedisAddr)
	if addr == "" {
		return nil, ErrNotFoundRedisAddr
	}
	db := 0
	dbStr := os.Getenv(RedisDb)
	if dbStr != "" {
		db1, err := strconv.ParseInt(dbStr, 10, 64)
		if err != nil {
			log.Errorf("parsing env %s err: %s", RedisDb, err.Error())
			return nil, err
		}
		db = int(db1)
	}
	return &RedisConfig{
		Addr:       addr,
		DB:         db,
		ClientName: os.Getenv(RedisClientName),
		Username:   os.Getenv(RedisUser),
		Password:   os.Getenv(RedisPwd),
	}, nil
}
func GetRedisCli(config *RedisConfig) (*redis.Client, error) {
	if config == nil {
		var err error
		config, err = GetConfig4Env()
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
	}
	return redis.NewClient(config), nil
}

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, key string) error
	Exist(ctx context.Context, key string) (bool, error)
}
