package app

import (
	"example.com/go-shop-orders/internal/web"
	"github.com/gin-gonic/gin"
)

func Handler(s *Service, version string) *gin.Engine {
	router := web.NewRouter()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, map[string]string{"status": "ok", "version": version})
	})

	products := router.Group("/products")
	products.GET("", func(c *gin.Context) {
		result, err := web.Paginate(c.Request, s.Products())
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, result)
	})
	products.POST("", func(c *gin.Context) {
		var in Product
		if err := web.Decode(c.Writer, c.Request, &in); err != nil {
			web.ReplyError(c, err)
			return
		}
		p, err := s.AddProduct(in)
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(201, p)
	})
	products.POST("/:id/restock", func(c *gin.Context) {
		var in struct {
			Quantity int `json:"quantity"`
		}
		if err := web.Decode(c.Writer, c.Request, &in); err != nil {
			web.ReplyError(c, err)
			return
		}
		p, err := s.Restock(c.Param("id"), in.Quantity)
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, p)
	})

	orders := router.Group("/orders")
	orders.GET("", func(c *gin.Context) {
		items, err := s.Orders(c.Query("status"))
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		page, err := web.Paginate(c.Request, items)
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, page)
	})
	orders.GET("/:id", func(c *gin.Context) {
		o, err := s.Get(c.Param("id"))
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, o)
	})
	orders.POST("", func(c *gin.Context) {
		var in CreateOrder
		if err := web.Decode(c.Writer, c.Request, &in); err != nil {
			web.ReplyError(c, err)
			return
		}
		o, err := s.Create(in)
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(201, o)
	})
	orders.PATCH("/:id/status", func(c *gin.Context) {
		var in struct {
			Status OrderStatus `json:"status"`
		}
		if err := web.Decode(c.Writer, c.Request, &in); err != nil {
			web.ReplyError(c, err)
			return
		}
		o, err := s.Transition(c.Param("id"), in.Status)
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, o)
	})
	return router
}
