//go:build ignore

package main

/*
练习 1:更新接口中的可选分页参数

背景:前端未填写 pageSize 时用默认值, 但输入 0 也属于需要纠正的值
本题用 *int 区分“未提供 nil”与“明确提供了一个数值”

要求:
  - nil 返回默认值 20
  - 指向的值小于 1 时, 将原变量改成 1, 并返回 1
  - 大于 100 时, 将原变量改成 100, 并返回 100
  - 其他数值保持原样并返回
  - 不解引用 nil, 也不通过重新指向局部变量假装修改调用方

建议签名:func normalizePageSize(size *int) int

手动核对:
  - size = 30:返回 30, 原变量仍为 30
  - size = 0:返回 1, 原变量变成 1
  - size = 200:返回 100, 原变量变成 100
  - nil:返回 20, 不发生 panic
*/

/*
练习 2:为批量导入任务制作独立的重试次数列表

背景:界面展示原始重试次数, 任务准备阶段需要每项加 1, 不能修改展示数据

要求:
  - 创建并返回独立的 []int, 每个值为原值加 1
  - 不修改传入的 slice;返回后修改结果的元素也不能影响原 slice
  - 输入 nil 或空 slice 时, 返回长度为 0 的 slice;不限定是否为 nil
  - 说明直接执行 result := source 为什么不能隔离元素修改

建议签名:func nextAttempts(source []int) []int

手动核对:
  - 输入 [0, 2, 4], 返回 [1, 3, 5], 输入仍为 [0, 2, 4]
  - 把结果第一个元素改成 99, 输入第一个元素仍为 0
  - 输入 [] 或 nil, 结果长度为 0
*/

/*
练习 3:合并用户的界面设置

背景:后端根据 PATCH 请求更新少量设置, 没有提交的键保持原值

要求:
  - 将 patch 中所有键值写入 current;已存在的键覆盖, 不存在的键新增
  - current 为 nil 时创建 map, 返回最终可用的 map;调用方接收返回值
  - current 非 nil 时原地更新, 同一 map 的其他引用能看到更改
  - patch 为 nil 时保留已有数据;即使两者都为 nil, 返回值也可写入

建议签名:
func applySettings(current, patch map[string]string) map[string]string

手动核对:
  - current={theme:light, locale:zh}, patch={theme:dark}:结果为 {theme:dark, locale:zh}
  - current=nil, patch={locale:en}:结果包含 locale=en
  - 两者 nil:结果长度为 0, 执行 result["theme"] = "dark" 不会 panic
  - 打印 map 时不要依赖键的顺序
*/

func main() {
}
