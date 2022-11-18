package kernel

import (
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
	"github.com/YunzeGao/fire/framework/gin"
)

type FireKernelProvider struct {
	HttpEngine *gin.Engine
}

func (fire *FireKernelProvider) Name() string {
	return contract.KernelKey
}

func (fire *FireKernelProvider) Params(container framework.IContainer) []interface{} {
	return []interface{}{fire.HttpEngine}
}

func (fire *FireKernelProvider) Register(container framework.IContainer) framework.NewInstance {
	return NewFireKernelService
}

func (fire *FireKernelProvider) IsDefer() bool {
	return false
}

// Boot 启动的时候判断是否由外界注入了Engine，如果注入的化，用注入的，如果没有，重新实例化
func (fire *FireKernelProvider) Boot(container framework.IContainer) error {
	if fire.HttpEngine == nil {
		fire.HttpEngine = gin.Default()
	}
	fire.HttpEngine.SetContainer(container)
	return nil
}
