package middlerware

import (
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"com/logging"
)

//请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 路径
		var path string
		// Start timer
		start := time.Now()

		//交给request
		c.Next()

		//请求日志信息
		latency := time.Now().Sub(start)
		urlPath := c.Request.URL.Path
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMsg := c.Errors.String()
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = urlPath + "?" + raw
		} else {
			path = urlPath
		}
		clientIP := c.ClientIP()
		//请求信息
		reqInfo := fmt.Sprintf("[StartTime] %v [IP] %s [Path] %s [Method] %s [LtcyTime] %v [StatusCode] %d [ErrorInfo] %s\n",
			start, clientIP, path, method, latency, statusCode, errorMsg)
		logging.Info(reqInfo)
	}
}
