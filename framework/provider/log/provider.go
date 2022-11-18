package log

import (
	"io"
	"strings"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/provider/log/formatter"
	"github.com/YunzeGao/fire/framework/provider/log/service"
)

type FireLogProvider struct {
	Driver     string
	Level      contract.LogLevel   // 日志级别
	Formatter  contract.Formatter  // 日志输出格式方法
	CtxFielder contract.CtxFielder // 日志context上下文信息获取函数
	Output     io.Writer           // 日志输出信息
}

func (provider *FireLogProvider) Name() string {
	return contract.FireLogKey
}

func (provider *FireLogProvider) Params(container framework.IContainer) []interface{} {
	// 获取configService
	configService := container.MustMake(contract.ConfigKey).(contract.IConfig)
	// 设置参数formatter
	if provider.Formatter == nil {
		provider.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				provider.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				provider.Formatter = formatter.TextFormatter
			}
		}
	}

	if provider.Level == contract.UnknownLevel {
		provider.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			provider.Level = logLevel(configService.GetString("log.level"))
		}
	}

	// 定义5个参数
	return []interface{}{container, provider.Level, provider.CtxFielder, provider.Formatter, provider.Output}
}

func (provider *FireLogProvider) Register(container framework.IContainer) framework.NewInstance {
	if provider.Driver == "" {
		tcs, err := container.Make(contract.ConfigKey)
		if err != nil {
			// 默认使用console
			return service.NewFireConsoleLog
		}

		cs := tcs.(contract.IConfig)
		provider.Driver = strings.ToLower(cs.GetString("log.Driver"))
	}
	// 根据driver的配置项确定
	switch provider.Driver {
	case "single":
		return service.NewFireSingleLog
	case "rotate":
		return service.NewFireRotateLog
	case "console":
		return service.NewFireConsoleLog
	case "custom":
		return service.NewFireCustomLog
	default:
		return service.NewFireConsoleLog
	}
}

func (provider *FireLogProvider) IsDefer() bool {
	return false
}

func (provider *FireLogProvider) Boot(container framework.IContainer) error {
	return nil
}

// logLevel get level from string
func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "warn":
		return contract.WarnLevel
	case "info":
		return contract.InfoLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	}
	return contract.UnknownLevel
}
