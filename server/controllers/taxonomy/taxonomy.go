package taxonomy

import (
	"darco/proto/controllers"
	"darco/proto/models/taxonomy"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary List anchor taxa
// @Description Anchors are taxa that were imported as the root of a subtree in the taxonomy.
// @id TaxonAnchors
// @tags Taxonomy
// @Success 200 {array} taxonomy.TaxonDB "Get anchor taxa list success"
// @Router /taxonomy/anchors [get]
func ListAnchors(ctx *gin.Context, db *edgedb.Client) {
	anchors, err := taxonomy.ListAnchorTaxa(db)
	if err != nil {
		ctx.Error(err).SetMeta(gin.H{
			"msg": "Failed to fetch taxonomy data",
		})
	} else {
		ctx.JSON(http.StatusOK, anchors)
	}
}

// @Summary List taxa
// @Description Lists taxa, optionally filtered by name, rank and status
// @id TaxonomyList
// @tags Taxonomy
// @Success 200 {array} taxonomy.TaxonSelect "Get taxon success"
// @Router /taxonomy/ [get]
// @Param pattern query string false "Name search pattern" minlength(2)
// @Param rank query taxonomy.TaxonRank false "Taxonomic rank"
// @Param status query taxonomy.TaxonStatus false "Taxonomic status"
func ListTaxa(ctx *gin.Context, db *edgedb.Client) {
	pattern := ctx.Query("pattern")
	rank := taxonomy.TaxonRank(ctx.Query("rank"))
	status := taxonomy.TaxonStatus(ctx.Query("status"))
	taxa, err := taxonomy.ListTaxa(db, pattern, rank, status)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, taxa)
	}
}

// @Summary Get a taxon by its code
// @Description
// @id GetTaxon
// @tags Taxonomy
// @Success 200 {object}  taxonomy.TaxonSelect "Get taxon success"
// @Failure 404
// @Router /taxonomy/{code} [get]
// @Param code path string true "Taxon code" minlength(3)
func GetTaxon(ctx *gin.Context, db *edgedb.Client) {
	code := ctx.Param("code")
	taxon, err := taxonomy.FindByCode(db, code)
	if err != nil {
		logrus.Debug(err)
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, taxon)
	}
}

// @Summary Delete a taxon by its code
// @Description
// @id DeleteTaxon
// @tags Taxonomy
// @Success 200
// @Failure 403
// @Failure 404
// @Router /taxonomy/{code} [delete]
// @Param code path string true "Taxon code" minlength(3)
func DeleteTaxon(ctx *gin.Context, db *edgedb.Client) {
	code, err := controllers.ParseCodeURI(ctx)
	if err != nil {
		return
	}
	err = taxonomy.Delete(code)
	if err != nil {
		ctx.Error(err)
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary Update a taxon by its code
// @Description
// @id UpdateTaxon
// @tags Taxonomy
// @Success 200 {object} taxonomy.TaxonSelect
// @Failure 403
// @Failure 404
// @Router /taxonomy/{code} [patch]
// @Param code path string true "Taxon code" minlength(3)
// @Param data body taxonomy.TaxonUpdate true "Taxon"
// func UpdateTaxon(ctx *gin.Context, db *edgedb.Client) {
// 	taxon, err := controllers.BindUpdateByCode[taxonomy.TaxonSelect](ctx, db, taxonomy.FindByCode)
// 	if err != nil {
// 		return
// 	}
// 	updatedTaxon, err := taxon.Update(db)
// 	if err != nil {
// 		ctx.Error(err)
// 	} else {
// 		ctx.JSON(http.StatusOK, updatedTaxon)
// 	}
// }
