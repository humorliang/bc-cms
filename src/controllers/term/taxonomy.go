package term

import (
	"github.com/gin-gonic/gin"
	"controllers"
	"net/http"
	"com/logging"
	"com/e"
	"db"
)

type TaxonomyTerm struct {
	TermName     string `json:"term_name" binding:"required"`
	Taxonomy     string `json:"taxonomy" binding:"required"`
	Description  string `json:"description"`
	TermParentId int `json:"term_parent_id"`
}

func AdminAddTaxonomyTerm(c *gin.Context) {
	ctx := controllers.Context{c}
	var taxTerm TaxonomyTerm
	if err := ctx.BindJSON(&taxTerm); err != nil {
		logging.Error(err)
		ctx.Response(http.StatusBadRequest, e.INVALID_PARAMS, "")
	} else {
		//添加分类项
		rowNum, rowId, err := db.QRUDExec("INSERT INTO bc_terms (term_name) VALUES (?)", taxTerm.TermName)
		if err != nil || rowNum == 0 {
			if rowNum == 0 {
				ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
			} else {
				logging.Error(err)
				ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
			}
		} else {
			//将创建分类法并添加分类项
			rowNum2, rowId2, err := db.QRUDExec("INSERT INTO bc_term_taxonomy (term_id,taxonomy,description,term_parent_id) VALUES (?,?,?,?)",
				rowId, taxTerm.Taxonomy, taxTerm.Description, taxTerm.TermParentId)
			if err != nil || rowNum2 == 0 {
				if rowNum2 == 0 {
					ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
				} else {
					logging.Error(err)
					ctx.Response(http.StatusInternalServerError, e.ERROR_TAXONOMY_TERM_FAIL, "")
				}
			} else {
				ctx.Response(http.StatusOK, e.SUCCESS, gin.H{
					"term_id":          rowId,
					"term_taxonomy_id": rowId2,
				})
			}
		}
	}
}
