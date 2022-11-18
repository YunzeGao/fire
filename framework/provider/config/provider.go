package config

import (
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireConfigProvider struct {
	container framework.IContainer
	folder    string
	env       string
	envMaps   map[string]string
}

func (provider *FireConfigProvider) Name() string {
	return contract.ConfigKey
}

func (provider *FireConfigProvider) Params(container framework.IContainer) []interface{} {
	return []interface{}{provider.folder, provider.envMaps, provider.env, provider.container}
}

func (provider *FireConfigProvider) Register(container framework.IContainer) framework.NewInstance {
	return NewFireConfig
}

func (provider *FireConfigProvider) IsDefer() bool {
	return true
}

func (provider *FireConfigProvider) Boot(container framework.IContainer) error {
	provider.folder = container.MustMake(contract.AppKey).(contract.App).ConfigFolder()
	provider.env = container.MustMake(contract.EnvKey).(contract.Env).AppEnv()
	provider.envMaps = container.MustMake(contract.EnvKey).(contract.Env).All()
	provider.container = container
	return nil
}
