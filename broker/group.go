package broker

import (
	"github.com/ml444/gcache/config"
)

type Group struct {
	Name       string
	Address    string
	shardCount uint64
	shardMap   map[int]*Shard
	hashFunc   func(key string) uint64
	loadFunc   func(key string) ([]byte, error)
}

func NewGroup(cfg *config.GroupConfig, address string, shardSerialNoList []int) *Group {
	c := &Group{
		Name:       cfg.Name,
		Address:    address,
		shardCount: uint64(cfg.ShardCount),
		shardMap:   make(map[int]*Shard, len(shardSerialNoList)),
		hashFunc:   cfg.HashFunc,
		loadFunc:   cfg.LoadFunc,
	}
	for i, serialNo := range shardSerialNoList {
		c.shardMap[i] = newShard(serialNo, cfg.Shard.MaxCount, cfg.Shard.MaxSize, cfg.Shard.Strategy)
	}
	return c
}

func (c *Group) Set(key string, value []byte) error {
	hash := c.hashFunc(key)
	shard, ok := c.shardMap[int(hash%c.shardCount)]
	if !ok {
		return ErrNotFound
	}
	return shard.Set(hash, key, value)
}

func (c *Group) Get(key string) ([]byte, error) {
	hash := c.hashFunc(key)
	shard, ok := c.shardMap[int(hash%c.shardCount)]
	if !ok {
		return nil, ErrNotFound
	}
	if !shard.isMaster {
		return nil, ErrNoMaster
	}
	return shard.Get(hash, key)
}
