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
	defer ch.Close()

	// 准备两个独立队列，每个消费者监听自己的队列。
	queues, err := rabbitmq.DeclareQueues(ch, []string{"fanout.queue1", "fanout.queue2"})
	if err != nil {
		return err
	}

	exchange := "fanout.exchange"
	if err := rabbitmq.DeclareExchange(ch, exchange, "fanout", true, false, false, nil); err != nil {
		return err
	}

	// fanout 忽略绑定键，将两个队列都绑定到这个交换机。
	for _, queue := range queues.Queues {
		if err := rabbitmq.BindQueue(ch, queue.Name, exchange, "", nil); err != nil {
			return err
		}
	}

	// 只发送一次，由交换机向两个绑定队列各投递一份。
	body := "你好，这是一条 fanout 广播消息"
	if err := rabbitmq.SendMessage(ch, exchange, "", body); err != nil {
		return err
	}
	log.Printf("已向 %s 发送一条消息：%s", exchange, body)
	return nil
}
