//go:build ignore

package main

// 练习 1：导出文件的首次备份。
// 建议：func backupNew(sourcePath, destinationPath string) (int64, error)。
// 要求：源为普通文件；目标使用 O_CREATE|O_EXCL|O_WRONLY，不覆盖已有文件。
// 常规：源为 5 字节 "hello"，新目标获得同样内容，返回 5, nil。
// 边界：空源成功且返回 0；源缺失返回 error。
// 目标已存在返回可被 errors.Is(err, os.ErrExist) 识别的错误，且原内容不变。
// 注意：先确认源可读取，再创建目标；若复制失败，明确返回错误，调用者不能使用半成品。

// 练习 2：复制上传流，同时计算 SHA-256，供内容一致性核验使用。
// 建议：func copyAndHash(destination io.Writer, source io.Reader) (int64, string, error)。
// 要求：使用 io.Copy、io.MultiWriter 和 sha256.New，一次遍历完成；返回十六进制摘要。
// 常规：输入 "abc"，写入 3 字节；摘要等于 fmt.Sprintf("%x", sha256.Sum256([]byte("abc")))。
// 边界：空输入写入 0 字节且有合法摘要；目标写入失败应返回 error，此时摘要返回空串。
// 用 bytes.Buffer 演示正常路径，自己写一个始终返回错误的 Writer 演示故障路径。
// 本函数不关闭调用者传入的接口；创建文件的一方负责关闭文件。

// 练习 3：覆盖更新本地构建产物。
// 建议：func replaceArtifact(sourcePath, destinationPath string) (int64, error)。
// 要求：自己实现，先检查已打开文件的身份，再截断目标，不能直接调用讲解中的 copyFile。
// 常规：用 "new" 覆盖较长的 "old-content" -> 目标恰好为 "new"，返回 3。
// 边界：源与目标同路径 -> error 且源不变；目标是源的硬链接也必须拒绝；源缺失时目标不变。
// 硬链接用 os.Link 在临时目录创建；环境不支持时清楚输出跳过原因。
// 写一句注释：复制中途失败时，本函数为什么不能保证目标保留旧内容？

func main() {
}
