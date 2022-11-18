package service

import (
	"os"
	"path/filepath"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/util"

	"github.com/pkg/errors"
)

type FireSingleLog struct {
	FireLog
	folder string
	file   string
	fd     *os.File
}

// NewFireSingleLog params sequence: level, ctxFielder, Formatter, map[string]interface(folder/file)
func NewFireSingleLog(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := container.MustMake(contract.AppKey).(contract.App)
	configService := container.MustMake(contract.ConfigKey).(contract.IConfig)

	log := &FireSingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}
	log.folder = folder
	if !util.Exists(folder) {
		_ = os.MkdirAll(folder, os.ModePerm)
	}
	log.file = "file.log"
	if configService.IsExist("log.file") {
		log.file = configService.GetString("log.file")
	}
	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}
	log.SetOutput(fd)
	log.container = container
	return log, nil
}
