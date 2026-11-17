// 初始化 rabbitMQ 实例的公共方法包
package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// 创建一个 rabbitMQ 实例，连接对应主机下账号的虚拟主机
func NewClient(host, port, virtualHost, username, password string) (*amqp.Connection, error) {
	address := "amqp://" + username + ":" + password + "@" + host + ":" + port + virtualHost

	client, err := amqp.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("连接 RabbitMQ 失败: %w", err)
	}

	return client, nil
}

// 创建 Channel，go 每一个协程都要单独使用一个 Channel
func NewChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("创建 Channel 失败: %w", err)
	}

	return channel, nil
}

// 声明交换机
// 声明规则为：声明时会判断是否存在同名交换机，如果存在则检查参数是否一致，不一致则报错和关闭 Channel；如果不存在则创建
func DeclareExchange(ch *amqp.Channel, name, kind string, durable, autoDelete, internal bool, args amqp.Table) error {
	err := ch.ExchangeDeclare(
		name,       // 交换机名称，例如 go.direct
		kind,       // 交换机类型：direct、fanout、topic、headers
		durable,    // 是否持久化：重启后保留交换机定义
		autoDelete, // 是否在曾经绑定后，最后一个绑定被移除时自动删除
		internal,   // 是否为内部交换机：true 时生产者不能直接向它发布
		false,      // noWait：等待服务器返回声明结果
		args,       // 额外配置，没有时传 nil
	)
	if err != nil {
		return fmt.Errorf("声明交换机 %q 失败: %w", name, err)
	}
	return nil
}

type QueueResult struct {
	Queues []amqp.Queue
}

// 声名 Queue
// 声明规则为：声明时会判断是否存在同名 Queue，如果存在则检查参数是否一致，不一致则报错和关闭 Channel；如果不存在则创建
func DeclareQueues(ch *amqp.Channel, queueNames []string) (*QueueResult, error) {
	result := &QueueResult{Queues: make([]amqp.Queue, 0, len(queueNames))}

	for _, name := range queueNames {
		queue, err := ch.QueueDeclare(
			name,  // 队列名称
			true,  // durable：RabbitMQ 重启后保留队列定义
			false, // autoDelete：不随最后一个消费者离开而自动删除
			false, // exclusive：不限定为当前连接独占
			false, // noWait：等待服务端返回声明结果
			nil,   // 不指定 x-queue-type，使用 vhost 的默认队列类型
		)
		if err != nil {
			// 参数冲突等服务端错误会关闭当前 Channel，但不影响其他 Channel。
			// 但已成功声明的队列和记录会保留，不会自动回滚。
			return result, fmt.Errorf("声明队列 %q 失败: %w", name, err)
		}

		result.Queues = append(result.Queues, queue)
	}

	return result, nil
}

func BindQueue(ch *amqp.Channel, queueName, exchange, bindingKey string, args amqp.Table) error {
	err := ch.QueueBind(
		queueName,  // 队列名称
		bindingKey, // 绑定键，与发送消息时的 routing key 配合使用
		exchange,   // 交换机名称
		false,      // noWait：等待服务器返回绑定结果
		args,       // 额外绑定配置
	)
	if err != nil {
		return fmt.Errorf("绑定队列 %q 到交换机 %q 失败: %w", queueName, exchange, err)
	}
	return nil
}

// 发送消息
func SendMessage(ch *amqp.Channel, exchange, routingKey, body string) error {
	err := ch.Publish(
		exchange,   // 交换机名称
		routingKey, // 路由键，由交换机决定消息进入哪个队列
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         []byte(body),
		},
	)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}
	return nil
}

// 订阅队列：autoAck 为 true 时自动确认，为 false 时由调用方处理成功后手动 Ack。
func ConsumeMessages(ch *amqp.Channel, queueName string, autoAck bool) (<-chan amqp.Delivery, error) {
	messages, err := ch.Consume(
		queueName, // 队列名称
		"",        // 自动生成消费者标识
		autoAck,   // 是否自动确认，由调用方决定
		false,     // exclusive：允许其他消费者共同消费
		false,     // noLocal
		false,     // noWait：等待服务器确认订阅成功
		nil,       // 没有额外配置
	)
	if err != nil {
		return nil, fmt.Errorf("订阅队列 %q 失败: %w", queueName, err)
	}
	return messages, nil
}
