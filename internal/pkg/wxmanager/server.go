package wxmanager

import (
	"github.com/gin-gonic/gin"
)

type WxManager struct {
}

func InitManger() {
	r := gin.Default()

	r.Use(SignatureMiddleware)
	r.GET("/home", Signature)
	r.POST("/home", RecvMessage)
	r.Run(":80")
}
