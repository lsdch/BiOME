package person

import (
	"context"
	"darco/proto/models"
	"encoding/json"
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

// Only used for OpenAPI specification
type PersonUpdate struct {
	ID        edgedb.UUID        `json:"id"`
	FirstName string             `json:"first_name" binding:"omitempty,alpha,min=2,max=32"`
	LastName  string             `json:"last_name" binding:"omitempty,alpha,min=2,max=32"`
	Contact   edgedb.OptionalStr `json:"contact"`
} // @name PersonUpdate

type Person struct {
	ID          edgedb.UUID `edgedb:"id" json:"id"`
	PersonInput `edgedb:"$inline"`
	FullName    string      `json:"full_name" edgedb:"full_name" binding:"required"`
	Meta        models.Meta `json:"meta" edgedb:"meta" binding:"required"`
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

func Find(db *edgedb.Client, id edgedb.UUID) (person Person, err error) {
	query := `select people::Person { *, meta: { * } } filter .id = <uuid>$0;`
	err = db.QuerySingle(context.Background(), query, &person, id)
	return person, err
}

func List(db *edgedb.Client) (people []Person, err error) {
	query := `select people::Person {
			id, first_name, last_name, full_name, contact, meta: { * }
		} order by .last_name;`
	err = db.Query(context.Background(), query, &people)
	return
}

func (person PersonInput) Create(db *edgedb.Client) (created Person, err error) {
	query := `with data := <json>$0
		select (insert people::Person {
			first_name := <str>data['first_name'],
			last_name := <str>data['last_name']
	}) { *, meta: { * }}`
	args, _ := json.Marshal(person)
	err = db.QuerySingle(context.Background(), query, &created, args)
	return created, err
}

func (person Person) Update(db *edgedb.Client) (updated Person, err error) {
	query := `with data := <json>$0
		select( update people::Person
			filter .id = <uuid>data['id']
			set {
				first_name := <str>data['first_name'],
				last_name := <str>data['last_name'],
			}
		) { *, meta: { * }};`
	args, _ := json.Marshal(person)
	err = db.QuerySingle(context.Background(), query, &updated, args)
	return updated, err
}

func Delete(db *edgedb.Client, id edgedb.UUID) (deleted Person, err error) {
	query := `select(delete people::Person filter .id = <uuid>$0){ *, meta:{ * }};`
	err = db.QuerySingle(context.Background(), query, &deleted, id)
	return
}
