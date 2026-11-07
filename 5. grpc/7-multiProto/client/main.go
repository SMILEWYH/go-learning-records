package main

import (
	"context"
	"fmt"
	mall "grpc-learning/7-multiProto/proto/proto_go"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 本地学习使用明文连接，不启用 TLS
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("创建 gRPC 客户端失败：%v", err)
	}
	defer conn.Close()

	// 两个服务共用同一个连接，客户端和消息类型都来自 mall 包。
	productClient := mall.NewProductServiceClient(conn)
	productCtx, cancelProduct := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelProduct()
	productRes, err := productClient.Find(productCtx, &mall.Request{
		ID: 1,
	})
	if err != nil {
		log.Printf("查询商品失败：%v", err)
		return
	}
	fmt.Println("商品查询结果：", productRes)

	orderClient := mall.NewOrderServiceClient(conn)
	orderCtx, cancelOrder := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelOrder()
	orderRes, err := orderClient.Find(orderCtx, &mall.Request{
		ID: 1,
	})
	if err != nil {
		log.Printf("查询订单失败：%v", err)
		return
	}
	fmt.Println("订单查询结果：", orderRes)
}
