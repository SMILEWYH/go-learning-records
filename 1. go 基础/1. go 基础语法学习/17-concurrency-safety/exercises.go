//go:build ignore

package main

// 练习 1：库存扣减。多个订单会同时购买同一个 SKU。
// 要求：定义 Inventory，使用 Mutex 保护库存；检查库存和扣减必须在同一把锁内。
// 建议：func (i *Inventory) Reserve(count int) bool；func (i *Inventory) Remaining() int。
// 初始库存 10，让 20 个 goroutine 各扣 1，WaitGroup 等待结束：恰好 10 次成功，剩余 0。
// 边界：扣减 0 或负数失败，库存不变；库存不足失败，库存不能变为负数。
// 提示：不要无锁地共享“成功次数”，可使用 channel 收集每个 goroutine 的结果。

// 练习 2：用户权限缓存。缓存的每个用户对应一个 []string 权限列表。
// 要求：用 map + RWMutex；Set 和 Get 都复制切片，避免调用者从锁外修改内部数据。
// 建议：func (c *PermissionCache) Set(userID string, permissions []string)
//
//	func (c *PermissionCache) Get(userID string) ([]string, bool)
//
// 常规：存入 ["read"] 后读出 ["read"], true；同时多次 Set/Get 不出现竞争。
// 边界：修改 Set 的输入切片或 Get 的返回切片，缓存仍为 ["read"]；不存在返回 nil, false。

// 练习 3：Webhook 事件去重。这里只实现进程内演示，重启后记录不保留。
// 要求：使用 sync.Map.LoadOrStore 实现“记录事件 ID 并判断是不是第一次”。
// 建议：func (d *EventDeduplicator) FirstSeen(eventID string) bool。
// 常规：20 个 goroutine 同时传 "evt-42"，恰好一次得到 true；新 ID "evt-43" 得到 true。
// 边界：空 ID 返回 false 且不存入；重复 ID 返回 false。不能把 Load 与 Store 拆开。
// 思考：如果事件处理失败，应如何设计重试？本题只判断首次出现，不代表业务处理成功。

func main() {
}
