package env

import (
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireEnvProvider struct {
	Folder string
}

func (fire *FireEnvProvider) Name() string {
	return contract.EnvKey
}

func (fire *FireEnvProvider) Params(container framework.IContainer) []interface{} {
	return []interface{}{fire.Folder}
}

func (fire *FireEnvProvider) Register(container framework.IContainer) framework.NewInstance {
	return NewFireEnv
}

func (fire *FireEnvProvider) IsDefer() bool {
	return false
}

func (fire *FireEnvProvider) Boot(container framework.IContainer) error {
	app := container.MustMake(contract.AppKey).(contract.App)
	fire.Folder = app.BaseFolder()
	return nil
}
