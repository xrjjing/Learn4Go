// gRPC 服务端示例
//
// 本示例展示如何实现一个简单的 gRPC 服务器，包含：
//   - Unary RPC: SayHello（一问一答）
//   - Server Streaming RPC: WatchTime（服务端持续推送）
//
// 运行前需要先生成 protobuf 代码：
//
//	cd examples/grpc && make proto
//
// 启动服务器：
//
//	go run ./examples/grpc/server
//
// 测试（使用 grpcurl）：
//
//	grpcurl -plaintext localhost:50051 hello.Greeter/SayHello
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/xrjjing/Learn4Go/examples/grpc/proto/hellopb"
)

// greeterServer 实现了 hellopb.GreeterServer 接口
// 在 Go 中，只要实现了接口定义的所有方法，就自动实现了该接口
type greeterServer struct {
	// 嵌入 UnimplementedGreeterServer 是 gRPC 的推荐做法
	// 这样即使 proto 文件新增方法，编译也不会出错（向前兼容）
	hellopb.UnimplementedGreeterServer
}

// SayHello 实现 Unary RPC
//
// Unary RPC 是最简单的 RPC 类型：
//   - 客户端发送一个请求
//   - 服务端处理后返回一个响应
//   - 类似于普通的函数调用
//
// 参数说明：
//   - ctx: 上下文，用于传递截止时间、取消信号等
//   - req: 客户端发来的请求消息
//
// 返回值：
//   - *hellopb.HelloReply: 响应消息
//   - error: 错误信息，nil 表示成功
func (s *greeterServer) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloReply, error) {
	// 检查上下文是否已取消（比如客户端超时断开）
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// 构造响应消息
	name := req.GetName()
	if name == "" {
		name = "匿名用户"
	}
	msg := fmt.Sprintf("你好，%s！欢迎学习 gRPC", name)

	log.Printf("收到 SayHello 请求: name=%s", name)
	return &hellopb.HelloReply{Message: msg}, nil
}

// WatchTime 实现 Server Streaming RPC
//
// Server Streaming RPC 特点：
//   - 客户端发送一个请求
//   - 服务端返回多个响应（流式）
//   - 适用于：实时数据推送、日志流、进度更新等场景
//
// 参数说明：
//   - req: 客户端请求（本例未使用）
//   - stream: 流对象，用于向客户端发送响应
//
// 流式 RPC 的生命周期：
//  1. 客户端发起调用
//  2. 服务端循环调用 stream.Send() 发送数据
//  3. 服务端返回 nil 表示流结束
//  4. 服务端返回 error 表示异常终止
func (s *greeterServer) WatchTime(_ *hellopb.Empty, stream hellopb.Greeter_WatchTimeServer) error {
	log.Println("收到 WatchTime 请求，开始推送时间...")

	// 创建定时器，每秒触发一次
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop() // 确保退出时停止定时器，释放资源

	// 推送 5 次时间数据
	for i := 0; i < 5; i++ {
		select {
		case t := <-ticker.C:
			// 定时器触发，发送当前时间
			resp := &hellopb.TimeReply{
				Time: t.Format(time.RFC3339), // RFC3339 是标准时间格式
			}
			if err := stream.Send(resp); err != nil {
				// 发送失败（比如客户端断开连接）
				log.Printf("发送失败: %v", err)
				return err
			}
			log.Printf("推送时间: %s", resp.Time)

		case <-stream.Context().Done():
			// 客户端取消了请求（超时或主动取消）
			log.Println("客户端取消请求")
			return stream.Context().Err()
		}
	}

	log.Println("时间推送完成")
	return nil // 正常结束流
}

func main() {
	// 创建 TCP 监听器
	// ":50051" 表示监听所有网卡的 50051 端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	// 创建 gRPC 服务器实例
	// 可以通过选项配置拦截器、TLS 等
	srv := grpc.NewServer()

	// 注册服务实现
	// 这告诉 gRPC 服务器如何处理 Greeter 服务的请求
	hellopb.RegisterGreeterServer(srv, &greeterServer{})

	// 注册反射服务（可选，但推荐）
	// 反射允许 grpcurl、evans 等工具在不知道 proto 文件的情况下调用服务
	reflection.Register(srv)

	log.Println("gRPC 服务器启动，监听 :50051")
	log.Println("使用 grpcurl 测试:")
	log.Println("  grpcurl -plaintext -d '{\"name\":\"Go学习者\"}' localhost:50051 hello.Greeter/SayHello")
	log.Println("  grpcurl -plaintext localhost:50051 hello.Greeter/WatchTime")

	// 开始服务（阻塞直到服务停止）
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("服务器退出: %v", err)
	}
}
