package taxonomy

import (
	"darco/proto/models/taxonomy"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary List anchor taxa
// @Description Anchors are taxa that were imported as the root of a subtree in the taxonomy.
// @id TaxonAnchors
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200 {array}  taxonomy.TaxonDB "Get anchor taxa list success"
// @Router /taxonomy/anchors [get]
func GetAnchors(ctx *gin.Context) {
	anchors, err := taxonomy.GetAnchorTaxa()
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
// @Accept json
// @Produce json
// @Success 200 {array} taxonomy.TaxonSelect "Get taxon success"
// @Router /taxonomy/ [get]
// @Param pattern query string false "Name search pattern" minlength(2)
// @Param rank query taxonomy.TaxonRank false "Taxonomic rank"
// @Param status query taxonomy.TaxonStatus false "Taxonomic status"
func ListTaxa(ctx *gin.Context) {
	pattern := ctx.Query("pattern")
	rank := taxonomy.TaxonRank(ctx.Query("rank"))
	status := taxonomy.TaxonStatus(ctx.Query("status"))
	taxa, err := taxonomy.ListTaxa(pattern, rank, status)
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
// @Accept json
// @Produce json
// @Success 200 {object}  taxonomy.TaxonSelect "Get taxon success"
// @Failure 404
// @Router /taxonomy/{code} [get]
// @Param code path string true "Taxon code" minlength(3)
func GetTaxon(ctx *gin.Context) {
	code := ctx.Param("code")
	taxon, err := taxonomy.GetTaxon(code)
	if err != nil {
		logrus.Debug(err)
		ctx.Error(err)
	} else if taxon == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	} else {
		ctx.JSON(http.StatusOK, taxon)
	}
}

// @Summary Delete a taxon by its code
// @Description
// @id DeleteTaxon
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200
// @Failure 403
// @Failure 404
// @Router /taxonomy/{code} [delete]
// @Param code path string true "Taxon code" minlength(3)
func DeleteTaxon(ctx *gin.Context) {
	code := ctx.Param("code")
	_, err := taxonomy.DeleteTaxon(code)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.Status(http.StatusOK)
	}
}

// @Summary Update a taxon by its code
// @Description
// @id UpdateTaxon
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200 {object} taxonomy.TaxonSelect
// @Failure 403
// @Failure 404
// @Router /taxonomy/{code} [patch]
// @Param code path string true "Taxon code" minlength(3)
// @Param data body taxonomy.TaxonInput true "Taxon"
func UpdateTaxon(ctx *gin.Context) {
	code := ctx.Param("code")
	taxon := taxonomy.TaxonInput{}
	err := ctx.ShouldBindJSON(&taxon)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	updatedTaxon, err := taxonomy.UpdateTaxon(code, taxon)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, updatedTaxon)
	}
}
