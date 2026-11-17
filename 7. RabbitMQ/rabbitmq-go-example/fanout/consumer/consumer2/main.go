package main

import (
	"log"

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
	defer ch.Close()

	queueName := "fanout.queue2"
	if _, err := rabbitmq.DeclareQueues(ch, []string{queueName}); err != nil {
		return err
	}

	// 预取一条，打印后手动确认；只消费 fanout.queue2。
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}
	messages, err := rabbitmq.ConsumeMessages(ch, queueName, false)
	if err != nil {
		return err
	}

	log.Printf("消费者 2 已就绪，监听 %s，按 Ctrl+C 退出", queueName)
	for message := range messages {
		log.Printf("消费者 2 收到广播消息：%s", message.Body)
		if err := message.Ack(false); err != nil {
			return err
		}
	}
	return nil
}
