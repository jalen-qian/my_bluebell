package router

import (
	"github.com/gin-gonic/gin"
	"github.com/web_app_base/controller"
	"github.com/web_app_base/logger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", controller.IndexHandler)
	return r
}
