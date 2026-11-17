# Fanout 发布/订阅示例

复用 `rabbitmq/init.go` 中已有的方法。生产者声明交换机、两个队列和绑定，因此无需提前在管理页面手动创建。

```text
fanout/
├── publisher/main.go
├── consumer/
│   ├── consumer1/main.go
│   └── consumer2/main.go
└── README.md
```

消息路径：

```text
                              ┌→ fanout.queue1 → 消费者 1
生产者 → fanout.exchange ──────┤
          （fanout 类型）       └→ fanout.queue2 → 消费者 2
```

生产者每次运行只调用一次 `SendMessage`，两个队列各收到一份。fanout 忽略绑定键和 Routing key，示例中都传空字符串。

## 运行

打开三个终端，分别进入项目根目录：

```bash
cd "/Users/vvyh/Desktop/Go/7. RabbitMQ/rabbitmq-go-example"
```

终端 1：

```bash
go run ./fanout/consumer/consumer1
```

终端 2：

```bash
go run ./fanout/consumer/consumer2
```

等两个消费者输出“已就绪”后，在终端 3 发送一条消息：

```bash
go run ./fanout/publisher
```

两个消费者都会打印同样的消息正文：

```text
消费者 1 收到广播消息：你好，这是一条 fanout 广播消息
消费者 2 收到广播消息：你好，这是一条 fanout 广播消息
```

消费者分别调用 `ConsumeMessages(ch, queueName, false)`，监听自己的队列。沿用手动 Ack 和预取 1 条，打印后确认，不设置模拟处理耗时。消费者会持续等待消息，按 Ctrl+C 退出。

生产者会先准备好两个队列和绑定，再发送消息，因此先运行生产者、后运行消费者也能收到保存在各自队列中的消息。再次运行生产者会再广播一条。

Work Queue 是多个消费者分担同一个队列的任务；这里是两个消费者分别监听两个队列，同一条发布消息会被两边各处理一次。
