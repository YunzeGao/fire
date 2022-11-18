package http

import (
	"github.com/YunzeGao/fire/app/http/module/demo"
	"github.com/YunzeGao/fire/framework/gin"
	"github.com/YunzeGao/fire/framework/middleware"
)

// Routes 绑定业务层路由
func Routes(engine *gin.Engine) {
	engine.Static("/dist/", "./dist/")
	engine.Use(middleware.Trace())
	_ = demo.Register(engine)
}
