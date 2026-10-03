<!--{
  "Title": "避免 SQL 注入风险"
}-->

将 SQL 参数值作为 `sql` 包函数的参数传入，可以避免 SQL 注入风险。`sql` 包中的许多函数分别接受 SQL 语句和语句所需的参数值；另一些函数则用于执行预处理语句并传入参数。

下面的示例使用 `?` 作为 `id` 参数的占位符，将 `id` 的值作为函数参数传入：

```
// Correct format for executing an SQL statement with parameters.
rows, err := db.Query("SELECT * FROM user WHERE id = ?", id)
```

执行数据库操作的 `sql` 包函数会根据传入的参数创建预处理语句。在运行时，`sql` 包将 SQL 语句转换为预处理语句，并将语句与参数分别发送。

**注意：** 参数占位符的形式取决于所用的数据库管理系统和驱动。例如，Postgres 的 [pq 驱动](https://pkg.go.dev/github.com/lib/pq) 使用 `$1` 这样的占位符，而不是 `?`。

你可能会想到用 `fmt` 包的函数，将 SQL 语句和参数拼接成一个字符串，例如：

```
// SECURITY RISK!
rows, err := db.Query(fmt.Sprintf("SELECT * FROM user WHERE id = %s", id))
```

这种做法不安全！Go 会先用参数值替换 `%s` 格式化占位符，拼接出完整的 SQL 语句，再将其发送给数据库管理系统。这会带来 [SQL 注入](https://en.wikipedia.org/wiki/SQL_injection)风险，因为调用方可能将意料之外的 SQL 片段作为 `id` 参数传入。该片段可能以不可预测的方式改变 SQL 语句，对应用造成危害。

例如，如果向 `%s` 传入特定值，最终可能得到如下语句，从而返回数据库中的全部用户记录：

```
SELECT * FROM user WHERE id = 1 OR 1=1;
```
