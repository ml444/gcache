package gcache

import (
	"errors"
	"sync"

	"github.com/ml444/gcache/strategy"
)

var ErrNotFound = errors.New("not found")
var ErrHashCollision = errors.New("hashing collision")

type Shard struct {
	sync.RWMutex
	sequence int
	maxCount int
	maxSize  int
	indexBuf []byte
	kvBuf    []byte
	indexMap map[uint64]uint64
	stats    *Stats
	encoder  *Encode
	strategy strategy.IStrategy
}

func newShard(sequence, maxCount, maxSize int, strategy strategy.IStrategy) *Shard {
	return &Shard{
		maxSize:  maxSize,
		maxCount: maxCount,
		sequence: sequence,
		indexBuf: make([]byte, maxSize),
		kvBuf:    make([]byte, maxSize),
		indexMap: make(map[uint64]uint64),
		stats:    &Stats{},
		encoder:  &Encode{},
		strategy: strategy,
	}
}

func (c *Shard) Set(hash uint64, key string, value []byte) error {
	c.Lock()
	defer c.Unlock()
	if len(c.indexBuf)+len(value) > c.maxSize {
		if c.strategy != nil {
			// Calling the phase-out strategy
			v, ok := c.strategy.Evict()
			if ok {
				delete(c.indexMap, v)
				// TODO: delete kvBuf and indexBuf
			}
		} else {
			return errors.New("indexBuf is full")
		}
	}
	if len(c.indexMap)+1 > c.maxCount {
		if c.strategy != nil {
			// Calling the phase-out strategy
			v, ok := c.strategy.Evict()
			if ok {
				delete(c.indexMap, v)
			}
		} else {
			return errors.New("indexMap is full")
		}
	}
	entry := Entry{
		Key:        key,
		Value:      value,
		DataOffset: uint64(len(c.kvBuf)),
	}
	indexBuf, dataBuf, err := c.encoder.Marshal(&entry)
	if err != nil {
		return err
	}
	c.indexBuf = append(c.indexBuf, indexBuf...)
	c.kvBuf = append(c.kvBuf, dataBuf...)
	c.indexMap[hash] = uint64(len(c.indexBuf))
	return nil
}

func (c *Shard) Get(hash uint64, key string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	offset, ok := c.indexMap[hash]
	if !ok {
		return nil, ErrNotFound
	}
	if offset > uint64(len(c.indexBuf)) {
		return nil, ErrNotFound
	}
	indexBuf := c.indexBuf[offset : offset+itemIndexByteLen]
	var entry Entry
	err := c.encoder.UnmarshalIndex(indexBuf, &entry)
	if err != nil {
		return nil, err
	}
	dataBuf := c.kvBuf[entry.DataOffset : entry.DataOffset+uint64(entry.keyLen+entry.valueLen+2+2)]
	err = c.encoder.UnmarshalKV(dataBuf, &entry)
	if err != nil {
		return nil, err
	}
	if entry.Key != key {
		return nil, ErrHashCollision
	}
	return entry.Value, nil
}
