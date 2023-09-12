package taxonomy

import (
	"darco/proto/models/taxonomy"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary List anchor taxa
// @Description Anchors are taxa that were imported as the root of a subtree in the taxonomy.
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200 {array}  taxonomy.Taxon{authorship=string} "Get anchor taxa list success"
// @Failure 500 {string} gin.Error
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

// @Summary Get a taxon by its code
// @Description
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200 {object}  taxonomy.TaxonSelect{taxon=taxonomy.Taxon{children=[]taxonomy.Taxon}} "Get taxon success"
// @Failure 404
// @Failure 500 {string} gin.Error
// @Router /taxa/{code} [get]
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
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200
// @Failure 403
// @Failure 404
// @Failure 500 {string} gin.Error
// @Router /taxa/{code} [delete]
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
// @tags Taxonomy
// @Accept json
// @Produce json
// @Success 200 {object} taxonomy.TaxonSelect
// @Failure 403
// @Failure 404
// @Failure 500 {string} gin.Error
// @Router /taxa/{code} [patch]
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
