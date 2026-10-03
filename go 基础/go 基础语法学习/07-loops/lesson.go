/*
章节：07-循环遍历
	本章掌握清单：
	1. 三段式 for
	2. 只写条件的 for
	3. for {} 与明确退出
	4. range 切片：索引与元素副本
	5. range map：需要顺序时先排序键
	6. range 字符串：字节偏移与 rune
	7. continue 与 break
	8. 标签：退出指定的外层循环
	9. 遍历期间修改数组或切片
	10. 遍历 map 时增删条目
	11. range 整数与每轮变量
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./07-loops
// 运行练习：go run ./07-loops exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：Go 用 for 覆盖普通循环、while 风格循环及集合遍历，没有 while 关键字。
// 所有演示均有明确终止条件，不会发请求、等待或无限循环。
// 后续再学：循环变量与闭包的组合见第 08 章；channel 遍历见第 15 章；自定义迭代器按需学习。

package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：三段式 for
	// ═══════════════════════════════════════════════════════
	// 执行顺序是：初始化一次 -> 检查条件 -> 循环体 -> 更新 -> 再检查条件。
	// i <= 3 包含 3；改成 i < 3 后会少执行一轮。

	// 初始化；条件；每轮结束后的操作。i 的作用域局限于当前循环。
	for i := 1; i <= 3; i++ {
		fmt.Println("批次：", i)
	}
	// i++ 是语句，不能写 value := i++，也没有 ++i。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：只写条件的 for
	// ═══════════════════════════════════════════════════════
	// remaining-- 让剩余数量逐轮减少，是这个循环能终止的关键。
	// 如果忘记更新条件依赖的变量，条件就可能一直成立。

	// 只保留条件，就是 while 风格的循环。
	remaining := 3
	for remaining > 0 {
		fmt.Println("剩余任务：", remaining)
		remaining--
	}

	// ═══════════════════════════════════════════════════════
	// 第 3 节：for {} 与明确退出
	// ═══════════════════════════════════════════════════════
	// 没有循环条件时，把退出条件放进循环体。
	// 本例先把 attempt 加一再判断，所以最终执行两轮。

	// for {} 本身没有条件，要在循环体内确保能 break 或 return。
	// 先执行一次，再检查条件，可表达 do-while 的意图。
	attempt := 0
	for {
		attempt++
		fmt.Println("模拟尝试：", attempt)
		if attempt >= 2 {
			break
		}
	}

	// ═══════════════════════════════════════════════════════
	// 第 4 节：range 切片：索引与元素副本
	// ═══════════════════════════════════════════════════════
	// for index, id := range ids 中，index 是位置，id 是该位置的值副本。
	// 对 int 副本加 5 不影响原切片；通过 scores[index] 才能写回原位置。

	// 遍历数组或切片：range 的第一个值是索引，第二个值是元素副本。
	ids := []string{"u-01", "u-02", "u-03"}
	for index, id := range ids {
		fmt.Println("用户索引与 ID：", index, id)
	}
	for _, id := range ids { // 不需要索引时用 _。
		fmt.Println("只读取 ID：", id)
	}
	scores := []int{10, 20}
	for _, score := range scores {
		score += 5
		fmt.Println("修改的元素副本：", score)
	}
	fmt.Println("原切片仍为：", scores)
	for index := range scores {
		scores[index] += 5 // 要修改原来的 int 元素，需要通过索引赋值。
	}
	fmt.Println("按索引修改后：", scores)
	// 元素若本身是 map、切片或指针，复制元素后仍可能访问共享数据。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：range map：需要顺序时先排序键
	// ═══════════════════════════════════════════════════════
	// map 遍历适合处理全部键值，但不承诺先看到 GET 还是 POST。
	// 排序 keys 后逐个读取，可以得到稳定的展示顺序。

	// map 的 range 顺序不保证稳定，nil 或空 map 会循环零次。
	counts := map[string]int{"GET": 3, "POST": 2}
	for method, count := range counts {
		fmt.Println("无顺序保证的统计：", method, count)
	}
	// 需要确定顺序时，收集 key 并排序，再按 key 读取。
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		fmt.Println("按 key 排序：", key, counts[key])
	}

	// ═══════════════════════════════════════════════════════
	// 第 6 节：range 字符串：字节偏移与 rune
	// ═══════════════════════════════════════════════════════
	// range 会按 UTF-8 解码，所以中文和 emoji 各作为一个码点返回。
	// 第一个返回值仍是字节位置，因此从 1 跳到 4。

	// range 字符串返回的是字节偏移和 rune，不是“第几个字符”的编号。
	for byteOffset, character := range "A中😀" {
		fmt.Printf("字节偏移=%d rune=%c\n", byteOffset, character)
	} // 字节偏移依次为 0、1、4。

	// ═══════════════════════════════════════════════════════
	// 第 7 节：continue 与 break
	// ═══════════════════════════════════════════════════════
	// continue 跳过本轮剩余部分后继续循环；break 直接结束所在循环。
	// 嵌套 switch 时先找最近一层可被 break 退出的结构。

	// continue 跳过当前轮剩余内容；break 退出最内层的 for 或 switch。
	statuses := []int{200, 404, 503, 502, 200}
	serverErrors := 0
	for _, status := range statuses {
		if status < 500 || status >= 600 {
			continue
		}
		serverErrors++
		if serverErrors == 2 {
			fmt.Println("发现两次服务端错误，停止检查")
			break
		}
	}
	fmt.Println("已统计的服务端错误：", serverErrors)
	// 当 switch 位于 for 内部时，switch 中的 break 只退出 switch；要退出外层可用标签。
	for _, state := range []string{"skip", "run"} {
		switch state {
		case "skip":
			break // 退出 switch；下方“循环尾部”仍然会输出 skip。
		default:
			fmt.Println("处理状态：", state)
		}
		fmt.Println("循环尾部：", state)
	}

	// ═══════════════════════════════════════════════════════
	// 第 8 节：标签：退出指定的外层循环
	// ═══════════════════════════════════════════════════════
	// search: 给外层循环命名，break search 可以一次退出两层查找。
	// 找到 failed 后，foundGroup 与 foundIndex 分别记录组号和组内位置。

	// 嵌套循环中使用标签，明确退出哪一层；不加标签只会退出最近一层。
	groups := [][]string{{"queued", "done"}, {"running", "failed"}}
	foundGroup, foundIndex := -1, -1
search:
	for groupIndex, group := range groups {
		for index, state := range group {
			if state == "failed" {
				foundGroup, foundIndex = groupIndex, index
				break search
			}
		}
	}
	fmt.Println("首个失败任务的位置：", foundGroup, foundIndex)
	// continue 标签则开始指定外层循环的下一轮；标签必须对应包围当前语句的循环。

	// ═══════════════════════════════════════════════════════
	// 第 9 节：遍历期间修改数组或切片
	// ═══════════════════════════════════════════════════════
	// 先区分复制的是整个数组还是切片描述，再判断下一轮从哪里取值。
	// queue 的遍历长度在开始时确定为 2；后来追加的元素不会加入本轮 range。

	// range 通常在开始时求值一次：数组表达式会产生数组副本，切片则复制切片描述。
	// 这里读取的 value 来自数组副本，即使修改了原数组的下一项，下一轮仍读到旧值。
	array := [3]int{10, 20, 30}
	for index, value := range array {
		if index == 0 {
			array[1] = 200
		}
		fmt.Println("数组 range 的副本值：", index, value)
	}
	fmt.Println("原数组的实际内容：", array)
	// range 切片固定本次遍历的长度，不会自动遍历后续 append 追加的元素。
	queue := []int{1, 2}
	for _, item := range queue {
		queue = append(queue, item+10)
	}
	fmt.Println("只遍历初始两项后追加：", queue) // [1 2 11 12]，不会无限追加。

	// ═══════════════════════════════════════════════════════
	// 第 10 节：遍历 map 时增删条目
	// ═══════════════════════════════════════════════════════
	// 删除过期键后，最终 map 只保留 active。
	// 如果边遍历边新增任务，新增项可能被访问，也可能不会，不能用它保证任务处理完整。

	// 遍历 map 时可以删除条目；还没走到的条目若被删掉，就不会再遍历到。
	// 遍历中新增的条目是否在本次被遍历到不保证，所以不要依赖边遍历边新增来处理完整队列。
	ttlByKey := map[string]int{"active": 30, "expired": 0}
	for key, ttl := range ttlByKey {
		if ttl <= 0 {
			delete(ttlByKey, key)
		}
	}
	fmt.Println("清理过期项后的 map：", ttlByKey)
	// 与 map 不同，range 切片时直接删除元素容易跳过数据或使用旧长度。
	// 一般应另建结果切片，append 要保留的元素，或使用 slices.DeleteFunc 等明确的操作。

	// ═══════════════════════════════════════════════════════
	// 第 11 节：range 整数与每轮变量
	// ═══════════════════════════════════════════════════════
	// range 3 依次产生 0、1、2，适合只需要次数的循环。
	// 闭包捕获循环变量时，要区分 := 创建每轮变量与 = 复用外部变量。

	// Go 1.22 起，可以 range 一个整数，得到 0 到 n-1；n <= 0 时执行零次。
	for index := range 3 {
		fmt.Println("固定重复次数：", index)
	}
	// Go 1.22 及更新语言版本中，:= 声明的循环变量每轮都是独立变量。
	// 如果使用 = 给循环外已存在的变量赋值，仍然复用那个变量；第 08 章用闭包演示区别。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把三段式循环的 <= 改成 <，先数清会执行几轮。
// 对比修改 range 得到的 score 和修改 scores[index]；解释为什么只有后者写回。
