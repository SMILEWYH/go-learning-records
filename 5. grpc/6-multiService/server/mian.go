package main

import (
	"context"
	"fmt"
	"grpc_study/grpc_proto/duo_proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type VideoServer struct {
}

func (VideoServer) Look(ctx context.Context, request *duo_proto.Request) (res *duo_proto.Response, err error) {
	fmt.Println("video:", request)
	return &duo_proto.Response{
		Name: "fengfeng",
	}, nil
}

type OrderServer struct {
}

func (OrderServer) Buy(ctx context.Context, request *duo_proto.Request) (res *duo_proto.Response, err error) {
	fmt.Println("order:", request)
	return &duo_proto.Response{
		Name: "fengfeng",
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	duo_proto.RegisterVideoServiceServer(s, &VideoServer{})
	duo_proto.RegisterOrderServiceServer(s, &OrderServer{})

	err = s.Serve(listen)
}
