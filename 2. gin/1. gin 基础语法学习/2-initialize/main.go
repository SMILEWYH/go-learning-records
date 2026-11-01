package main

import "github.com/gin-gonic/gin"

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Index(c *gin.Context) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  "success",
		Data: "2026-11-01",
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
