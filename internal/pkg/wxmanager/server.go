package wxmanager

import (
	"wxrobot/internal/app/common"
	"wxrobot/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type WxManager struct {
}

func NewWxMangerRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	initSignature(r)
	return r
}

func InitManger() {
	common.Initconfig()
	common.InitLogger()

	r := NewWxMangerRouter()
	r.Run(":80")
}
