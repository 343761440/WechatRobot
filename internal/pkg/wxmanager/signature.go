package wxmanager

import (
	"wxrobot/internal/pkg/message"
	"wxrobot/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Signature(c *gin.Context) {
	echostr := c.Query("echostr")
	c.Writer.Write([]byte(echostr))
}

func initSignature(r *gin.Engine) {
	g := r.Group("/signature")
	{
		sign := g.Group("/")
		sign.Use(middleware.SignatureMiddleware)
		{
			sign.GET("home", Signature)
			sign.POST("home", message.HandleMessage)
		}
	}
}
