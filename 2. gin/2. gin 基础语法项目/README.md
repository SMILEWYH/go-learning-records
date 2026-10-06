# Gin 实战项目

这三个独立项目从同级 `go 基础/go 基础语法项目` 升级而来，使用 Gin v1.12.0。模块要求 Go 1.27.0 或更高版本。每个项目都有独立的 `go.mod`、`go.sum`、中文说明和测试。

| 项目 | 内容 | 启动命令（进入对应目录后） | 默认地址 |
| --- | --- | --- | --- |
| [go-shop-orders](go-shop-orders/README.md) | 商品、库存、订单状态流转 | `go run .` | http://127.0.0.1:8091 |
| [go-room-booking](go-room-booking/README.md) | 房间、空闲查询、预约和取消 | `go run .` | http://127.0.0.1:8092 |
| [go-log-report](go-log-report/README.md) | 日志分析和 JSON 报告 | `go run . -serve` | http://127.0.0.1:8093 |

## 运行

分别在三个终端中执行：

```sh
cd ~/Desktop/Go/gin/go-shop-orders
go run .
```

```sh
cd ~/Desktop/Go/gin/go-room-booking
go run .
```

```sh
cd ~/Desktop/Go/gin/go-log-report
go run . -serve
```

首次运行会下载 Gin 依赖。默认依赖代理无法访问时，可在对应项目目录执行 `GOPROXY=https://goproxy.cn,direct go mod download` 后重试。

订单和预约项目可运行 `go run ./cmd/demo`，自动验证完整业务流程。日志项目保留 `go run .` 命令行模式，`go run . -out ./reports/summary.json` 将报告写入文件。

## 学习重点

1. 从 `gin.New()`、GET/POST/PATCH 路由和 `*gin.Context` 开始。
2. 对照订单与预约项目的 `router.Group`、`c.Param`、`c.Query`、`c.JSON`。
3. 日志接口使用 `ShouldBindQuery` 与 `binding` 标签校验并发数。
4. 读 `internal/web/router.go`：Logger、CustomRecovery、NoRoute、NoMethod。
5. 理解 `c.Request.Context()` 如何把取消和超时传递给日志分析协程。
6. 用 `httptest` 测试 Gin 的参数、状态码和错误响应。

写接口沿用严格 JSON 解码，保留未知字段检查和请求大小限制。业务继续使用原有锁、校验和 JSON 文件存储，不需要安装 MySQL。

## 验证与构建

在每个项目目录执行：

```sh
go test -race -cover ./...
go vet ./...
go build ./...
```

需要独立程序时，执行 `go build -o ./bin/server .`，再运行 `./bin/server`（日志 HTTP 服务加 `-serve`）。版本注入、参数和接口示例见各项目 README。`docs/验收记录.md` 记录本次升级的验证，`docs/原版验收记录.md` 保留原项目历史。

## 参考

- [Gin v1.12.0 发布说明](https://github.com/gin-gonic/gin/releases/tag/v1.12.0)
- [Gin 中间件说明](https://gin-gonic.com/en/docs/middleware/)
- [Gin 自定义异常恢复](https://gin-gonic.com/en/docs/middleware/custom-recovery/)
