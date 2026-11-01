package main

import (
	"example.com/hello/3-gin-json-response/res"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	res.OkWithData(c, map[string]any{
		"name": "枫枫",
	})
}

func main() {
	// 1. 初始化
	gin.SetMode("release") // 去除调试日志
	r := gin.Default()
	// 2. 挂载路由
	r.GET("/index", Index)
	// 3. 绑定端口,  运行
	r.Run(":8080")
}
