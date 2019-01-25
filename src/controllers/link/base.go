package link

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"net/http"
	"com/e"
	"db"
	"com/logging"
	"strconv"
	"com/gmysql"
	"com/setting"
)

//链接
type Link struct {
	LinkUrl         string `json:"link_url" binding:"required"`
	LinkName        string `json:"link_name" binding:"required"`
	LinkImageUrl    string `json:"link_image_url"`
	LinkDescription string `json:"link_description"`
}
type EditeLink struct {
	LinkId      int64 `json:"link_id" binding:"required"`
	LinkVisible int64 `json:"link_visible"`
}

type DelLink struct {
	LinkId int64 `json:"link_id" binding:"required"`
}

//添加链接
func AdminAddLink(c *gin.Context) {
	ctx := controllers.Context{c}
	var link Link
	if err := ctx.BindJSON(&link); err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		//获取用户ID
		userId, ok := ctx.Get("userId")
		if !ok {
			ctx.Response(http.StatusUnauthorized, e.ERROR_AUTH_GET_USER_FAIL, "")
			return
		}
		//数据插入
		num, _, err := db.QRUDExec("INSERT INTO bc_links (link_url,link_name,link_image_url,"+
			"link_description,link_owner) VALUES (?,?,?,?,?) ",
			link.LinkUrl, link.LinkName, link.LinkImageUrl, link.LinkDescription, userId)
		if err != nil || num == 0 {
			if num == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_LINK, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_ADD_LINK, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "添加成功")
		}
	}
}

//获取列表
func AdminGetLinks(c *gin.Context) {
	ctx := controllers.Context{c}
	pageNum, err := strconv.Atoi(ctx.Query("page_num"))
	if err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		//产生事务
		tx, err := gmysql.Con.Begin()
		if err != nil {
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_LINKS, "")
			return
		}
		pageSize := setting.AppSetting.CommentPageSize
		offsetSize := int64(pageNum-1) * pageSize
		rows, err := tx.Query("SELECT link_id,link_name,link_url,link_description,"+
			"user_nicename,link_visible,link_image_url FROM bc_users u,bc_links l "+
			"WHERE l.link_owner=u.user_id LIMIT ?,?", offsetSize, pageSize)
		res, err := db.Querys(rows)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_LINKS, "")
			return
		}
		rows2, err := tx.Query("SELECT count(link_id) as row_counts FROM bc_links LIMIT ?", 1)
		res2, err := db.Querys(rows2)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_LINKS, "")
			return
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_LINKS, "")
			return
		}
		ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
			"link_total": res2[0]["row_counts"],
			"page_num":   pageNum,
			"link_list":  res,
		})
	}
}

//编辑链接
func AdminEditLink(c *gin.Context) {
	ctx := controllers.Context{c}
	var link EditeLink
	if err := ctx.BindJSON(&link); err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	//数据插入
	num, _, err := db.QRUDExec("UPDATE bc_links SET link_visible=? WHERE link_id=?",
		link.LinkVisible, link.LinkId)
	if err != nil || num == 0 {
		if num == 0 {
			ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_LINK, "")
		} else {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_LINK, "")
		}
	} else {
		ctx.Response(http.StatusOK, e.SUCCESS, "编辑成功")
	}

}

//删除链接
func AdminDelLink(c *gin.Context) {
	ctx := controllers.Context{c}
	var link DelLink
	if err := ctx.BindJSON(&link); err != nil {
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	}
	//数据删除
	num, _, err := db.QRUDExec("DELETE FROM bc_links WHERE link_id=?", link.LinkId)
	if err != nil || num == 0 {
		if num == 0 {
			ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_LINK, "")
		} else {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_EDIT_LINK, "")
		}
	} else {
		ctx.Response(http.StatusOK, e.SUCCESS, "删除成功")
	}
}
