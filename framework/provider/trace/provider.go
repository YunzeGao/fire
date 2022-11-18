package trace

import (
	"github.com/YunzeGao/fire/framework"
	"github.com/YunzeGao/fire/framework/contract"
)

type FireTraceProvider struct {
	container framework.IContainer
}

// Register registe a new function for make a service instance
func (provider *FireTraceProvider) Register(container framework.IContainer) framework.NewInstance {
	return NewFireTraceService
}

// Boot will called when the service instantiate
func (provider *FireTraceProvider) Boot(container framework.IContainer) error {
	provider.container = container
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *FireTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *FireTraceProvider) Params(container framework.IContainer) []interface{} {
	return []interface{}{provider.container}
}

// Name define the name for this service
func (provider *FireTraceProvider) Name() string {
	return contract.TraceKey
}
