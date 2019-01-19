package middlerware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//404中间件
func NotFoundPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println(c.Writer.Status())
		if c.Writer.Status() == 404 {
			c.JSON(http.StatusNotFound, gin.H{
				"code": http.StatusNotFound,
				"msg":  "请求的URL有误",
				"data": "",
			})
		} else {
			c.Next()
		}
	}
}
