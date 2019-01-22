package post

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"com/logging"
	"net/http"
	"com/e"
	"db"
	"com/gmysql"
	"com/setting"
	"strconv"
)

type AddPostJson struct {
	PostTitle     string `json:"post_title" binding:"required"`
	PostExcerpt   string `json:"post_excerpt" binding:"required"`
	PostContent   string `json:"post_content" binding:"required"`
	PostPreImgUrl string `json:"post_pre_img_url"`
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
		userId, ok := ctx.Get("userId")
		if !ok {
			ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_GET_USER_FAIL, "")
			return
		}
		//数据插入
		_, res, err := db.QRUDExec("INSERT INTO bc_posts (post_author,post_title,post_excerpt,"+
			"post_content,post_pre_img_url) VALUES (?,?,?,?,?) ", userId, post.PostTitle,
			post.PostExcerpt, post.PostContent, post.PostPreImgUrl)
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

//获取文章列表
func AdminGETPosts(c *gin.Context) {
	ctx := controllers.Context{c}
	//对查询的参数进行类型转化
	pageNum, err := strconv.Atoi(ctx.Query("page_num"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	pageSize := setting.AppSetting.PostPageSize
	offsetSize := int64(pageNum-1) * pageSize
	//开启一个事务
	tx, err := gmysql.Con.Begin()
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	rows, err := tx.Query("SELECT post_id,user_nicename,post_date,"+
		"post_title,post_excerpt,post_status,comment_status,post_modified,"+
		"comment_count FROM bc_users u,bc_posts p "+
		"WHERE u.user_id=p.post_author Limit ?,?",
		offsetSize, pageSize)
	res, err := db.Querys(rows)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	rows2, err := tx.Query("SELECT FOUND_ROWS() AS row_counts LIMIT ?", 1)
	res2, err := db.Querys(rows2)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
		"post_total": res2[0]["row_counts"],
		"page_num":   pageNum,
		"post_list":  res,
	})
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

//获取分类文章列表
func GetTermPosts(c *gin.Context) {
	ctx := controllers.Context{c}
	termId, err := strconv.Atoi(ctx.Query("term_id"))
	pageNum, err := strconv.Atoi(ctx.Query("page_num"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	//展示数目
	pageSize := setting.AppSetting.PostPageSize
	//偏移量
	offsetSize := (pageNum - 1) * int(pageSize)
	//开启事物
	tx, err := gmysql.Con.Begin()
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	//查询
	rows, err := tx.Query("SELECT post_id,post_author,post_date,post_content,post_title,"+
		"post_excerpt,post_pre_img_url,comment_count "+
		"FROM bc_posts WHERE bc_posts.post_id IN (SELECT object_id "+
		"FROM bc_term_relationships tr,bc_term_taxonomy ta "+
		"WHERE tr.term_taxonomy_id=ta.term_taxonomy_id AND ta.term_id=?)"+
		"Limit ?,?", termId, offsetSize, pageSize)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	res, err := db.Querys(rows)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	rows2, err := tx.Query("SELECT FOUND_ROWS() AS row_counts Limit ?", 1)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	res2, err := db.Querys(rows2)
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	//提交事务
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.ERROR_GET_ARTICLES_FAIL, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
		"post_total": res2[0]["row_counts"],
		"page_num":   pageNum,
		"post_list":  res,
	})
}

//获取文章详情
func GetPost(c *gin.Context) {
	ctx := controllers.Context{c}
	postId, err := strconv.Atoi(ctx.Query("post_id"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	//开启事务
	tx, err := gmysql.Con.Begin()
	rows, err := tx.Query("SELECT user_nicename,post_date,post_content,"+
		"post_title,post_excerpt,comment_count FROM bc_posts p,bc_users u "+
		"WHERE post_id=? AND p.post_author=u.user_id", postId)
	res, err := db.Querys(rows)
	if err != nil || len(res) == 0 {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, "")
		return
	}
	rows2, err := tx.Query("SELECT comment_id,comment_author,comment_date,"+
		"comment_content,comment_parent FROM bc_comments "+
		"WHERE comment_post_id=? AND comment_approved=?", postId, 0)
	res2, err := db.Querys(rows2)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, "")
		return
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
		"post_author":   res[0]["user_nicename"],
		"post_date":     res[0]["post_date"],
		"post_content":  res[0]["post_content"],
		"post_title":    res[0]["post_title"],
		"post_excerpt":  res[0]["post_excerpt"],
		"comment_count": res[0]["comment_count"],
		"comment_list":  res2,
	})
}
