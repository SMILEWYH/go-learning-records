# 工作队列：100 条消息、两个消费者

生产者声明 `work.exchange`（direct 类型），将 `work.queue` 队列绑定到它，绑定和发送的键都使用空字符串 `""`，再以约每秒 50 条的速度发送 100 条消息。暂时不设置具体路由键；direct 交换机仍按键匹配，空路由键匹配空绑定键。

RabbitMQ 保留 `amq.` 前缀，当前容器会拒绝新建 `amq.work`，因此示例使用 `work.exchange`。

```text
生产者 → work.exchange → work.queue ┬→ 消费者 1：约 30 条/秒，预取 1 条
                                   └→ 消费者 2：约 10 条/秒，预取 1 条
```

## 运行

打开三个终端，每个终端先进入项目根目录：

```bash
cd "/Users/vvyh/Desktop/Go/7. RabbitMQ/rabbitmq-go-example"
```

先启动两个消费者，等各自输出“已就绪”：

```bash
# 终端 1
go run ./work-queues/consumer/consumer1
```

```bash
# 终端 2
go run ./work-queues/consumer/consumer2
```

然后在第三个终端发送 100 条消息：

```bash
go run ./work-queues/publisher
```

`consumer/consumer1/main.go` 和 `consumer/consumer2/main.go` 独立运行，各自建立连接和 Channel，订阅同一个队列。

- 生产者：`time.NewTicker(time.Second / 50)` 每 20 毫秒通知一次，收到通知后发送一条，总共 100 条。
- 消费者 1：每条消息处理时 `time.Sleep(time.Second / 30)`，约每秒处理 30 条。
- 消费者 2：每条消息处理时 `time.Sleep(time.Second / 10)`，约每秒处理 10 条。

使用简单等待模拟速度，实际还会受到日志、网络和调度影响。

两个消费者都在订阅前设置 `ch.Qos(1, 0, false)`，并传 `autoAck = false` 开启手动确认：

```go
if err := ch.Qos(1, 0, false); err != nil {
    return err
}
messages, err := rabbitmq.ConsumeMessages(ch, "work.queue", false)
if err != nil {
    return err
}
for message := range messages {
    time.Sleep(time.Second / 30) // 消费者 2 改成 / 10
    log.Printf("收到消息：%s", message.Body)
    if err := message.Ack(false); err != nil {
        return err
    }
}
```

`ConsumeMessages` 的第三个参数是 `autoAck`：`true` 为自动确认，`false` 为手动确认。
每个消费者最多持有一条未确认消息，处理并 Ack 后才继续获得消息，因此快消费者通常会处理更多任务，但不保证严格按 75/25 分配。
生产速度约 50 条/秒，高于两个消费者合计约 40 条/秒，发送期间可能出现短暂积压，生产者发送完后消费者会继续处理。

生产者不开启发布确认；本次开启的是消费者的处理确认。没有消息时，`for range` 会等待；消息通道关闭且剩余消息读完后，循环结束。按 Ctrl+C 结束进程，未确认的消息在连接关闭后可以重新投递。
