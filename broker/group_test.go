package broker

import "testing"

var shards []*Shard

func lookupShard(no int) *Shard {
	if len(shards) == 0 {
		for i := 0; i < 1000; i++ {
			shards = append(shards, &Shard{SerialNo: i + 1})
		}
	}
	for _, s := range shards {
		if s.SerialNo == no {
			return s
		}
	}
	return nil
}

var shardMap = make(map[int]*Shard)

func lookupShard2(no int) *Shard {
	if len(shardMap) == 0 {
		for i := 0; i < 1000; i++ {
			shardMap[i] = &Shard{SerialNo: i}
		}
	}
	return shardMap[no]
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupShard2(i % 9)
	}
}

func BenchmarkLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lookupShard(i % 9)
	}
}
