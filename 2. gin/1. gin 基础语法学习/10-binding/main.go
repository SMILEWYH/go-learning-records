package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 使用 bind 来代替之前的方法

func main() {
	r := gin.Default()

	// 查询参数
	r.GET("user", func(c *gin.Context) {
		type User struct {
			Name string `form:"name"`
			Age  int    `form:"age"`
		}
		var user User

		err := c.ShouldBindQuery(&user)
		fmt.Println(user, err)
	})

	// 动态参数 (路径参数)
	r.GET("user/:id/:name", func(c *gin.Context) {
		type User struct {
			Name string `uri:"name"`
			Age  int    `uri:"age"`
		}
		var user User

		err := c.ShouldBindUri(&user)
		fmt.Println(user, err)
	})

	// 表单参数, 不能解析 x-www-form-urlencoded的格式
	r.POST("user", func(c *gin.Context) {
		type User struct {
			Name string `form:"name"`
			Age  int    `form:"age"`
		}
		var user User

		err := c.ShouldBind(&user)
		fmt.Println(user, err)
	})

	// JSON 参数
	r.POST("product", func(c *gin.Context) {
		type Product struct {
			Name string `json:"name"`
			Size int    `json:"size"`
		}
		var product Product

		err := c.ShouldBindJSON(&product)
		fmt.Println(product, err)
	})

	// header 参数
	r.POST("header", func(c *gin.Context) {
		type Header struct {
			UserAgent   string `header:"User-Agent"`
			ContentType string `header:"Content-Type"`
		}
		var header Header

		err := c.ShouldBindJSON(&header)
		fmt.Println(header, err)
	})

	r.Run(":8080")
}
