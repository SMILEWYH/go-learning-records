package main

import (
	"context"
	mall "grpc-learning/7-multiProto/proto/proto_go"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ProductService struct {
	mall.UnimplementedProductServiceServer
}

func (ProductService) Find(ctx context.Context, request *mall.Request) (res *mall.Response, err error) {
	// request.ID：客户端传入的商品 ID

	return &mall.Response{
		Price: 10,
	}, nil
}

type OrderService struct {
	mall.UnimplementedOrderServiceServer
}

func (OrderService) Find(ctx context.Context, request *mall.Request) (res *mall.Response, err error) {
	// request.ID：客户端传入的订单 ID

	return &mall.Response{
		Price: 10,
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	// 使用 mall 包中的两个注册函数，注册到同一个 gRPC 服务器。
	mall.RegisterProductServiceServer(s, &ProductService{})
	mall.RegisterOrderServiceServer(s, &OrderService{})

	log.Println("gRPC 服务已启动，监听 :8080")
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
