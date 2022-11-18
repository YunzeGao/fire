package framework

type (
	NewInstance      func(...interface{}) (interface{}, error)
	IServiceProvider interface {
		// Name 这个服务提供者的凭证
		Name() string
		// Params 传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数
		Params(IContainer) []interface{}
		// Register 在服务容器中注册了一个实例化服务的方法，是否在注册的时候就实例化这个服务，需要参考 IsDefer 接口。
		Register(IContainer) NewInstance
		// IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那就是在第一次 make 的时候进行实例化操作
		//false 表示不需要延迟实例化，在注册的时候就实例化。true 表示延迟实例化
		IsDefer() bool
		// Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这个里面。
		//	整个服务实例化会实例化失败，返回 error
		Boot(IContainer) error
	}
)
