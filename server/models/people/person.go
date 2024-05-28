package people

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type PersonIdentity struct {
	FirstName string `json:"first_name" edgedb:"first_name" minLength:"2" maxLength:"32" fake:"{firstname}"`
	LastName  string `json:"last_name" edgedb:"last_name" minLength:"2" maxLength:"32" fake:"{lastname}"`
}

// PersonInner contains all properties defining a person, excluding links to related entities
type PersonInner struct {
	PersonIdentity `edgedb:"$inline"`
	ID             edgedb.UUID        `edgedb:"id" json:"id" binding:"required" format:"uuid"`
	FullName       string             `json:"full_name" edgedb:"full_name" binding:"required"`
	Alias          string             `json:"alias" edgedb:"alias" binding:"required"`
	Role           OptionalUserRole   `json:"role,omitempty" edgedb:"role"`
	Contact        edgedb.OptionalStr `json:"contact" edgedb:"contact" format:"email"`
	Comment        edgedb.OptionalStr `json:"comment" edgedb:"comment"`
}

// PersonUser is PersonInner with optional user informations attached
type PersonUser struct {
	PersonInner `edgedb:"$inline" json:",inline"`
	User        OptionalUserInner `edgedb:"user" json:"user,omitempty"`
}

// Person is the complete informations about a person, including related entities
type Person struct {
	PersonUser   `edgedb:"$inline" json:",inline"`
	Institutions []InstitutionInner `json:"institutions" edgedb:"institutions"`
	Meta         Meta               `json:"meta" edgedb:"meta"`
} // @name Person

type OptionalPerson struct {
	edgedb.Optional
	PersonInner `edgedb:"$inline" json:",inline"`
}

func FindPerson(db edgedb.Executor, id edgedb.UUID) (person Person, err error) {
	query := `select people::Person { *, ** } filter .id = <uuid>$0;`
	err = db.QuerySingle(context.Background(), query, &person, id)
	return person, err
}

func ListPersons(db edgedb.Executor) (people []Person, err error) {
	err = db.Query(context.Background(),
		`select people::Person { ** } order by .last_name;`,
		&people)
	return
}

func DeletePerson(db edgedb.Executor, id edgedb.UUID) (deleted Person, err error) {
	logrus.Infof("Deleting person: %v", id)
	query := `select(
		delete (<people::Person><uuid>$0)
	){ ** };`
	err = db.QuerySingle(context.Background(), query, &deleted, id)
	return
}

func (person Person) Delete(db edgedb.Executor) (Person, error) {
	return DeletePerson(db, person.ID)
}

type PersonInput struct {
	PersonIdentity
	Institutions []string                     `json:"institutions" binding:"omitempty,exist_all=people::Institution.code" fakesize:"2"`
	Alias        models.OptionalInput[string] `json:"alias,omitempty" binding:"unique_str=people::Person.alias" fake:"-"`
	Contact      models.OptionalInput[string] `json:"contact,omitempty" format:"email"`
	Comment      models.OptionalInput[string] `json:"comment,omitempty"`
}

func (p *PersonIdentity) GenerateAlias() string {
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
		return ""
	}
	if conflicts > 0 {
		alias = alias + fmt.Sprint(conflicts)
	}
	return alias
}

//go:embed queries/create_person.edgeql
var personCreateQuery string

func (person PersonInput) Create(db edgedb.Executor) (created Person, err error) {
	logrus.Infof("Creating person %+v", person)
	if !person.Alias.IsSet {
		person.Alias.Value = person.GenerateAlias()
	}
	args, _ := json.Marshal(person)
	err = db.QuerySingle(context.Background(), personCreateQuery, &created, args)
	return created, err
}

type PersonUpdate struct {
	FirstName    models.OptionalInput[string]   `json:"first_name,omitempty" minLength:"2" maxLength:"32"`
	LastName     models.OptionalInput[string]   `json:"last_name,omitempty" minLength:"2" maxLength:"32"`
	Contact      models.OptionalNull[string]    `json:"contact,omitempty" binding:"omitnil,nullemail" `
	Institutions models.OptionalInput[[]string] `json:"institutions,omitempty" binding:"omitnil,exist_all=people::Institution.code" fakesize:"3"` // Institution codes
	Alias        models.OptionalInput[string]   `json:"alias,omitempty" binding:"omitnil,min=3"`
	Comment      models.OptionalNull[string]    `json:"comment,omitempty" binding:"omitnil"`
}

//go:embed queries/update_person.edgeql
var personUpdateQuery string

func (person PersonUpdate) Update(db edgedb.Executor, id edgedb.UUID) (uuid edgedb.UUID, err error) {
	logrus.Infof("Updating person %+v", person)
	args, _ := json.Marshal(person)
	err = db.Execute(context.Background(), personUpdateQuery, id, args)
	return id, err
}
