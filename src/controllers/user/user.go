package user

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"net/http"
	"com/e"
	"com/gmysql"
	"com/logging"
	"utils"
	"com/setting"
)

//登录路由
func Login(c *gin.Context) {
	ctx := controllers.Context{c}
	// 这个将通过 content-type 头去推断绑定器使用哪个依赖。
	// application/x-www-form-urlencoded
	// application/json
	//fmt.Println(ctx.ContentType())
	if ctx.ContentType() == "application/x-www-form-urlencoded" {
		//查询数据
		var user controllers.User
		var userList []controllers.User
		rows, err := gmysql.Con.Query("SELECT user_id,user_login,user_nicename,user_email,user_registered,user_status "+
			"FROM bc_users WHERE user_login=? AND user_pass=? ", ctx.PostForm("user_login"), ctx.PostForm("user_pass"))
		defer rows.Close()
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR, "")
			return
		}
		//遍历数据
		for rows.Next() {
			err = rows.Scan(&user.UserId, &user.UserLogin,
				&user.UserNicename, &user.UserEmail, &user.UserRegistered, &user.UserStatus)
			if err != nil {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR, "")
				return
			}
			userList = append(userList, user)
		}
		//结果判断
		switch len(userList) {
		case 0:
			ctx.Response(http.StatusBadRequest, e.ERROR_EXITS_USER, "")
			return
		case 1:
			//token中用户信息
			userInfo := make(map[string]interface{})
			userInfo["user_id"] = userList[0].UserId
			//生成token
			token, err := utils.CreateToken(setting.AppSetting.JwtKey, userInfo)
			if err != nil {
				logging.Error(err)
				ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_TOKEN, "")
				return
			}
			ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
				"user_id":         userList[0].UserId,
				"user_nicename":   userList[0].UserNicename,
				"user_email":      userList[0].UserEmail,
				"user_registered": userList[0].UserRegistered,
				"user_status":     userList[0].UserStatus, // 0 正常 1 异常
				"token":           token,
			})
			return
		default:
			ctx.Response(http.StatusBadRequest, e.ERROR_EXITS_USER_REPEAT, "")
		}

	} else {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
}

//注册路由
func Register(c *gin.Context) {
	ctx := controllers.Context{c}
	if ctx.ContentType() == "application/x-www-form-urlencoded" {

	} else {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
}
