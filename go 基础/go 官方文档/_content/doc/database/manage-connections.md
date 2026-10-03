<!--{
  "Title": "管理连接"
}-->

对于绝大多数程序，无须修改 `sql.DB` 连接池的默认设置。不过，一些复杂程序可能需要调整连接池参数，或显式使用连接。本文介绍这些操作。

[`sql.DB`](https://pkg.go.dev/database/sql#DB) 数据库句柄可供多个 goroutine（轻量级协程）并发安全地使用，也就是其他语言常说的“线程安全”。另一些数据库访问库基于单个连接工作，每个连接一次只能执行一个操作。为解决这一限制，每个 `sql.DB` 都维护着一组到底层数据库的活动连接，并在需要时创建新连接，支持 Go 程序并行执行数据库操作。

连接池能够满足大多数数据访问需求。调用 `sql.DB` 的 `Query` 或 `Exec` 方法时，其实现会从连接池中获取可用连接，必要时创建新连接。不再需要连接后，包会将连接归还连接池。这种机制支持高度并行的数据库访问。

### 设置连接池属性 {#connection_pool_properties}

可以设置属性，控制 `sql` 包管理连接池的方式。要了解这些属性产生的影响，可以通过 [`DB.Stats`](https://pkg.go.dev/database/sql#DB.Stats) 获取统计信息。

#### 设置最大打开连接数 {#max_open_connections}

[`DB.SetMaxOpenConns`](https://pkg.go.dev/database/sql#DB.SetMaxOpenConns) 限制同时打开的连接数量。达到上限后，新的数据库操作会等待已有操作结束，直到可以获取连接或创建新的连接。默认情况下，只要需要连接而现有连接都在使用中，`sql.DB` 就会创建新连接。

请注意，设置上限后，使用数据库连接的过程类似于获取锁或信号量，因此应用可能在等待数据库连接时发生死锁。

#### 设置最大空闲连接数 {#max_idle_connections}

[`DB.SetMaxIdleConns`](https://pkg.go.dev/database/sql#DB.SetMaxIdleConns) 用于修改 `sql.DB` 保留的最大空闲连接数。

一个连接上的 SQL 操作结束后，通常不会立即关闭该连接：应用可能很快再次需要它，保留连接可以避免下一次操作时重新连接数据库。默认情况下，`sql.DB` 最多保留两个空闲连接。对于并行程度较高的程序，提高这一上限可以减少频繁重连。

#### 设置连接的最长空闲时间 {#max_idle_time}

[`DB.SetConnMaxIdleTime`](https://pkg.go.dev/database/sql#DB.SetConnMaxIdleTime) 用于设置连接在关闭前允许空闲的最长时间。`sql.DB` 会关闭空闲时间超过指定时长的连接。

默认情况下，空闲连接加入连接池后，会一直保留到再次需要使用它。当使用 `DB.SetMaxIdleConns` 增大空闲连接上限以应对突发的并行负载时，也可以配合 `DB.SetConnMaxIdleTime`，使这些连接在系统恢复空闲后得到释放。

#### 设置连接的最长生命周期 {#max_connection_lifetime}

[`DB.SetConnMaxLifetime`](https://pkg.go.dev/database/sql#DB.SetConnMaxLifetime) 用于设置连接在关闭前允许保持打开状态的最长时间。

默认情况下，只要满足上述限制，连接就可以一直使用和复用。在某些系统中，例如使用负载均衡数据库服务器时，限制同一个连接的使用时间并定期重新连接会有所帮助。

### 使用专用连接 {#dedicated_connections}

某些数据库会赋予同一连接上连续执行的一组操作特定的隐含语义。`database/sql` 包提供了可用于这种场景的函数。

最常见的例子是事务：通常以 `BEGIN` 命令开始，以 `COMMIT` 或 `ROLLBACK` 命令结束，期间在该连接上执行的所有命令都属于同一个事务。对于这种情况，应使用 `sql` 包提供的事务支持。请参阅[执行事务](/doc/database/execute-transactions)。

对于其他需要让一系列操作全部在同一连接上执行的情况，`sql` 包提供了专用连接。通过 [`DB.Conn`](https://pkg.go.dev/database/sql#DB.Conn) 可以获取一个 [`sql.Conn`](https://pkg.go.dev/database/sql#Conn) 专用连接。`sql.Conn` 的 `BeginTx`、`ExecContext`、`PingContext`、`PrepareContext`、`QueryContext` 和 `QueryRowContext` 方法，与 `DB` 上的对应方法行为类似，但只使用这个专用连接。使用完毕后，必须调用 `Conn.Close` 释放它。
