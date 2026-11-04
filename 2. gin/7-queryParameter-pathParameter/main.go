package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Query("name")
		age := c.DefaultQuery("age", "18") // 如果没有传递age参数, 默认值是18
		keyList := c.QueryArray("manyKey")

		fmt.Println(name, age, keyList)
	})

	r.GET("/user/:id/:name", func(c *gin.Context) {
		id := c.Param("id")
		name := c.Param("name")

		fmt.Println(id, name)
	})

	r.Run(":8080")
}
