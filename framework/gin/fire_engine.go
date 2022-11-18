package gin

import "github.com/YunzeGao/fire/framework"

func (engine *Engine) SetContainer(container framework.IContainer) {
	engine.container = container
}

// Bind engine实现container的绑定封装
func (engine *Engine) Bind(provider framework.IServiceProvider) error {
	return engine.container.Bind(provider)
}

// IsBind 关键字凭证是否已经绑定服务提供者
func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}
