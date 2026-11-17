package main

import (
	"log"
	"time"

	"rabbitmq-go-example/rabbitmq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	client, err := rabbitmq.NewClient(
		"127.0.0.1",
		"5672",
		"/",
		"admin",
		"local-learning-password",
	)
	if err != nil {
		return err
	}
	defer client.Close()

	ch, err := rabbitmq.NewChannel(client)
	if err != nil {
		return err
	}
	defer ch.Close() // defer 后注册的先执行：先关闭 Channel，再关闭连接。

	queueName := "work.queue"
	queues, err := rabbitmq.DeclareQueues(ch, []string{queueName})
	if err != nil {
		return err
	}

	exchange := "work.exchange"
	if err := rabbitmq.DeclareExchange(ch, exchange, "direct", true, false, false, nil); err != nil {
		return err
	}

	if err := rabbitmq.BindQueue(ch, queues.Queues[0].Name, exchange, "", nil); err != nil {
		return err
	}

	// 定时器每隔 20 毫秒通知一次，约每秒发送 50 条。
	ticker := time.NewTicker(time.Second / 50)
	defer ticker.Stop()
	for i := 1; i <= 100; i++ {
		<-ticker.C
		if err := rabbitmq.SendMessage(ch, exchange, "", "消息"); err != nil {
			return err
		}
		log.Printf("已发送：%s", "消息")
	}
	log.Println("100 条消息发送调用完成")
	return nil
}
