package main

import (
	"context"
	"fmt"
	hello_grpc "grpc-learning/4-grpc-helloWord"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type HelloServer struct {
}

func (HelloServer) SayHello(ctx context.Context, request *hello_grpc.HelloRequest) (res *hello_grpc.HelloResponse, err error) {
	fmt.Println("入参：", request.Name, request.Message)
	res = new(hello_grpc.HelloResponse)
	res.Name = "你好"
	res.Message = "ok"
	return
}

func main() {
	// 监听端口
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// 创建一个gRPC服务器实例
	s := grpc.NewServer()

	// 将server结构体注册为gRPC服务
	server := HelloServer{}
	hello_grpc.RegisterHelloServiceServer(s, &server)

	// 开始处理客户端请求。
	err = s.Serve(listen)
}
