//go:build ignore

package main

// 练习 1：提取 API 模型的字段清单。
// 建议：func dbColumns(model any) ([]string, error)。只处理结构体或一层非 nil 指针。
// 要求：按声明顺序提取已导出字段的 db 标签；缺失、空标签、"-" 均跳过。
// 自行定义 Profile：ID int64、Name string、Age int、Active bool，db 标签依次为 id/name/age/active；
// 另加 Note string（db:"-"）和未导出的 secret string，用于核对跳过规则。
// 常规：Profile{} 与 &Profile{} -> ["id","name","age","active"]。
// 边界：nil、(*Profile)(nil)、123 -> error；没有标签的结构体 -> 长度 0。
// 在注释中说明为什么 Type.Field 与 Value.Field 不能互换。

// 练习 2：按字段名修改配置。
// 建议：func updateField(target any, fieldName string, next any) error。
// 要求：接收非 nil 结构体指针，只匹配该结构体直接声明的字段；拒绝未知/未导出字段和不兼容类型。
// 常规：&Profile{Name:"old"}, "Name", "new" -> Name 为 new。
// 边界：传结构体值、nil、"missing"、"secret"、给 Age 传 int64(20) -> error 且原值不变。
// 提示：先 Kind/IsNil，再 Elem/CanSet/AssignableTo；不要仅比较 Kind。

// 练习 3：为小 ORM 添加按 id 更新的语句生成器。
// 建议：func buildUpdate(table string, model any, id int64) (string, []any, error)。
// 要求：参考讲解中的字段规则，但跳过 db:"id"；id 必须大于 0；不能省略 WHERE。
// 只生成 SQL 和 args，不执行数据库操作；检查表名、列名、重复列、无可更新字段及不支持类型。
// 常规：模型只有 Name string `db:"name"`，值为 O'Reilly，id=7：
// UPDATE `users` SET `name` = ? WHERE `id` = ?；参数依次为 "O'Reilly"、int64(7)。
// 边界：id=0、空结构体、非法表名、nil 模型 -> error；用户输入不得出现在 SQL 字符串中。

func main() {
}
