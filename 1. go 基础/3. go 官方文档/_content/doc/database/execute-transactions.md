<!--{
  "Title": "执行事务"
}-->

可以通过表示事务的 [`sql.Tx`](https://pkg.go.dev/database/sql#Tx) 执行数据库事务。除表示事务特有语义的 `Commit` 和 `Rollback` 方法外，`sql.Tx` 还提供了执行常见数据库操作所需的各种方法。调用 `DB.Begin` 或 `DB.BeginTx` 可以获取 `sql.Tx`。

[数据库事务](https://en.wikipedia.org/wiki/Database_transaction)将多个操作组织起来，共同完成一个整体目标。这些操作要么全部成功，要么全部不生效，两种情况下都能保持数据完整性。事务流程通常包括：

1. 开始事务。
2. 执行一组数据库操作。
3. 如果没有出错，提交事务，使数据库更改生效。
4. 如果发生错误，回滚事务，撤销事务中的更改。

`sql` 包提供了开始和结束事务的方法，以及执行中间数据库操作的方法，分别对应上述流程。

*   开始事务。

    [`DB.Begin`](https://pkg.go.dev/database/sql#DB.Begin) 或 [`DB.BeginTx`](https://pkg.go.dev/database/sql#DB.BeginTx) 开始一个新事务，并返回表示该事务的 `sql.Tx`。
*   执行数据库操作。

    使用 `sql.Tx` 可以在同一个连接上执行一系列查询或更新操作。为此，`Tx` 导出了以下方法：

    *   [`Exec`](https://pkg.go.dev/database/sql#Tx.Exec) 和 [`ExecContext`](https://pkg.go.dev/database/sql#Tx.ExecContext)：通过 `INSERT`、`UPDATE` 和 `DELETE` 等 SQL 语句修改数据库。

        更多信息请参阅[执行不返回数据的 SQL 语句](/doc/database/change-data)。

    *   [`Query`](https://pkg.go.dev/database/sql#Tx.Query)、[`QueryContext`](https://pkg.go.dev/database/sql#Tx.QueryContext)、[`QueryRow`](https://pkg.go.dev/database/sql#Tx.QueryRow) 和 [`QueryRowContext`](https://pkg.go.dev/database/sql#Tx.QueryRowContext)：执行返回数据行的操作。

        更多信息请参阅[查询数据](/doc/database/querying)。

    *   [`Prepare`](https://pkg.go.dev/database/sql#Tx.Prepare)、[`PrepareContext`](https://pkg.go.dev/database/sql#Tx.PrepareContext)、[`Stmt`](https://pkg.go.dev/database/sql#Tx.Stmt) 和 [`StmtContext`](https://pkg.go.dev/database/sql#Tx.StmtContext)：定义预处理语句。

        更多信息请参阅[使用预处理语句](/doc/database/prepared-statements)。

*   通过以下方式之一结束事务：
    *   使用 [`Tx.Commit`](https://pkg.go.dev/database/sql#Tx.Commit) 提交事务。

        如果 `Commit` 成功（返回 nil 错误），所有查询结果都被确认为有效，已执行的更新会作为一次原子更改应用到数据库。如果 `Commit` 失败，应丢弃该 `Tx` 中所有 `Query` 和 `Exec` 操作的结果，将其视为无效。
    *   使用 [`Tx.Rollback`](https://pkg.go.dev/database/sql#Tx.Rollback) 回滚事务。

        即使 `Tx.Rollback` 失败，该事务也不再有效，并且不会提交到数据库。

### 最佳实践 {#best_practices}

事务有时涉及复杂的语义和连接管理。遵循以下实践，有助于正确处理这些问题。

*   使用本节介绍的 API 管理事务。不要直接执行 `BEGIN`、`COMMIT` 等与事务相关的 SQL 语句，否则可能让数据库处于难以预料的状态，在并发程序中尤其如此。
*   使用事务时，注意不要同时直接调用不属于事务的 `sql.DB` 方法。这些方法会在事务之外执行，导致代码看到不一致的数据库状态，甚至引发死锁。

### 示例 {#example}

下面的示例通过事务为客户创建一个唱片订单。代码依次执行以下操作：

1. 开始事务。
2. 使用 `defer` 延迟回滚事务。如果事务成功，会在函数退出前完成提交，此时延迟执行的回滚不会产生作用。如果事务失败，尚未提交，函数退出时便会执行回滚。
3. 检查客户所订购唱片的库存是否充足。
4. 如果充足，更新库存数量，扣除订单中的唱片数量。
5. 创建新订单，并获取为该订单生成的 ID，以返回给客户端。
6. 提交事务并返回 ID。

示例使用了接受 `context.Context` 参数的 `Tx` 方法，因此当运行时间过长或客户端连接关闭时，可以取消函数执行，包括其中的数据库操作。更多信息请参阅[取消正在进行的操作](/doc/database/cancel-operations)。

```
// CreateOrder creates an order for an album and returns the new order ID.
func CreateOrder(ctx context.Context, albumID, quantity, custID int) (orderID int64, err error) {

	// Create a helper function for preparing failure results.
	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("CreateOrder: %v", err)
	}

	// Get a Tx for making transaction requests.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// Confirm that album inventory is enough for the order.
	var enough bool
	if err = tx.QueryRowContext(ctx, "SELECT (quantity >= ?) from album where id = ?",
		quantity, albumID).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			return fail(fmt.Errorf("no such album"))
		}
		return fail(err)
	}
	if !enough {
		return fail(fmt.Errorf("not enough inventory"))
	}

	// Update the album inventory to remove the quantity in the order.
	_, err = tx.ExecContext(ctx, "UPDATE album SET quantity = quantity - ? WHERE id = ?",
		quantity, albumID)
	if err != nil {
		return fail(err)
	}

	// Create a new row in the album_order table.
	result, err := tx.ExecContext(ctx, "INSERT INTO album_order (album_id, cust_id, quantity, date) VALUES (?, ?, ?, ?)",
		albumID, custID, quantity, time.Now())
	if err != nil {
		return fail(err)
	}
	// Get the ID of the order item just created.
	orderID, err = result.LastInsertId()
	if err != nil {
		return fail(err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	// Return the order ID.
	return orderID, nil
}
```
