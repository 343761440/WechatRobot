package utils

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetLogLevel aaa
func GetLogLevel(logLevelConfig string) log.Level {
	level := log.InfoLevel

	logLevelConfig = strings.ToUpper(logLevelConfig)

	switch logLevelConfig {
	case "DEBUG":
		level = log.DebugLevel
	case "INFO":
		level = log.InfoLevel
	case "ERROR":
		level = log.ErrorLevel
	case "FATAL":
		level = log.FatalLevel
	case "TRACE":
		level = log.TraceLevel
	case "WARN":
		level = log.WarnLevel
	}
	return level
}
