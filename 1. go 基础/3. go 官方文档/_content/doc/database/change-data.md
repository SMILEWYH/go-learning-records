<!--{
  "Title": "执行不返回数据的 SQL 语句"
}-->

执行不返回数据的数据库操作时，可以使用 `database/sql` 包中的 `Exec` 或 `ExecContext` 方法，例如执行 `INSERT`、`DELETE` 和 `UPDATE` 语句。

如果查询可能返回数据行，应改用 `Query` 或 `QueryContext`。更多信息请参阅[查询数据库](/doc/database/querying)。

`ExecContext` 与 `Exec` 的工作方式相同，但额外接受一个 `context.Context` 参数，详见[取消正在进行的操作](/doc/database/cancel-operations)。

下面的示例使用 [`DB.Exec`](https://pkg.go.dev/database/sql#DB.Exec) 执行语句，向 `album` 表中添加一条新的唱片记录。

```
func AddAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist) VALUES (?, ?)", alb.Title, alb.Artist)
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}

	// Get the new album's generated ID for the client.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddAlbum: %v", err)
	}
	// Return the new album's ID.
	return id, nil
}
```

`DB.Exec` 返回一个 [`sql.Result`](https://pkg.go.dev/database/sql#Result) 和一个错误值。如果错误为 `nil`，就可以通过 `Result` 获取最后插入记录的 ID（如上例），或者获取该操作影响的行数。

**注意：** 预处理语句中参数占位符的形式取决于所用的数据库管理系统和驱动。例如，Postgres 的 [pq 驱动](https://pkg.go.dev/github.com/lib/pq) 要求使用 `$1` 这样的占位符，而不是 `?`。

如果需要反复执行相同的 SQL 语句，可以考虑通过 `sql.Stmt` 创建可复用的预处理语句。更多信息请参阅[使用预处理语句](/doc/database/prepared-statements)。

**警告：** 不要使用 `fmt.Sprintf` 等字符串格式化函数拼接 SQL 语句，否则可能引入 SQL 注入风险。更多信息请参阅[避免 SQL 注入风险](/doc/database/sql-injection)。

#### 执行不返回数据行的 SQL 语句的方法 {#no_rows_functions}

<table id="no-rows-functions-list" class="DocTable">
  <thead>
    <tr class="DocTable-head">
      <th class="DocTable-cell" width="20%">函数</th>
      <th class="DocTable-cell">说明</th>
    </tr>
  </thead>
  <tbody>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#DB.Exec">DB.Exec</a></code><br/>
        <code><a href="https://pkg.go.dev/database/sql#DB.ExecContext">DB.ExecContext</a></code>
      </td>
      <td class="DocTable-cell">独立执行一条 SQL 语句。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Tx.Exec">Tx.Exec</a></code><br/>
        <code><a href="https://pkg.go.dev/database/sql#Tx.ExecContext">Tx.ExecContext</a></code>
      </td>
      <td class="DocTable-cell">在一个事务中执行 SQL 语句。更多信息请参阅
          <a href="/doc/database/execute-transactions">执行事务</a>。
      </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Stmt.Exec">Stmt.Exec</a></code><br/>
        <code><a href="https://pkg.go.dev/database/sql#Stmt.ExecContext">Stmt.ExecContext</a></code>
      </td>
      <td class="DocTable-cell">执行已经准备好的 SQL 语句。更多信息请参阅
          <a href="/doc/database/prepared-statements">使用预处理语句</a>。
      </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Conn.ExecContext">Conn.ExecContext</a></code>
      </td>
      <td class="DocTable-cell">用于专用连接。更多信息请参阅
          <a href="/doc/database/manage-connections">管理连接</a>。
      </td>
    </tr>
  </tbody>
</table>
