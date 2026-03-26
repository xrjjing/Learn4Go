# 03 数组、切片、map

## 🎯 这一章解决什么问题

这章是 Go 基础里最实用的一章之一。

目标是把这 3 个东西彻底区分开：

- 数组：长度固定
- 切片：最常用，动态视图
- map：键值对集合

你只要把切片和 map 这章吃透，后面看绝大多数 Go 业务代码都会顺很多。

## 🧩 最小代码

```go
package main

import "fmt"

func main() {
arr := [3]int{1, 2, 3}
s := []int{1, 2, 3}
s = append(s, 4)

m := map[string]int{"go": 1}
m["java"] = 2

fmt.Println(arr)
fmt.Println(s)
fmt.Println(m)
}
```

## 1️⃣ 数组：长度是类型的一部分

```go
arr := [3]int{1, 2, 3}
```

这里的 `[3]int` 不是“3 个 int 的值”，而是一个完整类型。

所以：

- `[3]int` 和 `[4]int` 是不同类型
- 数组长度固定
- 数组在 Go 里不是最常用的数据结构

Java 的数组更常见，而 Go 业务代码里更常见的是切片。

## 2️⃣ 切片：最常用的数据结构

```go
s := []int{1, 2, 3}
```

切片可以先理解成：

- 一个“对底层数组的视图”
- 自带长度和容量信息
- 支持 `append`

所以 Go 里大多数“列表”场景都用切片，不直接用数组。

### 为什么切片比数组常用

因为它更灵活：

- 长度可以变
- 可以切片操作：`s[1:3]`
- 可以 `append`
- 可以作为函数参数更自然

### `append` 是什么

```go
s = append(s, 4)
```

这个操作可能发生两种事：

- 底层数组容量够，直接往后放
- 底层数组容量不够，重新分配新数组

这就是为什么切片有时像“修改原对象”，有时又像“变成了新对象”。

## 3️⃣ map：键值对集合

```go
m := map[string]int{"go": 1}
```

可以理解成 Java 的：

```java
Map<String, Integer>
```

但 Go 的 map 有几个非常高频的坑。

### 读取不存在的 key

```go
v := m["python"]
```

如果 key 不存在，不会报错，而是返回零值。

所以更稳的写法是：

```go
v, ok := m["python"]
if ok {
fmt.Println(v)
}
```

### 判断是否存在

`ok` 是 Go 非常典型的用法。

- `ok == true`：存在
- `ok == false`：不存在

## 🆚 Java 开发怎么理解

### 数组

和 Java 数组有相似点，但 Go 的数组更“严格类型化”。

### 切片

不能简单类比成 `ArrayList`。

更准确地说，切片像：

- 一段窗口
- 指向底层数组
- 带 length / capacity

### map

可以类比成 `HashMap`，但使用方式更轻：

- 不需要 new 出具体实现类名
- 不需要 `get/put`
- 直接 `m[key]`

## ⚠️ 注意点 / 易错点

### 1. map 零值可读不可写

```go
var m map[string]int
fmt.Println(m["x"]) // ✅ 读，返回 0
m["x"] = 1          // ❌ panic
```

因为零值 map 还没初始化。

正确写法：

```go
m := make(map[string]int)
```

### 2. 切片可能共享底层数组

```go
a := []int{1, 2, 3, 4}
b := a[1:3]
b[0] = 99
fmt.Println(a) // [1 99 3 4]
```

这是切片最容易让 Java 开发误判的点之一。

### 3. `append` 后别忘了接回去

```go
s := []int{1, 2}
append(s, 3) // ❌ 结果被丢掉
```

应该写：

```go
s = append(s, 3)
```

### 4. 遍历 map 是无序的

```go
for k, v := range m {
fmt.Println(k, v)
}
```

输出顺序不要依赖。

### 5. `len` 和 `cap` 不一样

- `len(s)`：当前长度
- `cap(s)`：底层容量

刚开始先知道“容量影响 append 是否触发扩容”就够了。

## ▶️ 本章建议运行命令

```bash
go run ./examples/collections
go run ./examples/basics/map_basics
go run ./examples/basics/slice_aliasing
```

## 📌 本章小结

你要记住：

1. 数组长度固定，切片最常用
2. 切片是对底层数组的视图
3. `append` 可能触发扩容
4. map 读不存在 key 会拿到零值
5. 零值 map 可读不可写

## ⏭️ 下一章

下一章进入 Go 的“类替代物”：

- `struct`
- 方法
- 指针接收者
- 组合与嵌入
