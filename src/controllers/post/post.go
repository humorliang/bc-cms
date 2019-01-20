package post

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
	"strconv"
)

type AddPostJson struct {
	PostTitle   string `json:"post_title" binding:"required"`
	PostExcerpt string `json:"post_excerpt" binding:"required"`
	PostContent string `json:"post_content" binding:"required"`
}
type EditPostTitleJson struct {
	PostId    int    `json:"post_id" binding:"required"`
	PostTitle string `json:"post_title" binding:"required"`
}
type EditPostStatusJson struct {
	PostId     int `json:"post_id" binding:"required"`
	PostStatus int `json:"post_status" binding:"required"`
}
type EditPostCommentStatusJson struct {
	PostId        int `json:"post_id" binding:"required"`
	CommentStatus int `json:"comment_status" binding:"required"`
}
type PostId struct {
	PostId int `json:"post_id" binding:"required"`
}

//添加文章
func AdminAddPost(c *gin.Context) {
	ctx := controllers.Context{c}
	var post AddPostJson
	if err := ctx.BindJSON(&post); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	} else {
		//获取用户ID
		//userId, ok := ctx.Get("userId")
		//if !ok {
		//	ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_GET_USER_FAIL, "")
		//	return
		//}
		userId := 1
		//数据插入
		_, res, err := db.QRUDExec("INSERT INTO bc_posts (post_author,post_title,post_excerpt,post_content) VALUES (?,?,?,?) ",
			userId, post.PostTitle, post.PostExcerpt, post.PostContent)
		if err != nil || res == 0 {
			if res == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "添加成功")
		}
	}
}

//获取文章  bug post_id类型bug 添加？
func AdminGETPost(c *gin.Context) {
	ctx := controllers.Context{c}
	var pageNum controllers.PageNum
	if err := ctx.ShouldBindQuery(&pageNum); err == nil {
		pageSize := setting.AppSetting.PostPageSize
		offsetSize := (pageNum.PageNum - 1) * pageSize
		//开启一个事务
		tx, err := gmysql.Con.Begin()
		sql1 := fmt.Sprintf("SELECT post_id,user_nicename,post_date,"+
			"post_title,post_excerpt,post_status,comment_status,post_modified,"+
			"comment_count FROM bc_users u,bc_posts p "+
			"WHERE u.user_id=p.post_author Limit %d,%d",
			offsetSize, pageSize)
		//事务查询
		resList, err := db.TranscationQuerys(tx,
			sql1, "SELECT FOUND_ROWS() AS row_counts")
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, "")
			return
		}
		resPostInfo := resList[0]
		resPosts := resList[1]
		ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
			"post_total": resPosts[0]["row_counts"],
			"page_num":   pageNum.PageNum,
			"post_list":  resPostInfo,
		})
	} else {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	}
}

//编辑文章标题
func AdminEditTitlePost(c *gin.Context) {
	ctx := controllers.Context{c}
	var eInfo EditPostTitleJson
	if err := ctx.BindJSON(&eInfo); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_EDIT_ARTICLE_FAIL, "")
	} else {
		res, _, err := db.QRUDExec("UPDATE bc_posts SET post_title=? WHERE post_id=?",
			eInfo.PostTitle, eInfo.PostId)
		if err != nil || res == 0 {
			if res == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "编辑成功")
		}
	}

}

//编辑文章状态
func AdminEditStatusPost(c *gin.Context) {
	ctx := controllers.Context{c}
	var eInfo EditPostStatusJson
	if err := ctx.BindJSON(&eInfo); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_EDIT_ARTICLE_FAIL, "")
	} else {
		res, _, err := db.QRUDExec("UPDATE bc_posts SET post_status=? WHERE post_id=?",
			eInfo.PostStatus, eInfo.PostId)
		if err != nil || res == 0 {
			if res == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "编辑成功")
		}
	}

}

//编辑文章评论状态
func AdminEditCommentStatusPost(c *gin.Context) {
	ctx := controllers.Context{c}
	var eInfo EditPostCommentStatusJson
	if err := ctx.BindJSON(&eInfo); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_EDIT_ARTICLE_FAIL, "")
	} else {
		res, _, err := db.QRUDExec("UPDATE bc_posts SET comment_status=? WHERE post_id=?",
			eInfo.CommentStatus, eInfo.PostId)
		if err != nil || res == 0 {
			if res == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "编辑成功")
		}
	}
}

//删除文章
func AdminDeletePost(c *gin.Context) {
	ctx := controllers.Context{c}
	var pId PostId
	if err := ctx.BindJSON(&pId); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_DELETE_ARTICLE_FAIL, "")
	} else {
		res, _, err := db.QRUDExec("DELETE FROM bc_posts WHERE post_id=?", pId.PostId)
		if err != nil || res == 0 {
			if res == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "删除成功")
		}
	}
}

//获取文章详情
func GetPost(c *gin.Context) {
	ctx := controllers.Context{c}
	termId, err := strconv.Atoi(ctx.Query("term_id"))
	pageNum, err := strconv.Atoi(ctx.Query("page_num"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, err)
		return
	}
	//展示数目
	pageSize := setting.AppSetting.PostPageSize
	//偏移量
	offsetSize := (pageNum - 1) * int(pageSize)
	//开启事物
	fmt.Println(termId,offsetSize)

}
