package main

import (
	"github.com/gin-gonic/gin"
)

// 使用 bind 来代替之前的方法

func main() {
	r := gin.Default()

	// 常见的六种请求方法:
	// r.GET()
	// r.POST()
	// r.PUT()
	// r.PATCH()
	// r.DELETE()
	// r.OPTIONS()
	// r.Any() // 上面方法的集合

	// 路由分组:
	var gameGroupConstructor = func(r *gin.RouterGroup) {
		r.POST("detail", func(c *gin.Context) {})    // game/detail
		r.POST("storePage", func(c *gin.Context) {}) // game/storePage
	}

	var gameGroupConstructor2 = func(r *gin.RouterGroup) {
		r.POST("price", func(c *gin.Context) {}) // game/price
	}

	gameGroup := r.Group("game")
	gameGroupConstructor(gameGroup)

	// 除了可以让同一类接口更加清晰, 我们还可以在分组上使用中间件, 这样就不需要在每个路由上都使用中间件了
	gameGroup.Use(func(c *gin.Context) {})

	// 可以创建多个同名的路由分组, 这样就可以区分哪些路由分组需要使用中间件, 哪些不需要使用中间件
	gameGroup2 := r.Group("game")
	gameGroupConstructor2(gameGroup2) // gameGroup2 不使用中间件

	// RESETFul Api规范: 尽量使用名称的复数来定义路由
	// 在没有resetful规范正确, 表示创建用户、删除用户
	// /api/user_create
	// /api/users/create
	// /api/users/add
	// /api/add_user
	// /api/user/delete
	// /api/user_remove

	// 使用resetful规范
	// GET /api/users  用户列表
	// POST /api/users  创建用户
	// PUT /api/users/:id 更新用户信息
	// DELETE /api/users 批量删除用户
	// DELETE /api/users/:id 删除单个用户

	// 但是resetful它只是一个api规范
	// 大家在公司里面的项目，基本上都是POST请求, 因为:
	// 1. 很早之前的，那个时候还没有RESETFul 规范这个说法
	// 2. 很多公司的防火墙会拦截除GET和POST之外的请求

	r.Run(":8080")
}
