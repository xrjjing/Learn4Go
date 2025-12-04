# gRPC 示例

本目录包含一个简单的 gRPC 服务示例，演示了：

- **Unary RPC**: 一问一答模式（SayHello）
- **Server Streaming RPC**: 服务端流式响应（WatchTime）

## 目录结构

```
examples/grpc/
├── proto/
│   └── hello.proto      # 服务定义文件
├── server/
│   └── main.go          # 服务端实现
├── client/
│   └── main.go          # 客户端实现
├── Makefile             # 构建脚本
└── README.md            # 本文档
```

## 前置要求

### 1. 安装 protoc 编译器

**macOS:**
```bash
brew install protobuf
```

**Linux (Ubuntu/Debian):**
```bash
apt-get install -y protobuf-compiler
```

### 2. 安装 Go 插件

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

确保 `$GOPATH/bin` 在你的 PATH 中：
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 3. 安装 gRPC 依赖

在项目根目录执行：
```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

## 快速开始

### 1. 生成 protobuf 代码

```bash
cd examples/grpc
make proto
```

这会在 `proto/hellopb/` 目录下生成：
- `hello.pb.go`: 消息结构体
- `hello_grpc.pb.go`: gRPC 服务接口

### 2. 启动服务端

```bash
# 在项目根目录
go run ./examples/grpc/server
```

输出：
```
gRPC 服务器启动，监听 :50051
```

### 3. 运行客户端

打开新终端：
```bash
go run ./examples/grpc/client
```

### 4. 使用 grpcurl 测试（可选）

安装 grpcurl：
```bash
brew install grpcurl  # macOS
```

测试 Unary RPC：
```bash
grpcurl -plaintext -d '{"name":"Go学习者"}' localhost:50051 hello.Greeter/SayHello
```

测试 Streaming RPC：
```bash
grpcurl -plaintext localhost:50051 hello.Greeter/WatchTime
```

## RPC 类型说明

### Unary RPC (一元调用)

```
客户端 --请求--> 服务端
客户端 <--响应-- 服务端
```

最简单的 RPC 类型，类似普通函数调用。

### Server Streaming RPC (服务端流)

```
客户端 --请求--> 服务端
客户端 <--响应1-- 服务端
客户端 <--响应2-- 服务端
客户端 <--响应N-- 服务端
客户端 <--结束-- 服务端
```

服务端可以发送多个响应，适用于：
- 实时数据推送
- 大数据集分批返回
- 进度更新通知

## 常见问题

### Q: `protoc-gen-go: program not found`

确保 Go 插件已安装且在 PATH 中：
```bash
which protoc-gen-go
# 如果找不到，添加到 PATH：
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Q: 连接被拒绝

确保服务端已启动，且端口 50051 未被占用：
```bash
lsof -i :50051
```

### Q: import 路径错误

检查 `hello.proto` 中的 `go_package` 是否与你的 `go.mod` 模块名一致。
