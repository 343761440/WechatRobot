package log

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type LastError struct {
	errorInfo []string
	sync.Mutex
}

var gError LastError

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logrus.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Error logs a message at level Error on the standard logger.
func ErrorWithRecord(args ...interface{}) {
	logrus.Error(args...)
	gError.Lock()
	gError.errorInfo = append(gError.errorInfo, time.Now().Format("2006-01-02 15:04:05")+"\n"+fmt.Sprint(args...))
	gError.Unlock()
}

func GetLastError() []string {
	gError.Lock()
	defer gError.Unlock()

	res := gError.errorInfo
	gError.errorInfo = gError.errorInfo[0:0]
	return res
}
