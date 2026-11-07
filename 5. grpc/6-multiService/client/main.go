package main

import (
	"context"
	"fmt"
	online_mall "grpc-learning/6-multiService/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc connect addr [%s] 连接失败 %s", "127.0.0.1:8080", err)
	}
	defer conn.Close()

	productClient := online_mall.NewProductServiceClient(conn)
	res1, err1 := productClient.Find(context.Background(), &online_mall.Request{
		ID: 1,
	})
	fmt.Println(res1, err1)

	orderClient := online_mall.NewOrderServiceClient(conn)
	res2, err2 := orderClient.Find(context.Background(), &online_mall.Request{
		ID: 1,
	})
	fmt.Println(res2, err2)

}
