package router

import (
	"github.com/bluebell/middlewares"
	"github.com/bluebell/settings"

	"github.com/bluebell/controller"
	"github.com/bluebell/logger"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//如果配置是release，就使用release模式启动，否则使用debug模式（默认）
	//注意这一行不能写在gin.New()之后
	if settings.Conf.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	//路由组v1
	v1Group := r.Group("/api/v1")
	//注册接口
	v1Group.POST("/signup", controller.SignUpHandler)
	//登录接口
	v1Group.POST("/login", controller.LoginHandler)
	v1Group.Use(middlewares.JWTAuthMiddleware())
	{
		//社区列表接口
		v1Group.GET("/community", controller.CommunityHandler)
		//社区详情接口
		v1Group.GET("/community/:id", controller.CommunityDetailHandler)
		//创建帖子接口
		v1Group.POST("/post", controller.AddPostHandler)
		//帖子详情接口
		v1Group.GET("/post/:postId", controller.PostDetailHandler)
	}

	return r
}
