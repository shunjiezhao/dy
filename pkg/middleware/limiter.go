package middleware

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"time"

	"github.com/axiaoxin-com/ratelimiter"
	"github.com/gin-gonic/gin"
)

func Limiter() gin.HandlerFunc {
	return ratelimiter.GinMemRatelimiter(ratelimiter.GinRatelimiterConfig{
		LimitKey: func(c *gin.Context) string {
			return c.ClientIP() // 针对客户端的ip
		},
		LimitedHandler: func(c *gin.Context) {
			c.JSON(200, handlers.BuildResponse(errno.RemoteErr)) //稍后重试
			c.Abort()
			return
		},
		TokenBucketConfig: func(*gin.Context) (time.Duration, int) {
			return time.Second * 2, 4
		},
	})
}
