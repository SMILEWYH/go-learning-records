<!--{
  "Title": "取消正在进行的操作"
}-->

可以使用 Go 的 [`context.Context`](https://pkg.go.dev/context#Context) 管理正在进行的操作。`Context` 是 Go 中的标准数据值，用于报告它所代表的整体操作是否已被取消、是否仍有必要继续执行。在应用的函数调用和服务之间传递 `context.Context`，可以让不再需要的处理提前结束并返回错误。有关 `Context` 的更多信息，请参阅 [Go 并发模式：Context](/blog/context)。

例如，你可能希望：

*   结束运行时间过长的操作，包括迟迟无法完成的数据库操作。
*   传递来自其他位置的取消请求，例如客户端关闭连接时发出的取消信号。

许多面向 Go 开发者的 API 都提供了接受 `Context` 参数的方法，便于在整个应用中使用 `Context`。

### 超时后取消数据库操作 {#timeout_cancel}

可以通过 `Context` 设置超时时长或截止时间，超过该时间后就取消操作。要派生一个具有超时时长或截止时间的 `Context`，可以调用 [`context.WithTimeout`](https://pkg.go.dev/context#WithTimeout) 或 [`context.WithDeadline`](https://pkg.go.dev/context#WithDeadline)。

下面的超时示例派生了一个 `Context`，并将其传给 `sql.DB` 的 [`QueryContext`](https://pkg.go.dev/database/sql#DB.QueryContext) 方法。

```
func QueryWithTimeout(ctx context.Context) {
	// Create a Context with a timeout.
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Pass the timeout Context with a query.
	rows, err := db.QueryContext(queryCtx, "SELECT * FROM album")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Handle returned rows.
}
```

从父上下文派生出的上下文，会在父上下文被取消时自动取消。上例中的 `queryCtx` 就是从 `ctx` 派生的。例如，在 HTTP 服务器中，`http.Request.Context` 返回与请求关联的上下文。如果 HTTP 客户端断开连接或取消 HTTP 请求（HTTP/2 支持这种取消方式），该上下文也会被取消。将 HTTP 请求的上下文传给上面的 `QueryWithTimeout`，无论是整个 HTTP 请求被取消，还是查询耗时超过五秒，都会使数据库查询提前停止。

**注意：** 创建带有超时时长或截止时间的 `Context` 后，始终使用 `defer` 调用其返回的 `cancel` 函数。这样，所在函数退出时就会释放新 `Context` 持有的资源。此调用也会取消 `queryCtx`，不过函数返回时，通常已经没有操作需要继续使用 `queryCtx`。
