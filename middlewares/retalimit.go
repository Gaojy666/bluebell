package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimitMiddleware令牌桶中间件
// fillInterval向桶里填充令牌的间隔时间，cap是总容量
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// bucket.TakeAvailable(1)返回取1个令牌等待后令牌桶移除的令牌数
		// bucket.Take(1)返回取1个令牌需要等待的时间
		// 如果取不到令牌就返回响应
		//if bucket.Take(1) > 0 {
		if bucket.TakeAvailable(1) != 1 {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		// 取到令牌就放行
		c.Next()
	}
}
