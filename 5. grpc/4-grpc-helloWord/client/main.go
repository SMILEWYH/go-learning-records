package main

import (
	"context"
	"fmt"
	"grpc_study/grpc_proto/hello_grpc"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 使用 grpc.Dial 创建一个到指定地址的 gRPC 连接, 此处使用不安全的证书来实现 SSL/TLS 连接
	// conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", ":8080", err))
	}
	defer conn.Close()

	// 初始化客户端
	client := hello_grpc.NewHelloServiceClient(conn)
	result, err := client.SayHello(context.Background(), &hello_grpc.HelloRequest{
		Name:    "枫枫",
		Message: "ok",
	})
	fmt.Println(result, err)
}

// 如果需要经常操作这个连接对象，我们可以复用
/*
func initClient() hello_grpc.HelloServiceClient {
	addr := ":8080"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	client := hello_grpc.NewHelloServiceClient(conn)
	return client
}

func main() {
	client := initClient()
	result, err := client.SayHello(context.Background(), &hello_grpc.HelloRequest{
		Name:    "枫枫",
		Message: "ok",
	})
	result, err = client.SayHello(context.Background(), &hello_grpc.HelloRequest{
		Name:    "哈哈哈",
		Message: "ok1",
	})

	fmt.Println(result, err)
}
*/
