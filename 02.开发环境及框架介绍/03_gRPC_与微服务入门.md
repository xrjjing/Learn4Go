# gRPC 与微服务入门

## 核心概念
- IDL：Protocol Buffers（.proto）
- 四种调用：Unary / Server streaming / Client streaming / Bidirectional
- Go 生成：`protoc --go_out=. --go-grpc_out=. service.proto`

## 最小示例（伪代码）
```proto
service Greeter {
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}
```

```go
// 服务器注册
grpcServer := grpc.NewServer()
pb.RegisterGreeterServer(grpcServer, &Greeter{})
```

## 演进
1. 在本仓库先用 HTTP+JSON 完成链路
2. 拓展为 gRPC 服务，并通过 grpc-gateway 暴露 HTTP（需网络装依赖）

## 练习
- 设计一个 `User` 服务 proto（GetUser/CreateUser）
- 对比 gRPC 与 REST 在类型安全、性能、前后端契约上的差异
