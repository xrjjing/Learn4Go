// gRPC 客户端示例
//
// 本示例展示如何调用 gRPC 服务，包含：
//   - Unary RPC 调用: SayHello
//   - Server Streaming RPC 调用: WatchTime
//
// 运行前需要：
//   1. 生成 protobuf 代码: cd examples/grpc && make proto
//   2. 启动服务器: go run ./examples/grpc/server
//
// 运行客户端：
//   go run ./examples/grpc/client
package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/xrjjing/Learn4Go/examples/grpc/proto/hellopb"
)

func main() {
	// ============ 建立连接 ============
	//
	// gRPC 连接是长连接，底层使用 HTTP/2
	// 一个连接可以复用多个请求（多路复用）

	// grpc.WithTransportCredentials 指定传输层安全设置
	// insecure.NewCredentials() 表示不使用 TLS（仅用于开发环境）
	// 生产环境应该使用 credentials.NewClientTLSFromFile() 或 credentials.NewTLS()
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close() // 程序退出时关闭连接

	// 创建客户端存根（stub）
	// 存根封装了底层的网络调用，让 RPC 调用看起来像本地函数调用
	client := hellopb.NewGreeterClient(conn)

	// ============ Unary RPC 调用 ============
	log.Println("=== 调用 SayHello (Unary RPC) ===")
	callSayHello(client)

	// ============ Server Streaming RPC 调用 ============
	log.Println("\n=== 调用 WatchTime (Server Streaming RPC) ===")
	callWatchTime(client)
}

// callSayHello 演示 Unary RPC 调用
func callSayHello(client hellopb.GreeterClient) {
	// 创建带超时的上下文
	// 如果服务端 3 秒内没有响应，调用会自动取消
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // 确保资源释放

	// 构造请求
	req := &hellopb.HelloRequest{Name: "Go 学习者"}

	// 发起 RPC 调用
	// 这看起来像本地函数调用，但实际上会：
	//   1. 将请求序列化为 protobuf 二进制格式
	//   2. 通过网络发送到服务端
	//   3. 等待服务端响应
	//   4. 将响应反序列化为 Go 结构体
	resp, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("SayHello 调用失败: %v", err)
	}

	log.Printf("服务端响应: %s", resp.GetMessage())
}

// callWatchTime 演示 Server Streaming RPC 调用
func callWatchTime(client hellopb.GreeterClient) {
	// 创建带超时的上下文
	// 7 秒超时，因为服务端会推送 5 次（每秒一次）
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	// 发起流式 RPC 调用
	// 返回的是一个流对象，而不是单个响应
	stream, err := client.WatchTime(ctx, &hellopb.Empty{})
	if err != nil {
		log.Fatalf("WatchTime 调用失败: %v", err)
	}

	// 循环接收服务端推送的数据
	for {
		// Recv() 从流中读取一条消息
		// 如果没有数据，会阻塞等待
		msg, err := stream.Recv()
		if err == io.EOF {
			// io.EOF 表示服务端正常结束了流
			log.Println("服务端完成推送")
			break
		}
		if err != nil {
			// 其他错误（网络问题、服务端错误、超时等）
			log.Fatalf("接收数据失败: %v", err)
		}

		// 处理收到的数据
		log.Printf("收到时间: %s", msg.GetTime())
	}
}
