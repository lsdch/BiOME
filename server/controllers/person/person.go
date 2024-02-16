package person

import (
	"darco/proto/controllers"

	"darco/proto/models/people"
	_ "darco/proto/models/validations"

	"github.com/edgedb/edgedb-go"
	"github.com/gin-gonic/gin"
)

// List persons
// @Summary List persons
// @Tags People
// @Success 200 {array} people.Person
// @Router /people/persons [get]
func List(ctx *gin.Context, db *edgedb.Client) {
	controllers.ListItems[people.Person](ctx, db, people.ListPersons)
}

// @Summary Create person
// @Description Register a new person
// @id Createperson
// @tags People
// @Success 201 {object} people.Person
// @Failure 400 {object} validations.FieldErrors
// @Router /people/persons [post]
// @Param data body people.PersonInput true "Created person"
func Create(ctx *gin.Context, db *edgedb.Client) {
	controllers.CreateItem[people.PersonInput, people.Person](ctx, db)
}

// @Summary Delete person
// @id DeletePerson
// @tags People
// @Success 200 {object} people.Person
// @Failure 404 "person does not exist"
// @Router /people/persons/{id} [delete]
// @Param id path string true "Item UUID"
func Delete(ctx *gin.Context, db *edgedb.Client) {
	controllers.DeleteByID(ctx, db, people.DeletePerson)
}

// @Summary Update person
// @id UpdatePerson
// @tags People
// @Success 200 {object} people.Person
// @Failure 400 {object} validations.FieldErrors
// @Router /people/persons/{id} [patch]
// @Param id path string true "Item UUID"
// @Param data body people.PersonUpdate true "Update infos"
func Update(ctx *gin.Context, db *edgedb.Client) {
	// uuid, err := controllers.ParseUUIDfromURI(ctx)
	// if err != nil {
	// 	logrus.Errorf("%v", err)
	// 	return
	// }
	// person := people.PersonUpdate{ID: uuid}
	// err = ctx.ShouldBindJSON(&person)
	// if err != nil {
	// 	logrus.Errorf("%v", err)
	// 	return
	// }
	// logrus.Debugf("%+v", person)
	// updated, err := person.Update(db)
	// if err != nil {
	// 	logrus.Errorf("%v", err)
	// 	return
	// }
	// logrus.Infof("%v", updated)
	// ctx.JSON(http.StatusOK, updated)

	// controllers.UpdateByID(ctx, db, people.PersonInitUpdate)
	controllers.UpdateItemByUUID[people.PersonUpdate](ctx, db, people.FindPerson)
}
