package app

import (
	"example.com/go-room-booking/internal/web"
	"github.com/gin-gonic/gin"
)

func Handler(s *Service, version string) *gin.Engine {
	router := web.NewRouter()
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, map[string]string{"status": "ok", "version": version})
	})

	rooms := router.Group("/rooms")
	rooms.GET("", func(c *gin.Context) { c.JSON(200, Rooms()) })
	rooms.GET("/:id/availability", func(c *gin.Context) {
		ok, err := s.Available(c.Param("id"), c.Query("start"), c.Query("end"))
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, map[string]bool{"available": ok})
	})

	bookings := router.Group("/bookings")
	bookings.POST("", func(c *gin.Context) {
		var in BookingInput
		if err := web.Decode(c.Writer, c.Request, &in); err != nil {
			web.ReplyError(c, err)
			return
		}
		b, err := s.Create(in)
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(201, b)
	})
	bookings.GET("", func(c *gin.Context) {
		items, err := s.List(c.Query("room_id"), c.Query("status"))
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
	bookings.GET("/:id", func(c *gin.Context) {
		b, err := s.Get(c.Param("id"))
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, b)
	})
	bookings.POST("/:id/cancel", func(c *gin.Context) {
		b, err := s.Cancel(c.Param("id"))
		if err != nil {
			web.ReplyError(c, err)
			return
		}
		c.JSON(200, b)
	})
	return router
}
