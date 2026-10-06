<!--{
  "Title": "查询数据"
}-->

执行会返回数据的 SQL 语句时，可以使用 `database/sql` 包提供的 `Query` 系列方法。这些方法返回 `Row` 或 `Rows`，通过其 `Scan` 方法可以将数据复制到变量中。例如，执行 `SELECT` 语句时就会使用这些方法。

执行不返回数据的语句时，应使用 `Exec` 或 `ExecContext`。更多信息请参阅[执行不返回数据的语句](/doc/database/change-data)。

`database/sql` 包提供了两种查询方式：

*   **查询单行数据**：`QueryRow` 最多从数据库返回一行 `Row`。更多信息请参阅[查询单行数据](#single_row)。
*   **查询多行数据**：`Query` 将所有匹配行作为 `Rows` 结构体返回，代码可以遍历这些行。更多信息请参阅[查询多行数据](#multiple_rows)。

如果需要反复执行相同的 SQL 语句，可以考虑使用预处理语句。更多信息请参阅[使用预处理语句](/doc/database/prepared-statements)。

**警告：** 不要使用 `fmt.Sprintf` 等字符串格式化函数拼接 SQL 语句，否则可能引入 SQL 注入风险。更多信息请参阅[避免 SQL 注入风险](/doc/database/sql-injection)。

### 查询单行数据 {#single_row}

`QueryRow` 最多获取一行数据库记录，例如根据唯一 ID 查询数据。如果查询返回多行，`Scan` 方法只保留第一行，丢弃其余行。

`QueryRowContext` 与 `QueryRow` 的工作方式相同，但额外接受一个 `context.Context` 参数。更多信息请参阅[取消正在进行的操作](/doc/database/cancel-operations)。

下面的示例通过查询判断库存是否足够完成购买。SQL 语句在库存充足时返回 `true`，否则返回 `false`。[`Row.Scan`](https://pkg.go.dev/database/sql#Row.Scan) 通过指针将这个布尔返回值复制到 `enough` 变量中。

```
func canPurchase(id int, quantity int) (bool, error) {
	var enough bool
	// Query for a value based on a single row.
	if err := db.QueryRow("SELECT (quantity >= ?) from album where id = ?",
		quantity, id).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("canPurchase %d: unknown album", id)
		}
		return false, fmt.Errorf("canPurchase %d: %v", id, err)
	}
	return enough, nil
}
```

**注意：** 预处理语句中参数占位符的形式取决于所用的数据库管理系统和驱动。例如，Postgres 的 [pq 驱动](https://pkg.go.dev/github.com/lib/pq) 要求使用 `$1` 这样的占位符，而不是 `?`。

#### 处理错误 {#single_row_errors}

`QueryRow` 本身不返回错误，查询和扫描过程中的错误统一由 `Scan` 报告。如果查询没有找到任何行，`Scan` 会返回 [`sql.ErrNoRows`](https://pkg.go.dev/database/sql#ErrNoRows)。

#### 返回单行数据的方法 {#single_row_functions}

<table id="single-row-functions-list" class="DocTable">
  <thead>
    <tr class="DocTable-head">
      <th class="DocTable-cell" width="20%">函数</th>
      <th class="DocTable-cell">说明</th>
    </tr>
  </thead>
  <tbody>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#DB.QueryRow">DB.QueryRow</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#DB.QueryRowContext">DB.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">独立执行单行查询。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Tx.QueryRow">Tx.QueryRow</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Tx.QueryRowContext">Tx.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">在一个事务中执行单行查询。更多信息请参阅
        <a href="/doc/database/execute-transactions">执行事务</a>。
      </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Stmt.QueryRow">Stmt.QueryRow</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Stmt.QueryRowContext">Stmt.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">使用预处理语句执行单行查询。更多信息请参阅 <a href="/doc/database/prepared-statements">使用预处理语句</a>。
      </td>
    </tr>
    <tr class="DocTable-row">
        <td class="DocTable-cell">
  <code><a href="https://pkg.go.dev/database/sql#Conn.QueryRowContext">Conn.QueryRowContext</a></code>
      </td>
      <td class="DocTable-cell">用于专用连接。更多信息请参阅
        <a href="/doc/database/manage-connections">管理连接</a>。
      </td>
    </tr>
  </tbody>
</table>

### 查询多行数据 {#multiple_rows}

可以使用 `Query` 或 `QueryContext` 查询多行数据，它们返回表示查询结果的 `Rows`。代码通过 [`Rows.Next`](https://pkg.go.dev/database/sql#Rows.Next) 遍历返回的行，并在每次迭代时调用 `Scan`，将列值复制到变量中。

`QueryContext` 与 `Query` 的工作方式相同，但额外接受一个 `context.Context` 参数。更多信息请参阅[取消正在进行的操作](/doc/database/cancel-operations)。

下面的示例查询指定艺人的唱片，结果保存在 `sql.Rows` 中。代码使用 [`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan)，将列值复制到各个指针所指向的变量中。

```
func albumsByArtist(artist string) ([]Album, error) {
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", artist)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An album slice to hold data from returned rows.
	var albums []Album

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist,
			&alb.Price, &alb.Quantity); err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}
	if err = rows.Err(); err != nil {
		return albums, err
	}
	return albums, nil
}
```

注意延迟执行的 [`rows.Close`](https://pkg.go.dev/database/sql#Rows.Close) 调用。无论函数以哪种方式返回，它都能释放结果行持有的资源。完整遍历所有行也会隐式关闭结果集，不过使用 `defer` 能确保各种情况下都会关闭 `rows`，因此更稳妥。

**注意：** 预处理语句中参数占位符的形式取决于所用的数据库管理系统和驱动。例如，Postgres 的 [pq 驱动](https://pkg.go.dev/github.com/lib/pq) 要求使用 `$1` 这样的占位符，而不是 `?`。

#### 处理错误 {#multiple_rows_errors}

遍历查询结果后，务必检查 `sql.Rows` 是否记录了错误。如果查询在迭代过程中失败，代码需要通过这种方式发现错误。

#### 返回多行数据的方法 {#multiple_rows_functions}

<table id="multiple-row-functions-list" class="DocTable">
  <thead>
    <tr class="DocTable-head">
      <th class="DocTable-cell" width="20%">函数</th>
      <th class="DocTable-cell">说明</th>
    </tr>
  </thead>
  <tbody>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#DB.Query">DB.Query</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#DB.QueryContext">DB.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">独立执行查询。</td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Tx.Query">Tx.Query</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Tx.QueryContext">Tx.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">在一个事务中执行查询。更多信息请参阅
        <a href="/doc/database/execute-transactions">执行事务</a>。
      </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Stmt.Query">Stmt.Query</a></code><br />
        <code><a href="https://pkg.go.dev/database/sql#Stmt.QueryContext">Stmt.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">使用预处理语句执行查询。更多信息请参阅
        <a href="/doc/database/prepared-statements">使用预处理语句</a>。
    </td>
    </tr>
    <tr class="DocTable-row">
      <td class="DocTable-cell">
        <code><a href="https://pkg.go.dev/database/sql#Conn.QueryContext">Conn.QueryContext</a></code>
      </td>
      <td class="DocTable-cell">用于专用连接。更多信息请参阅
        <a href="/doc/database/manage-connections">管理连接</a>。
      </td>
    </tr>
  </tbody>
</table>

### 处理可为 NULL 的列值 {#nullable_columns}

`database/sql` 包提供了几种特殊类型。当列值可能为 NULL 时，可以将这些类型的值传给 `Scan`。每种类型都包含一个 `Valid` 字段，表示值是否非 NULL，并包含一个用于保存实际值的字段。

下面的示例查询客户姓名。如果姓名为 NULL，代码就使用另一个值作为应用中的替代值。

```
var s sql.NullString
err := db.QueryRow("SELECT name FROM customer WHERE id = ?", id).Scan(&s)
if err != nil {
	log.Fatal(err)
}

// Find customer name, using placeholder if not present.
name := "Valued Customer"
if s.Valid {
	name = s.String
}
```

有关各个类型的详细说明，请参阅 `sql` 包参考文档：

*    [`NullBool`](https://pkg.go.dev/database/sql#NullBool)
*    [`NullFloat64`](https://pkg.go.dev/database/sql#NullFloat64)
*    [`NullInt32`](https://pkg.go.dev/database/sql#NullInt32)
*    [`NullInt64`](https://pkg.go.dev/database/sql#NullInt64)
*    [`NullString`](https://pkg.go.dev/database/sql#NullString)
*    [`NullTime`](https://pkg.go.dev/database/sql#NullTime)

### 获取列数据 {#column_data}

遍历查询返回的行时，可以使用 `Scan` 将一行中的列值复制到 Go 值中，具体说明参见 [`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan) 参考文档。

所有驱动都支持一组基本的数据转换，例如将 SQL `INT` 转换为 Go `int`。部分驱动还支持额外的转换，具体请参阅对应驱动的文档。

`Scan` 会将数据库列类型转换为相近的 Go 类型。例如，它可以将 SQL 的 `CHAR`、`VARCHAR` 和 `TEXT` 转换为 Go 的 `string`。此外，只要列值适合，也可以转换为其他 Go 类型。例如，如果某个 `VARCHAR` 列始终存储数字，可以指定 Go 的 `int` 等数值类型接收结果，`Scan` 会使用 `strconv.Atoi` 完成转换。

有关 `Scan` 支持的类型转换，请参阅 [`Rows.Scan`](https://pkg.go.dev/database/sql#Rows.Scan) 参考文档。

### 处理多个结果集 {#multiple_result_sets}

如果数据库操作可能返回多个结果集，可以通过 [`Rows.NextResultSet`](https://pkg.go.dev/database/sql#Rows.NextResultSet) 获取它们。例如，发送一段分别查询多个表的 SQL，并且每个查询都返回独立结果集时，就可以使用这个方法。

`Rows.NextResultSet` 会准备下一个结果集，使后续的 `Rows.Next` 调用读取该结果集的第一行。它返回一个布尔值，表示是否存在下一个结果集。

下面的示例使用 `DB.Query` 执行两条 SQL 语句。第一个结果集来自第一个查询，包含 `album` 表中的所有行；下一个结果集来自第二个查询，包含 `song` 表中的行。

```
rows, err := db.Query("SELECT * from album; SELECT * from song;")
if err != nil {
	log.Fatal(err)
}
defer rows.Close()

// Loop through the first result set.
for rows.Next() {
	// Handle result set.
}

// Advance to next result set.
rows.NextResultSet()

// Loop through the second result set.
for rows.Next() {
	// Handle second set.
}

// Check for any error in either result set.
if err := rows.Err(); err != nil {
	log.Fatal(err)
}
```
