package middleware

import (
	"log"
	"time"

	"github.com/YunzeGao/fire/framework/gin"
)

func Cost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		end := time.Now()
		cost := end.Sub(start)
		log.Printf("api uri: %v, cost: %.6f 秒", ctx.DefaultUri(), cost.Seconds())
		return
	}
}
