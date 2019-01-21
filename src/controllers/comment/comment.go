package comment

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"com/logging"
	"net/http"
	"com/e"
	"db"
	"com/gmysql"
	"com/setting"
	"fmt"
	"strconv"
)

//评论结构体
type CommentInfo struct {
	CommentAuthor      string `json:"comment_author" binding:"required"`
	CommentAuthorEmail string `json:"comment_author_email" binding:"required"`
	CommentContent     string `json:"comment_content" binding:"required"`
	CommentParentId    int64  `json:"comment_parent_id"`
	CommentPostId      int64  `json:"comment_post_id" binding:"required"`
}
type CommentEdit struct {
	CommentId       int64 `json:"comment_id" binding:"required"`
	CommentApproved int64 `json:"comment_approved"`
}

//添加评论
func AddComment(c *gin.Context) {
	ctx := controllers.Context{c}
	var comment CommentInfo
	if err := ctx.BindJSON(&comment); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		//获取用户ID
		userId, ok := ctx.Get("userId")
		if !ok {
			ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_GET_USER_FAIL, "")
			return
		}else {
			userId = 1
		}
		commentAuthorIP := ctx.ClientIP()
		res, err := gmysql.Con.Exec("INSERT INTO bc_comments "+
			"(comment_author,comment_author_email,comment_content,"+
			"comment_parent,comment_post_id,comment_author_IP,user_id) "+
			"VALUES (?,?,?,?,?,?,?)",
			comment.CommentAuthor, comment.CommentAuthorEmail, comment.CommentContent,
			comment.CommentParentId, comment.CommentPostId, commentAuthorIP, userId)
		num, err := res.RowsAffected()
		if err != nil || num == 0 {
			if num == 0 {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_COMMENT_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_COMMENT_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "评论审核中")
		}
	}
}

//获取全部评论
func AdminGetComments(c *gin.Context) {
	ctx := controllers.Context{c}
	//获取查询参数
	pageNum, err := strconv.Atoi(ctx.Query("page_num"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	}
	pageSize := setting.AppSetting.UserPageSize
	offsetSize := int64(pageNum-1) * pageSize
	//开启一个事务
	tx, err := gmysql.Con.Begin()
	rows, err := tx.Query("SELECT comment_id,comment_post_id,post_title,"+
		"comment_author,comment_author_IP,comment_content,comment_date,comment_approved"+
		" FROM bc_comments c,bc_posts p "+
		"WHERE c.comment_post_id=p.post_id AND 1=? Limit ?,?",
		1, offsetSize, pageSize)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMY, "")
		return
	}
	//查询数据
	data, err := db.Querys(rows)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMY, "")
		return
	}
	//查询数量
	rows2, err := tx.Query("SELECT FOUND_ROWS() AS row_counts Limit ?", 1)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMY, "")
		return
	}
	data2, err := db.Querys(rows2)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMY, "")
		return
	}
	//提交事务
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMY, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
		"user_total": data2[0]["row_counts"],
		"page_num":   pageNum,
		"user_list":  data,
	})

}

//编辑评论状态
func AdminEditComment(c *gin.Context) {
	ctx := controllers.Context{c}
	var cEdit CommentEdit
	if err := ctx.BindJSON(&cEdit); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	} else {
		//开启事务
		tx, err := gmysql.Con.Begin()
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
		}
		fmt.Println(cEdit)
		res, err := tx.Exec("UPDATE bc_comments SET comment_approved=? WHERE comment_id=?",
			cEdit.CommentApproved, cEdit.CommentId)
		num, err := res.RowsAffected()
		//id, err := res.LastInsertId()
		//fmt.Println(res,num)
		if err != nil || num == 0 {
			if num == 0 {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
			}
			tx.Rollback()
			return
		}
		//评论状态判断
		if cEdit.CommentApproved == 0 {
			res2, err := tx.Exec("UPDATE bc_posts SET comment_count=comment_count+1 "+
				"WHERE post_id=(SELECT comment_post_id FROM bc_comments WHERE comment_id=?)", cEdit.CommentId)
			num2, err := res2.RowsAffected()
			if err != nil || num2 == 0 {
				if num2 == 0 {
					logging.Error(err)
					ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
				} else {
					logging.Error(err)
					ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
				}
				tx.Rollback()
				return
			}
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
				return
			}
		} else {
			res2, err := tx.Exec("UPDATE bc_posts SET comment_count=comment_count-1 "+
				"WHERE post_id=(SELECT comment_post_id FROM bc_comments WHERE comment_id=?)", cEdit.CommentId)
			num2, err := res2.RowsAffected()
			if err != nil || num2 == 0 {
				if num2 == 0 {
					logging.Error(err)
					ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
				} else {
					logging.Error(err)
					ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
				}
				tx.Rollback()
				return
			}
			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_COMMENT, "")
				return
			}
		}
		ctx.Response(http.StatusOK, e.SUCCESS, "编辑成功")
	}
}

//删除评论
func AdminDeleteComment(c *gin.Context) {
	ctx := controllers.Context{c}
	var com CommentEdit
	if err := ctx.BindJSON(&com); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	} else {
		num, _, err := db.QRUDExec("DELETE FROM bc_comments WHERE comment_id=?",
			com.CommentId)
		if err != nil || num == 0 {
			if num == 0 {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_DELETE_COMMENT, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_DELETE_COMMENT, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "删除成功")
		}
	}
}
