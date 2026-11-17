package main

import (
	"log"

	"rabbitmq-go-example/rabbitmq"
)

// 最简单示例
// 初始化后并没有使用 rabbitMQ 收发消息
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, err := rabbitmq.NewClient(
		"127.0.0.1",               // host：主机名或 IP
		"5672",                    // port：AMQP 端口
		"/",                       // virtualHost：虚拟主机
		"admin",                   // username：用户名
		"local-learning-password", // password：密码
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("RabbitMQ 网络连接建立成功！")

	// 第二步：在连接上创建 Channel。
	ch, err := rabbitmq.NewChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	// 单独声明队列：不存在时创建，已存在时检查参数。
	queueResult, err := rabbitmq.DeclareQueues(ch, []string{"go.hello"})
	if err != nil {
		return err
	}

	log.Printf("Channel 创建成功！队列：%s", queueResult.Queues[0].Name)
	return nil
}
