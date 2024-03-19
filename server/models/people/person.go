package people

import (
	"context"
	"darco/proto/db"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type PersonIdentity struct {
	FirstName string `json:"first_name" edgedb:"first_name" binding:"required,alphaunicode,min=2,max=32" faker:"first_name"`
	LastName  string `json:"last_name" edgedb:"last_name" binding:"required,alphaunicode,min=2,max=32" faker:"last_name"`
}

type PersonInner struct {
	PersonIdentity `edgedb:"$inline"`
	ID             edgedb.UUID        `edgedb:"id" json:"id" binding:"required"`
	FullName       string             `json:"full_name" edgedb:"full_name" binding:"required"`
	Alias          string             `json:"alias" edgedb:"alias" binding:"required"`
	Role           OptionalUserRole   `json:"role" edgedb:"role"`
	Contact        edgedb.OptionalStr `json:"contact" edgedb:"contact"`
	Comment        edgedb.OptionalStr `json:"comment" edgedb:"comment"`
}

type Person struct {
	PersonInner  `edgedb:"$inline" json:",inline"`
	Institutions []InstitutionInner `json:"institutions" edgedb:"institutions"`
	Meta         Meta               `json:"meta" edgedb:"meta"`
} // @name Person

type OptionalPerson struct {
	edgedb.Optional
	PersonInner `edgedb:"$inline" json:",inline"`
} // @name OptionalPerson

// func PersonStructLevelValidation(sl validator.StructLevel) {
// 	person := sl.Current().Interface().(PersonInput)
// 	var exists = false
// 	query := `
// 		select exists (
// 			select people::Person
// 			filter .first_name = <str>$0 and .last_name = <str>$1
// 		);`

// 	err := db.Client().QuerySingle(context.Background(), query, &exists, person.FirstName, person.LastName)
// 	if err != nil {
// 		logrus.Errorf("Unique validation query failed: %v with query %s", err, query)
// 	}
// 	if exists {
// 		sl.ReportError(fmt.Sprintf("%s %s", person.FirstName, person.LastName), "*", "Person", "person_unique", "")
// 	}
// }

func FindPerson(db edgedb.Executor, id edgedb.UUID) (person Person, err error) {
	query := `select people::Person { *, institutions: { * }, meta: { * } } filter .id = <uuid>$0;`
	err = db.QuerySingle(context.Background(), query, &person, id)
	return person, err
}

func ListPersons(db edgedb.Executor) (people []Person, err error) {
	err = db.Query(context.Background(),
		`select people::Person { *, institutions: { * }, meta: { * } } order by .last_name;`,
		&people)
	return
}

func DeletePerson(db edgedb.Executor, id edgedb.UUID) (deleted Person, err error) {
	logrus.Infof("Deleting person: %v", id)
	query := `select(
		delete (<people::Person><uuid>$0)
	){ *, institutions: { * }, meta:{ * }};`
	err = db.QuerySingle(context.Background(), query, &deleted, id)
	return
}

func (person Person) Delete(db edgedb.Executor) (Person, error) {
	return DeletePerson(db, person.ID)
}

type PersonInput struct {
	PersonIdentity
	Institutions []string `json:"institutions" binding:"omitempty,exist_all=people::Institution.code" faker:"len=10,slice_len=2"`
	Alias        *string  `json:"alias,omitempty" binding:"unique_str=people::Person.alias" faker:"-"`
	Contact      *string  `json:"contact,omitempty" binding:"omitnil,nullemail" faker:"email"`
	Comment      *string  `json:"comment,omitempty" faker:"sentence"`
} // @name PersonInput

func (p *PersonIdentity) GenerateAlias() *string {
	first_initial := ""
	if len(p.FirstName) > 0 {
		first_initial = string(p.FirstName[0])
	}

	alias := strings.ToLower(fmt.Sprintf("%s%s", first_initial, p.LastName))

	var conflicts int64
	query := `select (count (people::Person
			filter str_trim(.alias, "0123456789") = <str>$0
		))`
	if err := db.Client().QuerySingle(context.Background(),
		query, &conflicts, alias,
	); err != nil {
		logrus.Errorf("Error while checking for Person.alias duplicates: %v", err)
		return nil
	}
	if conflicts > 0 {
		alias = alias + fmt.Sprint(conflicts)
	}
	return &alias
}

// func (p *PersonInput) UnmarshalJSON(data []byte) error {
// 	type TmpInput PersonInput
// 	if err := json.Unmarshal(data, (*TmpInput)(p)); err != nil {
// 		return err
// 	}
// 	if p.Alias == nil {
// 		p.Alias = p.generateAlias()
// 		logrus.Infof("Generated alias '%s' for person %+v", *p.Alias, *p)
// 	}
// 	return nil
// }

//go:embed queries/create_person.edgeql
var personCreateQuery string

func (person PersonInput) Create(db edgedb.Executor) (created Person, err error) {
	logrus.Infof("Creating person %+v", person)
	if person.Alias == nil || *person.Alias == "" {
		person.Alias = person.GenerateAlias()
	}
	args, _ := json.Marshal(person)
	err = db.QuerySingle(context.Background(), personCreateQuery, &created, args)
	return created, err
}

type PersonUpdate struct {
	FirstName    *string   `json:"first_name,omitempty" binding:"omitnil,min=2,alphaunicode,max=32" faker:"first_name"`
	LastName     *string   `json:"last_name,omitempty" binding:"omitnil,min=2,alphaunicode,max=32" faker:"last_name"`
	Contact      *string   `json:"contact,omitempty" binding:"omitnil,nullemail" faker:"email"`
	Institutions *[]string `json:"institutions,omitempty" binding:"omitnil,exist_all=people::Institution.code" faker:"len=10,slice_len=3"` // Institution codes
	Alias        *string   `json:"alias,omitempty" binding:"omitnil,min=3" faker:"word,len=5"`
	Comment      *string   `json:"comment,omitempty" binding:"omitnil" faker:"sentence"`
} // @name PersonUpdate

//go:embed queries/update_person.edgeql
var personUpdateQuery string

func (person PersonUpdate) Update(db edgedb.Executor, id edgedb.UUID) (uuid edgedb.UUID, err error) {
	logrus.Infof("Updating person %+v", person)
	args, _ := json.Marshal(person)
	err = db.Execute(context.Background(), personUpdateQuery, id, args)
	return id, err
}
