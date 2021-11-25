package main

import (
	"wxrobot/internal/app/common"
	"wxrobot/internal/pkg/wxmanager"
)

func init() {
	common.Initconfig()
	common.InitLogger()
}

func main() {
	wxmanager.InitManger()
}
