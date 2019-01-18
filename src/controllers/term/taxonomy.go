package term

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"net/http"
	"com/logging"
	"com/e"
	"db"
	"com/gmysql"
)

type TaxonomyTerm struct {
	TermName     string `json:"term_name" binding:"required"`
	Taxonomy     string `json:"taxonomy" binding:"required"`
	Description  string `json:"description"`
	TermParentId int    `json:"term_parent_id"`
}

type Taxonomy struct {
	Taxonomy string `json:"taxonomy" binding:"required"`
}

type Term struct {
	TermId int64 `json:"term_id" binding:"required"`
}

//添加分类
func AdminAddTaxonomyTerm(c *gin.Context) {
	ctx := controllers.Context{c}
	var taxTerm TaxonomyTerm
	//数据绑定
	if err := ctx.BindJSON(&taxTerm); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		//开启事务
		tx, err := gmysql.Con.Begin()
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
			return
		}
		//添加分类
		num, id, err := db.TxQRUDExec(tx, "INSERT INTO bc_terms (term_name) VALUES (?)", taxTerm.TermName)
		if err != nil || num == 0 {
			tx.Rollback()
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
			return
		} else {
			//添加分类法及分类关系
			num2, id2, err := db.TxQRUDExec(tx, "INSERT INTO bc_term_taxonomy (term_id,taxonomy,description,term_parent_id) VALUES (?,?,?,?)",
				id, taxTerm.Taxonomy, taxTerm.Description, taxTerm.TermParentId)
			if err != nil || num2 == 0 {
				tx.Rollback()
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
				return
			} else {
				tx.Commit()
				ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
					"term_id":          id,
					"term_taxonomy_id": id2,
				})
			}
		}
	}

}

//获取全部分类
func AdminGetTaxonomys(c *gin.Context) {
	ctx := controllers.Context{c}
	//如果没有？mysql会使用文本传输协议，有？就会进行预处理
	rows, err := gmysql.Con.Query("SELECT term_taxonomy_id,tm.term_id,term_name,taxonomy,"+
		"description,term_parent_id FROM bc_term_taxonomy tm,bc_terms tr "+
		"WHERE tm.term_id=tr.term_id AND 1=?", 1) //要添加一个mysql ？占位符 使mysql使用新的协议，进行预处理
	defer rows.Close()
	//结果数据
	data, err := db.Querys(rows)
	if err != nil {
		logging.Error(err)
		ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMYS, "")
		return
	}
	ctx.Response(http.StatusOK, e.SUCCESS, data)
}

//获取某个分类法
func AdminGetTaxonomy(c *gin.Context) {
	ctx := controllers.Context{c}
	var taxo Taxonomy
	if err := ctx.BindJSON(&taxo); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		rows, err := gmysql.Con.Query("SELECT term_taxonomy_id,tm.term_id,term_name,taxonomy,"+
			"description,term_parent_id FROM bc_term_taxonomy tm,bc_terms tr "+
			"WHERE tm.term_id=tr.term_id and tm.taxonomy=?", taxo.Taxonomy)
		defer rows.Close()
		//结果数据
		data, err := db.Querys(rows)
		if err != nil {
			logging.Error(err)
			ctx.Response(http.StatusInternalServerError, e.ERROR_GET_TAXONOMY, "")
			return
		}
		ctx.Response(http.StatusOK, e.SUCCESS, data)
	}
}

//删除分类
func AdminDelTaxonomy(c *gin.Context) {
	ctx := controllers.Context{c}
	var term Term
	if err := ctx.BindJSON(&term); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
		return
	} else {
		num, _, err := db.QRUDExec("DELETE FROM bc_terms WHERE term_id=?", term.TermId)
		if err != nil || num == 0 {
			if num == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_DELETE_TERM, "")
			} else {
				ctx.Response(http.StatusInternalServerError, e.ERROR_DELETE_TERM, "")
			}
		} else {
			ctx.Response(http.StatusOK, e.SUCCESS, "删除成功")
		}
	}
}
