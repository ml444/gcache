package gcache

import (
	"sync"
	"time"
)

type Cacher interface {
	// Get 根据Key 荻取
	Get(key string, ptrValue interface{}) (err error)
	// 根据Key列表获取，如果末命中可能触发加载动作
	GetKeys(keys ...string)
	// 根据Key列表荻取，如果未命中不会触发加载动作
	GetAllPresent(keys ...string)
	// 获取所有key 列表
	Keys()
	// 写人一个k/v
	Set(key string, value interface{})
	// 写人一个k/v，并设置过期时间
	SetEx(key string, value interface{}, expire time.Duration) (err error)
	// 将entries写人缓存
	SetAll(entries interface{})
	// 如果缓存中没有则写入
	SetIfAbsent(key string, value interface{})

	// 比较旧的值相同时才置换
	CompareAndSwap(key string, newValue, oldValue interface{})

	// 删除一个key
	Delete(keys string) (err error)
	// 匹配k/v刪除
	DeleteCompare(key string, oldValue interface{})
	// 根据key列表删除
	DeleteKeys(keys ...string) (affected int64, err error)

	// 清空缓存
	Clear()
	Flush() (err error)
	Close() (err error)
}

type cacheShard struct {
	sync.RWMutex
}

func (c *cacheShard) Set(key string, value []byte) {

}

type Cache struct {
	shardCount uint64
	shards     []*cacheShard
	locks      []*sync.RWMutex
	hashFunc   func(key string) uint64
	loadFunc   func(key string) (interface{}, error)
}

func (c *Cache) Set(key string, value []byte) {
	hash := c.hashFunc(key)
	shard := c.shards[hash%c.shardCount]
	shard.Set(key, value)
}

func (c *Cache) Get(key string) []byte {
	return nil
}
