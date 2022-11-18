package id

import (
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireIDProvider struct {
}

// Register registe a new function for make a service instance
func (provider *FireIDProvider) Register(c framework.IContainer) framework.NewInstance {
	return NewFireIDService
}

// Boot will called when the service instantiate
func (provider *FireIDProvider) Boot(c framework.IContainer) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FireIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FireIDProvider) Params(c framework.IContainer) []interface{} {
	return []interface{}{}
}

// Name define the name for this service
func (provider *FireIDProvider) Name() string {
	return contract.IDKey
}
