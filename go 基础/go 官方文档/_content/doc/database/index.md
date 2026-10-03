<!--{
  "Title": "访问关系型数据库",
  "Breadcrumb": true
}-->

使用 Go，可以在应用中集成各种数据库和数据访问方式。本节介绍如何使用标准库的 [`database/sql`](https://pkg.go.dev/database/sql) 包访问关系型数据库。

如果你希望从入门教程开始学习 Go 数据访问，请参阅[教程：访问关系型数据库](/doc/tutorial/database-access)。

Go 也支持其他数据访问技术，包括以更高抽象层次访问关系型数据库的 ORM 库，以及非关系型的 NoSQL 数据存储。

*   **对象关系映射（ORM）库。** `database/sql` 包提供了较底层的数据访问函数；你也可以使用 Go 在更高的抽象层次上访问数据存储。Go 中两个流行的对象关系映射（ORM）库是 [GORM](https://gorm.io/index.html)（[包参考文档](https://pkg.go.dev/gorm.io/gorm)）和 [ent](https://entgo.io/)（[包参考文档](https://pkg.go.dev/entgo.io/ent)）。
*   **NoSQL 数据存储。** Go 社区为大多数 NoSQL 数据存储开发了驱动，包括 [MongoDB](https://docs.mongodb.com/drivers/go/) 和 [Couchbase](https://docs.couchbase.com/go-sdk/current/hello-world/overview.html)。可以在 [pkg.go.dev](https://pkg.go.dev/) 上搜索更多驱动。

### 支持的数据库管理系统 {#supported_dbms}

Go 支持所有常见的关系型数据库管理系统，包括 MySQL、Oracle、Postgres、SQL Server、SQLite 等。

完整的驱动列表参见 [SQLDrivers](/wiki/SQLDrivers) 页面。

### 执行查询或修改数据库的函数 {#functions}

`database/sql` 包针对不同的数据库操作提供了专门的函数。例如，`Query` 和 `QueryRow` 都可以执行查询，但 `QueryRow` 专用于预期只返回一行的场景，可避免返回一个仅包含一行数据的 `sql.Rows` 所带来的开销。要通过 `INSERT`、`UPDATE` 或 `DELETE` 等 SQL 语句修改数据库，可以使用 `Exec`。

更多内容请参阅：

*   [执行不返回数据的 SQL 语句](/doc/database/change-data)
*   [查询数据](/doc/database/querying)

### 事务 {#transactions}

通过 `sql.Tx`，可以在事务中执行数据库操作。事务允许将多个操作组合起来，最后统一提交，以一个原子步骤应用所有更改；也可以回滚，放弃这些更改。

有关事务的更多信息，请参阅[执行事务](/doc/database/execute-transactions)。

### 取消查询 {#query_cancellation}

如果希望能够取消数据库操作，例如在客户端连接关闭或操作耗时过长时，可以使用 `context.Context`。

对于数据库操作，可以使用 `database/sql` 包中接受 `Context` 参数的函数。通过 `Context`，可以设置操作的超时时长或截止时间，也可以将取消请求沿应用调用链传递到执行 SQL 语句的函数，确保不再需要的资源得到释放。

更多信息请参阅[取消正在进行的操作](/doc/database/cancel-operations)。

### 自动管理的连接池 {#connection_pool}

使用 `sql.DB` 数据库句柄时，实际使用的是一个内置连接池。它会根据代码的需要创建和释放连接。通过 `sql.DB` 句柄访问数据库是 Go 中最常见的方式。更多信息请参阅[打开数据库句柄](/doc/database/open-handle)。

`database/sql` 包会自动管理连接池。对于更复杂的需求，可以按照[设置连接池属性](/doc/database/manage-connections#connection_pool_properties)中的说明调整连接池属性。

对于需要独占一个连接的操作，`database/sql` 包提供了 [`sql.Conn`](https://pkg.go.dev/database/sql#Conn)。当使用 `sql.Tx` 事务并不合适时，`Conn` 尤其有用。

例如，你的代码可能需要：

*   通过 DDL 修改数据库结构，其中的逻辑可能具有自己的事务语义。如[执行事务](/doc/database/execute-transactions)所述，将 `sql` 包的事务函数与 SQL 事务语句混用是不推荐的做法。
*   执行会创建临时表的查询锁定操作。

更多信息请参阅[使用专用连接](/doc/database/manage-connections#dedicated_connections)。
