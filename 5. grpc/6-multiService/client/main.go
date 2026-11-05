package main

import (
	"context"
	"fmt"
	"grpc_study/grpc_proto/duo_proto"
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

	orderClient := duo_proto.NewOrderServiceClient(conn)
	res, err := orderClient.Buy(context.Background(), &duo_proto.Request{
		Name: "枫枫",
	})
	fmt.Println(res, err)

	videoClient := duo_proto.NewVideoServiceClient(conn)
	res, err = videoClient.Look(context.Background(), &duo_proto.Request{
		Name: "枫枫",
	})
	fmt.Println(res, err)

}
