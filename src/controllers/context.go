package controllers

import (
	"github.com/gin-gonic/gin"
	"com/e"
)

type Context struct {
	*gin.Context
}

//统一响应码
func (ctx *Context) Response(httpCode int, msg int, data interface{}) {
	if httpCode == 200 {
		ctx.JSON(httpCode, gin.H{
			"code": 0,
			"msg":  e.GetMsg(msg),
			"data": data,
		})
	} else {
		ctx.JSON(httpCode, gin.H{
			"code": msg,
			"msg":  e.GetMsg(msg),
			"data": data,
		})
	}
}
