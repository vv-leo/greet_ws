package middleware

import (
	"bytes"
	"fmt"
	"git.exclouds.org/ww/greet_ws/global"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		// proto := c.Request.Proto
		ua := c.Request.UserAgent()

		token := c.Request.Header.Get("X-Token")
		if token == "" {
			token = c.Request.Header.Get("X-Secret")
		}

		productId := c.Request.Header.Get("X-ProductId")
		var params string
		if method != "OPTIONS" {
			ctype := c.ContentType()
			switch ctype {
			case "text/plain", "application/json":
				data, _ := ioutil.ReadAll(c.Request.Body)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				params = string(data)
			case "application/x-www-form-urlencoded", "multipart/form-data":
				c.Request.ParseMultipartForm(128)
				params = "{"
				for k, v := range c.Request.PostForm {
					params = params + fmt.Sprintf("\"%s\":\"%s\",", k, strings.Join(v, ","))
				}
				if c.Request.MultipartForm != nil {
					for _, v := range c.Request.MultipartForm.File {
						for i := 0; i < len(v); i++ {
							field := strings.Replace(strings.Split(v[i].Header["Content-Disposition"][0], ";")[1], " name=", "", 1)
							field = strings.Replace(field, "\"", "", 2)
							params = params + fmt.Sprintf("\"%s\":\"%s\",", field, v[i].Filename)
						}
					}
				}
				if len(params) > 1 {
					params = params[0:len(params)-1] + "}"
				} else {
					params = params + "}"
				}
			}
		}

		c.Next()

		latency := time.Since(requestTime)
		status := c.Writer.Status()
		ua = ""
		token = ""
		loginfo := fmt.Sprintf("%d %s %s %d %s %s \"%s\" %s %s %s", requestTime.Unix(), method, path, status, latency, clientIP, ua, token, productId, params)
		go func() {
			global.LOGGER.Info(loginfo)
		}()
	}
}
