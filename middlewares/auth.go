package middlewares

import (
	"bluebell/controller"
	"bluebell/dao/redis"
	"bluebell/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleWare 基于JWT的认证中间件
// 检查请求是否按要求携带了jwt token认证
func JWTAuthMiddleWare() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		// 用来校验用户携带的token是否与存储的userID-token一致
		token, err := redis.GetTokenFromID(mc.UserId)
		if err != nil {
			controller.ResponseError(c, controller.CodeServerBusy)
			c.Abort()
			return
		}
		// 如果不一致，则说明有第二个设备登陆了，当前第一个设备只能退出登录
		if token != parts[1] {
			controller.ResponseError(c, controller.CodeInvalidatedLogin)
			c.Redirect(http.StatusTemporaryRedirect, "/api/v1/login")
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set(controller.CtxUserIdKey, mc.UserId)
		c.Next() // 后续的处理函数的请求中，可以用过c.Get(CtxUserIdKey)来获取当前请求的用户信息
	}
}
