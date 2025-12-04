// 包与模块管理示例
// 对应章节: 09_包_模块管理.md
package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// 本示例展示如何使用标准库中的包
// 实际项目中，你会创建自己的包并用 go.mod 管理依赖

func main() {
	fmt.Println("=== 标准库包使用示例 ===")

	// ========== fmt 包 ==========
	fmt.Println("--- fmt 包 ---")
	name := "Go"
	age := 14
	fmt.Printf("  %s 语言已经 %d 岁了\n", name, age)

	// Sprintf 返回格式化字符串
	s := fmt.Sprintf("版本: %s, 年龄: %d", name, age)
	fmt.Println(" ", s)

	// ========== strings 包 ==========
	fmt.Println("\n--- strings 包 ---")
	text := "Hello, Go World!"

	fmt.Println("  Contains:", strings.Contains(text, "Go"))
	fmt.Println("  HasPrefix:", strings.HasPrefix(text, "Hello"))
	fmt.Println("  HasSuffix:", strings.HasSuffix(text, "!"))
	fmt.Println("  ToUpper:", strings.ToUpper(text))
	fmt.Println("  Replace:", strings.Replace(text, "World", "语言", 1))
	fmt.Println("  Split:", strings.Split("a,b,c", ","))
	fmt.Println("  Join:", strings.Join([]string{"a", "b", "c"}, "-"))
	fmt.Println("  TrimSpace:", strings.TrimSpace("  hello  "))

	// ========== math 包 ==========
	fmt.Println("\n--- math 包 ---")
	fmt.Println("  Pi:", math.Pi)
	fmt.Println("  Sqrt(16):", math.Sqrt(16))
	fmt.Println("  Pow(2, 10):", math.Pow(2, 10))
	fmt.Println("  Max(3, 5):", math.Max(3, 5))
	fmt.Println("  Abs(-10):", math.Abs(-10))
	fmt.Println("  Ceil(3.2):", math.Ceil(3.2))
	fmt.Println("  Floor(3.8):", math.Floor(3.8))

	// ========== time 包 ==========
	fmt.Println("\n--- time 包 ---")
	now := time.Now()
	fmt.Println("  当前时间:", now)
	fmt.Println("  年:", now.Year())
	fmt.Println("  月:", now.Month())
	fmt.Println("  日:", now.Day())
	fmt.Println("  格式化:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("  Unix时间戳:", now.Unix())

	// 时间解析
	t, _ := time.Parse("2006-01-02", "2024-01-01")
	fmt.Println("  解析时间:", t)

	// 时间计算
	future := now.Add(24 * time.Hour)
	fmt.Println("  明天:", future.Format("2006-01-02"))

	// ========== 包的可见性规则 ==========
	fmt.Println("\n--- 可见性规则 ---")
	fmt.Println("  大写开头 = 导出（public）")
	fmt.Println("  小写开头 = 私有（private）")
	fmt.Println("  例: fmt.Println 是导出的")
	fmt.Println("  例: strings 包内部的 makeASCIISet 是私有的")

	// ========== init 函数 ==========
	fmt.Println("\n--- init 函数 ---")
	fmt.Println("  每个包可以有 init() 函数")
	fmt.Println("  init() 在 main() 之前自动执行")
	fmt.Println("  常用于初始化配置、注册驱动等")
}
