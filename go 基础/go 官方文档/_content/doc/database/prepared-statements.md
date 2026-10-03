<!--{
  "Title": "使用预处理语句"
}-->

可以定义一个预处理语句并反复使用。这样能够避免每次执行数据库操作时重新创建语句的开销，从而略微提高代码运行效率。

**注意：** 预处理语句中参数占位符的形式取决于所用的数据库管理系统和驱动。例如，Postgres 的 [pq 驱动](https://pkg.go.dev/github.com/lib/pq) 要求使用 `$1` 这样的占位符，而不是 `?`。

### 什么是预处理语句？ {#what_prepared_statement}

预处理语句是由数据库管理系统解析并保存的 SQL，通常包含参数占位符，但不包含实际参数值。之后，可以传入一组参数值来执行该语句。

### 如何使用预处理语句 {#use_prepared_statement}

如果预计会反复执行相同的 SQL，可以使用 `sql.Stmt` 预先准备语句，再按需执行。

下面的示例创建了一个用于从数据库查询指定唱片的预处理语句。[`DB.Prepare`](https://pkg.go.dev/database/sql#DB.Prepare) 根据给定的 SQL 文本返回表示预处理语句的 [`sql.Stmt`](https://pkg.go.dev/database/sql#Stmt)。要执行语句，可以将 SQL 参数传给 `Stmt.Exec`、`Stmt.QueryRow` 或 `Stmt.Query`。

```
// AlbumByID retrieves the specified album.
func AlbumByID(id int) (Album, error) {
	// Define a prepared statement. You'd typically define the statement
	// elsewhere and save it for use in functions such as this one.
	stmt, err := db.Prepare("SELECT * FROM album WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var album Album

	// Execute the prepared statement, passing in an id value for the
	// parameter whose placeholder is ?
	err := stmt.QueryRow(id).Scan(&album.ID, &album.Title, &album.Artist, &album.Price, &album.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case of no rows returned.
		}
		return album, err
	}
	return album, nil
}
```

### 预处理语句的行为 {#behavior}

准备好的 [`sql.Stmt`](https://pkg.go.dev/database/sql#Stmt) 提供了常用的 `Exec`、`QueryRow` 和 `Query` 方法，用于执行语句。有关这些方法的使用方式，请参阅[查询数据](/doc/database/querying)和[执行不返回数据的 SQL 语句](/doc/database/change-data)。

由于 `sql.Stmt` 已经代表一条确定的 SQL 语句，它的 `Exec`、`QueryRow` 和 `Query` 方法只接受与占位符对应的 SQL 参数值，不再接受 SQL 文本。

根据使用场景，可以通过不同方式创建 `sql.Stmt`：

*   `DB.Prepare` 和 `DB.PrepareContext` 创建可在事务外独立执行的预处理语句，使用方式类似于 `DB.Exec` 和 `DB.Query`。
*   `Tx.Prepare`、`Tx.PrepareContext`、`Tx.Stmt` 和 `Tx.StmtContext` 创建用于特定事务的预处理语句。`Prepare` 和 `PrepareContext` 根据 SQL 文本定义语句；`Stmt` 和 `StmtContext` 则接收 `DB.Prepare` 或 `DB.PrepareContext` 的结果，也就是将非事务专用的 `sql.Stmt` 转换为当前事务使用的 `sql.Stmt`。
*   `Conn.PrepareContext` 通过表示专用连接的 `sql.Conn` 创建预处理语句。

使用完语句后，务必调用 `stmt.Close`，释放可能与之关联的数据库资源，例如底层连接。如果语句只是函数中的局部变量，使用 `defer stmt.Close()` 即可。

#### 创建预处理语句的方法 {#prepared_statement_functions}

<table id="prepared-statement-functions-list" class="DocTable">
    <thead>
        <tr class="DocTable-head">
            <th class="DocTable-cell" width="20%">函数</th>
            <th class="DocTable-cell">说明</th>
        </tr>
    </thead>
    <tbody>
        <tr class="DocTable-row">
            <td class="DocTable-cell">
                <code><a href="https://pkg.go.dev/database/sql#DB.Prepare">DB.Prepare</a></code><br />
                <code><a href="https://pkg.go.dev/database/sql#DB.PrepareContext">DB.PrepareContext</a></code>
            </td>
            <td class="DocTable-cell">准备一条语句，以便独立执行，或通过 Tx.Stmt 转换为事务内使用的预处理语句。</td>
        </tr>
        <tr class="DocTable-row">
            <td class="DocTable-cell">
                <code><a href="https://pkg.go.dev/database/sql#Tx.Prepare">Tx.Prepare</a></code><br />
                <code><a href="https://pkg.go.dev/database/sql#Tx.PrepareContext">Tx.PrepareContext</a></code><br />
                <code><a href="https://pkg.go.dev/database/sql#Tx.Stmt">Tx.Stmt</a></code><br />
                <code><a href="https://pkg.go.dev/database/sql#Tx.StmtContext">Tx.StmtContext</a></code>
            </td>
            <td class="DocTable-cell">准备一条用于特定事务的语句。更多信息请参阅
                <a href="/doc/database/execute-transactions">执行事务</a>。
            </td>
        </tr>
        <tr class="DocTable-row">
            <td class="DocTable-cell">
                <code><a href="https://pkg.go.dev/database/sql#Conn.PrepareContext">Conn.PrepareContext</a></code>
            </td>
            <td class="DocTable-cell">用于专用连接。更多信息请参阅 
                <a href="/doc/database/manage-connections">管理连接</a>。
            </td>
        </tr>
    </tbody>
</table>
