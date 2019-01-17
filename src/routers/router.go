package routers

import (
	"github.com/gin-gonic/gin"
	"controllers/user"
	"controllers/post"
	"controllers/term"
)

//路由初始化
func SetUp(r *gin.Engine) {

	//请求 api v1.0版本
	rV1 := r.Group("/v1")

	//界面API
	rV1.POST("/user/login", user.Login)
	rV1.POST("/user/register", user.Register)

	//后台API
	//用户操作
	rV1Admin := rV1.Group("/admin")
	rV1Admin.GET("/users", user.GetUsers)
	rV1Admin.DELETE("/user", user.DelUser)
	//文章操作
	rV1Admin.POST("/post", post.AdminAddPost)
	rV1Admin.GET("/posts",post.AdminGETPost)
	rV1Admin.PUT("/post/title",post.AdminEditTitlePost)
	rV1Admin.PUT("/post/status",post.AdminEditStatusPost)
	rV1Admin.PUT("/post/comment/status",post.AdminEditCommentStatusPost)
	rV1Admin.DELETE("/post",post.AdminDeletePost)
	//分类操作
	rV1Admin.POST("/taxonomy/term",term.AdminAddTaxonomyTerm)
}
