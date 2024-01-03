package routes

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"

	_ "bluebell/docs" // 千万不要忘了导入你上一步生成的docs

	"github.com/gin-contrib/pprof"
	gs "github.com/swaggo/gin-swagger"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.LoadHTMLFiles("./templates/index.html")
	// 加载静态文件，将HTML文件中所有有关/static文件的请求
	// 映射到./static文件夹中
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("api/v1")

	//v1.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "ok")
	//})
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlewares.JWTAuthMiddleWare())
	v1.GET("/posts", controller.GetPostListHandler)

	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	v1.POST("/post", controller.CreatePostHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	{
		// 根据时间或分数获取帖子列表
		v1.GET("/posts2", controller.GetPostListHandler2)
		// 投票
		v1.POST("/vote", controller.PostVoteController)
	}

	pprof.Register(r) //注册pprof相关路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
