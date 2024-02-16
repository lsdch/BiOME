package institution

import (
	"darco/proto/controllers"
	"darco/proto/models/people"

	_ "darco/proto/models/validations"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

// List institutions
// swagger:route GET /people/institutions/
// @Summary List Institutions
// @ID List Institutions
// @Tags People
// @Success 200 {array} people.Institution
// @Router /people/institutions [get]
func List(ctx *gin.Context, db *edgedb.Client) {
	controllers.ListItems(ctx, db, people.ListInstitutions)
}

// @Summary Create institution
// @Description Register a new institution that people work in.
// @id CreateInstitution
// @tags People
// @Success 201 {object} people.Institution
// @Failure 400 {object} validations.FieldErrors
// @Router /people/institutions [post]
// @Param data body people.InstitutionInput true "Institution informations"
func Create(ctx *gin.Context, db *edgedb.Client) {
	controllers.CreateItem[people.InstitutionInput, people.Institution](ctx, db)
}

// @Summary Update institution
// @id UpdateInstitution
// @tags People
// @Success 200 {object} people.Institution
// @Failure 400 {object} validations.FieldErrors
// @Router /people/institutions/{code} [patch]
// @Param code path string true "Institution code"
// @Param data body people.InstitutionUpdate true "Institution informations"
func Update(ctx *gin.Context, db *edgedb.Client) {
	controllers.UpdateItemByCode[people.InstitutionUpdate](ctx, db, people.FindInstitution)
}

// @Summary Delete institution
// @id DeleteInstitution
// @tags People
// @Success 200 {object} people.Institution "Deleted item"
// @Failure 404 "Institution does not exist"
// @Router /people/institutions/{code} [delete]
// @Param code path string true "Institution short name"
func Delete(ctx *gin.Context, db *edgedb.Client) {
	controllers.DeleteByCode(ctx, db, people.DeleteInstitution)
}
