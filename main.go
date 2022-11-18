package main

import (
	"github.com/YunzeGao/fire/app/console"
	"github.com/YunzeGao/fire/app/http"
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/provider/app"
	"github.com/YunzeGao/fire/framework/provider/config"
	"github.com/YunzeGao/fire/framework/provider/env"
	"github.com/YunzeGao/fire/framework/provider/id"
	"github.com/YunzeGao/fire/framework/provider/kernel"
	"github.com/YunzeGao/fire/framework/provider/log"
	"github.com/YunzeGao/fire/framework/provider/trace"
)

func main() {
	container := framework.NewFireContainer()
	_ = container.Bind(&app.FireAppProvider{})
	// 后续初始化需要绑定的服务提供者...
	_ = container.Bind(&env.FireEnvProvider{})
	_ = container.Bind(&config.FireConfigProvider{})
	_ = container.Bind(&id.FireIDProvider{})
	_ = container.Bind(&trace.FireTraceProvider{})
	_ = container.Bind(&log.FireLogProvider{})
	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(); err == nil {
		_ = container.Bind(&kernel.FireKernelProvider{HttpEngine: engine})
	}
	// 运行root命令
	_ = console.RunCommand(container)
}
