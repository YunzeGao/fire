package app

import (
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireAppProvider struct {
	BaseFolder string
}

func (fire *FireAppProvider) Name() string {
	return contract.AppKey
}

func (fire *FireAppProvider) Params(container framework.IContainer) []interface{} {
	return []interface{}{container, fire.BaseFolder}
}

func (fire *FireAppProvider) Register(container framework.IContainer) framework.NewInstance {
	return NewFireApp
}

func (fire *FireAppProvider) IsDefer() bool {
	return false
}

func (fire *FireAppProvider) Boot(container framework.IContainer) error {
	return nil
}
