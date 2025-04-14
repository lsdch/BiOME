package people

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"

	"github.com/sirupsen/logrus"
)

type PersonIdentity struct {
	FirstName string `json:"first_name" gel:"first_name" minLength:"2" maxLength:"32" fake:"{firstname}"`
	LastName  string `json:"last_name" gel:"last_name" minLength:"2" maxLength:"32" fake:"{lastname}"`
}

// PersonInner contains all properties defining a person, excluding links to related entities
type PersonInner struct {
	PersonIdentity `gel:"$inline"`
	ID             geltypes.UUID        `gel:"id" json:"id" binding:"required" format:"uuid"`
	FullName       string               `json:"full_name" gel:"full_name" binding:"required"`
	Alias          string               `json:"alias" gel:"alias" binding:"required"`
	Role           OptionalUserRole     `json:"role,omitempty" gel:"role"`
	Contact        geltypes.OptionalStr `json:"contact" gel:"contact" format:"email"`
	Comment        geltypes.OptionalStr `json:"comment" gel:"comment"`
}

// PersonUser is PersonInner with optional user informations attached
type PersonUser struct {
	PersonInner `gel:"$inline" json:",inline"`
	User        models.Optional[UserInner] `gel:"user" json:"user"`
}

// Person is the complete informations about a person, including related entities
type Person struct {
	PersonUser    `gel:"$inline" json:",inline"`
	Organisations []OrganisationInner `json:"organisations" gel:"organisations"`
	Meta          Meta                `json:"meta" gel:"meta"`
}

type OptionalPerson struct {
	geltypes.Optional
	PersonInner `gel:"$inline" json:",inline"`
}

func FindPerson(db geltypes.Executor, id geltypes.UUID) (person Person, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		select people::Person { *, ** } filter .id = <uuid>$0;
		`, &person, id)
	return person, err
}

func ListPersons(db geltypes.Executor) (people []Person, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select people::Person { ** } order by .last_name;
		`, &people)
	return
}

func DeletePerson(db geltypes.Executor, id geltypes.UUID) (deleted Person, err error) {
	logrus.Infof("Deleting person: %v", id)
	query := `#edgeql
		select(
			delete (<people::Person><uuid>$0)
		){ ** };`
	err = db.QuerySingle(context.Background(), query, &deleted, id)
	return
}

func (person Person) Delete(db geltypes.Executor) (Person, error) {
	return DeletePerson(db, person.ID)
}

type PersonInput struct {
	PersonIdentity
	Organisations []string                     `json:"organisations,omitempty" fakesize:"2"`
	Alias         models.OptionalInput[string] `json:"alias,omitempty" fake:"-"`
	Contact       models.OptionalInput[string] `json:"contact,omitempty" format:"email"`
	Comment       models.OptionalInput[string] `json:"comment,omitempty"`
}

func (p *PersonInput) WithOrganisationCodes(codes map[string]string) PersonInput {
	for i, code := range p.Organisations {
		if org, ok := codes[code]; ok {
			p.Organisations[i] = org
		}
	}
	return *p
}

func (p *PersonIdentity) GenerateAlias() string {
	first_initial := ""
	if len(p.FirstName) > 0 {
		first_initial = string(p.FirstName[0])
	}

	alias := strings.ToLower(fmt.Sprintf("%s%s", first_initial, p.LastName))

	var conflicts int64
	if err := db.Client().QuerySingle(context.Background(),
		`#edgeql
			select (count (people::Person
				filter str_trim(.alias, "0123456789") = <str>$0
			))
		`, &conflicts, alias,
	); err != nil {
		logrus.Errorf("Error while checking for Person.alias duplicates: %v", err)
		return ""
	}
	if conflicts > 0 {
		alias = alias + fmt.Sprint(conflicts)
	}
	return alias
}

func (person PersonInput) Save(db geltypes.Executor) (created Person, err error) {
	logrus.Infof("Creating person %+v", person)
	if !person.Alias.IsSet {
		person.Alias.Value = person.GenerateAlias()
	}
	args, _ := json.Marshal(person)
	logrus.Infof("Creating person with args: %s", string(args))
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select people::insert_person(<json>$0) { ** }
		`, &created, args)
	return created, err
}

type PersonUpdate struct {
	FirstName     models.OptionalInput[string]   `gel:"first_name" json:"first_name,omitempty" minLength:"2" maxLength:"32"`
	LastName      models.OptionalInput[string]   `gel:"last_name" json:"last_name,omitempty" minLength:"2" maxLength:"32"`
	Contact       models.OptionalNull[string]    `gel:"contact" json:"contact,omitempty" `
	Organisations models.OptionalInput[[]string] `gel:"organisations" json:"organisations,omitempty" fakesize:"3"` // Organisation codes
	Alias         models.OptionalInput[string]   `gel:"alias" json:"alias,omitempty"`
	Comment       models.OptionalNull[string]    `gel:"comment" json:"comment,omitempty"`
}

func (u PersonUpdate) Save(e geltypes.Executor, id geltypes.UUID) (updated Person, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update people::Person filter .id = <uuid>$0 set {
				%s
			}) { ** }
		`,
		Mappings: map[string]string{
			"first_name": "<str>item['first_name']",
			"last_name":  "<str>item['last_name']",
			"contact":    "<str>item['contact']",
			"alias":      "<str>item['alias']",
			"comment":    "<str>item['comment']",
			"organisations": `#edgeql
				(
					select people::Organisation
					filter .code in array_unpack(<array<str>>item['organisations'])
				)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, id, data)
	updated.Meta.Save(e)
	return
}
