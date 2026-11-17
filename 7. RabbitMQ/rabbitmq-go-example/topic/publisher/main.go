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

	// 准备四个独立队列，每个消费者监听自己的队列。
	queues, err := rabbitmq.DeclareQueues(ch, []string{"topic.queue1", "topic.queue2", "topic.queue3", "topic.queue4"})
	if err != nil {
		return err
	}

	exchange := "topic.exchange"
	if err := rabbitmq.DeclareExchange(ch, exchange, "topic", true, false, false, nil); err != nil {
		return err
	}

	// 四个队列分别绑定一条规则，顺序与上面的队列列表一一对应。
	bindingKeys := []string{"china.#", "japan.#", "#.weather", "#.news"}
	for i, bindingKey := range bindingKeys {
		if err := rabbitmq.BindQueue(ch, queues.Queues[i].Name, exchange, bindingKey, nil); err != nil {
			return err
		}
	}

	// 按本次要求，发送下面四个原样的 Routing key。
	// 注意：# 仅在绑定键中是通配符，在消息的 Routing key 中不做通配展开。
	// 常见业务消息可改用 china.weather、japan.news 等具体名称。
	for _, routingKey := range []string{"china.#", "japan.#", "#.weather", "#.news"} {
		body := "这是一条 " + routingKey + " 消息"
		if err := rabbitmq.SendMessage(ch, exchange, routingKey, body); err != nil {
			return err
		}
		log.Printf("已向 %s 发送消息：routingKey=%s，body=%s", exchange, routingKey, body)
	}
	return nil
}
