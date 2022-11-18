package framework

import (
	"errors"
	"fmt"
	"sync"
)

// IContainer 是一个服务容器，提供绑定服务和获取服务的功能
type IContainer interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回 error
	Bind(provider IServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务，
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会 panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
	PrintProviders() []string
}

type FireContainer struct {
	IContainer
	// providers 存储注册的服务提供者，key 为字符串凭证
	providers map[string]IServiceProvider
	// instance 存储具体的实例，key 为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

// NewFireContainer 创建一个服务容器
func NewFireContainer() *FireContainer {
	return &FireContainer{
		providers: map[string]IServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (fire *FireContainer) PrintProviders() []string {
	var ret []string
	for _, provider := range fire.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind 将服务容器和关键字做了绑定
func (fire *FireContainer) Bind(provider IServiceProvider) error {
	fire.lock.Lock()
	key := provider.Name()
	fire.providers[key] = provider
	fire.lock.Unlock()

	if provider.IsDefer() == false {
		if err := provider.Boot(fire); err != nil {
			return err
		}
		// 实例化方法
		params := provider.Params(fire)
		if instance, err := provider.Register(fire)(params...); err != nil {
			fmt.Println("bind service provider ", key, " error: ", err)
			return errors.New(err.Error())
		} else {
			fire.instances[key] = instance
		}
	}
	return nil
}

func (fire *FireContainer) IsBind(key string) bool {
	return fire.findServiceProvider(key) != nil
}

func (fire *FireContainer) Make(key string) (interface{}, error) {
	return fire.make(key, nil, false)
}

func (fire *FireContainer) MustMake(key string) interface{} {
	serv, err := fire.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

func (fire *FireContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return fire.make(key, params, true)
}

func (fire *FireContainer) findServiceProvider(key string) IServiceProvider {
	fire.lock.RLock()
	defer fire.lock.RUnlock()
	if sp, ok := fire.providers[key]; ok {
		return sp
	}
	return nil
}

func (fire *FireContainer) newInstance(sp IServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(fire); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(fire)
	}
	method := sp.Register(fire)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

func (fire *FireContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	fire.lock.RLock()
	defer fire.lock.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := fire.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}
	if forceNew {
		return fire.newInstance(sp, params)
	}
	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	if ins, ok := fire.instances[key]; ok {
		return ins, nil
	}
	// 容器中还未实例化，则进行一次实例化
	ins, err := fire.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	fire.instances[key] = ins
	return ins, nil
}

func (fire *FireContainer) NameList() []string {
	ret := make([]string, len(fire.providers))
	index := 0
	for _, provider := range fire.providers {
		ret[index] = provider.Name()
		index++
	}
	return ret
}
