package taxonomy

import (
	"darco/proto/controllers"
	"darco/proto/models/taxonomy"
	_ "darco/proto/models/validations"
	"net/http"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Summary List taxa
// @Description Lists taxa, optionally filtered by name, rank and status
// @id TaxonomyList
// @tags Taxonomy
// @Success 200 {array} taxonomy.TaxonSelect "Get taxon success"
// @Router /taxonomy/ [get]
// @Param filter body taxonomy.ListFilters false "Query filters"
func ListTaxa(ctx *gin.Context, db *edgedb.Client) {
	var filters = new(taxonomy.ListFilters)
	ctx.ShouldBindJSON(filters)
	taxa, err := taxonomy.ListTaxa(db, filters)
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
// @Success 200 {object} taxonomy.TaxonSelect "Get taxon success"
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

// @Summary Create taxon
// @Description This provides a way to register new unclassified taxa,
// @Description that have not yet been published to GBIF.
// @Description Importing taxa directly from GBIF should be preferred otherwise.
// @id CreateTaxon
// @tags Taxonomy
// @Success 201 {object} taxonomy.TaxonSelect
// @Failure 400 {object} validations.FieldErrors
// @Router /taxonomy/ [post]
// @Param data body taxonomy.TaxonInput true "New taxon"
func CreateTaxon(ctx *gin.Context, db *edgedb.Client) {
	controllers.CreateItem[taxonomy.TaxonInput, taxonomy.TaxonSelect](ctx, db)
}

// @Summary Delete a taxon by its code
// @Description
// @id DeleteTaxon
// @tags Taxonomy
// @Success 200 {object} taxonomy.TaxonSelect
// @Failure 403
// @Failure 404
// @Router /taxonomy/{code} [delete]
// @Param code path string true "Taxon code" minlength(3)
func DeleteTaxon(ctx *gin.Context, db *edgedb.Client) {
	controllers.DeleteByCode(ctx, db, taxonomy.Delete)
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
func UpdateTaxon(ctx *gin.Context, db *edgedb.Client) {
	controllers.UpdateItemByCode[taxonomy.TaxonUpdate](ctx, db, taxonomy.FindByID)
}
