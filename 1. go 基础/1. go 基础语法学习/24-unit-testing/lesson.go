/*
章节:24-单元测试
*/

// TS 对照:testing 对应 Vitest/Jest 的基础测试能力;t.Run 类似具名子测试
// Go 通常使用标准库测试, 不需要先安装第三方测试框架
// 直接运行 lesson.go 只执行讲解;下面的测试代码是注释示例
// 需要实际运行测试时, 可自行新建 lesson_test.go（package main, 导入 testing）并写入测试
// 在本章目录明确指定文件:go test -v lesson.go lesson_test.go
// 指定用例:go test -v -run '^TestFreeStandard$' lesson.go lesson_test.go
// 竞争与覆盖率:go test -race -count=1 -cover lesson.go lesson_test.go

package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrEmptyReceipt = errors.New("收据不能为空")

// ShippingFee 用整数“分”计算运费, 适合用纯函数单元测试验证业务规则
// 订单不足 100 元:普通运费 6 元;满 100 元免普通运费;加急额外收 10 元
func ShippingFee(subtotalCents int, express bool) (int, error) {
	if subtotalCents < 0 {
		return 0, errors.New("订单金额不能为负数")
	}
	fee := 600
	if subtotalCents >= 10000 {
		fee = 0
	}
	if express {
		fee += 1000
	}
	return fee, nil
}

// SaveReceipt 用于演示有文件副作用的测试对象, 可在测试中用 t.TempDir 隔离文件
func SaveReceipt(path, receipt string) error {
	if receipt == "" {
		return ErrEmptyReceipt
	}
	return os.WriteFile(path, []byte(receipt+"\n"), 0o600)
}

func main() {
	// ═══════════════════════════════════════════════════════
	// 第 1 节:测试对象:先明确输入和期望结果
	// ═══════════════════════════════════════════════════════
	// ShippingFee 是纯函数, 同样的金额与加急选项应返回同样的运费, 适合直接测试
	// 本例订单 120 元满足免普通运费条件, 加急再加 1000 分, 因此最终结果是 1000

	fee, err := ShippingFee(12000, true)
	if err != nil {
		fmt.Println("运费计算失败：", err)
		return
	}
	fmt.Println("订单 120 元、加急配送，运费（分）：", fee) // 1000

	// ═══════════════════════════════════════════════════════
	// 第 2 节:测试文件、测试函数与断言
	// ═══════════════════════════════════════════════════════
	// 测试函数形如 func TestShippingFee(t *testing.T), 需要放在自行创建的 _test.go 文件中
	// 一个最小断言可写成下面这样, 放在 _test.go 文件内并导入 testing:
	//
	//   func TestFreeStandard(t *testing.T) {
	//       got, err := ShippingFee(10000, false)
	//       if err != nil {
	//           t.Fatal(err)
	//       }
	//       if got != 0 {
	//           t.Errorf("got %d, want 0", got)
	//       }
	//   }
	//
	// 测试必须检查期望;只打印结果, 即使算错也可能显示 PASS

	fmt.Println("测试写入自己的 lesson_test.go 后，执行 go test -v lesson.go lesson_test.go")
	// _test.go 文件只在 go test 时参与编译, 不会被 go run 的 main 调用
	// 测试函数必须形如 TestXxx(t *testing.T), 失败时用 t.Error/t.Fatal 等报告
	// t.Fatal 停止当前测试, t.Error 记录失败后继续;不要用 Println 代替断言

	// ═══════════════════════════════════════════════════════
	// 第 3 节:表驱动测试、子测试与边界
	// ═══════════════════════════════════════════════════════
	// 表驱动测试的 cases 把 name、输入、want、wantErr 放到每一行数据中
	// 循环负责调用函数, t.Run 给每条用例起名, 失败时能直接定位是哪种情况
	// 类似 Vitest 的 test.each:数据表列出情况, t.Run 分别报告结果
	// 放在 _test.go 中的最小表驱动示例:
	// func TestShippingBoundary(t *testing.T) {
	//     cases := []struct { name string; subtotal, want int }{
	//         {"未满门槛", 9999, 600},
	//         {"刚好门槛", 10000, 0},
	//     }
	//     for _, c := range cases {
	//         t.Run(c.name, func(t *testing.T) {
	//             got, err := ShippingFee(c.subtotal, false)
	//             if err != nil { t.Fatal(err) }
	//             if got != c.want { t.Errorf("got %d, want %d", got, c.want) }
	//         })
	//     }
	// }
	// 先覆盖 -1、0、9999、10000, 再分别测试加急与普通配送
	// 9999 与 10000 紧挨免运费阈值, 比只检查 12000 更容易发现 >= 写成 > 的错误

	// ═══════════════════════════════════════════════════════
	// 第 4 节:错误、副作用与资源隔离
	// ═══════════════════════════════════════════════════════
	// 编写 TestSaveReceipt 时, 可用 t.TempDir 创建独立目录, 写入后读回并比较精确内容
	// 传空收据时既检查 errors.Is(err, ErrEmptyReceipt), 也检查旧文件没有被覆盖
	// 这两类断言分别回答“是否报告失败”和“失败后数据是否保持正确”
	// 测试辅助函数用 t.Helper 把失败位置指向调用处, t.Cleanup 在测试及子测试结束后关闭资源
	// TestMain 的 m.Run 负责执行整个包的测试;单个测试资源通常直接用 TempDir/Cleanup 管理

	// ═══════════════════════════════════════════════════════
	// 第 5 节:并行、竞争检查与覆盖率
	// ═══════════════════════════════════════════════════════
	// t.Parallel 表示此测试可与其他允许并行的测试一起执行, 数据和临时目录应互不影响
	// 竞争检测寻找实际执行到的未同步访问;覆盖率统计执行到的语句, 两者都不能替代断言
	// 自行编写测试后, 用 go test -count=1 -cover lesson.go lesson_test.go 运行
	// 显式列出文件, 避免把独立练习程序 exercises.go 的另一个 main 一起编译

	// go test -race 只检查实际运行到的路径;通过并不代表所有并发路径都已被证明安全
	// -run 选择测试/子测试;-count=1 绕过成功测试缓存;-cover 衡量语句执行覆盖率
	// 覆盖率高不意味着断言充分, 例如执行了函数却没有检查结果仍可能“覆盖”

	// ═══════════════════════════════════════════════════════
	// 第 6 节:Benchmark、Fuzz 与修改验证
	// ═══════════════════════════════════════════════════════
	// 可以编写 BenchmarkShippingFee 测量重复调用的成本, -benchmem 额外显示内存分配
	// 可以编写 FuzzShippingFee 检查业务性质:非负订单的加急费用必须比普通费用多 1000 分
	// 普通 go test 只跑有限种子;指定 -fuzz 与 -fuzztime 后才会探索更多输入
	// 可临时把 ShippingFee 的 >= 改成 >, 运行自己编写的 TestShippingBoundary, 核对阈值用例报错后恢复

	// BenchmarkXxx 用于测量;FuzzXxx 根据性质断言探索输入, 普通 go test 只运行它的有限种子
	// 单元测试聚焦隔离的业务行为;真实数据库、外部 API 协作需要另做集成测试
	// 不为测试而暴露所有内部细节;优先断言输入输出、错误类型和必要副作用
}
