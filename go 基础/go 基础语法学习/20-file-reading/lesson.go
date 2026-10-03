/*
章节：20-文件读取
	本章掌握清单：
	1. 准备临时文件与资源清理
	2. ReadFile：一次读取全部字节
	3. Read：循环读取数据块
	4. Scanner：逐行扫描
	5. Reader：保留分隔符与末行
	6. Seek 与路径、编码
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./20-file-reading
// 运行练习：go run ./20-file-reading exam
// 阅读顺序：从 main 进入 demoReading，按六个小节运行；readChunks/readLines/readDelimited 是对应的读取实现。
// TS 对照：os.ReadFile 类似 fs.promises.readFile；io.Reader 是流读取能力的接口。
// Go 文件操作通常直接返回结果和 error；需要并发时由调用方安排 goroutine。
// 相对路径相对于进程当前工作目录 CWD，不是 lesson.go 所在的目录。
// 后续再学：有界流式解析、字符编码转换、压缩数据和随机访问大文件。

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// readChunks 接受接口，既能读文件，也能读 strings.Reader 或 HTTP 响应体。
// 它不负责关闭 reader：创建资源的一方需要明确承担关闭责任。
func readChunks(reader io.Reader) (string, error) {
	buffer := make([]byte, 7)
	var result strings.Builder
	for {
		n, err := reader.Read(buffer)
		// 同一次 Read 可以同时返回 n > 0 和 io.EOF，必须先处理数据。
		if n > 0 {
			result.Write(buffer[:n])
		}
		if errors.Is(err, io.EOF) {
			return result.String(), nil // EOF 表示读完，通常不是业务失败。
		}
		if err != nil {
			return "", fmt.Errorf("读取数据块: %w", err)
		}
	}
	// 分块边界是字节边界，可能截断某个 UTF-8 字符，不能逐块当作完整字符处理。
	// 本示例最后聚合为字符串，仍占用与文件大小相应的内存。
	// 真正处理大文件时，可以读一块处理一块，不把全部内容保留在内存。
}

func readLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开行文本: %w", err)
	}
	defer func() {
		// errors.Join 保留读取与关闭两个错误；两个参数都是 nil 时结果也是 nil。
		err = errors.Join(err, file.Close())
	}()
	scanner := bufio.NewScanner(file)
	// Scanner 有 token 长度限制；按业务设置上限，超长行必须检查 Err。
	// 上限是缓冲上限，可能包含换行等额外字节，不等同于正文长度。
	scanner.Buffer(make([]byte, 1024), 64*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) // 默认按行分割，去掉换行符。
		// scanner.Bytes() 能减少分配，但内容可能在下次 Scan 时被覆盖；要长期保存就复制。
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("按行读取: %w", err)
	}
	return lines, nil
}

// Reader.ReadString 保留分隔符；没有 Scanner 的默认 token 上限，但超长行仍可能占很多内存。
// 对不可信内容应自行限制大小，而不是为了绕过 Scanner 限制就无界读取。
func readDelimited(reader io.Reader) ([]string, error) {
	buffered := bufio.NewReader(reader)
	var lines []string
	for {
		line, err := buffered.ReadString('\n') // ReadBytes('\n') 返回 []byte，错误规则相同。
		if len(line) > 0 {
			lines = append(lines, line) // 最后一行可能没有换行，同时返回数据与 EOF。
		}
		if errors.Is(err, io.EOF) {
			return lines, nil
		}
		if err != nil {
			return nil, fmt.Errorf("按分隔符读取: %w", err)
		}
	}
}

func demoSeek(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	if _, err := file.Seek(7, io.SeekStart); err != nil {
		return err
	}
	buffer := make([]byte, 5)
	if _, err := io.ReadFull(file, buffer); err != nil {
		return err
	}
	fmt.Printf("从字节偏移 7 读取 5 字节：%q\n", buffer) // ready
	// ReadFull 要求填满指定缓冲区；完全没有数据时返回 EOF，读到部分时返回 ErrUnexpectedEOF。
	// SeekStart/SeekCurrent/SeekEnd 分别相对开头/当前位置/结尾；单位都是字节。
	// Seek 适用于支持定位的文件；管道、网络流未必支持，不要把它假定为所有 Reader 的能力。
	// 已用 bufio.Reader 缓冲过内容后，不要直接对底层文件 Seek 再继续读旧缓冲。
	return nil
}

func demoReading() (err error) {

	// ═══════════════════════════════════════════════════════
	// 第 1 节：准备临时文件与资源清理
	// ═══════════════════════════════════════════════════════
	// 先创建只属于本次演示的目录和固定内容，后面的四种读取方式都用它核对结果。
	// defer 在退出 demoReading 时删除临时目录；真正读取的函数各自关闭打开的文件。

	dir, err := os.MkdirTemp("", "go-reading-*")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()
	path := filepath.Join(dir, "events.log")
	if err := os.WriteFile(path, []byte("status=ready\nuser=u-1\n"), 0o600); err != nil {
		return err
	}

	// ═══════════════════════════════════════════════════════
	// 第 2 节：ReadFile：一次读取全部字节
	// ═══════════════════════════════════════════════════════
	// os.ReadFile 已负责打开、读取和关闭，调用方只需检查返回的 error。
	// 适合大小可控的文件；content 是 []byte，%q 可以看到其中的换行。

	content, err := os.ReadFile(path) // 一次性读入 []byte，适合大小可控的配置等文件。
	if err != nil {
		return err
	}
	fmt.Printf("一次性读取：%q\n", content)

	// ═══════════════════════════════════════════════════════
	// 第 3 节：Read：循环读取数据块
	// ═══════════════════════════════════════════════════════
	// 跳到 readChunks：每次 Read 返回本次有效字节数 n，只能处理 buffer[:n]。
	// 先处理数据再判断 EOF，防止丢掉最后一批同时带 EOF 的数据。

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	chunked, readErr := readChunks(file)
	if err := errors.Join(readErr, file.Close()); err != nil {
		return err
	}
	fmt.Printf("分块读取：%q\n", chunked)

	// ═══════════════════════════════════════════════════════
	// 第 4 节：Scanner：逐行扫描
	// ═══════════════════════════════════════════════════════
	// readLines 用 Scan 推进，Text 取当前行，循环结束后 Err 判断是否失败。
	// 正常读完与超长行都可能让 Scan 返回 false，因此必须再检查 Err。

	lines, err := readLines(path)
	if err != nil {
		return err
	}
	fmt.Println("逐行读取：", lines)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：Reader：保留分隔符与末行
	// ═══════════════════════════════════════════════════════
	// ReadString 读到分隔符时把它包含在结果中；最后一行无换行仍需要保存。
	// 对比 Scanner 返回的行内容和这里带 \n 的结果。

	delimited, err := readDelimited(strings.NewReader("first\nlast"))
	if err != nil {
		return err
	}
	fmt.Printf("Reader 保留分隔符且保留末行：%q\n", delimited)

	// ═══════════════════════════════════════════════════════
	// 第 6 节：Seek 与路径、编码
	// ═══════════════════════════════════════════════════════
	// Seek 按字节定位，ReadFull 继续读取准确数量的字节。
	// 文件路径按 CWD 解释；读到的字节也不会自动转换成指定文本编码。

	if err := demoSeek(path); err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("当前工作目录：", cwd)
	// ReadFile/Read/ReadAll 读取的都是字节，不会自动把 GBK 等编码转换成 UTF-8。
	// 二进制文件也能读取，但不应按行拆分；解析固定宽度字段可继续学习 encoding/binary。
	// io.ReadAll(reader) 和 ReadFile 一样会累积全部结果，使用前需要确认内容大小可控。
	return nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := demoReading(); err != nil {
		fmt.Fprintln(os.Stderr, "读取演示失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 把 readChunks 的缓冲区从 7 改为 2，观察最终文本是否一致。
// 给 readDelimited 的最后一行加上换行，比较最后一次返回数据与 EOF 的时机。
