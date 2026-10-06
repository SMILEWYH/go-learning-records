/*
章节：21-文件写入
	本章掌握清单：
	1. 临时目录与写入前的准备
	2. WriteFile 与截断覆盖
	3. O_APPEND：向末尾追加
	4. bufio.Writer：先 Flush，再 Close
	5. 临时文件替换与持久化边界
	6. 打开标志与权限的选择
*/

// 请从 go 基础语法学习 项目根目录运行以下命令，同目录文件共同组成 main 包。
// 运行讲解：go run ./21-file-writing
// 运行练习：go run ./21-file-writing exam
// 阅读顺序：先读 demoWriting 的六个小节，再回看 writeWithFlags、writeBuffered、replaceWithTemp 的实现。
// TS 对照：WriteFile 类似 fs.writeFile，O_APPEND 类似 appendFile 的用途。
// Go 使用显式打开标志与权限位，写入错误和关闭错误都需要处理。
// 后续再学：掉电一致性、目录 Sync、跨平台替换策略、跨进程文件写入协议。

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func writeWithFlags(path, content string, flags int) (err error) {
	file, err := os.OpenFile(path, flags, 0o600)
	if err != nil {
		return fmt.Errorf("打开写入文件: %w", err)
	}
	defer func() {
		// 延迟关闭不能直接 defer file.Close() 后完全丢弃关闭错误。
		err = errors.Join(err, file.Close())
	}()
	n, err := file.WriteString(content)
	if err != nil {
		return fmt.Errorf("写入内容: %w", err)
	}
	if n != len(content) {
		return io.ErrShortWrite
	}
	return nil
}

func writeBuffered(path string, lines []string) (err error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return fmt.Errorf("写入缓冲区: %w", err)
		}
	}
	// Close 底层文件不会自动刷新 bufio.Writer！必须先 Flush，再 Close，并检查两者错误。
	// Flush 只是把 Go 的用户态缓冲交给底层 Writer，不等于文件已经持久化到磁盘。
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("刷新缓冲区: %w", err)
	}
	return nil
}

// replaceWithTemp 先写完整临时文件，再提交目标；示例只接收程序自己控制的临时目录路径。
// 创建在同目录，避免跨文件系统 Rename；新文件沿用 CreateTemp 的 0600 权限，不继承旧权限。
// Unix 的同文件系统 rename 通常提供原子名称替换；非 Unix 平台即使同目录也不承诺原子性。
// 原子可见性和掉电持久性是两件事：这里 Sync 文件，但没有同步父目录，不承诺掉电不丢更新。
func replaceWithTemp(path string, content []byte) (err error) {
	file, err := os.CreateTemp(filepath.Dir(path), ".replacement-*")
	if err != nil {
		return err
	}
	temporaryPath := file.Name()
	closed := false
	defer func() {
		if !closed {
			err = errors.Join(err, file.Close())
		}
		// 成功 Rename 后临时路径已不存在；失败时删除本次创建的半成品。
		if removeErr := os.Remove(temporaryPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			err = errors.Join(err, removeErr)
		}
	}()
	if _, err := file.Write(content); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	closeErr := file.Close()
	closed = true
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(temporaryPath, path)
}

