# 企业会议室预约服务

面向小团队的会议室排期：查看房间、查询空闲时段、创建预约、筛选预约、取消预约。重点是时间区间、冲突检测、并发争抢和文件保存。

项目只使用当前 `go-demo` 的知识和 Go 标准库。组织者姓名是业务字段，没有实现登录或权限；没有循环会议、邮件邀请、外部日历同步或多地域时区配置。

## 1. 立即体验

```sh
cd ~/Desktop/go-room-booking
go run ./cmd/demo
```

演示自动选择明天的时间并使用临时目录：查询会议室 → 创建预约 → 再预约同一时段得到 409 → 取消 → 再次预约成功 → 查询有效预约。可重复运行，不受今天日期或已有数据影响。

启动正式练习服务：

```sh
go run .
```

默认地址 `http://127.0.0.1:8092`。另开一个 macOS 终端：

```sh
curl -s http://127.0.0.1:8092/rooms
# macOS date 的 -v+1d 表示明天；这里用 UTC 时间，避免夏令时歧义。
DAY=$(date -u -v+1d '+%Y-%m-%d')
curl -s -X POST http://127.0.0.1:8092/bookings \
  -H 'Content-Type: application/json' \
  -d "{\"room_id\":\"room-a\",\"organizer\":\"小林\",\"title\":\"需求评审\",\"attendees\":3,\"start\":\"${DAY}T09:00:00Z\",\"end\":\"${DAY}T10:00:00Z\"}"
```

新数据中的预约 ID 是 `booking-000001`。完整输入和输出字段见 [接口示例](examples/requests.http)。重复提交创建请求会得到冲突，不会返回原预约；网络超时后先查询列表确认结果，不要假设创建失败。

```sh
curl -s 'http://127.0.0.1:8092/bookings?room_id=room-a&status=active'
curl -s -X POST http://127.0.0.1:8092/bookings/booking-000001/cancel
```

## 2. 业务规则

| 房间 | 名称 | 最大人数 |
| --- | --- | --- |
| room-a | 青竹会议室 | 4 |
| room-b | 云杉会议室 | 10 |
| room-c | 远程会议间 | 2 |

- 开始、结束时间使用 RFC3339，必须包含时区，例如 `2026-09-28T09:00:00+08:00`。服务内部和响应统一为 UTC。
- 时间需要对齐 UTC 的 15 分钟刻度，时长 15 分钟到 4 小时，开始时间晚于当前时刻且在未来 90 天内。
- 同一会议室的有效预约不能重叠。同一组织者即使换会议室，也不能同时主持两个会议；组织者按照去除首尾空白后的姓名精确比较。
- 区间采用 `[开始,结束)`：09:00–10:00 和 10:00–11:00 可以相邻。
- 重叠条件为 `新开始 < 旧结束 && 旧开始 < 新结束`，既覆盖部分重叠，也覆盖完全包含。
- 人数为 1..房间容量；组织者 1..40 个 rune，标题 1..100 个 rune。
- 状态为 active 或 canceled。active 表示预约未被取消，历史已结束会议仍保留 active，不启动后台任务修改状态。
- 只有尚未开始的预约可以取消；重复取消成功。取消后释放原时段，历史记录继续保留。
- 空闲查询只检查房间，不检查组织者和人数，也不会预占资源。创建时会再次校验全部条件，检查和保存由同一把锁保护。

## 3. API

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /healthz | 健康检查和构建版本 |
| GET | /rooms | 固定房间列表 |
| GET | /rooms/{id}/availability?start=...&end=... | 返回 `{"available":true/false}` |
| POST | /bookings | 创建预约，成功 201 |
| GET | /bookings?room_id=room-a&status=active&page=1&size=20 | 过滤和分页，按开始时间、ID 排序 |
| GET | /bookings/{id} | 预约详情 |
| POST | /bookings/{id}/cancel | 取消，不需要请求体，成功 200 |

分页 page 默认 1，size 默认 20、范围 1..100，超出页数返回空列表。写请求（创建）必须为单个 JSON 对象，最多 16 KiB，拒绝未知字段、null、多个值。

查询空闲示例（沿用上面的 DAY）：

```sh
curl -sG http://127.0.0.1:8092/rooms/room-a/availability \
  --data-urlencode "start=${DAY}T09:00:00Z" \
  --data-urlencode "end=${DAY}T10:00:00Z"
```

`+08:00` 中的 `+` 放进 URL 时要编码；`--data-urlencode` 可以处理。400 表示输入问题，404 表示资源不存在，409 表示排期或状态冲突，415 表示请求类型错误，413 表示请求过大，500 表示保存失败。业务错误为 `{"error":"原因"}`；路由默认的 404/405 为标准库文本。

## 4. 代码阅读路线

1. `internal/app/model.go`：结构体嵌入让 Booking 复用 BookingInput 字段，对应第 11 章；状态定义对应第 12 章。
2. `service.go` 中的 `parseInterval`、`overlaps`、`validateInput`：先读纯规则，再读 Create 和 Cancel。
3. `Service.now` 是一个 `func() time.Time` 函数字段（第 08 章），运行时用 time.Now，测试时固定时间；测试不依赖当天日期或 Sleep。
4. `Saver` 小接口对应第 13 章，`clone → Save → state = next` 配合第 17 章锁，确保保存失败不占用时段。
5. `internal/filedata` 对应第 20、21、23 章；持久化文件损坏会报错，不会清空重建。
6. `http.go` 和 `internal/web` 对应第 19、27 章的分页及路由；`main.go` 对应第 30 章的配置和优雅退出。
7. `service_test.go`、`http_test.go` 对应第 24 章，理解为什么并发测试与时间边界测试要分别编写。

`time.Parse/Format` 用固定布局表达时间格式，是标准库 time 的业务应用；这里不引入日历框架。源代码中会解释 UTC 归一化和区间比较。

## 5. 验证与后续练习

```sh
go test -race -cover ./...
go vet ./...
go test -v ./internal/app -run TestConcurrentBooking
```

测试覆盖 20 人并发争抢同一房间只能成功一人、不同时区表示同一时段、首尾相接、容量限制、取消释放、过去时间、90 天限制、保存失败和重启恢复。

可以继续做三个范围内的练习：

1. 给会议室增加投影仪字段并添加查询筛选。
2. 新增“某个组织者的预约列表”，保持分页与排序约定。
3. 用 map、切片、循环统计每个房间已经预约的分钟数，排除 canceled。

详细配置、备份与构建说明见 [运行与部署](docs/运行与部署.md)。单个数据文件只供一个进程使用；每次保存完整快照，适合当前阶段的小数据项目。
