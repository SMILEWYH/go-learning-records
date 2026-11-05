package main

import (
	"fmt"
	"grpc_study/stream_proto/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ServiceStream struct{}

func (ServiceStream) Fun(request *proto.Request, stream proto.ServiceStream_FunServer) error {
	fmt.Println(request)
	for i := 0; i < 10; i++ {
		stream.Send(&proto.Response{
			Text: fmt.Sprintf("第%d轮数据", i),
		})
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	proto.RegisterServiceStreamServer(server, &ServiceStream{})

	server.Serve(listen)
}
