// 流程控制示例
// 对应章节: 03_流程控制.md
package main

import "fmt"

func main() {
	// if-else
	score := 85
	if score >= 90 {
		fmt.Println("优秀")
	} else if score >= 60 {
		fmt.Println("及格")
	} else {
		fmt.Println("不及格")
	}

	// if 带初始化语句
	if n := 10; n%2 == 0 {
		fmt.Println(n, "是偶数")
	}

	// switch（无需 break，自动终止）
	day := 3
	switch day {
	case 1:
		fmt.Println("周一")
	case 2:
		fmt.Println("周二")
	case 3:
		fmt.Println("周三")
	case 4, 5:
		fmt.Println("周四或周五")
	default:
		fmt.Println("周末")
	}

	// switch 无条件（替代 if-else 链）
	hour := 14
	switch {
	case hour < 12:
		fmt.Println("上午")
	case hour < 18:
		fmt.Println("下午")
	default:
		fmt.Println("晚上")
	}

	// for 循环（Go 只有 for，没有 while）
	// 基本形式
	for i := 0; i < 3; i++ {
		fmt.Println("循环:", i)
	}

	// 类似 while
	j := 0
	for j < 3 {
		fmt.Println("while形式:", j)
		j++
	}

	// 无限循环
	// for { ... }

	// range 遍历
	nums := []int{10, 20, 30}
	for idx, val := range nums {
		fmt.Printf("index=%d, value=%d\n", idx, val)
	}

	// break 和 continue
	for i := 0; i < 5; i++ {
		if i == 2 {
			continue // 跳过 2
		}
		if i == 4 {
			break // 到 4 终止
		}
		fmt.Println("break/continue:", i)
	}
}
