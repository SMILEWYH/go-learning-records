package main

import (
	"github.com/gin-gonic/gin"
)

// 使用 bind 来代替之前的方法

func main() {

	r := gin.Default()

	r.POST("product", func(c *gin.Context) {
		type Product struct {
			Name string `json:"name" binding:"required, min=5"` // 这样去使用, 如果有多个规则, 使用逗号分隔
			Size int    `json:"size"`
		}
		var product Product

		err := c.ShouldBindJSON(&product)
		if err != nil {
			// do someting
		}

		// next step
	})

	r.Run(":8080")
	// 不能为空，并且不能没有这个字段
	// required： 必填字段，如：binding:"required"

	// 针对字符串的长度
	// min 最小长度，如：binding:"min=5"
	// max 最大长度，如：binding:"max=10"
	// len 长度，如：binding:"len=6"

	// 针对数字的大小
	// eq 等于，如：binding:"eq=3"
	// ne 不等于，如：binding:"ne=12"
	// gt 大于，如：binding:"gt=10"
	// gte 大于等于，如：binding:"gte=10"
	// lt 小于，如：binding:"lt=10"
	// lte 小于等于，如：binding:"lte=10"

	// 针对同级字段的
	// eqfield 等于其他字段的值，如：PassWord string `binding:"eqfield=其他字段名称"`
	// nefield 不等于其他字段的值

	// - 忽略字段，如：binding:"-" 或者不写

	// 枚举  只能是red 或green
	// oneof=red green

	// 字符串
	// contains=fengfeng  // 包含fengfeng的字符串
	// excludes // 不包含
	// startswith  // 字符串前缀
	// endswith  // 字符串后缀

	// 数组
	// dive  // dive后面的验证就是针对数组中的每一个元素

	// 网络验证 (其实就是检查数据结构是否符合要求)
	// ip
	// ipv4
	// ipv6
	// uri
	// url
	// uri 在于I(Identifier)是统一资源标示符，可以唯一标识一个资源。
	// url 在于Locater，是统一资源定位符，提供找到该资源的确切路径

	// 日期验证  要求字段值是一个符合 "年-月-日" 格式的有效日期字符串
	// datetime=2006-01-02
}
