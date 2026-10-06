/*
章节：04-数组切片与映射
	本章掌握清单：
	1. 数组：长度固定，赋值复制元素
	2. 切片：分清长度、容量和 nil
	3. 共享底层数组与追加覆盖
	4. copy：复制元素与浅复制
	5. 排序：明确是否修改原数据
	6. map：键值、存在性和初始化
	7. 切片的增删查改与内容比较
	8. map 的复制、清空与共享
	9. nil 与空集合在 JSON 中的区别
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./04-collections
// 运行练习：go run ./04-collections exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：数组 [N]T 有固定长度；[]T 更接近常用的动态列表，但赋值不等于复制元素。
// 后续再学：集合遍历见第 07 章；指针及结构体见第 09/11 章；并发访问 map 见第 17 章。

package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：数组：长度固定，赋值复制元素
	// ═══════════════════════════════════════════════════════
	// [3]string 中的 3 是类型的一部分；初始化不足三个元素时，其余位置为 ""。
	// 这里的元素都是字符串，修改副本的一项不会影响原数组。

	// 数组长度属于类型：[2]string 与 [3]string 是不同类型。
	environments := [3]string{"development", "staging", "production"}
	arrayCopy := environments // 数组赋值会复制每个元素。
	arrayCopy[0] = "local"
	fmt.Println("原数组：", environments, "副本：", arrayCopy)
	fmt.Println("长度：", len(environments), "最后一项：", environments[len(environments)-1])
	// 访问下标前要确保范围有效；空列表没有最后一项。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：切片：分清长度、容量和 nil
	// ═══════════════════════════════════════════════════════
	// len 决定当前能访问哪些下标；cap 表示当前切片从起点最多能扩展到的长度。
	// make([]string, 0, 4) 不能直接写 items[0]，必须先 append 或合法地扩展长度。

	// 切片描述一段底层数组，包含指向数据的位置、长度 len、容量 cap。
	var nilItems []string
	emptyItems := []string{}
	fmt.Println("nil 切片：", nilItems == nil, len(nilItems), cap(nilItems))
	fmt.Println("空切片：", emptyItems == nil, len(emptyItems), cap(emptyItems))
	// 两者都可直接 append；切片只能与 nil 比较，不能彼此用 == 比较。
	items := make([]string, 0, 4) // 长度 0、容量 4，并不是已经有四个可索引元素。
	items = append(items, "login", "profile")
	fmt.Println("追加后：", items, "len：", len(items), "cap：", cap(items))
	// make([]string, 4) 则会创建四个元素，它们初始都是 ""。
	// append 返回新的切片描述，必须接收；容量不足时会分配新的底层数组。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：共享底层数组与追加覆盖
	// ═══════════════════════════════════════════════════════
	// alias、first 与 items 起初指向同一组元素；first 长度虽短，容量仍允许继续写入。
	// 因此 append(first, "settings") 会改到 items[1]，而不只是得到一个独立的新列表。

	// 切片赋值与截取通常共享底层数组；截取范围是左闭右开。
	alias := items
	first := items[:1]
	alias[0] = "sign-in"
	fmt.Println("修改 alias 后：", items, first)
	// 此处容量足够，向 first 追加会覆盖底层数组的第二项。
	first = append(first, "settings")
	fmt.Println("向子切片追加后：", items, first)
	// 容量增长策略是实现细节，不要依赖 append 后容量恰好翻倍。

	// 三下标切片 s[low:high:max] 能限制容量，让后续 append 分配新数组。
	limited := items[:1:1]
	limited = append(limited, "isolated")
	fmt.Println("限制容量后再追加：", items, limited)
	// 注意：追加之前 limited[0] 仍然与 items[0] 共享。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：copy：复制元素与浅复制
	// ═══════════════════════════════════════════════════════
	// 先用 make 建好目标长度，再 copy；只有容量、没有长度的目标会复制零项。
	// 外层元素如果还是切片，复制的只是内部切片描述，所以 nestedCopy 仍能改到原数据。

	// copy 只复制 min(len(dst), len(src)) 个元素，不会自动扩展目标长度。
	backup := make([]string, len(items))
	n := copy(backup, items)
	backup[0] = "backup-only"
	fmt.Println("复制项数：", n, "原数据：", items, "副本：", backup)
	// 复制嵌套切片、map、指针时，内部数据仍可能共享，这叫浅复制。
	nested := [][]int{{1, 2}, {3, 4}}
	nestedCopy := make([][]int, len(nested))
	copy(nestedCopy, nested)
	nestedCopy[0][0] = 99
	fmt.Println("浅复制共享内部切片：", nested, nestedCopy)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：排序：明确是否修改原数据
	// ═══════════════════════════════════════════════════════
	// 排序是原地操作；本节先复制 latencies，再排序 sorted。
	// 运行后应看到原顺序 [80 20 50] 和升序 [20 50 80] 同时保留。

	// slices.Sort 原地排序；想保留原顺序，应先复制再排序。
	latencies := []int{80, 20, 50}
	sorted := make([]int, len(latencies))
	copy(sorted, latencies)
	slices.Sort(sorted)
	fmt.Println("原始耗时：", latencies, "升序：", sorted)

	// ═══════════════════════════════════════════════════════
	// 第 6 节：map：键值、存在性和初始化
	// ═══════════════════════════════════════════════════════
	// users[12] 得到空字符串，无法单凭它判断用户是否存在；第二个返回值 exists 才能区分。
	// 切片的 nil 可直接 append，map 的 nil 却不能直接写入，二者要分别记忆。

	// map 是键值映射，遍历顺序不保证稳定，不应依赖其输出顺序。
	// key 必须支持 == / != 比较，即 comparable，不局限于基本类型。
	// 例如元素可比较的数组也能作为 key；切片、map、函数不能作为 key。
	users := map[int]string{11: "alice", 13: ""}
	fmt.Println("存在的用户：", users[11], "不存在的用户：", users[12])
	missing, exists := users[12]
	empty, emptyExists := users[13]
	fmt.Printf("缺失：%q %t；空值：%q %t\n", missing, exists, empty, emptyExists)
	delete(users, 13) // 删除不存在的 key 也合法。
	users[15] = "bob"
	fmt.Println("用户数：", len(users))

	// nil map 可以读取、len、delete，但写入会 panic；写入前必须初始化。
	var cache map[string]int
	fmt.Println("nil map 读取：", cache["count"], "长度：", len(cache))
	cache = make(map[string]int)
	cache["count"] = 1
	cacheAlias := cache // map 赋值仍共享同一份映射内容。
	cacheAlias["count"] = 2
	fmt.Println("通过别名修改 map：", cache["count"])

	coordinates := map[[2]int]string{{3, 4}: "node-a"}
	fmt.Println("数组作为 key：", coordinates[[2]int{3, 4}])
	// 数字字面量 0011 是旧式八进制，十进制值为 9，不是带格式的编号 11。
	// 建议八进制写 0o11；需要保留前导零的业务 ID 应使用字符串。
	fmt.Println("0011、0o11、11：", 0011, 0o11, 11)
	userCodes := map[string]string{"0011": "alice"}
	fmt.Println("保留前导零的业务 ID：", userCodes["0011"])

	// ═══════════════════════════════════════════════════════
	// 第 7 节：切片的增删查改与内容比较
	// ═══════════════════════════════════════════════════════
	// 插入和删除会改变切片长度，返回值代表操作后的有效范围。
	// 查找失败返回 -1，不能把这个结果直接用于下标访问。

	// 日常切片操作：追加另一切片要使用 ...，像把一批元素展开传入。
	routes := []string{"/health"}
	extraRoutes := []string{"/users", "/orders"}
	routes = append(routes, extraRoutes...)
	// slices.Insert/Delete 可能复用底层数组，都必须接收返回的切片。
	routes = slices.Insert(routes, 1, "/metrics")
	fmt.Println("插入路由：", routes)
	routes = slices.Delete(routes, 2, 3) // 删除 [2:3)，即 /users；下标必须合法。
	fmt.Println("删除后：", routes, "查找 orders：", slices.Index(routes, "/orders"))
	fmt.Println("包含 health：", slices.Contains(routes, "/health"), "查找缺失项：", slices.Index(routes, "/missing"))
	// Delete 会清零被移出有效范围的尾部元素，帮助释放其中引用的数据。
	// 不要继续把操作前的旧切片长度当作有效长度；别名也可能看到元素被移动或清零。
	clonedRoutes := slices.Clone(routes)
	fmt.Println("浅复制后内容相等：", slices.Equal(routes, clonedRoutes))
	fmt.Println("nil 与空切片的内容相等：", slices.Equal(nilItems, emptyItems))
	// Equal 比较元素，nil 与空切片都含零个元素，所以内容相等；与 == nil 的判断不同。

	// ═══════════════════════════════════════════════════════
	// 第 8 节：map 的复制、清空与共享
	// ═══════════════════════════════════════════════════════
	// cacheAlias := cache 共享映射；maps.Clone(cache) 则建立独立的一层映射。
	// 本例值是 int，修改副本的 count 不影响原 map；嵌套容器需要进一步复制。

	// maps.Clone 复制映射本身，仍然是浅复制；嵌套切片、map 等值依旧可能共享。
	cacheBackup := maps.Clone(cache)
	cacheBackup["count"] = 99
	fmt.Println("复制 map 后修改：", cache["count"], cacheBackup["count"])
	clear(cacheBackup) // 删除全部键；对 nil map 也安全。
	fmt.Println("清空复制后的 map：", len(cacheBackup), "原 map 仍有：", len(cache))
	// map 的元素不能直接取地址：不能写 &cache["count"]。
	// 若 map 的值为结构体，修改字段通常需取出、修改、写回，或让值存指针（第 09/11 章）。
	// 遍历与“收集 key 后排序输出”留到第 07 章；map 的遍历顺序本身不可依赖。

	// ═══════════════════════════════════════════════════════
	// 第 9 节：nil 与空集合在 JSON 中的区别
	// ═══════════════════════════════════════════════════════
	// Go 中 len 都为 0 的两个切片，对前端却可能分别表现为 null 和 []。
	// 决定是否初始化空集合时，要同时考虑 Go 操作和 API 的返回约定。

	// 面向前端的一个实际区别：encoding/json 默认把 nil 切片编码为 null，空切片为 []。
	// 这里仅观察集合编码结果，完整 JSON 结构会在第 11 章学习。
	nilJSON, nilJSONErr := json.Marshal(nilItems)
	emptyJSON, emptyJSONErr := json.Marshal(emptyItems)
	fmt.Printf("nil 切片 JSON=%s err=%v；空切片 JSON=%s err=%v\n", nilJSON, nilJSONErr, emptyJSON, emptyJSONErr)
	// nil map 与已初始化的空 map 对应 null 与 {}；后端应按接口约定选择初始化方式。
	// clear(slice) 会把元素变成零值，长度不变；clear(map) 则删键，长度变为 0。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 make([]string, 0, 4) 改成 make([]string, 4)，观察为什么 append 后前面多出空字符串。
// 把 copy 的目标改成 make([]string, 0, len(items))，预测复制数量；再恢复正确写法。
