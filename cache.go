package gcache

import (
	"time"
)

type ICache interface {
	// Get 根据Key 荻取
	Get(key string) (value []byte, err error)
	// GetKeys 根据Key列表获取，如果末命中可能触发加载动作
	GetKeys(keys ...string)
	// GetKeysPresent 根据Key列表荻取，如果未命中不会触发加载动作
	GetKeysPresent(keys ...string) (map[string][]byte, error)
	// Keys 获取所有key列表
	Keys()
	// Set 写入一个k/v
	Set(key string, value []byte) (err error)
	// SetEx 写入一个k/v，并设置过期时间
	SetEx(key string, value []byte, expire time.Duration) (err error)
	// SetAll 将entries写人缓存
	SetAll(entries interface{})
	// SetIfAbsent 如果缓存中没有则写入
	SetIfAbsent(key string, value []byte) (err error)

	// CompareAndSwap 比较旧的值相同时才置换
	CompareAndSwap(key string, newValue, oldValue []byte) (err error)

	// Delete 删除一个key
	Delete(key string) (err error)
	// DeleteCompare 匹配k/v刪除
	DeleteCompare(key string, oldValue []byte) (err error)
	// DeleteKeys 根据key列表删除
	DeleteKeys(keys ...string) (affected int64, err error)

	// Clear 清空缓存
	Clear()
	Flush() (err error)
	Close() (err error)
}
