package service

import (
	"io"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireCustomLog struct {
	FireLog
}

func NewFireCustomLog(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)
	output := params[4].(io.Writer)

	log := &FireCustomLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetOutput(output)
	log.container = container
	return log, nil
}
