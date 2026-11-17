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
	conn, err := rabbitmq.NewClient(
		"127.0.0.1",
		"5672",
		"/",
		"admin",
		"local-learning-password",
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := rabbitmq.NewChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close() // defer 后注册的先执行：先关闭 Channel，再关闭连接。

	queueName := "work.queue"
	queues, err := rabbitmq.DeclareQueues(ch, []string{queueName})
	if err != nil {
		return err
	}

	// 最多接收一条未确认消息；确认后再接收下一条。
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	// false 表示关闭自动确认，改为处理成功后手动 Ack。
	messages, err := rabbitmq.ConsumeMessages(ch, queues.Queues[0].Name, false)
	if err != nil {
		return err
	}

	log.Println("消费者 1 已就绪，约每秒处理 30 条，预取 1 条，按 Ctrl+C 退出")
	count := 0
	for message := range messages {
		// 用等待模拟业务处理耗时，约每秒处理 30 条。
		time.Sleep(time.Second / 30)
		log.Printf("消费者 1 收到消息：%s", message.Body)

		// 打印完成后确认当前消息，false 表示只确认这一条。
		if err := message.Ack(false); err != nil {
			return err
		}
		count++
		log.Printf("消费者 1 已处理并确认第 %d 条", count)
	}
	return nil
}
