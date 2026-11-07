package main

import (
	"context"
	online_mall "grpc-learning/6-multiService/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ProductService struct {
	online_mall.UnimplementedProductServiceServer
}

func (ProductService) Find(ctx context.Context, request *online_mall.Request) (res *online_mall.Response, err error) {
	// request.ID

	return &online_mall.Response{
		Price: 10,
	}, nil
}

type OrderService struct {
	online_mall.UnimplementedOrderServiceServer
}

func (OrderService) Find(ctx context.Context, request *online_mall.Request) (res *online_mall.Response, err error) {
	// request.ID

	return &online_mall.Response{
		Price: 10,
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	online_mall.RegisterProductServiceServer(s, &ProductService{})
	online_mall.RegisterOrderServiceServer(s, &OrderService{})

	err = s.Serve(listen)
}
