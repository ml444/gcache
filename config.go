package gcache

import (
	"hash/fnv"

	"github.com/ml444/gcache/strategy"
)

type Config struct {
	ShardCount int
	HashFunc   func(key string) uint64
	LoadFunc   func(key string) ([]byte, error)
	Shard      struct {
		MaxSize  int
		MaxCount int
		Strategy strategy.IStrategy
	}
}

func DefaultConfig() *Config {
	return &Config{
		ShardCount: 1024,
		HashFunc: func(key string) uint64 {
			hasher := fnv.New64a()
			hasher.Write([]byte(key))
			return hasher.Sum64()
		},
		LoadFunc: nil,
		Shard: struct {
			MaxSize  int
			MaxCount int
			Strategy strategy.IStrategy
		}{
			MaxSize:  1024 * 1024 * 1024,
			MaxCount: 1000000,
			Strategy: nil,
		},
	}
}
