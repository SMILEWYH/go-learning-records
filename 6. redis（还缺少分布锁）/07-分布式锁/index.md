# 9. Redis的分布式锁

分布式锁：多个进程、多个服务实例访问同一个共享资源时，通过共同访问一个锁服务，协调谁可以进入临界区。

临界区：需要避免并发冲突的那一段业务代码，例如读取订单状态、判断能否更新、再更新订单。

> 学习顺序：为什么需要锁 → SET NX EX 加锁 → 唯一标识 → Lua 解锁 → 过期与续期 → 故障处理 → Go 实战。

## 9.1 为什么需要分布式锁

### 9.1.1 本地锁和分布式锁

Go 中的 `sync.Mutex` 可以协调同一个进程内、使用同一把锁的多个 goroutine。

如果服务部署了多个实例，每个实例都有自己的内存和 Mutex，那么实例 A 的 Mutex 无法阻止实例 B 执行业务。

~~~text
服务实例 A：goroutine → A 的 Mutex → 共享数据库
服务实例 B：goroutine → B 的 Mutex → 共享数据库

A 和 B 的 Mutex 不是同一把锁，两个实例仍然可能同时修改数据。
~~~

使用 Redis 时，多个实例竞争同一个锁 key，由 Redis 决定谁能够获取它。

~~~text
服务实例 A ─┐
           ├─→ Redis 中的 lock:order:1001 → 获得锁的实例执行业务
服务实例 B ─┘
~~~

**锁 key 相同，才是在竞争同一个资源。** 所有需要协调的业务入口都必须遵守相同的加锁规则。

### 9.1.2 常见使用场景

- 缓存重建：热点缓存过期后，减少同时查询数据库的请求。
- 定时任务：多个服务实例部署了同一个任务，协调某一轮由哪个实例执行。
- 资源更新：协调针对同一个订单、商品等资源的并发操作。

> 分布式锁不等于业务幂等，也不能单独保证任务只执行一次。涉及业务正确性时，还要设计数据库约束、条件更新或其他校验机制。

### 9.1.3 一把锁需要考虑什么

| 要求 | 含义 |
| --- | --- |
| 互斥 | 在锁有效且系统假设成立时，同一个资源只允许一个持锁者进入临界区 |
| 避免永久占用 | 持锁进程崩溃后，其他客户端最终仍有机会获取锁 |
| 正确释放 | 只能释放自己这一轮获取的锁 |
| 有限等待 | 获取锁和访问 Redis 都要有超时，不能无限阻塞 |
| 故障处理 | 能处理网络超时、锁过期、续期失败、Redis 切换等情况 |

## 9.2 使用 SET NX EX 加锁

### 9.2.1 命令含义

~~~shell
SET lock:order:1001 token-A NX EX 30
~~~

| 参数 | 含义 |
| --- | --- |
| `lock:order:1001` | 锁的 key，表示订单 1001 对应的锁 |
| `token-A` | 本次加锁的唯一标识，实际程序应生成随机 token |
| `NX` | 只有 key 不存在时，才允许设置 |
| `EX 30` | key 在 30 秒后过期 |
| `PX 30000` | 另一种写法，过期时间单位为毫秒 |

