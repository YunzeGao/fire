package http

import "github.com/YunzeGao/fire/framework/gin"

func NewHttpEngine() (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	// 默认启动一个Web引擎
	engine := gin.New()
	engine.Use(gin.Recovery())
	Routes(engine)
	return engine, nil
}
