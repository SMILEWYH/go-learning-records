package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter 为每个服务安装相同的日志、异常恢复和错误响应。
func NewRouter() *gin.Engine {
	router := gin.New()
	// 直接在本机监听，不信任客户端传入的代理地址头。
	if err := router.SetTrustedProxies(nil); err != nil {
		panic(err)
	}
	router.HandleMethodNotAllowed = true
	router.Use(gin.Logger(), gin.CustomRecovery(func(c *gin.Context, recovered any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "服务内部错误"})
	}))
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "接口不存在"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "请求方法不支持"})
	})
	return router
}
