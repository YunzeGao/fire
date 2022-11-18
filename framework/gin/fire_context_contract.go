package gin

import "github.com/YunzeGao/fire/framework/contract"

func (c *Context) MustMakeAPP() contract.App {
	return c.MustMake(contract.AppKey).(contract.App)
}

func (c *Context) MustMakeKernel() contract.IKernel {
	return c.MustMake(contract.KernelKey).(contract.IKernel)
}

// MustMakeConfig 从容器中获取配置服务
func (c *Context) MustMakeConfig() contract.IConfig {
	return c.MustMake(contract.ConfigKey).(contract.IConfig)
}

// MustMakeLog 从容器中获取日志服务
func (c *Context) MustMakeLog() contract.ILog {
	return c.MustMake(contract.FireLogKey).(contract.ILog)
}
