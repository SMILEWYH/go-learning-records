package main

import (
	"fmt"
	"grpc_study/stream_proto/proto"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

type ServiceStream struct{}

func (ServiceStream) DownLoadFile(request *proto.Request, stream proto.ServiceStream_DownLoadFileServer) error {
	fmt.Println(request)
	file, err := os.Open("static/1.gvb_web项目搭建.mp4")
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		buf := make([]byte, 2048)
		_, err = file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		stream.Send(&proto.FileResponse{
			Content: buf,
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
