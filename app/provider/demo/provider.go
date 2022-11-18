package demo

import "github.com/YunzeGao/fire/framework"

type DemoProvider struct {
	framework.IServiceProvider

	c framework.IContainer
}

func (sp *DemoProvider) Name() string {
	return DemoKey
}

func (sp *DemoProvider) Register(c framework.IContainer) framework.NewInstance {
	return NewService
}

func (sp *DemoProvider) IsDefer() bool {
	return false
}

func (sp *DemoProvider) Params(c framework.IContainer) []interface{} {
	return []interface{}{sp.c}
}

func (sp *DemoProvider) Boot(c framework.IContainer) error {
	sp.c = c
	return nil
}
