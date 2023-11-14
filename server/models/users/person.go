package users

import (
	"context"
	"darco/proto/models"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type PersonInput struct {
	FirstName string             `json:"first_name" edgedb:"first_name" binding:"required,alpha,min=2,max=32"`
	LastName  string             `json:"last_name" edgedb:"last_name" binding:"required,alpha,min=2,max=32"`
	Contact   edgedb.OptionalStr `json:"contact" edgedb:"contact"`
} // @name PersonInput

type Person struct {
	ID          edgedb.UUID `edgedb:"id" json:"id"`
	PersonInput `edgedb:"$inline"`
	FullName    string `json:"full_name" edgedb:"full_name"`
} // @name Person

func PersonStructLevelValidation(sl validator.StructLevel) {
	person := sl.Current().Interface().(PersonInput)
	var exists = false
	query := "select exists (select people::Person filter .first_name = <str>$0 and .last_name = <str>$1)"
	err := models.DB().QuerySingle(context.Background(), query, &exists, person.FirstName, person.LastName)
	if err != nil {
		logrus.Errorf("Unique validation query failed: %v with query %s", err, query)
	}
	if exists {
		sl.ReportError(fmt.Sprintf("%s %s", person.FirstName, person.LastName), "*", "Person", "person_unique", "")
	}
}
