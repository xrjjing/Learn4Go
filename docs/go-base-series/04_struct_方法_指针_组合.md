# 04 struct、方法、指针、组合

## 🎯 这一章解决什么问题

如果你来自 Java，这一章最关键。

因为这里要完成一次真正的心智切换：

- Go 没有 class
- Go 没有继承
- Go 用 `struct + 方法 + 组合` 组织业务对象

你要开始适应 Go 的“数据和行为分离定义”。

## 🧩 最小代码

```go
package main

import "fmt"

type User struct {
Name string
Age  int
}

func (u User) SayHello() {
fmt.Println("hello, I am", u.Name)
}

func (u *User) Grow() {
u.Age++
}

func main() {
u := User{Name: "Tom", Age: 18}
u.SayHello()
u.Grow()
fmt.Println(u.Age)
}
```

## 1️⃣ `struct` 是什么

`struct` 是结构体，用来组织数据字段。

```go
type User struct {
Name string
Age  int
}
```

你可以把它理解成 Java 类里“只保留字段定义的那一部分”。

但 Go 不把方法写进结构体内部，而是定义在外面。

## 2️⃣ 方法为什么写在外面

```go
func (u User) SayHello() {
fmt.Println("hello, I am", u.Name)
}
```

`(u User)` 叫 **接收者**。

你可以把它理解成 Java 的 `this`，但它是显式写出来的。

- `User` 是接收者类型
- `u` 是这个接收者变量名

所以 Go 里方法长得像“带接收者的函数”。

## 3️⃣ 值接收者 vs 指针接收者

### 值接收者

```go
func (u User) SayHello() {}
```

表示拿到的是一个副本。

一般适合：

- 不修改对象状态
- 结构体较小
- 语义偏只读

### 指针接收者

```go
func (u *User) Grow() {
u.Age++
}
```

表示要改原对象。

一般适合：

- 需要修改字段
- 结构体较大，不想复制
- 结构体里含有锁等不能安全复制的成员

## 4️⃣ 组合与嵌入

Go 没有继承，但支持组合。

```go
type Address struct {
City string
}

type User struct {
Name string
Address
}
```

这里 `Address` 被匿名嵌入后，`User` 可以直接访问：

```go
u.City
```

这叫方法/字段被“提升”。

但要注意：

**嵌入不是继承。**

`User` 不会因为嵌入了 `Address`，就变成 `Address` 的子类型。

## 🆚 Java 开发怎么理解

### Java

```java
class User {
private String name;
private int age;
void grow() { age++; }
}
```

### Go

```go
type User struct {
Name string
Age  int
}

func (u *User) Grow() {
u.Age++
}
```

可以先这样理解：

- Java：类把字段和方法包在一起
- Go：结构体存数据，方法挂在结构体外部

## ⚠️ 注意点 / 易错点

### 1. 方法不能写在 struct 里面

Java 开发最常犯的语法迁移错误之一。

Go 不支持：

```go
type User struct {
func Grow() {} // ❌
}
```

### 2. 修改对象状态时要用指针接收者

```go
func (u User) Grow() {
u.Age++
}
```

这只会改副本，不会改原对象。

### 3. 一种类型的方法接收者风格尽量统一

如果一个结构体大多数方法都要改状态，通常统一用指针接收者更稳。

### 4. 结构体字段大小写决定可见性

```go
type User struct {
name string // 包内可见
Age  int    // 导出
}
```

Go 没有 `public/private` 关键字，而是靠首字母大小写。

### 5. 嵌入不是多态继承

这点一定要反复提醒自己，否则后面接口那章会混。

## ▶️ 本章建议运行命令

```bash
go run ./examples/structs
go run ./examples/basics/pointer_receiver
```

## 📌 本章小结

你要记住：

1. Go 用 `struct` 表达数据对象
2. 方法写在结构体外，通过接收者绑定
3. 要改状态时通常用指针接收者
4. Go 没有继承，核心是组合和嵌入
5. 大小写决定导出与否

## ⏭️ 下一章

下一章进入 Go 风格最强烈的部分：

- `interface`
- `error`
- `defer`
- `panic / recover`
