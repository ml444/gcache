package hash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type ConsistentRing struct {
	sync.RWMutex
	defaultVirtualCount int
	sortedVirtualNodes  []uint32

	virtualNodeMap    map[uint32]string //The real node information corresponding to the virtual node
	realNodeMap       map[string]bool
	hashFunc          func(key string) uint32
	customVirtualKeys func(key string) []uint32
}

func NewConsistent(opts ...OptionFunc) *ConsistentRing {
	ring := &ConsistentRing{}
	for _, opt := range opts {
		opt(ring)
	}
	if ring.hashFunc == nil {
		ring.hashFunc = func(key string) uint32 {
			return crc32.ChecksumIEEE([]byte(key))
		}
	}
	return ring
}

func (c *ConsistentRing) hashKey(key string) uint32 {
	return c.hashFunc(key)
}

func (c *ConsistentRing) Add(node string) error {
	return c.AddWithVirtualCount(node, c.defaultVirtualCount)
}

func (c *ConsistentRing) AddWithVirtualCount(node string, virtualCount int) error {
	if node == "" {
		return nil
	}

	if c.virtualNodeMap == nil {
		c.virtualNodeMap = map[uint32]string{}
	}
	if c.realNodeMap == nil {
		c.realNodeMap = map[string]bool{}
	}

	if ok := c.realNodeMap[node]; ok {
		return errors.New("node already existed: " + node)
	}

	c.Lock()
	defer c.Unlock()

	c.realNodeMap[node] = true
	//add virtual node
	if c.customVirtualKeys != nil {
		virtualKeys := c.customVirtualKeys(node)
		for _, virtualKey := range virtualKeys {
			c.virtualNodeMap[virtualKey] = node
			c.sortedVirtualNodes = append(c.sortedVirtualNodes, virtualKey)
		}
		return nil
	} else {
		for i := 0; i < virtualCount; i++ {
			virtualKey := c.hashKey(node + strconv.Itoa(i))
			c.virtualNodeMap[virtualKey] = node
			c.sortedVirtualNodes = append(c.sortedVirtualNodes, virtualKey)
		}
	}

	//虚拟结点排序
	sort.Slice(c.sortedVirtualNodes, func(i, j int) bool {
		return c.sortedVirtualNodes[i] < c.sortedVirtualNodes[j]
	})

	return nil
}

func (c *ConsistentRing) Remove(node string) error {
	if node == "" {
		return nil
	}

	if !c.realNodeMap[node] {
		return errors.New("node not existed: " + node)
	}

	c.Lock()
	defer c.Unlock()

	delete(c.realNodeMap, node)
	var removeMap = map[uint32]bool{}
	//remove virtual node
	if c.customVirtualKeys != nil {
		virtualKeys := c.customVirtualKeys(node)
		for _, virtualKey := range virtualKeys {
			removeMap[virtualKey] = true
			delete(c.virtualNodeMap, virtualKey)
		}
	} else {
		for i := 0; i < c.defaultVirtualCount; i++ {
			virtualKey := c.hashKey(node + strconv.Itoa(i))
			removeMap[virtualKey] = true
			delete(c.virtualNodeMap, virtualKey)
		}
	}

	var newSortedVirtualNodes []uint32
	for _, virtualKey := range c.sortedVirtualNodes {
		if removeMap[virtualKey] {
			continue
		}
		newSortedVirtualNodes = append(newSortedVirtualNodes, virtualKey)
	}
	c.sortedVirtualNodes = newSortedVirtualNodes
	return nil
}

func (c *ConsistentRing) Get(key string) string {
	hash := c.hashKey(key)

	c.RLock()
	defer c.RUnlock()

	i := c.getPosition(hash)
	return c.virtualNodeMap[c.sortedVirtualNodes[i]]
}

func (c *ConsistentRing) getPosition(hash uint32) int {
	i := sort.Search(len(c.sortedVirtualNodes), func(i int) bool { return c.sortedVirtualNodes[i] >= hash })
	if i < len(c.sortedVirtualNodes) {
		return i
	} else {
		return 0
	}
}
