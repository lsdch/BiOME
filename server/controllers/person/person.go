package person

import (
	"darco/proto/controllers"
	"darco/proto/models/person"

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
	controllers.ListItems[person.Person](ctx, db, person.List)
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
// @id DeletePerson
// @tags People
// @Success 200 {object} person.Person
// @Failure 404 "person does not exist"
// @Router /people/persons/{id} [delete]
// @Param id path string true "Item UUID"
func Delete(ctx *gin.Context, db *edgedb.Client) {
	controllers.DeleteByID[person.Person](ctx, db, person.Delete)
}

// @Summary Update person
// @id UpdatePerson
// @tags People
// @Success 200 {object} person.Person
// @Failure 400 {object} validations.FieldErrors
// @Router /people/persons/{id} [patch]
// @Param id path string true "Item UUID"
// @Param data body person.PersonUpdate true "Update infos"
func Update(ctx *gin.Context, db *edgedb.Client) {
	controllers.UpdateByID[person.Person](ctx, db, person.Find)
}
