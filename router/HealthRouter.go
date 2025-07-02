package router

import (
	"git.exclouds.org/ww/greet_ws/global"
	"github.com/gin-gonic/gin"
)

func HealthRouter(Router *gin.RouterGroup) {
	// 健康检查
	Router.GET("/health", func(c *gin.Context) {

		// 检查 Redis 连接（如果启用）
		redisHealth := true
		if global.SERVER_CONFIG.SystemConfig.UseMultipoint && global.REDIS != nil {
			if _, err := global.REDIS.Ping().Result(); err != nil {
				redisHealth = false
			}
		}

		// 所有依赖都健康时才返回200
		if redisHealth {
			c.JSON(200, gin.H{
				"status": "healthy",
				"detail": gin.H{
					"database": "ok",
					"redis":    "ok",
				},
			})
			return
		}

		// 否则返回503表示服务不可用
		detail := gin.H{}

		if !redisHealth {
			detail["redis"] = "error"
		} else {
			detail["redis"] = "ok"
		}

		c.JSON(503, gin.H{
			"status": "unhealthy",
			"detail": detail,
		})
	})
}
