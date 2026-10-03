# 电商订单与库存服务

为小型商家管理商品、可售库存和订单履约。重点是把真实业务规则写完整：服务端定价、并发防超卖、重复请求处理、取消退库存、失败时保持数据一致。

这是基于 `go-demo` 当前 30 章知识的完整可运行练习。只用标准库；订单表示商家已接受的履约单，不处理付款、退款、物流平台和用户登录。

## 1. 先运行一个完整流程

```sh
cd ~/Desktop/go-shop-orders
go run ./cmd/demo
```

演示会使用临时目录和空闲端口，完成：查询商品 → 下单 → 重复提交同一下单请求 → 取消 → 尝试非法发货 → 再次查询库存。每一步核对 HTTP 状态码；正常结束时显示“演示通过”。不会改动正式运行的 `data`。

随后运行长期服务：

```sh
go run .
```

默认地址 `http://127.0.0.1:8091`。另开终端：

```sh
curl -s http://127.0.0.1:8091/products
curl -s -X POST http://127.0.0.1:8091/orders \
  -H 'Content-Type: application/json' \
  -d '{"request_id":"demo-001","customer":"小林","items":[{"product_id":"p-1001","quantity":2}]}'
```

新数据文件中的首笔订单为 `ord-000001`，商品金额是 **3980 分**，`p-1001` 库存从 **20** 变为 **18**。已有数据时，请使用实际返回的订单 ID。

```sh
curl -s http://127.0.0.1:8091/orders/ord-000001
curl -s -X PATCH http://127.0.0.1:8091/orders/ord-000001/status \
  -H 'Content-Type: application/json' -d '{"status":"canceled"}'
curl -s http://127.0.0.1:8091/products
```

重复取消不会再次增加库存。取消后重启服务，订单状态和库存仍然保留。

## 2. 需求与业务规则

- 初始商品：`p-1001` 办公笔记本，1990 分/本、20 本；`p-1002` 桌面收纳盒，3500 分/个、10 个。
- 金额用 `int64` 派生的 `Cents`，所有单价和总额单位都是分，不用浮点数计算。
- 下单提交商品 ID 和数量；商品名称、单价及总额由服务端生成，并保存当时快照。
- 每单 1..100 种商品；相同商品不允许重复行，每项 1..1000 件；任何一项不存在或库存不足，整单失败。
- `request_id` 为 1..64 位英文字母、数字、`-`、`_`。相同编号和相同内容返回原订单（仍为 201），不会再次扣库存；相同编号但内容不同返回 409。明细行顺序不影响比较。
- 已取消订单的 `request_id` 也不能再创建新订单；重新购买必须使用新编号。客户名称在去除首尾空白后比较。
- 状态仅允许 `confirmed → shipped` 或 `confirmed → canceled`。重复到达同一目标状态成功；终态间不能互转。
- 创建与补货都限制库存规模；补货时为待履约订单的取消保留容量。取消不会导致库存超过上限。
- 数据检查、库存变更、订单写入都在同一把锁下完成；保存失败时不替换内存状态。

## 3. API

写请求使用 `Content-Type: application/json`，请求体最多 16 KiB，拒绝未知字段、`null`、数组和多个 JSON。业务错误返回 `{"error":"原因"}`；路由本身产生的 404/405 使用标准库默认文本。

| 方法 | 路径 | 行为 |
| --- | --- | --- |
| GET | /healthz | 健康检查与版本 |
| GET | /products?page=1&size=20 | 按商品 ID 排序分页 |
| POST | /products | 新增商品，成功 201 |
| POST | /products/{id}/restock | 补货，成功 200 |
| POST | /orders | 创建或返回重复订单，成功 201 |
| GET | /orders?status=confirmed&page=1&size=20 | 状态筛选、按订单 ID 排序分页 |
| GET | /orders/{id} | 查询订单详情 |
| PATCH | /orders/{id}/status | 发货或取消 |

分页默认 page=1、size=20，size 为 1..100，越界页返回空 `items`，保留总数。空列表返回 `[]`。`status` 可为 confirmed、shipped、canceled，不填查询全部。

新增商品示例：

```json
{"id":"p-1003","name":"签字笔","price_cents":250,"stock":100}
```

价格为 1..100000000 分，库存为 0..1000000，名称 1..80 个 rune。补货示例 `{"quantity":10}`，单次为 1..100000。客户名称为 1..80 个 rune。不存在资源返回 404，冲突返回 409，非法输入返回 400，保存失败返回 500。

请求模板见 [examples/requests.http](examples/requests.http)，无需安装编辑器插件，也可以直接使用上述 curl。

## 4. 按这个顺序阅读

| 文件 | 阅读目的 | 原学习章节 |
| --- | --- | --- |
| internal/app/model.go | 定义商品、订单、金额、状态 | 03、04、11、12 |
| internal/app/service.go | 业务校验、快照、锁、错误、保存接口 | 05–10、13、17、18 |
| internal/filedata/file.go | 读取和完整保存 JSON | 20、21、23 |
| internal/app/http.go | 把业务方法映射为 HTTP 接口 | 27 |
| internal/web/web.go | 严格 JSON、泛型分页、服务退出 | 16、18、19、27、30 |
| main.go | 显式组装依赖与配置 | 01、10、30 |
| internal/app/*_test.go | 业务边界、并发、真实 HTTP 测试 | 14、17、24、27 |
| cmd/demo/main.go | 客户端发请求并核对结果 | 27 |

调用链：`HTTP 请求 → handler 解码 → Service 校验和修改候选副本 → Saver 保存 → 替换内存 → JSON 响应`。

`internal` 是 Go 的包组织约定，用来限制包被项目之外导入；核心仍是第 01 章的 module/package 和导出规则。这里按业务、HTTP、文件职责分文件，不需要学习复杂分层框架。

## 5. 测试与练习

```sh
go test -race -cover ./...
go vet ./...
go test -v ./internal/app -run TestConcurrentStockAndRequestDeduplication
```

已有测试会核对 40 个请求抢 20 件库存、20 个相同请求只创建一单、整单失败、重复取消、保存失败回滚、重启恢复、JSON 和分页边界。命令行入口与演示通过实际运行另外验证；单元覆盖率不代表端到端保证。

建议读懂后自行完成：

1. 给商品增加可选分类字段，扩展创建校验和列表筛选。
2. 增加按客户名筛选订单，先写空输入和无匹配测试。
3. 用已经学过的切片与循环统计各状态的订单数与金额；取消单不能计入履约金额。

以上是后续练习，当前交付的主业务已经实现，无需补 TODO 才能运行。

配置、构建、数据备份和故障排查见 [运行与部署](docs/运行与部署.md)。