func demoWriting() (err error) {

	// ═══════════════════════════════════════════════════════
	// 第 1 节：临时目录与写入前的准备
	// ═══════════════════════════════════════════════════════
	// demoWriting 创建独立临时目录，再依次展示覆盖、追加、缓冲和替换。
	// 辅助函数把打开、写入和关闭放在一次调用内，失败时把错误交回上层。

	dir, err := os.MkdirTemp("", "go-writing-*")
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, os.RemoveAll(dir)) }()

	// ═══════════════════════════════════════════════════════
	// 第 2 节：WriteFile 与截断覆盖
	// ═══════════════════════════════════════════════════════
	// WriteFile 适合一次写入完整内容；OpenFile 的 flags 则让你明确控制打开方式。
	// O_WRONLY 表示写，O_TRUNC 表示打开时清空旧内容，两者用 | 组合。

	configPath := filepath.Join(dir, "config.json")
	// WriteFile 不存在就创建，存在就截断并覆盖；不是追加，也不是原子替换。
	// 0o 是八进制前缀。Unix 下 0600 为所有者可读写；0644 再允许其他人读取。
	// 实际新建权限还受 umask 影响；mode 不会重设已有文件权限。
	// Windows 的权限语义与 Unix 不同，不能用这些位代替 Windows ACL 配置。
	if err := os.WriteFile(configPath, []byte("{\"mode\":\"development\"}\n"), 0o600); err != nil {
		return err
	}
	if err := writeWithFlags(configPath, "{\"mode\":\"test\"}\n", os.O_WRONLY|os.O_TRUNC); err != nil {
		return err
	}
	config, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	fmt.Printf("截断覆盖后：%s", config) // 更短的新内容不会残留旧文件尾部。

	// ═══════════════════════════════════════════════════════
	// 第 3 节：O_APPEND：向末尾追加
	// ═══════════════════════════════════════════════════════
	// O_CREATE 在文件不存在时创建，O_APPEND 让每次写入落到末尾。
	// 因此第二条事件保留第一条记录，适合本例的顺序追加日志。

	logPath := filepath.Join(dir, "audit.log")
	for _, event := range []string{"user.created\n", "user.updated\n"} {
		if err := writeWithFlags(logPath, event, os.O_WRONLY|os.O_CREATE|os.O_APPEND); err != nil {
			return err
		}
	}
	log, err := os.ReadFile(logPath)
	if err != nil {
		return err
	}
	fmt.Printf("追加写入：\n%s", log)

	// ═══════════════════════════════════════════════════════
	// 第 4 节：bufio.Writer：先 Flush，再 Close
	// ═══════════════════════════════════════════════════════
	// 先写内存缓冲，最后 Flush 交给底层文件；Close 文件不会代你刷新 Go 缓冲。
	// 进入 writeBuffered，注意 Write、Flush、Close 各自都可能返回错误。

	bufferedPath := filepath.Join(dir, "buffered.txt")
	if err := writeBuffered(bufferedPath, []string{"first", "second"}); err != nil {
		return err
	}
	buffered, err := os.ReadFile(bufferedPath)
	if err != nil {
		return err
	}
	fmt.Printf("刷新缓冲后的内容：%q\n", buffered)

	// ═══════════════════════════════════════════════════════
	// 第 5 节：临时文件替换与持久化边界
	// ═══════════════════════════════════════════════════════
	// replaceWithTemp 依次创建临时文件、写内容、Sync、Close，最后 Rename 提交目标名称。
	// 失败清理由本次创建的临时文件承担；是否能原子替换还取决于操作系统。

	if err := replaceWithTemp(configPath, []byte("{\"mode\":\"ready\"}\n")); err != nil {
		return err
	}
	replaced, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	fmt.Printf("同目录临时文件替换后：%s", replaced)

	// ═══════════════════════════════════════════════════════
	// 第 6 节：打开标志与权限的选择
	// ═══════════════════════════════════════════════════════
	// 权限描述谁可以访问，flags 描述本次怎样打开；这两组参数承担不同职责。
	// 对已有文件写入时，决定是否清空的是 O_TRUNC，不能只看传入的权限位。

	// O_RDONLY / O_WRONLY / O_RDWR 选择访问方式；O_CREATE 创建不存在的文件；
	// O_TRUNC 打开时清空；O_APPEND 每次写入定位到末尾。
	// O_CREATE|O_EXCL 可要求目标必须不存在，避免意外覆盖已有文件。
	// 不加 O_TRUNC 的普通写入从开头覆盖，可能留下原文件更长的尾部。
	// Close 不等于保证掉电不丢数据；需要持久化保证时还要理解 Sync 和文件系统行为。
	// 权限控制访问，不能替代路径校验；filepath.Join 并不会把用户路径限制在指定目录内。
	return nil
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "exam" {
		runExercises()
		return
	}
	if err := demoWriting(); err != nil {
		fmt.Fprintln(os.Stderr, "写入演示失败：", err)
		os.Exit(1)
	}
}

// ── 学习方式：先预测，再运行，最后动手修改 ──
// 在临时目录示例里去掉 O_TRUNC，观察较短的新配置后面是否残留旧内容。
// 沿 writeBuffered 找到 Flush 和 Close，解释为什么只关闭底层文件还不够。
