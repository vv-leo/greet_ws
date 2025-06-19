package initialize

import (
	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/middleware"
	"git.exclouds.org/ww/greet_ws/router"

	"github.com/gin-gonic/gin"
)

func RoutersInitialize() *gin.Engine {
	var Router = gin.Default()

	// 跨域
	Router.Use(middleware.CorsMiddleware())
	global.LOGGER.Info("use middleware cors")
	Router.Use(middleware.LoggerMiddleware())
	global.LOGGER.Info("use middleware logger")

	// 健康检查路由 - 直接挂载在根路径上，不需要认证
	router.HealthRouter(Router.Group(""))
	global.LOGGER.Info("health check router initialized")

	global.LOGGER.Info("router init success")
	return Router
}
