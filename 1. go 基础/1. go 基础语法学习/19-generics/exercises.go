//go:build ignore

package main

// 练习 1:接口结果筛选写一个通用函数供订单、商品等列表使用
// 建议:func Filter[T any](items []T, keep func(T) bool) []T
// 要求:保持原顺序, 不修改输入切片;回调约定非 nil
// 常规:[1,2,3,4] 按偶数筛选 -> [2,4];再用含 Active 字段的用户切片筛选启用用户
// 边界:nil 输入、全部不满足条件 -> 长度为 0;筛选后修改结果元素不能改变输入元素
// 本题元素用 int 或只含值字段的结构体即可, 不要求深拷贝嵌套指针

// 练习 2:把接口返回的数组转换为按 ID 查找的索引
// 建议:func IndexBy[T any, K comparable](items []T, keyOf func(T) K) map[K]T
// 要求:返回非 nil 的 map;相同 key 后出现的记录覆盖前面的;回调约定非 nil
// 常规:[{ID:"u-1",Name:"旧名"},{ID:"u-1",Name:"新名"}] -> key "u-1" 对应 "新名"
// 边界:空输入 -> 可直接写入的空 map;空字符串 ID 是合法 key, 按同样覆盖规则处理
// 再用整数 ID 验证一次;说明为什么 K 需要 comparable

// 练习 3:内存列表分页自行定义 Page[T any], 字段为 Items []T、Total int, 表达返回结果
// 建议:func Paginate[T any](items []T, page, pageSize int) (Page[T], error)
// 要求:page 从 1 开始;pageSize 为 1..100;Total 始终为输入总数;结果 Items 复制一份
// 常规:[1,2,3,4,5], page=2, pageSize=2 -> Items=[3,4], Total=5
// 边界:越界页和空输入 -> Items 长度 0;page<=0 或非法 pageSize -> error
// page 极大也应得到空页, 避免直接计算 (page-1)*pageSize 导致整数溢出

func main() {
}
