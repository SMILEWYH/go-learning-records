// 章节：04-数组切片与映射
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// exam 自动检查三道练习；前四章填写 exercise1/2/3，输出要求见 exercises_test.go。

// 练习 1：构造批量请求列表。
// 背景：前端把选中的用户 ID 交给后端批量查询。
// 要求：
//  1. 用 make 创建长度为 0、容量为 4 的 []string。
//  2. 依次追加 "u-01"、"u-02"，再把另一个切片 {"u-03", "u-04"} 整批追加进去。
//  3. 用 [:2] 取出第一页，用 [2:] 取出第二页。
//  4. 打印完整切片的 len、cap，以及两页内容。
//
// 常规：最终 len=4、cap=4；第一页 [u-01 u-02]，第二页 [u-03 u-04]。
// 边界：另建一个 nil 的 []string 并追加 "u-01"，应得到长度 1 的有效切片。
// 提示：整批追加使用 append(target, source...)；不要对长度 0 的切片直接写 [0]。
// 建议实现位置：exercise1()，无需 helper；手动处理这组确定长度的数据即可。
// 思考：如果修改第一页中的元素，完整切片会不会改变？请实际验证。
func exercise1() {
	// TODO：创建、追加、切片，并打印结果。
}

// 练习 2：为监控面板排序耗时，同时保留原始数据。
// 背景：监控面板既需要按请求顺序展示耗时，也需要展示从低到高的耗时列表。
// 要求：
//  1. 创建 source := []int{120, 35, 80, 35}。
//  2. 用 make 分配同长度的目标切片，再用 copy 复制，记录复制的元素数量。
//  3. 在 import 中加入 slices，用 slices.Sort 对目标切片升序排序。
//  4. 把目标第一项改为 0，确认 source 完全不变。
//
// 常规：复制数量为 4；排序后为 [35 35 80 120]；source 始终为 [120 35 80 35]。
// 边界：对 []int{} 重复操作，复制数量为 0，排序结果仍为空；不要再写第一项。
// 建议实现位置：exercise2()；学完函数后可抽取 sortedCopy(values []int) []int。
// 思考：如果仅写 target := source，为什么无法满足“保留原始数据”的要求？
func exercise2() {
	// TODO：分配独立底层存储，再复制、排序、验证原切片。
}

// 练习 3：区分用户昵称缺失与昵称为空。
// 背景：接口响应里，“这个用户不存在”与“用户还没设置昵称”是两个状态。
// 要求：
//  1. 创建 map[string]string，数据为 "0011":"alice"、"0013":""。
//  2. 用 value, ok 读取 "0011"、"0013"、"0012"，分别打印值和是否存在。
//  3. 删除 "0011" 并再次查询；再新增 "0012":"bob"。
//
// 常规："0011" -> "alice", true；新增后的 "0012" -> "bob", true。
// 边界："0013" -> "", true；尚未新增的 "0012" -> "", false。
// 删除后的 "0011" -> "", false。
// 建议实现位置：exercise3()，每个 key 直接查询，不需要遍历或 helper。
// 思考：把 "0011" 改成数字 0011 时，它的十进制值是多少？为何业务 ID 要用字符串？
func exercise3() {
	// TODO：创建 map、执行查询与增删，并输出 value 和 ok。
}

// ── 自动检查启动入口：填写习题时无需修改 ──
// runExercises 只启动本章的检查入口；检查规则、测试和参考答案都在 exercises_test.go。
// Go 的普通运行不会编译 _test.go，因此通过 go test 进入 TestExam。
func runExercises() {
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法定位本章源码，检查未执行")
		return
	}
	dir := filepath.Dir(source)
	if !filepath.IsAbs(dir) { // 兼容以 -trimpath 构建后在课程根目录或章节目录运行。
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("无法定位章节：", err)
			return
		}
		dir = cwd
		if _, err := os.Stat(filepath.Join(dir, "exercises.go")); err != nil {
			dir = filepath.Join(cwd, filepath.Base(filepath.Dir(source)))
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "test", "-json", "-count=1", "-run", "^TestExam$", "-timeout", "5m", ".")
	command.Dir = dir
	command.Env = append(os.Environ(), "GO_LESSON_EXAM=1")
	command.WaitDelay = 2 * time.Second
	output, runErr := command.Output()
	decoder := json.NewDecoder(bytes.NewReader(output))
	var diagnostics strings.Builder
	passed := false
	for {
		var event struct{ Action, Test, Output string }
		if err := decoder.Decode(&event); err != nil {
			break
		}
		diagnostics.WriteString(event.Output)
		if event.Test == "TestExam" {
			if event.Action == "pass" {
				passed = true
			}
			if event.Action == "output" && !strings.HasPrefix(event.Output, "=== RUN") && !strings.HasPrefix(event.Output, "--- PASS") {
				fmt.Print(event.Output)
			}
		}
	}
	if runErr != nil || !passed {
		fmt.Println("练习检查未正常结束：", runErr)
		fmt.Print(diagnostics.String())
		if failure, ok := runErr.(*exec.ExitError); ok {
			fmt.Print(string(failure.Stderr))
		}
	}
}
