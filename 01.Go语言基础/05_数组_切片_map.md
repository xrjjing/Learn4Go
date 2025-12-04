# 05 数组、切片、map

## 数组 vs 切片
- 数组长度固定，值语义；切片为动态视图，包含指针/长度/容量
- 切片扩容：`append`，可能触发底层数组重分配

## map
- 无序，读取不存在键返回零值；使用 `v, ok := m[k]` 判断

## 示例
```go
arr := [3]int{1,2,3}
s := []int{1,2,3}
s = append(s, 4)

m := map[string]int{"a":1}
m["b"] = 2
if v, ok := m["c"]; !ok {
    fmt.Println("not found")
}
```

## 练习
- 实现切片去重（返回新切片）
- 统计一段文本的词频，使用 map + strings.Fields
