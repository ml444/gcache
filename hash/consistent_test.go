package hash

import (
	"fmt"
	"sort"
	"testing"
)

func TestSortSearch(t *testing.T) {
	// 已排序的切片
	numbers := []int{1, 3, 5, 7, 9, 11, 13, 15}

	// 要搜索的元素
	target := 0

	// 使用 sort.Search 查找元素的索引
	index := sort.Search(len(numbers), func(i int) bool {
		fmt.Printf("i: %d, numbers[i]: %d\n", i, numbers[i])
		return numbers[i] >= target
	})

	// 检查结果
	if index < len(numbers) && numbers[index] == target {
		fmt.Printf("找到了 %d，索引是 %d\n", target, index)
	} else {
		fmt.Printf("%d 不存在于切片中 %d \n", target, index)
	}
}
