package comment

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"com/logging"
	"net/http"
	"com/e"
	"db"
	"com/gmysql"
	"fmt"
	"com/setting"
)

//评论结构体
type CommentInfo struct {
	CommentAuthor      string `json:"comment_author" binding:"required"`
	CommentAuthorEmail string `json:"comment_author_email" binding:"required"`
	CommentContent     string `json:"comment_content" binding:"required"`
	CommentParentId    int    `json:"comment_parent_id"`
	CommentPostId      int    `json:"comment_post_id" binding:"required"`
}

//添加留言或者评论
func AddComment(c *gin.Context) {
	ctx := controllers.Context{c}
	var comment CommentInfo
	if err := ctx.BindJSON(&comment); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		//获取用户ID
		//userId, ok := ctx.Get("userId")
		//if !ok {
		//	ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_GET_USER_FAIL, "")
		//	return
		//}
		userId := 1
		num, id, err := db.QRUDExec("INSERT INTO bc_comments "+
			"(comment_author,comment_author_email,comment_content,"+
			"comment_parent,comment_post_id,user_id) VALUES (?,?,?,?,?,?)",
			comment.CommentAuthor, comment.CommentAuthorEmail, comment.CommentContent,
			comment.CommentParentId, comment.CommentPostId, userId)
		if err != nil || num == 0 {
			if num == 0 {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_COMMENT_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_COMMENT_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
				"comment_id": id,
			})
		}
	}
}

//获取全部评论
func AdminGetComments(c *gin.Context) {
	ctx := controllers.Context{c}
	var pageNum controllers.PageNum
	if err := ctx.ShouldBindQuery(&pageNum); err == nil {
		pageSize := setting.AppSetting.UserPageSize
		offsetSize := (pageNum.PageNum - 1) * pageSize
		//开启一个事务
		tx, err := gmysql.Con.Begin()
		sql1 := fmt.Sprintf("SELECT comment_id,comment_post_id,post_title,"+
			"comment_author,comment_author_IP,comment_content,comment_date,comment_approved"+
			" FROM bc_comments c,bc_posts p "+
			"WHERE c.comment_post_id=p.post_id AND ?=? Limit %d,%d",
			offsetSize, pageSize)
		resList, err := db.TranscationQuerys(tx,
			sql1, "SELECT FOUND_ROWS() AS row_counts")

		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_COMMENTS, "")
			return
		}
		commentInfo := resList[0]
		comments := resList[1]
		ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
			"user_total": comments[0]["row_counts"],
			"page_num":   pageNum.PageNum,
			"user_list":  commentInfo,
		})
	} else {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	}
}
