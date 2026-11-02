package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 之前使用 binding 自带规则, 如: min=5 这样会导致这个字段不仅被长度被限制为最小5, 同时它也为一个必填字段
// 如果我们需要 min=5 的能力, 但又不想受到必填字段的影响, 那么就需要使用到 binding 自定义规则

// 补充：这个需求也可以用 binding:"omitempty,min=5" 实现，这里用于演示自定义规则。

// optionalMin5 允许空字符串，非空字符串至少需要 5 个字符。
func optionalMin5(fl validator.FieldLevel) bool {
	fmt.Println("fl.Field(): ", fl.Field())
	fmt.Println("fl.FieldName(): ", fl.FieldName())
	fmt.Println("fl.StructFieldName(): ", fl.StructFieldName())
	fmt.Println("fl.Parent(): ", fl.Parent())
	fmt.Println("fl.Top(): ", fl.Top())
	fmt.Println("fl.Param(): ", fl.Param())

	// 获取当前正在校验的字符串字段。
	value := fl.Field().String()

	// 缺少 name 时字段为零值空字符串，直接通过校验。
	if value == "" {
		return true
	}

	// 按 Unicode 字符数量校验，避免把一个中文汉字按多个字节计算。
	return utf8.RuneCountInString(value) >= 5
}

func main() {
	// 获取 Gin 使用的校验器，在开始处理请求之前注册规则。
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("无法获取 Gin 校验器")
	}

	// 第一个参数是规则名，第二个参数是校验函数；注册成功不代表请求校验成功。
	if err := v.RegisterValidation("optional_min5", optionalMin5); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.POST("product", func(c *gin.Context) {
		type Product struct {
			Name string `json:"name" binding:"optional_min5"` // 这样去使用, 如果有多个规则, 使用逗号分隔
			Size int    `json:"size"`
		}
		var product Product

		err := c.ShouldBindJSON(&product)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// next step
		c.JSON(200, product)
	})

	r.Run(":8080")
}
