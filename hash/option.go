package hash

type OptionFunc func(*ConsistentRing)

func WithHashFunc(hashFunc func(key string) uint32) OptionFunc {
	return func(c *ConsistentRing) {
		c.hashFunc = hashFunc
	}
}

func WithVirtualCount(virtualCount int) OptionFunc {
	return func(c *ConsistentRing) {
		c.defaultVirtualCount = virtualCount
	}
}

func WithCustomVirtualKeysFunc(customVirtualKeysFunc func(key string) []uint32) OptionFunc {
	return func(c *ConsistentRing) {
		c.customVirtualKeys = customVirtualKeysFunc
	}
}
