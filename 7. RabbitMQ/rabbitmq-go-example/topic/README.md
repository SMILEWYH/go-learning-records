# Topic 主题路由：四个独立消费者

使用 `topic.exchange`（topic 类型），复用 `rabbitmq/init.go`，四个消费者分别通过 `ConsumeMessages` 监听自己的队列。

```text
topic/
├── publisher/main.go
├── consumer/
│   ├── consumer1/main.go
│   ├── consumer2/main.go
│   ├── consumer3/main.go
│   └── consumer4/main.go
└── README.md
```

| 消费者 | 队列 | Binding key |
| --- | --- | --- |
| consumer1 | topic.queue1 | china.# |
| consumer2 | topic.queue2 | japan.# |
| consumer3 | topic.queue3 | #.weather |
| consumer4 | topic.queue4 | #.news |

## 发送与匹配

按当前要求，生产者发送四个原样的 Routing key：`china.#`、`japan.#`、`#.weather`、`#.news`。在以上绑定下，四个消费者分别收到一条对应消息。

`#` 只有在 Binding key 中才表示零个或多个单词，消息的 Routing key 不会做通配展开。实际业务更常发送具体主题，例如：

| 具体的 Routing key | 收到消息的消费者 |
| --- | --- |
| china.weather | consumer1、consumer3 |
| japan.news | consumer2、consumer4 |
| china.sports | consumer1 |
| france.news | consumer4 |

可修改生产者的发送列表来体验这些具体主题。规则见 [RabbitMQ 官方 Topic 教程](https://www.rabbitmq.com/tutorials/tutorial-five-go)。

## 运行

打开五个终端，每个终端进入项目根目录：

```bash
cd "/Users/vvyh/Desktop/Go/7. RabbitMQ/rabbitmq-go-example"
```

前四个终端分别运行一个消费者：

```bash
# 终端 1
go run ./topic/consumer/consumer1
```

```bash
# 终端 2
go run ./topic/consumer/consumer2
```

```bash
# 终端 3
go run ./topic/consumer/consumer3
```

```bash
# 终端 4
go run ./topic/consumer/consumer4
```

等四个消费者都输出“已就绪”，在终端 5 发送：

```bash
go run ./topic/publisher
```

生产者声明四个队列、交换机和对应绑定，再发送消息。消费者预取 1 条，打印后手动 Ack，按 Ctrl+C 退出。

如果之前运行过旧的两队列版本，旧绑定不会因修改代码自动消失。应移除旧的 `topic.queue1 → japan.#`、`topic.queue2 → #.weather`、`topic.queue2 → #.news` 三条绑定，再使用上表的绑定。当前电脑上的这些旧绑定已在本次调整中清理。
