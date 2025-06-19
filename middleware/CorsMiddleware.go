package middleware

import (
	"fmt"
	"git.exclouds.org/ww/greet_ws/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 处理跨域请求,支持options访问
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		add := addOrigin(c)
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			if add {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusUnavailableForLegalReasons)
			}
		}
		// 处理请求
		c.Next()
	}
}

func addOrigin(c *gin.Context) bool {
	add := false
	origin := c.Request.Header.Get("Origin")
	allowOrigins := global.SERVER_CONFIG.SystemConfig.AllowOrigins
	if len(allowOrigins) == 0 || len(origin) == 0 { //未配置
		add = true
	} else {
		for _, oriSuf := range allowOrigins {
			if strings.Contains(origin, oriSuf) {
				add = true
				break
			}
		}
	}

	if add {
		c.Header("Access-Control-Allow-Origin", origin) //添加origin
	} else {
		global.LOGGER.Warn(fmt.Sprintf("origin不匹配：%s", origin))
	}
	return add
}
