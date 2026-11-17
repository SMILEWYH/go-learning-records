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
	queues, err := rabbitmq.DeclareQueues(ch, []string{"direct.queue1", "direct.queue2"})
	if err != nil {
		return err
	}

	exchange := "direct.exchange"
	if err := rabbitmq.DeclareExchange(ch, exchange, "direct", true, false, false, nil); err != nil {
		return err
	}

	// 同一个队列可以有多个绑定键，每个键分别绑定一次。
	// 队列 1 接收 blue 和 red。
	for _, bindingKey := range []string{"blue", "red"} {
		if err := rabbitmq.BindQueue(ch, queues.Queues[0].Name, exchange, bindingKey, nil); err != nil {
			return err
		}
	}

	// 队列 2 接收 yellow 和 red。
	for _, bindingKey := range []string{"yellow", "red"} {
		if err := rabbitmq.BindQueue(ch, queues.Queues[1].Name, exchange, bindingKey, nil); err != nil {
			return err
		}
	}

	// direct 精确匹配绑定键：blue 到队列 1，yellow 到队列 2，red 到两个队列。
	for _, routingKey := range []string{"blue", "yellow", "red"} {
		body := "这是一条 " + routingKey + " 消息"
		if err := rabbitmq.SendMessage(ch, exchange, routingKey, body); err != nil {
			return err
		}
		log.Printf("已向 %s 发送消息：routingKey=%s，body=%s", exchange, routingKey, body)
	}
	return nil
}
