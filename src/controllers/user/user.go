package user

import (
	"github.com/gin-gonic/gin"
	"controllers"
)

type LoginForm struct {
	UserLogin string `json:"user_login",form:"user_login",binding:"required"`
	UserPass  string `json:"user_pass",form:"user_login",binding:"required"`
	AuthCode  string `json:"auth_code",form:"auth_code"`
}

//登录路由
func Login(c *gin.Context) {
	ctx := controllers.Context{c}
	userLogin:=ctx.PostForm("user_login")
}
