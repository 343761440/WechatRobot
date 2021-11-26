package wxmanager

import (
	"wxrobot/internal/app/common"
	"wxrobot/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
	"gorm.io/gorm"
)

type WxManager struct {
	rds    *redis.Client
	db     *gorm.DB
	router *gin.Engine
}

func NewWxMangerRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	initSignature(r)
	return r
}

func (wm *WxManager) Run() {
	wm.router.Run(":80")
}

func InitManger() {
	common.Initconfig()
	common.InitLogger()

	r := NewWxMangerRouter()
	rdb, err := common.GetRedisClient()
	if err != nil {
		log.Fatal("init redis failed,err=", err)
	}

	db, err := common.InitDB(false, false)
	if err != nil {
		log.Fatal("init db failed,err=", err)
	}

	wxmgr := WxManager{
		router: r,
		rds:    rdb,
		db:     db,
	}

	wxmgr.Run()
}
