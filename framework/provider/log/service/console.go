package service

import (
	"os"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireConsoleLog struct {
	FireLog
}

func NewFireConsoleLog(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	log := &FireConsoleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.container = container
	return log, nil
}
