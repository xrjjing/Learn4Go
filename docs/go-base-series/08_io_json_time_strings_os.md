# 08 常用标准库：io、json、time、strings、os

## 🎯 这一章解决什么问题

到这里开始，你要从“只会语法”进入“能写小程序”。

这章聚焦最常用的标准库能力：

- 时间处理
- 字符串处理
- 文件与路径
- IO 读写
- JSON 编解码

## 1️⃣ `time`

最常见用法：

```go
now := time.Now()
fmt.Println(now)
```

构造一个时间：

```go
t := time.Date(2026, time.March, 24, 10, 30, 0, 0, time.Local)
```

这里最后一个参数是 `*time.Location`，表示时区。

常见值：

- `time.Local`
- `time.UTC`
- `time.LoadLocation("Asia/Shanghai")`

### 你之前问过的点

`time.Date(..., loc *time.Location)` 里最后一个参数不是随便写整数，而是传时区对象。

## 2️⃣ `strings`

常见函数：

```go
strings.Contains(s, "go")
strings.Split(s, ",")
strings.Join(parts, "-")
strings.TrimSpace(s)
```

这些会在文本处理和配置处理里高频出现。

## 3️⃣ `os`

`os` 负责和操作系统打交道。

常见用法：

```go
os.ReadFile("a.txt")
os.WriteFile("a.txt", data, 0644)
os.Getenv("HOME")
```

## 4️⃣ `io`

`io` 更像一层抽象能力。

你会经常看到：

- `io.Reader`
- `io.Writer`
- `io.Copy`

Go 很多库都围绕 Reader / Writer 设计，这是一条非常重要的主线。

## 5️⃣ `encoding/json`

这是 Go 里最常用的包之一。

序列化：

```go
data, err := json.Marshal(v)
```

反序列化：

```go
err := json.Unmarshal(data, &v)
```

结构体字段通常要写 json tag：

```go
type User struct {
Name string `json:"name"`
Age  int    `json:"age"`
}
```

## 🧩 一个最小示例

```go
package main

import (
"encoding/json"
"fmt"
"time"
)

type User struct {
Name string `json:"name"`
Age  int    `json:"age"`
}

func main() {
u := User{Name: "Tom", Age: 18}
data, err := json.Marshal(u)
if err != nil {
    fmt.Println("marshal failed:", err)
    return
}

fmt.Println(string(data))
fmt.Println(time.Now().Format(time.RFC3339))
}
```

## 🆚 Java 开发怎么理解

- `time` 类似 `java.time`，但 API 风格不同
- `strings` 类似 `String` / `StringUtils` 的常用子集
- `os` 类似 `Files`、环境变量和进程相关 API 的集合
- `io.Reader/Writer` 类似 Java IO 流模型，但 Go 更偏接口化
- `encoding/json` 类似 Jackson 的最基础使用层

## ⚠️ 注意点 / 易错点

### 1. `time.Date` 最后一个参数必须是时区对象

这就是你之前问到的高频点。

### 2. JSON 反序列化要传指针

```go
json.Unmarshal(data, &u)
```

不是：

```go
json.Unmarshal(data, u) // ❌
```

### 3. 结构体字段小写时，json 包通常看不到

```go
type User struct {
name string `json:"name"`
}
```

这通常不会按你期望工作，因为字段未导出。

### 4. `os.ReadFile` 很方便，但大文件别无脑一次性全读

小文件练习没问题，大文件再考虑流式读取。

### 5. `strings.Split` 结果可能包含空串

文本清洗时记得处理边界。

## ▶️ 本章建议运行命令

```bash
go run ./examples/java_compare/bufio_wordcount
go run ./examples/java_compare/urlencode
go run ./examples/basics/rune_utf8
```

如果你想继续扩展，可以参考：

- `docs/IO操作指南.md`
- `examples/io/*`

## 📌 本章小结

你要记住：

1. `time` 负责时间、时区、格式化
2. `strings` 是文本处理高频工具包
3. `os` 负责文件、环境变量、系统交互
4. `io.Reader/Writer` 是 Go IO 体系核心抽象
5. `encoding/json` 是写 API 时的核心标准库

## ⏭️ 下一章

下一章开始收工程化尾：

- `go test`
- `go build`
- `go fmt`
- `go vet`
- 什么时候用 `go run`
