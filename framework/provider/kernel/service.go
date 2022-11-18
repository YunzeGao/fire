package kernel

import (
	"net/http"

	"github.com/YunzeGao/fire/framework/gin"
)

type FireKernelService struct {
	engine *gin.Engine
}

// NewFireKernelService 初始化web引擎服务实例
func NewFireKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &FireKernelService{
		engine: httpEngine,
	}, nil
}

// HttpEngine 返回web引擎
func (s *FireKernelService) HttpEngine() http.Handler {
	return s.engine
}
