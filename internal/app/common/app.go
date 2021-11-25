package common

import (
	"wxrobot/internal/app/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initconfigImpl(panic bool) {
	viper.SetConfigName(kCommonConfigName)
	viper.AddConfigPath(kConfigPath)
	// viper.MergeInConfig()
	err := viper.ReadInConfig()
	if err != nil {
		if panic {
			log.Panicf("viper read config err: %s", err.Error())
		} else {
			log.Errorf("viper read config err: %s", err.Error())
		}
	}
	viper.AutomaticEnv()
}

// Initconfig 初始化配置viper
func Initconfig() {
	initconfigImpl(true)
}

// InitLogger 初始化日志模块
func InitLogger() {
	level := utils.GetLogLevel(viper.GetString("LOG_LEVEL"))
	log.SetLevel(level)
	formatter := &log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	log.SetFormatter(formatter)
	log.Debug("debug log level")
	log.Info("start")
}
