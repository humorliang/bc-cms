package middlerware

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"com/e"
	"utils"
	"com/setting"
	"com/logging"
)

//jwt 中间件
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := controllers.Context{c}

		//获取token
		token := ctx.Request.Header.Get("authorization")
		if token == "" {
			ctx.Response(http.StatusBadRequest, e.ERROR_AUTH, "")
		} else {
			//判断token
			claims, err := utils.ParseToken(token, setting.AppSetting.JwtKey)
			if err != nil {
				logging.Error(err)
				//token错误判断
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ctx.Response(http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, "")
					ctx.Abort()
					return

				default:
					ctx.Response(http.StatusBadRequest, e.ERROR_AUTH, "")
					ctx.Abort()
					return
				}
			}
			//获取userId
			if userId, ok := claims["user_id"]; ok {
				ctx.Set("userId", userId)
			} else {
				ctx.Response(http.StatusBadRequest, e.ERROR_AUTH_GET_USER_FAIL, "")
				ctx.Abort()
				return
			}
		}
		//交给下一个中间件
		ctx.Next()

	}
}
