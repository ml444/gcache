package gcache

import (
	"hash/fnv"
	"testing"
)

func TestShard(t *testing.T) {
	//for _, s := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
	//	t.Log(s, ":", int(s[0])%4)
	//}
	m := map[int]int{}
	for i := 0; i < 100; i++ {
		//s := strconv.Itoa(i)
		a := i % 9
		t.Log(i, ":", a)
		m[a] = m[a] + 1
		//t.Log(i, ":", s[0]%4)
	}
	t.Log(m)

}

func hashString(input string) uint32 {
	// 创建 FNV-1a 哈希对象
	hasher := fnv.New32a()

	// 写入要计算哈希值的数据
	hasher.Write([]byte(input))

	// 获取哈希值并转换为整数
	hashValue := hasher.Sum32()

	return hashValue
}
