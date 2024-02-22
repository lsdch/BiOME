package people

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/users/user_role"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type PersonInner struct {
	FirstName   string             `json:"first_name" edgedb:"first_name" binding:"required,alphaunicode,min=2,max=32"`
	MiddleNames edgedb.OptionalStr `json:"middle_names" edgedb:"middle_names" binding:"omitempty,alphaunicode,max=32"`
	LastName    string             `json:"last_name" edgedb:"last_name" binding:"required,alphaunicode,min=2,max=32"`
}

type PersonInput struct {
	PersonInner
	Institutions []string `json:"institutions" binding:"omitempty,exist_all=people::Institution.code"`
	Alias        *string  `json:"alias,omitempty" binding:"unique_str=people::Person.alias"`
	Contact      *string  `json:"contact,omitempty" binding:"omitnil,nullemail"`
} // @name PersonInput

func (p *PersonInput) generateAlias() string {
	middle_initials := ""
	middle_names, isSet := p.MiddleNames.Get()
	if isSet && len(middle_names) > 0 {
		re := regexp.MustCompile(`\W`)
		split := re.Split(middle_names, -1)
		for _, word := range split {
			middle_initials = middle_initials + string(word[0])
		}
	}
	first_initial := ""
	if len(p.FirstName) > 0 {
		first_initial = string(p.FirstName[0])
	}

	alias := strings.ToLower(fmt.Sprintf("%s%s%s", first_initial, middle_initials, p.LastName))

	var conflicts []string
	query := `select (people::Person filter str_trim(.alias, "0123456789") = <str>$0).alias`
	if err := db.Client().Query(context.Background(), query, &conflicts, alias); err != nil {
		logrus.Errorf("Error while checking for Person.alias duplicates: %v", err)
		return ""
	}
	if len(conflicts) > 0 {
		alias = alias + fmt.Sprint(len(conflicts))
	}
	return alias
}

func (p *PersonInput) UnmarshalJSON(data []byte) error {
	type TmpInput PersonInput
	if err := json.Unmarshal(data, (*TmpInput)(p)); err != nil {
		return err
	}
	if p.Alias == nil {
		alias := p.generateAlias()
		logrus.Infof("Generated alias %s for person %+v", alias, *p)
		p.Alias = &alias
	}
	return nil
}

type Person struct {
	PersonInner  `edgedb:"$inline"`
	Institutions []Institution              `json:"institutions" edgedb:"institutions"`
	ID           edgedb.UUID                `edgedb:"id" json:"id" binding:"required"`
	FullName     string                     `json:"full_name" edgedb:"full_name" binding:"required"`
	Alias        string                     `json:"alias" edgedb:"alias" binding:"required"`
	Role         user_role.OptionalUserRole `json:"role" edgedb:"role"`
	Contact      edgedb.OptionalStr         `json:"contact" edgedb:"contact"`
	Meta         models.Meta                `json:"meta" edgedb:"meta"`
} // @name Person

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

func FindPerson(db *edgedb.Client, id edgedb.UUID) (person Person, err error) {
	query := `select people::Person { *, institutions: { * }, meta: { * } } filter .id = <uuid>$0;`
	err = db.QuerySingle(context.Background(), query, &person, id)
	return person, err
}

func ListPersons(db *edgedb.Client) (people []Person, err error) {
	query := `select people::Person { *, institutions: { * }, meta: { * } } order by .last_name;`
	err = db.Query(context.Background(), query, &people)
	return
}

func DeletePerson(db *edgedb.Client, id edgedb.UUID) (deleted Person, err error) {
	logrus.Infof("Deleting person: %v", id)
	query := `select(
		delete people::Person filter .id = <uuid>$0
	){ *, institutions: { * }, meta:{ * }};`
	err = db.QuerySingle(context.Background(), query, &deleted, id)
	return
}

//go:embed queries/create_person.edgeql
var personCreateQuery string

func (person PersonInput) Create(db *edgedb.Client) (created Person, err error) {
	logrus.Infof("Creating person %+v", person)
	args, _ := json.Marshal(person)
	err = db.QuerySingle(context.Background(), personCreateQuery, &created, args)
	return created, err
}

type PersonUpdate struct {
	FirstName    *string   `json:"first_name,omitempty" binding:"omitnil,min=2,alphaunicode,max=32"`
	MiddleNames  *string   `json:"middle_names,omitempty" edgedb:"middle_names" binding:"omitnil,nullalphaunicode,max=32"`
	LastName     *string   `json:"last_name,omitempty" binding:"omitnil,min=2,alphaunicode,max=32"`
	Contact      *string   `json:"contact,omitempty" binding:"omitnil,nullemail"`
	Institutions *[]string `json:"institutions,omitempty" binding:"omitnil,exist_all=people::Institution.code"` // Institution codes
	Alias        *string   `json:"alias,omitempty" binding:"omitnil,min=3"`
} // @name PersonUpdate

//go:embed queries/update_person.edgeql
var personUpdateQuery string

func (person PersonUpdate) Update(db *edgedb.Client, id edgedb.UUID) (uuid edgedb.UUID, err error) {
	logrus.Infof("Updating person %+v", person)
	args, _ := json.Marshal(person)
	err = db.Execute(context.Background(), personUpdateQuery, id, args)
	return id, err
}
