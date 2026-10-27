//go:build ignore

package main

/*
练习 1：购物车条目与金额计算。

背景：结算页面需要按“分”计算金额，避免用浮点数表示人民币小数。

要求：
  - 定义 CartItem，含 SKU string、UnitPriceCents int、Quantity int。
  - 用值接收者 TotalCents() int 返回单价乘数量，不修改条目。
  - 用指针接收者 SetQuantity(quantity int) bool 更新数量。
  - quantity<0 或接收者为 nil 时返回 false，其他情况更新并返回 true。
  - 为聚焦结构体，假设单价非负且乘积在 int 范围内。

建议签名：
func (item CartItem) TotalCents() int
func (item *CartItem) SetQuantity(quantity int) bool

手动核对：
  - 单价 1990，数量 2：金额 3980；修改数量为 3 后金额 5970。
  - 数量改成 -1：返回 false，原数量保持 3。
  - 数量改成 0：返回 true，金额为 0；nil 指针调用更新方法返回 false。
*/

/*
练习 2：组织信息与成员信息的组合。

背景：成员卡片同时展示成员昵称和所在组织名称，两者都可能叫 Name。

要求：
  - 定义 TeamInfo{Name string} 与 TeamMember{TeamInfo; Name string}。
  - 给 TeamInfo 定义 TeamLabel() string，返回“团队：”加团队名。
  - 给 TeamMember 定义 CardTitle() string，返回“成员名 / 团队名”。
  - 在 CardTitle 中明确访问嵌入字段的 Name，观察外层字段的遮蔽。
  - 演示直接调用 member.TeamLabel() 的方法提升。

建议签名：func (member TeamMember) CardTitle() string

手动核对：
  - 成员小林、团队平台组：CardTitle 为“小林 / 平台组”。
  - member.Name 为“小林”，member.TeamInfo.Name 为“平台组”。
  - 成员名为空、团队平台组：结果为“ / 平台组”，不隐式填充默认值。
  - member.TeamLabel() 为“团队：平台组”。
*/

/*
练习 3：安全地构造登录后的响应数据。

背景：后端需要给前端返回 id、displayName、可选 avatarURL，但不能发送密码摘要。

要求：
  - 定义 LoginResponse，字段为 ID int、DisplayName、AvatarURL、PasswordHash string。
  - 配置正确 JSON tag：id、displayName、avatarURL,omitempty、-。
  - 新增 encoding/json 导入，构造数据后调用 json.Marshal，并检查 err。
  - 输出 string(data)，不要直接依赖字节切片的数字打印形式。
  - 用注释说明：JSON 忽略字段不等于可以安全地直接打印整个结构体。

建议签名：func printLoginResponse(response LoginResponse)

手动核对（按 JSON 内容核对，不依赖对象键序）：
  - ID=1、DisplayName=Lin、AvatarURL=""：仅含 id=1 与 displayName=Lin。
  - AvatarURL=/avatars/1.png：响应包含 avatarURL。
  - PasswordHash 无论是否为空，JSON 中都不能有该字段或它的值。
*/

func main() {
}
