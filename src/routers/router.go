package routers

import (
	"github.com/gin-gonic/gin"
)

//路由初始化
func SetUp(r *gin.Engine) {

	//请求 api v1.0版本
	rV1 := r.Group("/v1")

	//界面API
	rV1.GET("/user/login",)


	//后台API
	//rV1Admin := rV1.Group("/admin", middleware.JWTAuth())


}
