package wxmanager

import "github.com/gin-gonic/gin"

func initHealth(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.Status(200)
	})
}
