/*
章节：03-基本数据类型
	本章掌握清单：
	1. 整数类型与显式转换
	2. 浮点数与整数除法
	3. 布尔值：条件必须明确
	4. byte 与 rune：字节和码点
	5. 字符串长度与下标
	6. 字符串字面量与转义
	7. 算术、比较与位运算
	8. strconv：在文本与数值之间转换
	9. 修改文本：先转换，再构造新字符串
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./03-basic-types
// 运行练习：go run ./03-basic-types exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：Go 不把整数和小数统一成 number，也不会把 0 或 "" 隐式转成 false。
// 后续再学：切片下标和共享见第 04 章；错误分支见第 05/18 章；复杂 Unicode 分词按业务单独学习。

package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode/utf8"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：整数类型与显式转换
	// ═══════════════════════════════════════════════════════
	// 类型控制可表示的数值范围。int64(pageSize) 是显式转换，原来的 pageSize 仍是 int。
	// 300 的二进制超出 8 位，因此变量转 uint8 后只保留低 8 位，结果为 44。

	// 有符号整数可表示负数；无符号整数只能表示非负数。
	// int8 / int16 / int32 / int64，uint8 / uint16 / uint32 / uint64 指定位数。
	// int 和 uint 的位数与目标架构有关，可能为 32 或 64，不是永远等于 int64。
	count := 42 // 整数字面量赋给无类型约束的新变量时，通常推断为 int。
	var recordID int64 = 9001
	var statusCode uint16 = 200
	fmt.Printf("count=%d 类型=%T；当前 int 位数=%d\n", count, count, strconv.IntSize)
	fmt.Println("记录 ID：", recordID, "状态码：", statusCode)

	// 不同整数类型的变量不能直接相加，需要明确转换。
	// 转到更窄的类型可能丢失高位，生产代码应先验证范围。
	pageSize := 20
	var offset int64 = 100
	fmt.Println("下一页偏移：", offset+int64(pageSize))
	var large int = 300
	fmt.Println("300 转 uint8 后：", uint8(large)) // 44，并非 300。

	// ═══════════════════════════════════════════════════════
	// 第 2 节：浮点数与整数除法
	// ═══════════════════════════════════════════════════════
	// 表达式先按操作数类型计算，再赋给结果。先做 3/4 再转 float64，只会得到 0。
	// 要得到 0.75，必须在除法之前把操作数转为浮点数。

	// Go 只有 float32 / float64，没有名为 float 的内置类型。
	duration := 12.5 // 默认推断为 float64。
	var ratio float32 = 0.75
	fmt.Printf("耗时=%v 类型=%T；比例类型=%T\n", duration, duration, ratio)
	// 整数除法会截去小数部分；应在相除之前转换类型。
	completed, total := 3, 4
	fmt.Println("整数除法：", completed/total)
	fmt.Println("浮点除法：", float64(completed)/float64(total))
	// 浮点数是近似值。金额入门例子宜用整数分存储，如 1999 表示 19.99 元。
	priceCents := 1999
	fmt.Printf("金额：%d.%02d 元\n", priceCents/100, priceCents%100)

	// ═══════════════════════════════════════════════════════
	// 第 3 节：布尔值：条件必须明确
	// ═══════════════════════════════════════════════════════
	// count > 0 的结果才是 bool，count 本身仍是 int。
	// Go 要求把“数量非零”写成 count != 0，不会自动把数值当作条件。

	// bool 只有 true / false，零值为 false。
	var enabled bool
	hasRecords := count > 0
	fmt.Println("功能开启：", enabled, "存在记录：", hasRecords)
	// 不能写 bool(count)、if count 或 if "abc"；条件必须是布尔表达式。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：byte 与 rune：字节和码点
	// ═══════════════════════════════════════════════════════
	// 同一个值既可用 %d 看数值，也可用 %c 看字符。
	// 一个中文码点能放在 rune 中，但它的 UTF-8 编码通常需要多个 byte。

	// byte 是 uint8 的别名，常用于原始字节；rune 是 int32 的别名，常表示 Unicode 码点。
	// 单引号是字符字面量；双引号是字符串，不是可以自由互换的两种引号。
	var marker byte = 'A'
	var character rune = '界'
	fmt.Printf("byte：数值=%d 字符=%c 类型=%T\n", marker, marker, marker)
	fmt.Printf("rune：数值=%d 字符=%c 类型=%T\n", character, character, character)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：字符串长度与下标
	// ═══════════════════════════════════════════════════════
	// Go语言 的字节数是 8，码点数是 4；name[2] 取到“语”的第一个字节。
	// 先转 []rune 再索引，才能按码点取得完整的“语”。

	// 字符串是不可变的字节序列，不保证所有字符串都包含有效 UTF-8。
	// 常见的文本字符串采用 UTF-8；len 返回字节数，索引 s[i] 得到单个 byte。
	name := "Go语言"
	fmt.Printf("文本=%s 字节数=%d 码点数=%d\n", name, len(name), utf8.RuneCountInString(name))
	fmt.Printf("第一个字节：%c；中文首字节：0x%x\n", name[0], name[2])
	runes := []rune(name)
	fmt.Printf("第三个码点：%c\n", runes[2])
	// 对照 TS：JavaScript 的 string.length 按 UTF-16 代码单元计数，也不等于“可见字符数”。
	// emoji 组合、组合重音等可由多个 rune 构成；rune 数量也不是用户感知的字符数量。

	// ═══════════════════════════════════════════════════════
	// 第 6 节：字符串字面量与转义
	// ═══════════════════════════════════════════════════════
	// 双引号里的 \n 是一个换行字符；反引号中的反斜杠保留原样。
	// %q 会把换行显示成转义形式，适合核对字符串实际包含什么。

	// 双引号支持转义，反引号创建原始字符串，可包含换行但不解释 \n。
	jsonText := `{"service":"api","ready":true}`
	multiline := "first\nsecond"
	rawPath := `C:\logs\api.txt`
	fmt.Println("JSON 文本：", jsonText)
	fmt.Printf("转义换行：%q；原始路径：%s\n", multiline, rawPath)

	// 数字转十进制文本用 strconv.Itoa；string(65) 表示码点 U+0041，即 "A"。
	fmt.Println("数字转文本：", strconv.Itoa(65), "码点转文本：", string(rune(65)))

	// ═══════════════════════════════════════════════════════
	// 第 7 节：算术、比较与位运算
	// ═══════════════════════════════════════════════════════
	// 先用括号理解运算顺序，再看权限位：01 | 10 得到 11，表示同时拥有两个标记。
	// permissions & readFlag != 0 是检查读标记是否存在，结果是 bool。

	// 运算符：乘除取余比加减优先，比较比 && 优先，&& 比 || 优先。
	// 关系复杂时使用括号表达意图；Go 没有 ** 幂运算符，^ 是按位异或。
	fmt.Println("优先级：", 2+3*4, (2+3)*4)
	fmt.Println("负数整除与取余：", -7/3, -7%3) // -2、-1，商向零截断，余数与被除数同号或为 0。
	// 整数除数为零会失败；浮点转整数同样向零截断，不等于四舍五入。
	fraction := -3.8
	fmt.Println("浮点转 int：", int(fraction))
	// 常用整数位运算：& 与、| 或、^ 异或、&^ 清除位、<< 左移、>> 右移。
	const readFlag = 1 << 0
	const writeFlag = 1 << 1
	permissions := readFlag | writeFlag
	fmt.Printf("权限位=%02b 可读=%t 去掉写权限=%02b\n", permissions, permissions&readFlag != 0, permissions&^writeFlag)
	// 正常的 Web 字段无需强行编码成位标记；认识位运算即可。

	// ═══════════════════════════════════════════════════════
	// 第 8 节：strconv：在文本与数值之间转换
	// ═══════════════════════════════════════════════════════
	// 类型转换处理已经有类型的值，解析则要判断文本是否符合数字格式。
	// Atoi 返回 (int, error)；文本非法时必须看错误，不能把返回的 0 当作有效输入。

	// 请求参数是字符串，转换成数字需要解析，不是简单的类型转换。
	parsedSize, parseErr := strconv.Atoi("25") // 十进制字符串 -> int。
	badSize, badErr := strconv.Atoi("25px")
	parsedID, idErr := strconv.ParseInt("9000000000", 10, 64)
	fmt.Println("解析 pageSize：", parsedSize, parseErr, "错误输入：", badSize, badErr)
	fmt.Println("解析 int64 ID：", parsedID, idErr, "格式化：", strconv.FormatInt(parsedID, 10))
	// base=10 明确按十进制；base=0 会识别 0x、0o 以及前导 0 等进制前缀。
	// Atoi/ParseInt 对非法文本或溢出返回错误；不能只用返回的数值而忽略错误。
	parsedRatio, ratioErr := strconv.ParseFloat("0.75", 64)
	parsedBool, boolErr := strconv.ParseBool("true")
	fmt.Println("解析浮点与布尔：", parsedRatio, ratioErr, parsedBool, boolErr)

	// ═══════════════════════════════════════════════════════
	// 第 9 节：修改文本：先转换，再构造新字符串
	// ═══════════════════════════════════════════════════════
	// 字符串不能写 name[0] = ...；字节切片或 rune 切片才允许修改元素。
	// 修改完成后用 string(...) 生成结果，原来的字符串值不会被原地改变。

	// 字符串不能原地修改。转为 []byte 适合字节处理，转为 []rune 适合按码点处理。
	textBytes := []byte("api")
	textBytes[0] = 'A'
	textRunes := []rune("猫服务")
	textRunes[0] = '云'
	fmt.Println("修改后的文本：", string(textBytes), string(textRunes))
	invalidText := string([]byte{0xff})
	fmt.Println("是否有效 UTF-8：", utf8.ValidString(invalidText), "转换 rune：", []rune(invalidText))
	// 无效 UTF-8 转为 rune 时会得到替换字符 U+FFFD；[]rune 不是任意二进制数据的无损容器。
	// 复数类型还有 complex64 / complex128；一般 Web 业务很少直接使用。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 分别计算 float64(3/4) 和 float64(3)/4，解释结果为什么不同。
// 把文本换成 "A中😀"，先猜字节数和码点数，再运行核对。