命令格式和返回值参见 [SET 官方文档](https://redis.io/docs/latest/commands/set/)。

~~~shell
127.0.0.1:6379> SET lock:order:1001 token-A NX EX 30
OK
127.0.0.1:6379> SET lock:order:1001 token-B NX EX 30
(nil)
127.0.0.1:6379> GET lock:order:1001
"token-A"
127.0.0.1:6379> TTL lock:order:1001
(integer) 27
~~~

以上操作假设在锁过期前连续执行；TTL 的具体值取决于操作耗时。

- 返回 `OK`：本次设置成功。
- 返回 `(nil)`：key 已存在，本次没有拿到锁。
- 网络错误或超时：结果可能不确定，不能当作加锁成功。

注意：独立的 `SETNX` 命令返回 `1` 或 `0`；这里的 `SET ... NX EX ...` 返回 `OK` 或空回复，不要混淆。

### 9.2.2 为什么必须同时设置过期时间

下面的两条命令存在风险：

~~~shell
SETNX lock:order:1001 token-A
EXPIRE lock:order:1001 30
~~~

如果第一条执行成功，程序在执行第二条之前崩溃，锁就没有过期时间，可能一直占用资源。

**使用一条 SET 命令，把“仅不存在时创建”和“设置过期时间”一起完成。** 普通 Pipeline 只是批量发送命令，不能消除这个问题。

### 9.2.3 锁粒度

~~~text
lock:order:1001       订单 1001 的锁
lock:order:1002       订单 1002 的锁
lock:cache:product:8  商品 8 的缓存重建锁
~~~

不同资源可以使用不同的锁，减少不必要的等待。如果所有订单都使用 `lock:order`，订单之间也会相互阻塞。

锁 key 应带上项目或业务前缀，避免不同业务意外共用。锁 key 和实际缓存数据 key 也应分开。

## 9.3 唯一标识和误删锁

### 9.3.1 为什么不能直接 DEL

~~~text
第 0 秒：A 获取锁，TTL = 30 秒。
第 30 秒：A 的业务还没有完成，锁已经过期。
第 31 秒：B 获取同一个 key 的锁，开始执行业务。
第 35 秒：A 业务完成，执行 DEL，把 B 的锁删除。
第 36 秒：C 又获取锁，此时 B 和 C 可能同时执行业务。
~~~

> 释放锁前，需要确认当前锁仍然属于自己这一轮操作。

### 9.3.2 为什么 value 不能固定为 1

如果 A、B 加锁时都写入 `1`，A 无法分辨当前 key 是自己的旧锁，还是 B 后来创建的新锁。

因此每次新的加锁操作都应生成独立的随机 token，例如使用 UUID 或安全随机字节。

- key：标识需要保护的资源，竞争者之间相同。
- value：标识本次加锁操作，不同轮次应不同。

不能只使用固定的进程 ID、机器 IP 或 goroutine 编号，同一进程也会多次加锁。

## 9.4 使用 Lua 安全释放锁

### 9.4.1 为什么 GET 判断后再 DEL 仍然不安全

~~~text
A 执行 GET，确认 value 是自己的 token。
A 暂停，锁恰好过期。
B 获取锁，写入 B 的 token。
A 恢复，执行 DEL，仍然删除了 B 的锁。
~~~

问题在于：客户端发出的 GET 和 DEL 是两个独立操作，中间可能发生变化。

### 9.4.2 Lua 解锁脚本

~~~lua
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
end
return 0
~~~

| 写法 | 含义 |
| --- | --- |
| `redis.call(...)` | 在脚本内调用 Redis 命令 |
| `KEYS[1]` | 传入的第一个 key，Lua 数组下标从 1 开始 |
| `ARGV[1]` | 传入的第一个普通参数，此处为自己的 token |
| `return 1` | 删除成功，来自 DEL 的返回值 |
| `return 0` | key 不存在或 token 不一致，本次没有删除 |

Redis 在执行脚本期间，不会让其他客户端的命令插入脚本中间，因此可以原子地完成“比较并删除”。脚本应尽量短，避免长时间阻塞 Redis。[Lua 官方文档](https://redis.io/docs/latest/develop/programmability/eval-intro/)

> 这里的原子执行强调不被其他命令穿插，不代表脚本出错会自动回滚之前的写入。Redis 内的 Lua 也不能把外部 MySQL 操作包成同一个原子操作。

### 9.4.3 在 redis-cli 中调用

~~~shell
127.0.0.1:6379> SET lock:demo token-A NX EX 30
OK
127.0.0.1:6379> EVAL 'if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0' 1 lock:demo token-B
(integer) 0
127.0.0.1:6379> GET lock:demo
"token-A"
127.0.0.1:6379> EVAL 'if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0' 1 lock:demo token-A
(integer) 1
~~~

`EVAL 脚本 1 lock:demo token-A` 中，`1` 表示后面有一个 key，剩余参数进入 `ARGV`。以上示例需在 30 秒内完成。

### 9.4.4 新版本的条件删除命令

Redis 8.4 起支持下面的条件删除方式：

~~~shell
DELEX lock:demo IFEQ token-A
~~~

它也能完成“值相同才删除”。较早版本使用 Lua；学习 Lua 仍有助于理解续期和其他条件更新。[DELEX 官方文档](https://redis.io/docs/latest/commands/delex/)

## 9.5 锁过期与自动续期

### 9.5.1 锁过期不代表业务停止

Redis 的过期机制只负责让锁 key 失效，不会终止客户端的 goroutine，也不会撤销已经发出的数据库请求。

~~~text
A 获取锁 → A 执行业务较慢 → 锁过期 → B 获取锁
                                      ↓
                              A、B 的业务可能重叠
~~~

唯一 token 和 Lua 解锁防止的是误删锁，不能解决这里的业务重叠。

### 9.5.2 TTL 如何选择

- 太短：正常业务还没结束，锁就可能过期。
- 太长：持锁者崩溃后，其他请求要等待较久。
- 根据业务耗时、网络延迟和暂停风险选择，并留出余量。
- 获取锁的网络往返也会消耗有效期，收到成功回复不等于还剩完整 TTL。

带过期时间的锁也称为租约（lease），客户端拥有的是一定时间范围内的使用权。

### 9.5.3 续期与看门狗

看门狗：客户端在业务执行期间定期检查并延长锁的有效期。它是客户端实现的机制，Redis 不会自动判断业务是否完成。

例如 TTL 为 30 秒，可以设计为每隔约 10 秒尝试续期；这只是示例参数，需要结合实际耗时调整。

续期也必须确认锁属于自己：

~~~lua
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0
~~~

调用示例：

~~~shell
EVAL 'if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("PEXPIRE", KEYS[1], ARGV[2]) end return 0' 1 lock:demo token-A 30000
~~~

这会把剩余有效期重设为从现在起 30000 毫秒，不是在原来的到期时间上再加 30000 毫秒。

不能直接 `EXPIRE`，否则可能给别人的锁续期；锁已丢失时也不能用普通 `SET` 强行覆盖。

### 9.5.4 续期失败怎么办

1. 返回 `0`：当前 key 不存在或已不属于自己，应认为锁已丢失。
2. 网络超时：无法确认续期成功，不能无限假定锁仍然有效。
3. 无法确认仍在有效期内时，停止发起后续受保护操作，并尽力取消正在执行的业务。
4. 正常完成后停止续期 goroutine，再尝试按 token 释放锁。

还应设置业务总时限，防止业务卡死而续期任务一直运行。进程暂停、断网等也会让看门狗失效；取消信号不能保证已经发出的写请求被撤回。

续期机制与其限制参见 [Redis 分布式锁文档](https://redis.io/docs/latest/develop/clients/patterns/distributed-locks/#making-the-algorithm-more-reliable-extending-the-lock)。

## 9.6 获取失败、重试和异常处理

### 9.6.1 没拿到锁怎么处理

| 方式 | 适用思路 |
| --- | --- |
| 立即返回 | 返回处理中或繁忙，让调用方决定后续行为 |
| 有限重试 | 在明确的等待期限内，间隔一段时间再次尝试 |
| 等待已有结果 | 缓存重建时，稍后重新查询缓存 |
| 降级或排队 | 根据业务规则返回旧数据、进入任务队列等 |

不要写无间隔的死循环抢锁。重试可以使用退避和随机抖动，避免大量请求同时反复访问 Redis。

需要区分三个时间：**单次网络请求超时、获取锁的总等待时间、拿到锁后的 TTL**，它们不是同一个概念。

### 9.6.2 加锁超时不等于一定失败

~~~text
客户端发送 SET → Redis 设置成功 → 响应在网络中丢失 → 客户端收到超时
~~~

客户端此时不能直接进入临界区。保守处理可以是放弃本次业务，并尝试按该 token 清理可能存在的锁，清理失败则由 TTL 兜底。

即使清理成功，也要考虑原加锁请求可能延迟到达；此时 TTL 仍是限制占用时间的重要机制。复杂的重试需要围绕一次加锁尝试设计状态，不能只凭一次 GET 就假定长期持锁。

### 9.6.3 解锁失败与客户端退出

- 解锁返回 `0`：没有删除锁，不能报告为释放成功；锁可能已经过期或易主。
- 解锁网络错误：记录异常，可有限重试同一个比较删除操作。
- Go 使用 `defer` 释放锁，但进程被强制终止时，defer 不一定执行。
- 解锁需要独立且有超时的 context，避免业务 context 已取消，导致清理请求直接失败。
- TTL 是异常兜底；正常完成仍应及时释放，减少其他请求等待。

## 9.7 主从复制、哨兵和 Cluster 的影响

### 9.7.1 主从切换为什么可能丢锁

~~~text
1. A 在主节点 M 获取锁。
2. M 尚未把锁复制给从节点 R，就发生故障。
3. R 被提升为新主节点，但没有这个锁 key。
4. B 在 R 上获取锁。
5. A 仍在执行业务，A、B 可能同时操作资源。
~~~

根本原因是主从复制存在异步窗口。哨兵提高故障后的服务可用性，不会自动消除这个窗口。[官方故障切换说明](https://redis.io/docs/latest/develop/clients/patterns/distributed-locks/#why-failover-based-implementations-are-not-enough)

### 9.7.2 Cluster 也不是自动的锁安全保证

Cluster 按哈希槽分配 key。同一个锁 key 仍由对应槽的主节点处理，不会自动到所有主节点上投票加锁。

该分片发生主从切换时，也要考虑复制丢失。Lua 解锁和续期只操作一个 key；如果以后编写操作多个 key 的脚本，在 Cluster 中要满足这些 key 位于同一个槽的要求。

Redis 重启后锁能否恢复，也与持久化和故障时机有关。锁 key 被内存淘汰策略提前移除，同样会破坏原本预期的有效期，应避免让锁参与可淘汰缓存的管理策略。

## 9.8 Redlock（进阶了解）

Redlock 使用多个相互独立的 Redis 主节点，客户端向这些节点申请同一个资源的锁。

以 5 个独立节点为例：

1. 使用相同的 key、相同的本次 token，向各节点尝试加锁。
2. 至少在 3 个节点上成功，并且扣除获取耗时和时钟漂移余量后仍有有效期，才认为获取成功。
3. 业务必须在剩余有效期内完成，或按照协议成功续期。
4. 获取失败时，向所有尝试过的节点发起按 token 解锁，包括响应超时的节点。

**这些独立节点不是普通的一主多从，也不是直接拿 Cluster 的多个分片当作 Redlock。**

Redlock 仍需要分析时钟、进程暂停、网络延迟和节点重启等假设。它不会自动解决过期客户端继续写入的问题，也不会自动生成下面介绍的 fencing token。

入门先掌握单实例锁的生命周期和故障边界，再阅读 [Redlock 算法与安全性讨论](https://redis.io/docs/latest/develop/clients/patterns/distributed-locks/#the-redlock-algorithm)。

## 9.9 Fencing token 与业务正确性

### 9.9.1 如何拒绝过期持锁者的写入

Fencing token（隔离令牌）：与每次成功获取锁关联的严格递增编号，写入资源时必须携带，由资源端校验。

~~~text
A 获得编号 101，随后暂停，租约过期。
B 获得编号 102，数据库接受 B 的写入并记住 102。
A 恢复后带着 101 写入，数据库发现编号更旧，拒绝它。
~~~

关键是数据库或其他资源端要把“检查编号”和“执行写入”原子地完成。仅在客户端检查一次，检查后仍可能暂停。

### 9.9.2 随机 token 和 fencing token 的区别

| 类型 | 要求 | 主要作用 |
| --- | --- | --- |
| 随机 token | 不同加锁轮次能够区分 | 避免误删、误续期别人的锁 |
| Fencing token | 与获取锁顺序关联，严格递增 | 让资源端拒绝比已接受编号更旧的请求 |

随机 UUID 没有递增顺序，不能直接承担 fencing 的作用。也不能随意在加锁后补一个独立 `INCR` 就认为完成了设计：还要保证编号与获取锁顺序一致，并处理故障切换后的计数回退问题。

资源端如果不支持校验，这个机制就无法仅靠客户端实现。具体原理参见 [Martin Kleppmann 对分布式锁与 fencing 的分析](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html)。

### 9.9.3 锁和幂等不是一回事

~~~text
A 获取锁 → 扣款成功 → 释放锁
B 稍后获取锁 → 同一个请求又扣款一次
~~~

这两次操作没有并发重叠，但仍然发生了重复处理。

因此要根据业务建立请求幂等键、唯一约束、订单状态流转和数据库事务等机制。若业务已经能用数据库条件更新安全完成，也应先判断是否还需要额外的 Redis 锁。

## 9.10 可重入、公平性和死锁

### 9.10.1 基础锁不可重入

同一个持锁者再次执行 `SET ... NX ...` 时，key 已存在，仍然会失败。

可重入锁需要额外维护所有者身份、重入次数和对应的释放次数。不能简单看到 token 相同就随意重复进入，还要定义哪些调用属于同一个所有者。

### 9.10.2 基础锁不保证公平

Redis 按到达并被处理的命令竞争锁，不保证最早开始等待的客户端先获得锁。定时重试也不等于严格排队。

### 9.10.3 多把锁与死锁

~~~text
A 持有锁 1，等待锁 2。
B 持有锁 2，等待锁 1。
~~~

可以统一获取顺序，例如都先获取锁 1，再获取锁 2；设置有限等待；后续锁获取失败时释放已经获得的锁。

TTL 可以限制部分占用时间，但不能代替正确的业务设计。如果续期任务一直运行，错误的等待关系仍可能持续。

## 9.11 缓存击穿中的完整流程

~~~text
查询缓存
  ├─ 命中 → 返回数据
  └─ 未命中 → 尝试获取该资源的缓存重建锁
                 ├─ 未获得 → 有限等待后重新读缓存，或按规则降级
                 └─ 获得 → 再次查询缓存
                              ├─ 命中 → 释放锁，返回数据
                              └─ 未命中 → 查数据库 → 回写缓存
                                                   → 释放锁 → 返回数据
~~~

**拿到锁后要再次检查缓存。** 因为从第一次缓存未命中到成功获取锁之间，可能已经有其他请求完成了重建。

查询或回写失败时也需要进入清理流程。缓存重建特别慢、锁过期时，仍可能重复回源；如果旧数据回写会覆盖新数据，还要结合数据版本等机制处理。

## 9.12 Go 加锁与解锁示例

### 9.12.1 示例范围和准备

下面使用 `go-redis/v9` 演示一次加锁、模拟业务、按 token 解锁。示例不包含自动续期和 fencing，不应直接作为资金等业务的完整正确性方案。

准备一个可访问的 Redis；示例默认连接本机 `127.0.0.1:6379`、数据库 0，不配置密码。其他环境自行修改客户端连接参数。

在独立练习目录中，将下面的 Go 代码保存为 `main.go`，执行：

~~~shell
go mod init redis-lock-demo
go get github.com/redis/go-redis/v9
go run .
~~~

如果目录已有 `go.mod`，跳过 `go mod init`。`go get` 可能随版本更新提高 Go 版本要求，以实际安装的依赖为准。

### 9.12.2 完整示例

~~~go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var unlockScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
end
return 0
`)

func release(rdb *redis.Client, key, token string) {
	// 使用独立的超时 context，避免业务取消后无法发出解锁请求。
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	n, err := unlockScript.Run(ctx, rdb, []string{key}, token).Int64()
	if err != nil {
		log.Printf("解锁结果不确定：%v", err)
		return
	}
	if n == 0 {
		log.Print("未删除锁：锁已不存在或不属于本次操作")
		return
	}
	log.Print("释放锁成功")
}

func run() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:                  "127.0.0.1:6379",
		MaxRetries:            -1, // 示例禁用自动重试，显式处理不确定结果。
		DialTimeout:           2 * time.Second,
		ReadTimeout:           2 * time.Second,
		WriteTimeout:          2 * time.Second,
		ContextTimeoutEnabled: true,
	})
	defer rdb.Close()

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Errorf("生成 token：%w", err)
	}
	token := hex.EncodeToString(b)
	key := "learning:lock:order:1001"
	ttl := 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	started := time.Now()
	ok, err := rdb.SetNX(ctx, key, token, ttl).Result()
	cancel()
	if err != nil {
		// Redis 可能已设置成功，只是响应丢失；清理只比较自己的 token。
		release(rdb, key, token)
		return fmt.Errorf("加锁结果不确定，放弃本次业务：%w", err)
	}
	if !ok {
		log.Print("没有获得锁，本次不执行业务")
		return nil
	}
	defer release(rdb, key, token)

	// 保守扣除加锁耗时与余量。此检查不能消除之后的进程暂停风险。
	remaining := ttl - time.Since(started) - time.Second
	if remaining <= 0 {
		return fmt.Errorf("剩余锁有效期不足")
	}
	businessCtx, businessCancel := context.WithTimeout(context.Background(), remaining)
	defer businessCancel()

	log.Print("获得锁，开始模拟业务")
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	select {
	case <-businessCtx.Done():
		return fmt.Errorf("业务超时：%w", businessCtx.Err())
	case <-timer.C:
		log.Print("模拟业务完成")
		return nil
	}
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
	}
}
~~~

### 9.12.3 代码要点

- `SetNX(ctx, key, token, ttl)` 中传入正的 TTL，客户端会组合条件设置和过期时间；方法名不代表它只发送裸 `SETNX`。[go-redis 实现](https://github.com/redis/go-redis/blob/master/string_commands.go)
- TTL 不要传 `0`，否则可能创建不自动过期的锁。
- 只有 `err == nil` 且 `ok == true`，才继续评估有效期并开始业务。
- 获取失败不进入业务，结果不确定时仅按自己的 token 尝试清理。
- `defer release(...)` 在 `rdb.Close()` 之前执行，因为 defer 后注册的先执行。
- 超时 context 需要业务代码主动配合；它不会强制终止 goroutine，也不能撤销已被外部系统接受的写入。

## 9.13 动手验证

以下实验使用专门的 `learning:lock:*` 或 `lock:demo` key，不需要清空数据库。

| 实验 | 操作 | 应观察到的现象 |
| --- | --- | --- |
| 竞争同一把锁 | 两个终端在锁有效期内向同一个 key 执行 SET NX EX | 一个成功，另一个没有获得锁 |
| 错误 token 解锁 | 用 token-B 释放 token-A 的锁 | 返回 0，原锁仍存在 |
| 正确 token 解锁 | 在过期前使用 token-A 执行比较删除 | 返回 1，锁被删除 |
| 过期后的旧持锁者 | A 设置短 TTL；到期后 B 加锁；A 再用自己的 token 解锁 | 返回 0，B 的锁保留 |
| 续期 | 用正确 token 执行续期脚本，再查看 PTTL | 剩余有效期被重新设置 |
| 进程退出 | 运行 Go 示例，在模拟业务期间强制结束进程 | 即使未执行 defer，锁也会在 TTL 到期后失效 |
| 业务超过 TTL | 缩短 TTL，并故意让模拟任务忽略取消信号继续运行 | 锁可能已被别人获取，但旧任务仍在运行 |

最后一个实验用于理解风险，不要把忽略取消的写法用于实际业务。

## 9.14 常见面试问题

**1、Redis 分布式锁如何实现？**

使用 `SET key 唯一token NX EX/PX 有效期` 加锁，释放时原子比较 token 后删除；还需要考虑业务超时、续期、重试和 Redis 故障。

**2、为什么需要过期时间？**

避免客户端崩溃后永久占锁。正常完成仍应主动释放。

**3、为什么解锁需要 Lua？**

将检查归属和删除合成一个不可被其他命令插入的操作，避免误删。支持条件删除命令的版本也可使用对应命令。

**4、有了唯一 token 和 Lua 就绝对安全吗？**

它们解决误删问题；业务超出有效期、主从切换等仍可能导致业务并发重叠。

**5、看门狗能保证锁永远有效吗？**

不能。进程暂停、网络异常和续期失败都需要处理。

**6、哨兵、Cluster、Redlock 是同一个概念吗？**

不是。哨兵负责监控和故障切换；Cluster 负责分片并提供故障切换；Redlock 是客户端在多个独立 Redis 节点上协调获取锁的算法。

**7、分布式锁能保证不重复扣款吗？**

不能单独保证。即使请求依次执行，也可能重复扣款，还要设计幂等和数据库约束等业务机制。

**8、学习完成后应能独立解释什么？**

为什么要原子加锁、为什么需要唯一 token、为什么要原子解锁、锁过期后旧业务怎么办，以及主从切换为什么可能让两个客户端都认为自己持锁。
