package model

import (
	"sync"
	"wxrobot/internal/app/common"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ModelManager struct {
	db *gorm.DB
}

var (
	gModelMgr *ModelManager
	once      sync.Once
)

func GetInstance() *ModelManager {
	once.Do(func() {
		gModelMgr = &ModelManager{}
	})
	return gModelMgr
}

func init() {
	db, err := common.InitDB(false, false)
	if err != nil {
		log.Fatal("init db failed,err=", err)
	}
	//更新需要的数据表
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;").
		AutoMigrate(&WxUser{}, &Question{}, &TodoItem{}, &LoveJoke{}, &JcEvent{}); err != nil {
		log.Fatal("database init model error ", err)
		return
	}
	mm := GetInstance()
	mm.db = db
}
