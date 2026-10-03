<!--{
  "Title": "打开数据库句柄",
  "Breadcrumb": true
}-->

[`database/sql`](https://pkg.go.dev/database/sql) 包通过减少手动管理连接的需要，简化数据库访问。与许多数据访问 API 不同，使用 `database/sql` 时，不需要显式打开连接、执行操作，再关闭连接。代码打开的是一个代表连接池的数据库句柄，然后通过句柄执行数据访问操作。只有在需要释放资源时，例如释放查询结果行或预处理语句持有的资源，才调用相应的 `Close` 方法。

换句话说，由 [`sql.DB`](https://pkg.go.dev/database/sql#DB) 表示的数据库句柄负责管理连接，代替你的代码打开和关闭连接。通过句柄执行的数据库操作可以并发访问数据库。更多信息请参阅[管理连接](/doc/database/manage-connections)。

**注意：** 也可以预留一个数据库连接供特定操作使用。更多信息请参阅[使用专用连接](/doc/database/manage-connections#dedicated_connections)。

除了 `database/sql` 包提供的 API，Go 社区还为所有常见以及许多不常见的数据库管理系统（DBMS）开发了驱动。

打开数据库句柄的主要步骤如下：

1. 查找驱动。

    驱动负责在 Go 代码与数据库之间转换请求和响应。更多信息请参阅[查找并导入数据库驱动](#database_driver)。

2. 打开数据库句柄。

    导入驱动后，就可以为特定数据库打开句柄。更多信息请参阅[打开数据库句柄](#opening_handle)。

3. 确认连接可用。

    打开数据库句柄后，可以检查连接是否可用。更多信息请参阅[确认连接](#confirm_connection)。

通常，代码不会显式打开或关闭数据库连接，这些工作由数据库句柄完成。不过，代码仍应释放执行过程中获得的资源，例如包含查询结果的 `sql.Rows`。更多信息请参阅[释放资源](#free_resources)。

### 查找并导入数据库驱动 {#database_driver}

需要选择一个支持所用数据库管理系统的驱动。可以在 [SQLDrivers](/wiki/SQLDrivers) 中查找适合数据库的驱动。

与其他 Go 包一样，导入驱动包后，代码就能使用它。例如：

```
import "github.com/go-sql-driver/mysql"
```

注意，如果不会直接调用驱动包中的函数，例如仅由 `sql` 包隐式使用驱动，就需要使用空白导入，即在导入路径前加上下划线：

```
import _ "github.com/go-sql-driver/mysql"
```

**注意：** 推荐使用 `database/sql` 包中的函数执行数据库操作，避免直接使用驱动自身的 API。这样可以降低代码与特定数据库管理系统的耦合度，方便将来更换数据库。

### 打开数据库句柄 {#opening_handle}

`sql.DB` 数据库句柄支持读取和写入数据库，既可以独立执行操作，也可以在事务中执行。

调用 `sql.Open`（接受连接字符串）或 `sql.OpenDB`（接受 `driver.Connector`），都可以获取数据库句柄。两者都返回指向 [`sql.DB`](https://pkg.go.dev/database/sql#DB) 的指针。

**注意：** 不要将数据库凭据保存在 Go 源代码中。更多信息请参阅[存储数据库凭据](#store_credentials)。

#### 使用连接字符串打开 {#open_connection_string}

要通过连接字符串连接数据库，可以使用 [`sql.Open` 函数](https://pkg.go.dev/database/sql#Open)。字符串格式取决于所用的驱动。

下面是 MySQL 的示例：

```
db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/jazzrecords")
if err != nil {
	log.Fatal(err)
}
```

不过，以结构化方式设置连接属性，通常能让代码更易读。具体方式因驱动而异。

例如，可以将上面的示例替换为下面的代码：使用 MySQL 驱动的 [`Config`](https://pkg.go.dev/github.com/go-sql-driver/mysql#Config) 设置属性，再通过其 [`FormatDSN` 方法](https://pkg.go.dev/github.com/go-sql-driver/mysql#Config.FormatDSN) 构造连接字符串。

```
// Specify connection properties.
cfg := mysql.NewConfig()
cfg.User = username
cfg.Passwd = password
cfg.Net = "tcp"
cfg.Addr = "127.0.0.1:3306"
cfg.DBName = "jazzrecords"

// Get a database handle.
db, err = sql.Open("mysql", cfg.FormatDSN())
if err != nil {
	log.Fatal(err)
}
```

#### 使用 Connector 打开 {#open_connector}

如果希望使用连接字符串无法表达的驱动专属连接功能，可以使用 [`sql.OpenDB` 函数](https://pkg.go.dev/database/sql#OpenDB)。每种驱动都支持自己的连接属性，通常也提供针对特定数据库管理系统定制连接请求的方法。

将前面的 `sql.Open` 示例改为 `sql.OpenDB` 后，可以像下面这样创建句柄：

```
// Specify connection properties.
cfg := mysql.NewConfig()
cfg.User = username
cfg.Passwd = password
cfg.Net = "tcp"
cfg.Addr = "127.0.0.1:3306"
cfg.DBName = "jazzrecords"

// Get a driver-specific connector.
connector, err := mysql.NewConnector(&cfg)
if err != nil {
	log.Fatal(err)
}

// Get a database handle.
db = sql.OpenDB(connector)
```

#### 处理错误 {#handle_errors}

代码应检查创建句柄时返回的错误，例如 `sql.Open` 的返回值。这类错误不是建立连接时的错误，而是初始化句柄失败的错误。例如，无法解析指定的 DSN 时就可能发生这种情况。

### 确认连接 {#confirm_connection}

打开数据库句柄时，`sql` 包不一定立即创建数据库连接，而可能等代码实际需要连接时再创建。如果不会马上使用数据库，但希望确认能够建立连接，可以调用 [`Ping`](https://pkg.go.dev/database/sql#DB.Ping) 或 [`PingContext`](https://pkg.go.dev/database/sql#DB.PingContext)。

下面的代码通过 Ping 数据库来确认连接可用。

```
db, err = sql.Open("mysql", connString)

// Confirm a successful connection.
if err := db.Ping(); err != nil {
	log.Fatal(err)
}
```

### 存储数据库凭据 {#store_credentials}

不要将数据库凭据存储在 Go 源代码中，否则可能使他人获得数据库内容。应将凭据保存在代码之外、但代码能够访问的位置。例如，可以使用凭据管理服务来保存凭据，并通过其 API 获取用于数据库身份验证的信息。

一种常见方式是在程序启动前将凭据放入环境变量，例如从密钥管理服务中加载，然后在 Go 程序中通过 [`os.Getenv`](https://pkg.go.dev/os#Getenv) 读取：

```
username := os.Getenv("DB_USER")
password := os.Getenv("DB_PASS")
```

这种方式也便于在本地测试时手动设置环境变量。

### 释放资源 {#free_resources}

虽然使用 `database/sql` 包时无须显式管理或关闭连接，但代码仍应在不再需要时释放已获取的资源，例如表示查询结果的 `sql.Rows`，或表示预处理语句的 `sql.Stmt` 所持有的资源。

通常可以使用 `defer` 延迟调用 `Close`，确保在所在函数退出时释放资源。

下面的示例通过延迟调用 `Close`，释放 [`sql.Rows`](https://pkg.go.dev/database/sql#Rows) 持有的资源。

```
rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Loop through returned rows.
```
