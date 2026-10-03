/*
章节：02-输入输出
	本章掌握清单：
	1. Print、Println 与 Printf：直接输出
	2. 格式占位符：让值按指定形式显示
	3. Sprint 系列：把结果留在字符串里
	4. Fprint 系列：自己选择输出目标
	5. Append 系列：把文本追加到字节切片
	6. Scan 与 Fscan：按空白扫描输入
	7. 整行读取：保留空格与最后一行
	8. 标准输入、标准输出与标准错误
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./02-input-output
// 运行练习：go run ./02-input-output exam
// 阅读顺序：从 main 开始按小节阅读；遇到辅助函数时，再按调用名称找到它的定义。
// TS 对照：Printf 类似带格式约束的 console.log；Sprintf 类似生成字符串的模板表达式。
// 演示使用内存输入，不会停下来等待键盘输入。
// 后续再学：Reader/Writer 接口的定义见第 13 章；错误处理见第 18 章；大文件逐行读取见第 20 章。

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}

	// ═══════════════════════════════════════════════════════
	// 第 1 节：Print、Println 与 Printf：直接输出
	// ═══════════════════════════════════════════════════════
	// 先分清三个名字：Print 组合输出参数，Println 适合一行结果，Printf 适合固定格式。
	// 本节最后的 "\n" 是换行字符；Printf 也需要自己写换行。

	// Print 不自动换行，Println 用空格分隔参数并在结尾换行。
	// Print 的细节：相邻两个参数都不是字符串时，会自动加空格。
	// 因此“Print 永远不加空格”不准确，需要固定布局时使用 Printf。
	fmt.Print("服务：", "order-api", "\n")
	fmt.Print(1, 2, "\n")
	fmt.Println("请求：", "GET", "/orders")

	// ═══════════════════════════════════════════════════════
	// 第 2 节：格式占位符：让值按指定形式显示
	// ═══════════════════════════════════════════════════════
	// 格式字符串中的占位符按顺序对应后面的参数，%d 要求对应参数是整数。
	// %% 只输出一个百分号，不消耗参数。先预测 12.345 用 %.2f 显示时保留几位小数。

	// 占位符决定输出形式，不改变被输出变量的类型。
	status, elapsed, cached := 200, 12.345, true
	fmt.Printf("默认值 %%v：%v\n", status)
	fmt.Printf("类型 %%T：%T\n", status)
	fmt.Printf("整数 %%d：%d\n", status)
	fmt.Printf("浮点 %%f：%f；两位小数 %%.2f：%.2f\n", elapsed, elapsed)
	fmt.Printf("字符串 %%s：%s；布尔 %%t：%t\n", "/orders", cached)
	fmt.Printf("带引号 %%q：%q；Go 语法 %%#v：%#v\n", "a\nb", "hello")
	fmt.Printf("进度：100%%\n")
	// %+v 对结构体会额外显示字段名，第 11 章再看具体例子。
	// 占位符与类型不匹配时往往会输出 %!d(...) 一类诊断；go vet 也能检查常见问题。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：Sprint 系列：把结果留在字符串里
	// ═══════════════════════════════════════════════════════
	// 可以按名字记忆：Print 是输出，前面加 S 后就是生成字符串。
	// 例如日志正文先用 Sprintf 构造，再决定保存到哪里。

	// Sprint / Sprintln / Sprintf 返回字符串，不会直接打印。
	label := fmt.Sprint("service=", "order-api")
	line := fmt.Sprintln("ready", true) // 返回值自身包含结尾换行。
	message := fmt.Sprintf("GET %s -> %d (%.1fms)", "/orders", status, elapsed)
	fmt.Println(label)
	fmt.Print(line)
	fmt.Println(message)
	// 返回的字符串可以保存到变量。Sprintf 等普通函数调用不能用于初始化 Go 常量。

	// ═══════════════════════════════════════════════════════
	// 第 4 节：Fprint 系列：自己选择输出目标
	// ═══════════════════════════════════════════════════════
	// Fprintf 的第一个参数是写到哪里，第二个才是格式字符串。
	// 这里写到内存，所以调用结束后还要用 buffer.String() 才能取出内容。

	// Fprint / Fprintln / Fprintf 向指定目标写入。
	// bytes.Buffer 是内存缓冲区，后面也可以换成文件或 HTTP 响应目标。
	// &buffer 表示把 buffer 的地址交给函数，第 09 章会解释指针。
	var buffer bytes.Buffer
	written, writeErr := fmt.Fprintf(&buffer, "request_id=%s status=%d", "req-42", status)
	fmt.Println("写入字节数：", written, "错误：", writeErr)
	fmt.Println("缓冲区内容：", buffer.String())

	// ═══════════════════════════════════════════════════════
	// 第 5 节：Append 系列：把文本追加到字节切片
	// ═══════════════════════════════════════════════════════
	// []byte 表示一组字节，先把它看作可以继续追加的文本载体。
	// payload = ... 中的赋值很重要：函数返回的是追加后的切片。

	// Append / Appendln / Appendf 返回追加后的 []byte。
	// 必须接收返回值；切片语法和 append 在第 04 章深入学习。
	payload := []byte("log: ")
	payload = fmt.Appendf(payload, "status=%d", status)
	payload = fmt.Appendln(payload)
	fmt.Print(string(payload))

	// ═══════════════════════════════════════════════════════
	// 第 6 节：Scan 与 Fscan：按空白扫描输入
	// ═══════════════════════════════════════════════════════
	// &username、&pageSize 把变量地址交出去，扫描函数才能填入结果。
	// 输入 alice invalid 时，用户名可以先成功赋值，而第二项解析失败；观察 n=1。

	// Scan 从终端读，Fscan 从我们指定的输入源读。
	// 扫描按空白分隔，所以 "Ada Lovelace" 默认会被当作两个词。
	// n 是成功赋值的项目数，不是输入字符数。
	reader := strings.NewReader("alice 25")
	var username string
	var pageSize int
	n, err := fmt.Fscan(reader, &username, &pageSize)
	fmt.Printf("扫描：n=%d err=%v username=%q pageSize=%d\n", n, err, username, pageSize)

	// 错误输入不会自动变成有效值。应观察 n、err 及目标变量当前的值。
	badReader := strings.NewReader("alice invalid")
	var badName string
	var badSize int
	n, err = fmt.Fscan(badReader, &badName, &badSize)
	fmt.Printf("错误输入：n=%d err=%v name=%q size=%d\n", n, err, badName, badSize)

	// ═══════════════════════════════════════════════════════
	// 第 7 节：整行读取：保留空格与最后一行
	// ═══════════════════════════════════════════════════════
	// 按词读取适合“姓名代号 + 数量”，按行读取适合包含空格的一整段姓名。
	// 本例第二行没有换行，返回的错误是 EOF，但 lastLine 仍包含有效文本。

	// 需要保留姓名中的空格时，应按整行读取，而不是把一行切成词元。
	// Scanln 同样按空白扫描，只是到换行停止，不能直接代替“读取完整的一行”。
	lineReader := bufio.NewReader(strings.NewReader("Ada Lovelace\nGrace Hopper"))
	fullLine, lineErr := lineReader.ReadString('\n')
	fmt.Printf("首行包含分隔符：%q err=%v\n", fullLine, lineErr)
	lastLine, lastErr := lineReader.ReadString('\n')
	fmt.Printf("末行无换行：%q err=%v\n", lastLine, lastErr)
	// 最后一次同时返回有效的 "Grace Hopper" 和 EOF；处理错误时不能把已读到的数据丢掉。
	// 真实输入可能使用 CRLF 换行，是否去掉 \r\n 应由你的文本格式约定决定。

	// ═══════════════════════════════════════════════════════
	// 第 8 节：标准输入、标准输出与标准错误
	// ═══════════════════════════════════════════════════════
	// os.Stdin、os.Stdout、os.Stderr 是三个不同的流。
	// 将程序正常结果和诊断分开，方便在终端中分别重定向或接到其他程序。

	// 输入源和输出目标可以替换：Reader 提供读能力，Writer 提供写能力。
	// strings.Reader、bytes.Buffer、文件、HTTP 请求体等常通过统一接口参与 I/O。
	// 这里先理解“可替换的输入/输出目标”，接口的具体定义在第 13 章学习。
	fmt.Fprintln(os.Stdout, "标准输出：程序的正常结果")
	fmt.Fprintln(os.Stderr, "标准错误：演示诊断信息，不表示程序运行失败")
	// fmt.Print 默认写 os.Stdout；fmt.Scan 默认读 os.Stdin。
	// 标准输入不一定来自键盘，也可以来自管道或文件重定向。
	// ReadString 可换成 bufio.NewReader(os.Stdin)，但本例不这样做以避免等待输入。

	// 若想手动体验终端输入，可自己临时添加：
	// n, err := fmt.Scan(&username, &pageSize)
	// fmt.Println(n, err, username, pageSize)
	// fmt.Scan 会等待标准输入；正式服务一般从 HTTP 请求、配置或数据库接收数据。
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 elapsed 改成 1.2，比较 %f、%.2f 和 %g；它们如何显示同一个值？
// 把输入改成 "alice 0"，观察合法的 0 与 invalid 的区别；再试含空格的姓名。
