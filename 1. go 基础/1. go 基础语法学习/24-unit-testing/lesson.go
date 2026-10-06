/*
章节：24-单元测试
	本章掌握清单：
	1. 测试对象：先明确输入和期望结果
	2. 测试文件、测试函数与断言
	3. 表驱动测试、子测试与边界
	4. 错误、副作用与资源隔离
	5. 并行、竞争检查与覆盖率
	6. Benchmark、Fuzz 与修改验证
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./24-unit-testing
// 运行练习：go run ./24-unit-testing exam
// 阅读顺序：先读本文件的规则与六个小节，再在 exercises_test.go 中依次查找 TestShippingFee、TestSaveReceipt、BenchmarkShippingFee、FuzzShippingFee。
// TS 对照：testing 对应 Vitest/Jest 的基础测试能力；t.Run 类似具名子测试。
// Go 通常使用标准库测试，不需要先安装第三方测试框架。
// 测试：go test -v ./24-unit-testing
// 指定子测试：go test -v ./24-unit-testing -run 'TestShippingFee/free_standard'
// 检查竞争：go test -race ./24-unit-testing
// 重新执行并看覆盖率：go test -count=1 -cover ./24-unit-testing
// 基准入门：go test ./24-unit-testing -run '^$' -bench BenchmarkShippingFee -benchmem
// 有限时间模糊测试：go test ./24-unit-testing -run '^$' -fuzz '^FuzzShippingFee$' -fuzztime=1s -parallel=2
// 后续再学：HTTP/数据库集成测试、测试替身、性能剖析、CI。这里的基准只展示写法，不给性能结论。

package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrEmptyReceipt = errors.New("收据不能为空")

// ShippingFee 用整数“分”计算运费，适合用纯函数单元测试验证业务规则。
// 订单不足 100 元：普通运费 6 元；满 100 元免普通运费；加急额外收 10 元。
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

// SaveReceipt 的文件行为由 exercises_test.go 使用 t.TempDir 隔离验证。
func SaveReceipt(path, receipt string) error {
	if receipt == "" {
		return ErrEmptyReceipt
	}
	return os.WriteFile(path, []byte(receipt+"\n"), 0o600)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：测试对象：先明确输入和期望结果
	// ═══════════════════════════════════════════════════════
	// ShippingFee 是纯函数，同样的金额与加急选项应返回同样的运费，适合直接测试。
	// 本例订单 120 元满足免普通运费条件，加急再加 1000 分，因此最终结果是 1000。

	fee, err := ShippingFee(12000, true)
	if err != nil {
		fmt.Println("运费计算失败：", err)
		return
	}
	fmt.Println("订单 120 元、加急配送，运费（分）：", fee) // 1000

	// ═══════════════════════════════════════════════════════
	// 第 2 节：测试文件、测试函数与断言
	// ═══════════════════════════════════════════════════════
	// 打开同目录 exercises_test.go：测试函数形如 func TestShippingFee(t *testing.T)。
	// 一个最小断言可写成下面这样，放在 _test.go 文件内并导入 testing：
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
	// 测试必须检查期望；只打印结果，即使算错也可能显示 PASS。

	fmt.Println("打开 exercises_test.go 学习真实测试，再执行 go test -v ./24-unit-testing")
	// _test.go 文件只在 go test 时参与编译，不会被 go run 的 main 调用。
	// 测试函数必须形如 TestXxx(t *testing.T)，失败时用 t.Error/t.Fatal 等报告。
	// t.Fatal 停止当前测试，t.Error 记录失败后继续；不要用 Println 代替断言。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：表驱动测试、子测试与边界
	// ═══════════════════════════════════════════════════════
	// exercises_test.go 中的 tests 把 name、输入、want、wantErr 放到每一行数据中。
	// 循环负责调用函数，t.Run 给每条用例起名，失败时能直接定位是哪种情况。
	// 先覆盖 -1、0、9999、10000，再分别测试加急与普通配送。
	// 9999 与 10000 紧挨免运费阈值，比只检查 12000 更容易发现 >= 写成 > 的错误。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：错误、副作用与资源隔离
	// ═══════════════════════════════════════════════════════
	// TestSaveReceipt 用 t.TempDir 创建独立目录，写入后读回并比较精确内容。
	// 传空收据时既检查 errors.Is(err, ErrEmptyReceipt)，也检查旧文件没有被覆盖。
	// 这两类断言分别回答“是否报告失败”和“失败后数据是否保持正确”。
	// openTestFile 中 t.Helper 把失败位置指向调用处，t.Cleanup 在测试及子测试结束后关闭文件。
	// TestMain 的 m.Run 负责执行整个包的测试；单个测试资源通常直接用 TempDir/Cleanup 管理。

	// ═══════════════════════════════════════════════════════
	// 第 5 节：并行、竞争检查与覆盖率
	// ═══════════════════════════════════════════════════════
	// t.Parallel 表示此测试可与其他允许并行的测试一起执行，数据和临时目录应互不影响。
	// 竞争检测寻找实际执行到的未同步访问；覆盖率统计执行到的语句，两者都不能替代断言。
	// 只验证本章现有讲解测试，可运行：
	//   go test ./24-unit-testing -run "^(TestShippingFee|TestSaveReceipt|TestSaveReceiptMissingDirectory)$" -count=1 -cover
	// 本章 exercises_test.go 包含待完成练习；直接运行全部测试可能因未填写练习而失败。

	// go test -race 只检查实际运行到的路径；通过并不代表所有并发路径都已被证明安全。
	// -run 选择测试/子测试；-count=1 绕过成功测试缓存；-cover 衡量语句执行覆盖率。
	// 覆盖率高不意味着断言充分，例如执行了函数却没有检查结果仍可能“覆盖”。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：Benchmark、Fuzz 与修改验证
	// ═══════════════════════════════════════════════════════
	// BenchmarkShippingFee 测量重复调用的成本，-benchmem 额外显示内存分配。
	// FuzzShippingFee 检查业务性质：非负订单的加急费用必须比普通费用多 1000 分。
	// 普通 go test 只跑有限种子；指定 -fuzz 与 -fuzztime 后才会探索更多输入。
	// 可临时把 ShippingFee 的 >= 改成 >，运行 TestShippingFee，核对阈值用例报错后恢复。

	// BenchmarkXxx 用于测量；FuzzXxx 根据性质断言探索输入，普通 go test 只运行它的有限种子。
	// 单元测试聚焦隔离的业务行为；真实数据库、外部 API 协作需要另做集成测试。
	// 不为测试而暴露所有内部细节；优先断言输入输出、错误类型和必要副作用。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 先读 TestShippingFee 的阈值用例，再故意引入边界错误，确认测试能发现它并恢复代码。
// 区分“执行到了某行”“检查了正确结果”和“性能是否满足要求”，它们需要不同证据。
