# 运维日志分析工具

给维护多个服务的团队生成本地访问日志报告：统计每个服务的请求量、4xx/5xx 次数、平均与最大耗时、服务端错误率，并定位格式有问题的日志行。

这是命令行批处理业务，与订单和预约两个在线服务不同。输入为已导出的 JSON Lines 日志，不连接线上机器；固定 worker 数量并发读取文件，主 goroutine 汇总结果。只使用标准库。

## 1. 立即运行

```sh
cd ~/Desktop/go-log-report
go run .
# 写入报告文件：
go run . -input ./samples -out ./reports/summary.json -workers 2
```

内置两份模拟日志共 7 行，其中 6 行合法，1 行为故意保留的格式错误。预期：

| 服务 | 请求数 | 4xx | 5xx | 平均耗时 ms | 最大耗时 ms | 5xx 比例 |
| --- | --- | --- | --- | --- | --- | --- |
| booking | 3 | 0 | 1 | 100 | 120 | 1/3 |
| orders | 3 | 1 | 1 | 150 | 300 | 1/3 |

`server_error_rate` 为 0..1 的比例，不是百分数。示例输出约为 `0.3333333333333333`；4xx 单独计数，不并入服务端错误率。错误行会显示 `api-01.jsonl:4`，不将原始日志内容写入诊断。

```sh
# 严格模式：本例包含坏行，因此应该失败并非零退出，不写报告。
go run . -strict -out ./reports/strict.json
```

默认宽松模式保留有效行，并明确报告无效行数量。严格模式只要存在无效行就整批失败；已有目标报告保持不变。

## 2. 输入格式与业务约定

每行一个 JSON 对象：

```json
{"time":"2026-09-27T09:00:00+08:00","service":"orders","status":200,"duration_ms":100}
```

| 字段 | 规则 |
| --- | --- |
| time | RFC3339 时间，必须带时区 |
| service | 1..64 位小写字母、数字、`-`、`_` |
| status | 100..599 整数 |
| duration_ms | 必填的 0..3600000 整数，单位毫秒，不能是 null |

- 只扫描输入目录第一层的 `.jsonl` 普通文件，扩展名大小写不敏感；忽略软链接、子目录和其他扩展名。
- 每个文件最多 10 MiB、100000 行，每批最多 100 个文件；Scanner 的缓冲上限为 64 KiB，超长行视为文件读取失败。
- 文件读取失败、超长行、取消、超时都会导致整批失败，不写出部分成功报告。普通 JSON 错误按 strict 选项处理。
- 空行是无效日志；空文件合法，计入文件数，但请求量为 0；无 `.jsonl` 文件的目录视为参数错误。
- 拒绝未知字段和一行多个 JSON。文件末尾没有换行仍读取最后一行。
- 重复日志按多次请求计数，没有跨文件去重规则；分析覆盖输入文件的全部日期，不自动选择“今天”。
- 服务按名称排序，问题按文件名、行号排序，因此 worker 数不同仍输出相同报告。
- 只保留前 100 条问题详情，仍统计全部无效行，并通过 issues_truncated 告知截断。
- 请处理已经停止写入的日志副本。文件在分析期间被其他程序修改时，不提供快照一致性；os.Root 限制目录访问边界，但不是文件系统快照。

## 3. 参数

| 参数 | 默认值 | 约束 |
| --- | --- | --- |
| -input | ./samples | 输入目录，必须存在 |
| -out | 空 | 为空输出 JSON 到 stdout，否则保存 .json 报告 |
| -workers | 4 | 1..16，并发文件数 |
| -timeout | 10s | 大于 0，最多 10m；用于分析任务的协作取消 |
| -strict | false | 是否拒绝任何坏行 |
| -version | false | 显示构建版本 |

```sh
go run . -input ./samples -workers 1 -timeout 5s
go run . -h
```

错误输出到 stderr，退出码为 1；成功为 0。`-out` 文件允许覆盖旧报告，通过同目录临时文件完整替换；必须以 `.json` 结尾，且不能与输入日志是同一个文件。只读取输入，不修改或归档原始日志。相对路径按当前工作目录解析。

Ctrl+C / SIGTERM 和分析超时会传播给任务投递、任务接收、结果发送以及逐行处理。主函数等待各协程退出再关闭目录。这里针对普通本地文件；context 不能强制打断底层操作系统已经阻塞的磁盘读取，也不覆盖报告写入阶段。

## 4. 阅读顺序与知识对应

| 位置 | 重点 | go-demo 章节 |
| --- | --- | --- |
| internal/analyzer/analyzer.go: Event/Stats/Report | 数据模型、JSON tag、整数统计 | 03、04、11 |
| parseLine | Decoder、可选指针、业务错误 | 09、18、27 |
| readFile | Scanner、EOF、大小限制、关闭资源 | 10、13、20 |
| Analyze 前半段 | 目录过滤、os.Root | 23 |
| Analyze 中间部分 | jobs/results、WaitGroup、生产者关闭通道 | 14、15 |
| ctx.Done 与 ctx.Err | 每个阻塞位置响应取消 | 16 |
| Analyze 后半段 | 单个消费者汇总 map，排序后输出 | 07、08、19 |
| internal/filedata/file.go | 完整替换报告文件 | 21 |
| main.go | flags、信号、退出码、构建版本 | 02、18、30 |
| analyzer_test.go | 临时目录、稳定断言、错误输入 | 24 |

任务关系：主流程枚举文件 → 唯一生产者投递 jobs → 固定 worker 读取 → 唯一消费者汇总；协调者等 worker 全结束后关闭 results。worker 不共享统计 map，因此无需给每一次计数加 Mutex。

## 5. 测试与构建

```sh
go test -race -cover ./...
go vet ./...
mkdir -p bin
go build -trimpath -ldflags '-X main.version=v1.0.0' -o ./bin/go-log-report .
./bin/go-log-report -version
./bin/go-log-report -input ./samples -out ./reports/summary.json
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o ./bin/go-log-report-linux-amd64 .
```

沿用 Go 1.27.0 模块要求，本机 Go 1.27.1 验证。测试覆盖多 worker 结果一致、实际统计值、最后一行没有换行、坏行、strict、空目录、空文件、软链接忽略、取消、超长行和错误详情上限。不会访问外部网站或生产日志。

## 6. 当前边界和练习

这是一次性本地报告工具，没有日志实时采集、告警通知、仪表盘、数据库或分布式任务调度。所有输入和输出都由本机用户选择。为了控制内存和计算量，保留了文件数、文件大小和行数限制。

读懂后可以用现有知识继续练习：

1. 添加 `-service`，只统计指定服务，但保留整份输入的无效行统计。
2. 添加按状态码的分布统计，并按状态码排序输出。
3. 使用第 22 章的 io.Copy，将已完成报告备份到新文件，要求目标已存在时拒绝覆盖。

当前版本没有未实现的业务占位符，以上均为可选扩展。
