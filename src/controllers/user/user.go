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
	"github.com/gin-gonic/gin/binding"
	"db"
	"strconv"
)

//登录
type LoginInfo struct {
	UserLogin string `json:"user_login" form:"user_login" binding:"required"` //中间用空格隔开千万不要用 逗号
	UserPass  string `json:"user_pass" form:"user_pass" binding:"required"`
	AuthCode  string `json:"auth_code" form:"auth_code"`
}
type UserId struct {
	UserId int `json:"user_id" binding:"required"`
}

//注册
type RegisterInfo struct {
	UserLogin    string `json:"user_login" form:"user_login" binding:"required"` //中间用空格隔开千万不要用 逗号
	UserPass     string `json:"user_pass" form:"user_pass" binding:"required"`
	UserNicename string `json:"user_nicename" form:"user_nicename" `
	UserEmail    string `json:"user_email" form:"user_email"`
}

//登录路由
func Login(c *gin.Context) {
	ctx := controllers.Context{c}
	//数据绑定
	loginInfo := &LoginInfo{}
	if ctx.ContentType() == "application/x-www-form-urlencoded" {
		err := ctx.MustBindWith(loginInfo, binding.Form)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
			return
		}
	} else if ctx.ContentType() == "application/json" {
		err := ctx.BindJSON(loginInfo)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
			return
		}
	} else {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	//数据查询
	rows, err := gmysql.Con.Query("SELECT user_id,user_login,user_nicename,user_email,user_registered,user_status "+
		"FROM bc_users WHERE user_login=? AND user_pass=? ", loginInfo.UserLogin, utils.Md5Encrypt(loginInfo.UserPass))
	defer rows.Close()
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR, "")
		return
	}
	//返回查询结果集合
	userList, err := db.Querys(rows)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR, "")
		return
	}
	//结果判断
	switch len(userList) {
	case 0:
		ctx.Response(http.StatusBadRequest, e.ERROR_EXITS_USER, "")
		return
	case 1:
		//token中用户信息
		userInfo := make(map[string]interface{})
		userInfo["user_id"] = userList[0]["user_id"]
		//生成token
		token, err := utils.CreateToken(setting.AppSetting.JwtKey, userInfo)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_TOKEN, "")
			return
		}
		ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
			"user_id":         userList[0]["user_id"],
			"user_nicename":   userList[0]["user_nicename"],
			"user_email":      userList[0]["user_email"],
			"user_registered": userList[0]["user_registered"],
			"user_status":     userList[0]["user_status"], // 0 正常 1 异常
			"token":           token,
		})
		return
	default:
		ctx.Response(http.StatusBadRequest, e.ERROR_EXITS_USER_REPEAT, "")
		return
	}

}

//注册路由
func Register(c *gin.Context) {
	ctx := controllers.Context{c}
	regInfo := &RegisterInfo{}
	if ctx.ContentType() == "application/x-www-form-urlencoded" {
		err := ctx.MustBindWith(regInfo, binding.Form)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
			return
		}
	} else if ctx.ContentType() == "application/json" {
		err := ctx.BindJSON(regInfo)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
			return
		}
	} else {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	//用户判断
	//数据查询
	rows, err := gmysql.Con.Query("SELECT user_id,user_login,user_nicename,user_email,user_registered,user_status "+
		"FROM bc_users WHERE user_login=?", regInfo.UserLogin)
	defer rows.Close()
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_EXITS_REGISTER, "")
		return
	}
	//返回查询结果集合
	userList, err := db.Querys(rows)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_EXITS_REGISTER, "")
		return
	}
	if len(userList) != 0 {
		ctx.Response(http.StatusInternalServerError, e.ERROR_EXITS_REGISTER_USER, "")
		return
	}
	//数据插入
	res, err := gmysql.Con.Exec("INSERT INTO bc_users (user_login,user_pass,user_nicename,user_email) "+
		"VALUES (?,?,?,?) ", regInfo.UserLogin, utils.Md5Encrypt(regInfo.UserPass), regInfo.UserNicename, regInfo.UserEmail)
	if err != nil {
		ctx.Response(http.StatusInternalServerError, e.ERROR_EXITS_REGISTER_USER, "")
		return
	}
	//插入的Id
	uId, err := res.LastInsertId()
	if err != nil {
		ctx.Response(http.StatusInternalServerError, e.ERROR_EXITS_REGISTER_USER, "")
		return
	}
	//生成token
	//token中用户信息
	userInfo := make(map[string]interface{})
	userInfo["user_id"] = uId
	token, err := utils.CreateToken(setting.AppSetting.JwtKey, userInfo)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_TOKEN, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
		"user_id":       uId,
		"user_login":    regInfo.UserLogin,
		"user_nicename": regInfo.UserNicename,
		"user_email":    regInfo.UserEmail,
		"token":         token,
	})
}

//获取用户列表
func GetUsers(c *gin.Context) {
	ctx := controllers.Context{c}
	pageNum, err := strconv.Atoi(ctx.Query("page_num"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	pageSize := setting.AppSetting.UserPageSize
	offsetSize := int64(pageNum-1) * pageSize
	//开启一个事务
	tx, err := gmysql.Con.Begin()
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_USERS, "")
		return
	}
	rows, err := tx.Query("SELECT user_id,user_nicename,user_email,"+
		"user_registered,user_status FROM bc_users Limit %d,%d",
		offsetSize, pageSize)
	res, err := db.Querys(rows)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_USERS, "")
		return
	}
	rows2, err := tx.Query("SELECT FOUND_ROWS() AS row_counts LIMIT ?", 1)
	res2, err := db.Querys(rows2)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_USERS, "")
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_USERS, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
		"user_total": res2[0]["row_counts"],
		"page_num":   pageNum,
		"user_list":  res,
	})

}

//删除冻结用户
func DelUser(c *gin.Context) {
	ctx := controllers.Context{c}
	var user UserId
	if err := ctx.BindJSON(&user); err != nil {
		ctx.Response(http.StatusBadRequest, e.ERROR_EXITS_SUPER_USER, "")
	} else {
		//数据操作
		res, _, err := db.QRUDExec("DELETE FROM bc_users WHERE user_id=?", user.UserId)
		if err != nil || res == 0 {
			if res == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_NOT_USER, "")
			} else {
				ctx.Response(http.StatusInternalServerError, e.ERROR_DELETE_USER_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "删除成功")
		}
	}
}
