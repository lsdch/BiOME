package institution

import (
	"darco/proto/controllers"
	"darco/proto/models/institution"
	"net/http"

	_ "darco/proto/models/validations"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

// List institutions
// swagger:route GET /people/institutions/
// @Summary List Institutions
// @ID List Institutions
// @Tags People
// @Success 200 {array} institution.Institution
// @Router /people/institutions [get]
func List(ctx *gin.Context, db *edgedb.Client) {
	items, err := institution.List(db)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, items)
	}
}

// @Summary Create institution
// @Description Register a new institution that people work in.
// @id CreateInstitution
// @tags People
// @Success 201 {object} institution.Institution
// @Failure 400 {object} validations.FieldErrors
// @Router /people/institutions [post]
// @Param data body institution.InstitutionInput true "Institution informations"
func Create(ctx *gin.Context, db *edgedb.Client) {
	controllers.CreateItem[institution.Institution, institution.InstitutionInput](ctx, db)
}

// @Summary Update institution
// @id UpdateInstitution
// @tags People
// @Success 200 {object} institution.Institution
// @Failure 400 {object} validations.FieldErrors
// @Router /people/institutions/{code} [patch]
// @Param code path string true "Institution code"
// @Param data body institution.InstitutionUpdate true "Institution informations"
func Update(ctx *gin.Context, db *edgedb.Client) {
	controllers.UpdateByCode[institution.Institution](ctx, db, institution.Find)
}

// @Summary Delete institution
// @id DeleteInstitution
// @tags People
// @Success 200 "Delete successful"
// @Failure 404 "Institution does not exist"
// @Router /people/institutions/{code} [delete]
// @Param code path string true "Institution short name"
func Delete(ctx *gin.Context, db *edgedb.Client) {
	controllers.DeleteByCode[institution.Institution](ctx, db, institution.Delete)
}
