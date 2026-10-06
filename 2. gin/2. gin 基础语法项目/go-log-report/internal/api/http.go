// Package api 使用 Gin 暴露日志报告，复用 analyzer 的并发分析逻辑。
package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"example.com/go-log-report/internal/analyzer"
	"example.com/go-log-report/internal/web"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Input   string
	Workers int
	Strict  bool
	Timeout time.Duration
}

type reportQuery struct {
	Workers *int  `form:"workers" binding:"omitempty,min=1,max=16"`
	Strict  *bool `form:"strict"`
}

func Handler(config Config, version string) *gin.Engine {
	return newHandler(config, version, analyzer.Analyze)
}

// 分析函数作为依赖传入，测试可以稳定控制请求取消和并发时序。
func newHandler(config Config, version string, analyze func(context.Context, string, analyzer.Options) (analyzer.Report, error)) *gin.Engine {
	router := web.NewRouter()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "version": version})
	})
	// 一个服务同时只运行一次完整分析，避免请求数乘以 worker 数耗尽资源。
	busy := make(chan struct{}, 1)
	router.GET("/reports", func(c *gin.Context) {
		// ParseQuery 显式报告无效转义，不能静默忽略部分参数。
		values, err := url.ParseQuery(c.Request.URL.RawQuery)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "查询参数格式错误"})
			return
		}
		for key, value := range values {
			if (key != "workers" && key != "strict") || len(value) != 1 || value[0] == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "仅支持单个非空 workers、strict 参数"})
				return
			}
		}
		var query reportQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "workers 必须为 1..16，strict 必须为布尔值"})
			return
		}
		workers, strict := config.Workers, config.Strict
		if query.Workers != nil {
			workers = *query.Workers
		}
		if query.Strict != nil {
			strict = *query.Strict
		}
		select {
		case busy <- struct{}{}:
			defer func() { <-busy }()
		default:
			c.Header("Retry-After", "1")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "正在分析，请稍后重试"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()
		// 先取得完整统计，再按 HTTP 协议表达严格模式的拒绝原因。
		report, err := analyze(ctx, config.Input, analyzer.Options{Workers: workers})
		if err == nil {
			err = ctx.Err()
		}
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"error": "日志分析超时"})
			case errors.Is(err, context.Canceled):
				c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{"error": "请求已取消"})
			default:
				log.Printf("日志分析失败: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "日志读取或分析失败，请检查服务端日志目录"})
			}
			return
		}
		if strict && report.InvalidLines > 0 {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"error":         fmt.Sprintf("严格模式拒绝 %d 行无效日志", report.InvalidLines),
				"invalid_lines": report.InvalidLines, "issues": report.Issues,
			})
			return
		}
		c.JSON(http.StatusOK, report)
	})
	return router
}
