package main

import (
	"github.com/gin-gonic/gin"
)

// 全局中间件, 一般搭配分组一起使用

func main() {
	r := gin.Default()
	// 全局路由可以看 13-binding-router

	// 这章主要讲中间件传递参数:
	// 在中间件中 c.Set("key", "value")
	// 那么在接下来的流程中可以通过 key, ok :=c.Get("key") 获取到该值

	r.Run(":8080")
}
