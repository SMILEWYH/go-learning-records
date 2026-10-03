/*
章节：05-条件判断
	本章掌握清单：
	1. if / else：二选一
	2. else if：按顺序划分范围
	3. 初始化语句、作用域与遮蔽
	4. 短路求值与条件组合
	5. 默认值与条件覆盖
	6. 提前返回：先处理无法继续的情况
	7. 边界检查与执行顺序
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./05-if
// 运行练习：go run ./05-if exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：Go 的条件无需圆括号，但必须是 bool；&& / || 的结果也一定是 bool。
// 后续再学：函数参数和返回值见第 08 章；实际 HTTP 身份认证和授权需在服务开发时单独学习。

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：if / else：二选一
	// ═══════════════════════════════════════════════════════
	// if 后写 bool 表达式，{} 内写命中条件后要执行的语句。
	// 条件为 true 执行 if，为 false 执行 else，两边只执行一边。

	// 基本分支：条件不用括号，但花括号不能省略。
	loggedIn := true
	if loggedIn {
		fmt.Println("可以访问个人资料")
	} else {
		fmt.Println("请先登录")
	}
	// else 要与前一块结束的 } 放在同一行，避免 Go 自动插入分号造成语法错误。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：else if：按顺序划分范围
	// ═══════════════════════════════════════════════════════
	// 第一个条件不成立才会继续判断下一个，所以第二个分支隐含 responseMS >= 100。
	// 用 99、100、499、500 核对每个区间，避免只测试中间值。

	// else if 按顺序判断，首次匹配后跳过后面的分支。
	responseMS := 350
	if responseMS < 100 {
		fmt.Println("耗时等级：fast")
	} else if responseMS < 500 {
		fmt.Println("耗时等级：normal")
	} else {
		fmt.Println("耗时等级：slow")
	}
	// 这里 100ms 进入 normal，500ms 进入 slow；边界要与需求一致。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：初始化语句、作用域与遮蔽
	// ═══════════════════════════════════════════════════════
	// 分号前先取值，分号后再判断；role 和 exists 可供整段 if/else 使用。
	// 如果外层也有 role，初始化中的 := 会建立内层变量，结束后外层值保持不变。

	// if 可以先执行一条初始化语句，变量只在 if / else 整段里可用。
	roles := map[string]string{"u-01": "admin", "u-02": "viewer"}
	if role, exists := roles["u-01"]; !exists {
		fmt.Println("用户不存在")
	} else if role == "admin" {
		fmt.Println("管理员可以修改配置")
	} else {
		fmt.Println("只允许查看配置")
	}
	// 此处不能继续使用 role 或 exists；它们已超出作用域。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：短路求值与条件组合
	// ═══════════════════════════════════════════════════════
	// 把“能不能安全访问”的条件放前面：先检查长度，再访问第一个元素。
	// || 左侧为 true 时不执行右侧，所以缓存命中时不会调用 tracePolicyCheck。

	// && 和 || 会从左到右短路求值。
	// 左侧为 false，&& 不再计算右侧，因此空切片不会触发 [0] 越界。
	var selectedIDs []string
	if len(selectedIDs) > 0 && selectedIDs[0] != "" {
		fmt.Println("首个选中 ID：", selectedIDs[0])
	} else {
		fmt.Println("没有可用的首个 ID")
	}
	role := "editor"
	canEdit := role == "admin" || role == "editor"
	fmt.Println("能否编辑：", canEdit)
	// 与 TS 不同，不能写 name || "guest" 来选择一个字符串；Go 逻辑运算只用于 bool。
	// 条件组合要留意优先级：&& 比 || 优先，权限判断常用括号明确分组。
	// 已登录 && (管理员 || 资源所有者)，不能误写成“已登录且管理员，或者是资源所有者”。
	isOwner := true
	fmt.Println("组合权限判断：", loggedIn && (role == "admin" || isOwner))

	// 短路也会跳过函数调用及其副作用；不要把“必须执行的记录操作”藏在条件右边。
	// 这里用下面的小函数打印执行痕迹，函数语法在第 08 章展开。
	hasCachedDecision := true
	allowed := hasCachedDecision || tracePolicyCheck()
	fmt.Println("有缓存时跳过策略检查：", allowed)
	hasCachedDecision = false
	allowed = hasCachedDecision || tracePolicyCheck()
	fmt.Println("没有缓存时执行策略检查：", allowed)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：默认值与条件覆盖
	// ═══════════════════════════════════════════════════════
	// 先把默认结果写清楚，再在条件成立时替换。
	// 这种写法也适合默认排序方式、默认页面标题和缺省配置。

	// Go 没有三元表达式 condition ? a : b，通常先给默认值，再按条件覆盖。
	name := ""
	displayName := "guest"
	if name != "" {
		displayName = name
	}
	fmt.Println("显示名：", displayName)

	// ═══════════════════════════════════════════════════════
	// 第 6 节：提前返回：先处理无法继续的情况
	// ═══════════════════════════════════════════════════════
	// 把未登录、权限不足逐个排除，剩下的路径就可以直接执行业务。
	// 阅读下方 printAccessDecision 时，注意 return 只结束本次函数调用。

	// 嵌套适合体现从属关系；层数过多时可以用提前返回简化。
	if loggedIn {
		if canEdit {
			fmt.Println("已登录且具有编辑权限")
		}
	}
	// 下方小函数只是为演示 return 的作用范围；第 08 章会详细学习函数。
	printAccessDecision(false, "admin")
	printAccessDecision(true, "viewer")
	printAccessDecision(true, "admin")

	// ═══════════════════════════════════════════════════════
	// 第 7 节：边界检查与执行顺序
	// ═══════════════════════════════════════════════════════
	// 比较 remainingQuota < 0 与 <= 0：只有后者把恰好为 0 的情况也归为耗尽。
	// 先完成校验再扣减配额，才能避免无效输入已经触发业务操作。

	// 输入边界应在做危险操作之前处理；逻辑分支要覆盖零、负值、正常值和上限。
	// 例如“剩余配额”为 0 也需要拒绝，而不只是拒绝负数。
	remainingQuota := 0
	if remainingQuota <= 0 {
		fmt.Println("配额耗尽：不执行下一项任务")
	} else {
		fmt.Println("开始任务，当前配额：", remainingQuota)
	}
	// 典型检查顺序：先确认是否存在，再检查范围，最后下标访问、除法或执行业务操作。
}

func tracePolicyCheck() bool {
	fmt.Println("执行策略检查")
	return true
}

func printAccessDecision(loggedIn bool, role string) {
	// return 结束当前函数，不会结束调用它的 main。
	// 优先处理失败条件，成功路径就不必层层嵌套。
	if !loggedIn {
		fmt.Println("拒绝：未登录")
		return
	}
	if role != "admin" {
		fmt.Println("拒绝：需要管理员权限")
		return
	}
	fmt.Println("允许：管理员操作")
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 依次用 99、100、499、500 检查耗时分级。
// 在 if 前声明同名变量，分别用 := 和 =，核对第 01 章讲的遮蔽规则。
