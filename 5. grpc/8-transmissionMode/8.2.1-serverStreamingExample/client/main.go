package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc_study/stream_proto/proto"
	"io"
	"log"
	"os"

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

	stream, err := client.DownLoadFile(context.Background(), &proto.Request{
		Name: "张三",
	})

	file, err := os.OpenFile("static/1.gvb_web项目搭建1.mp4", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	var index int
	for {
		index++
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Printf("第%d 次， 写入 %d 数据\n", index, len(response.Content))
		writer.Write(response.Content)
	}
	writer.Flush()
}
