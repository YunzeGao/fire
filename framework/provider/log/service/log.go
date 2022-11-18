package service

import (
	"context"
	"io"
	pkgLog "log"
	"time"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/provider/log/formatter"
)

type FireLog struct {
	level      contract.LogLevel
	formatter  contract.Formatter
	ctxFielder contract.CtxFielder
	output     io.Writer
	container  framework.IContainer
}

func (log *FireLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

func (log *FireLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	if !log.IsLevelEnable(level) {
		return nil
	}
	// 使用ctxFielder 获取 context 中的信息
	fs := fields
	if log.ctxFielder != nil {
		t := log.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}
	// 如果绑定了trace服务，获取trace信息
	if log.container.IsBind(contract.TraceKey) {
		tracer := log.container.MustMake(contract.TraceKey).(contract.Trace)
		if tc := tracer.GetTrace(ctx); tc != nil {
			for k, v := range tracer.ToMap(tc) {
				fs[k] = v
			}
		}
	}
	// 将日志信息按照formatter序列化为字符串
	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}

	content, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}
	if level == contract.PanicLevel {
		pkgLog.Panic(string(content))
		return nil
	}
	_, _ = log.output.Write(content)
	_, _ = log.output.Write([]byte("\r\n"))
	return nil
}

// SetOutput 设置output
func (log *FireLog) SetOutput(output io.Writer) {
	log.output = output
}

// Panic 输出panic的日志信息
func (log *FireLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.PanicLevel, ctx, msg, fields)
}

// Fatal will add fatal record which contains msg and fields
func (log *FireLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.FatalLevel, ctx, msg, fields)
}

// Error will add error record which contains msg and fields
func (log *FireLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.ErrorLevel, ctx, msg, fields)
}

// Warn will add warn record which contains msg and fields
func (log *FireLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.WarnLevel, ctx, msg, fields)
}

// Info 会打印出普通的日志信息
func (log *FireLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.InfoLevel, ctx, msg, fields)
}

// Debug will add debug record which contains msg and fields
func (log *FireLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.DebugLevel, ctx, msg, fields)
}

// Trace will add trace info which contains msg and fields
func (log *FireLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	_ = log.logf(contract.TraceLevel, ctx, msg, fields)
}

// SetLevel set log level, and higher level will be recorded
func (log *FireLog) SetLevel(level contract.LogLevel) {
	log.level = level
}

// SetCtxFielder will get fields from context
func (log *FireLog) SetCtxFielder(handler contract.CtxFielder) {
	log.ctxFielder = handler
}

// SetFormatter will set formatter handler will covert data to string for recording
func (log *FireLog) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}
