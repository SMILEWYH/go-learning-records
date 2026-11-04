//go:build ignore

package main

// 练习 1:准备按年月归档的上传目录
// 建议:func ensureUploadDirectory(root string, year, month int) (string, error)
// 要求:year 为 2000..2100、month 为 1..12;得到 root/uploads/YYYY/MM;使用 filepath.Join/MkdirAll
// 常规:2026, 9 -> 目录 uploads/2026/09;重复调用成功, 已有文件保留
// 边界:month=0、13 或非法年份 -> error 且不创建子目录;目标路径上已有普通文件 -> error
// root 仅由本题的临时目录提供, 年份月份均验证后再转换成路径段

// 练习 2:列出可在前端展示的本地图片文件
// 建议:func listImages(root string) ([]string, error)
// 要求:只读根目录一层;返回普通文件的文件名;支持 .png/.jpg/.webp, 扩展名忽略大小写;按名称排序
// 常规:a.png、b.JPG、readme.txt -> ["a.png","b.JPG"]
// 边界:目录名 fake.png、符号链接 link.png 都不计入;空目录 -> 长度 0;根目录不存在 -> error
// 通过 DirEntry.Type/Info 判断类型并处理 Info 错误, 不递归进入子目录

// 练习 3:统计构建产物体积
// 建议:func artifactBytes(root string) (int64, error)
// 要求:WalkDir 递归累加普通文件的 Size;跳过名为 cache 的子目录;忽略符号链接
// 常规:a.txt 为 5 字节、nested/b.txt 为 3 字节、cache/x.txt 为 99 字节 -> 8
// 边界:空目录 -> 0;路径不存在或获取 Info 失败 -> error, 出错时不能报告成功统计
// 回调先检查传入 error, 再访问 DirEntry;遍历起点 root 必须是目录

func main() {
}
