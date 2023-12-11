//go:build go1.9

package strategy

import (
	"container/list"
	"sync"
)

type LRU struct {
	size      int
	innerList *list.List
	innerMap  sync.Map
}

type entry struct {
	key   interface{}
	value interface{}
}

func NewLRU(size int) *LRU {
	return &LRU{
		size:      size,
		innerList: list.New(),
		innerMap:  sync.Map{},
	}
}

func (lru *LRU) loadElement(key interface{}) (*list.Element, bool) {
	v, ok := lru.innerMap.Load(key)
	if !ok {
		return nil, false
	}
	el, ok := v.(*list.Element)
	if !ok {
		return nil, false
	}
	return el, true
}

func (lru *LRU) storeElement(key, value interface{}) {
	lru.innerMap.Store(key, value)
}

func (lru *LRU) deleteElement(key interface{}) {
	lru.innerMap.Delete(key)
}

func (lru *LRU) Get(key interface{}) (interface{}, bool) {
	if e, ok := lru.loadElement(key); ok {
		lru.innerList.MoveToFront(e)
		return e.Value.(*entry).value, true
	}
	return nil, false
}

func (lru *LRU) Put(key interface{}, value interface{}) (evicted bool) {
	if el, ok := lru.loadElement(key); ok {
		lru.innerList.MoveToFront(el)
		el.Value.(*entry).value = value
		return false
	} else {
		e := &entry{key, value}
		el = lru.innerList.PushFront(e)
		lru.storeElement(key, el)

		if lru.innerList.Len() > lru.size {
			last := lru.innerList.Back()
			lru.innerList.Remove(last)
			lru.deleteElement(last.Value.(*entry).key)
			return true
		}
		return false
	}
}
