package app

import (
	"errors"
	"path/filepath"

	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/util"
)

type FireApp struct {
	container  framework.IContainer // 服务容器
	baseFolder string               // 基础路径
	configMap  map[string]string    // 加载配置文件
}

// NewFireApp 初始化FireApp
func NewFireApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.IContainer)
	baseFolder := params[1].(string)
	return &FireApp{baseFolder: baseFolder, container: container}, nil
}

func (app *FireApp) Version() string {
	return "0.1"
}

func (app *FireApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}
	return util.GetExecDirectory()
}

func (app *FireApp) ConfigFolder() string {
	return filepath.Join(app.BaseFolder(), "config")
}

func (app *FireApp) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

func (app *FireApp) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

func (app *FireApp) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "http")
}

func (app *FireApp) ProviderFolder() string {
	if val, ok := app.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "provider")
}

func (app *FireApp) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

// ConsoleFolder 定义业务定义
func (app *FireApp) ConsoleFolder() string {
	if val, ok := app.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "console")
}

func (app *FireApp) CommandFolder() string {
	if val, ok := app.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

func (app *FireApp) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

func (app *FireApp) TestFolder() string {
	if val, ok := app.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "test")
}

func (app *FireApp) LoadAppConfig(kv map[string]string) {
	for k, v := range kv {
		app.configMap[k] = v
	}
}
