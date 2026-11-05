package main

import (
	"context"
	"fmt"
	"grpc_study/stream_proto/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := ":8080"
	// 使用 grpc.Dial 创建一个到指定地址的 gRPC 连接。
	// 此处使用不安全的证书来实现 SSL/TLS 连接
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	defer conn.Close()

	// 初始化客户端
	client := proto.NewServiceStreamClient(conn)
	stream, err := client.Fun(context.Background(), &proto.Request{
		Name: "张三",
	})
	for i := 0; i < 10; i++ {
		response, err := stream.Recv()
		fmt.Println(response, err)
	}

	// 客户端不知道服务端什么时候结束
	// for {
	// 	response, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	fmt.Println(response)
	// }

}
