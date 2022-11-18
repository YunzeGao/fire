package middleware

import (
	"context"
	"log"
	"time"

	"github.com/YunzeGao/fire/framework/gin"
)

func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), duration)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			ctx.Next()
			finish <- struct{}{}
		}()
		select {
		case p := <-panicChan:
			_ = ctx.ISetStatus(500).IJson(map[string]interface{}{
				"msg": "timeout",
			})
			log.Println(p)
		case <-finish:
			log.Println("finish = ", finish)
		case <-durationCtx.Done():
			_ = ctx.ISetStatus(500).IJson(map[string]interface{}{
				"msg": "timeout",
			})
		}
	}
}
