// 数组、切片、Map 示例
// 对应章节: 05_数组_切片_map.md
package main

import "fmt"

func main() {
	// ========== 数组 ==========
	// 数组长度固定，是值类型
	var arr1 [3]int                    // 零值初始化
	arr2 := [3]int{1, 2, 3}            // 字面量初始化
	arr3 := [...]int{1, 2, 3, 4, 5}    // 编译器推断长度

	fmt.Println("数组:")
	fmt.Println("  arr1:", arr1)
	fmt.Println("  arr2:", arr2)
	fmt.Println("  arr3:", arr3, "长度:", len(arr3))

	// ========== 切片 ==========
	// 切片是动态数组，是引用类型
	slice1 := []int{1, 2, 3}           // 字面量创建
	slice2 := make([]int, 3, 5)        // make(类型, 长度, 容量)

	fmt.Println("\n切片:")
	fmt.Println("  slice1:", slice1)
	fmt.Println("  slice2:", slice2, "len:", len(slice2), "cap:", cap(slice2))

	// 切片操作
	nums := []int{0, 1, 2, 3, 4, 5}
	fmt.Println("  原切片:", nums)
	fmt.Println("  nums[1:4]:", nums[1:4])  // [1, 2, 3]
	fmt.Println("  nums[:3]:", nums[:3])    // [0, 1, 2]
	fmt.Println("  nums[3:]:", nums[3:])    // [3, 4, 5]

	// append 追加元素
	slice1 = append(slice1, 4, 5)
	fmt.Println("  append后:", slice1)

	// copy 复制切片
	dst := make([]int, 3)
	copy(dst, slice1)
	fmt.Println("  copy后:", dst)

	// ========== Map ==========
	// map 是键值对集合
	m1 := make(map[string]int)
	m1["apple"] = 100
	m1["banana"] = 200

	m2 := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	fmt.Println("\nMap:")
	fmt.Println("  m1:", m1)
	fmt.Println("  m2:", m2)

	// 访问与检查
	val, ok := m1["apple"]
	if ok {
		fmt.Println("  apple =", val)
	}

	val, ok = m1["orange"]
	if !ok {
		fmt.Println("  orange 不存在")
	}

	// 删除
	delete(m1, "apple")
	fmt.Println("  删除apple后:", m1)

	// 遍历 map（顺序不固定）
	fmt.Println("  遍历m2:")
	for k, v := range m2 {
		fmt.Printf("    %s: %d\n", k, v)
	}

	// ========== 切片与Map的陷阱 ==========
	fmt.Println("\n注意事项:")

	// 切片是引用：修改会影响原数据
	original := []int{1, 2, 3}
	ref := original
	ref[0] = 999
	fmt.Println("  切片引用:", original) // [999, 2, 3]

	// nil 切片 vs 空切片
	var nilSlice []int
	emptySlice := []int{}
	fmt.Println("  nil切片:", nilSlice == nil)     // true
	fmt.Println("  空切片:", emptySlice == nil)    // false
}
