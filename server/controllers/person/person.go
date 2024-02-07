package person

import (
	"darco/proto/controllers"
	"darco/proto/models/person"
	"net/http"

	_ "darco/proto/models/validations"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

// List persons
// @Summary List persons
// @Tags People
// @Success 200 {array} person.Person
// @Router /people/persons [get]
func List(ctx *gin.Context, db *edgedb.Client) {
	items, err := person.List(db)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.JSON(http.StatusOK, items)
	}
}

// @Summary Create person
// @Description Register a new person
// @id Createperson
// @tags People
// @Success 201 {object} person.Person
// @Failure 400 {object} validations.FieldErrors
// @Router /people/persons [post]
// @Param data body person.PersonInput true "Created person"
func Create(ctx *gin.Context, db *edgedb.Client) {
	controllers.CreateItem[person.Person, person.PersonInput](ctx, db)
}

// @Summary Delete person
// @id Deleteperson
// @tags People
// @Success 204 "Deleted item"
// @Failure 404 "person does not exist"
// @Router /people/persons/{id} [delete]
// @Param id path string true "Item UUID"
func Delete(ctx *gin.Context, db *edgedb.Client) {
	controllers.DeleteByID[person.Person](ctx, db, person.Delete)
}

// @Summary Update person
// @id Updateperson
// @tags People
// @Success 200 {object} person.Person
// @Failure 400 {object} validations.FieldErrors
// @Router /people/persons/{id} [patch]
// @Param id path string true "Item UUID"
// @Param data body person.PersonUpdate true "Update infos"
func Update(ctx *gin.Context, db *edgedb.Client) {
	controllers.UpdateByID[person.Person](ctx, db, person.Find)
}
