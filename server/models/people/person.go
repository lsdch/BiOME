package people

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"

	"github.com/sirupsen/logrus"
)

// @mapstructure
type PersonIdentity struct {
	FirstName string `json:"first_name" gel:"first_name" minLength:"2" maxLength:"32" fake:"{firstname}" mapstructure:"first_name"`
	LastName  string `json:"last_name" gel:"last_name" minLength:"2" maxLength:"32" fake:"{lastname}" mapstructure:"last_name"`
}

// PersonInner contains all properties defining a person, excluding links to related entities
type PersonInner struct {
	PersonIdentity `gel:"$inline"`
	ID             geltypes.UUID        `gel:"id" json:"id" binding:"required" format:"uuid"`
	FullName       string               `json:"full_name" gel:"full_name" binding:"required"`
	Alias          string               `json:"alias" gel:"alias" binding:"required"`
	Role           OptionalUserRole     `json:"role,omitempty" gel:"role"`
	Contact        geltypes.OptionalStr `json:"contact" gel:"contact" format:"email"`
	Organisation   geltypes.OptionalStr `json:"organisation,omitempty" gel:"organisation"`
	Comment        geltypes.OptionalStr `json:"comment" gel:"comment"`
}

// PersonUser is PersonInner with optional user informations attached
type PersonUser struct {
	PersonInner `gel:"$inline"`
	User        models.Optional[UserInner] `gel:"user" json:"user"`
}

// Person is the complete informations about a person, including related entities
type Person struct {
	PersonUser `gel:"$inline"`
	Meta       Meta `json:"meta" gel:"meta"`
}

type OptionalPerson struct {
	geltypes.Optional
	PersonInner `gel:"$inline"`
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
			select people::Person { ** } order by (exists .user) desc then .last_name asc;
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
	Contact     models.OptionalInput[string] `json:"contact,omitzero" format:"email"`
	Comment     models.OptionalInput[string] `json:"comment,omitzero"`
	ForceCreate bool                         `json:"force_create,omitzero"`
}

func (person PersonInput) Save(db geltypes.Executor) (created Person, err error) {
	args, _ := json.Marshal(person)
	logrus.Debugf("Creating person with args: %s", string(args))
	err = db.QuerySingle(context.Background(),
		`#edgeql
			with args := <json>$0,
			select (
				insert Person {
					first_name := <str>args['first_name'],
					last_name := <str>args['last_name'],
					contact := <str>json_get(args, 'contact'),
					comment := <str>json_get(args, 'comment'),
					organisation := <str>json_get(args, 'organisation')
				}
			) { ** }
		`, &created, args)
	return created, err
}

type PersonUpdate struct {
	FirstName    models.OptionalInput[string] `gel:"first_name" json:"first_name,omitempty" minLength:"2" maxLength:"32"`
	LastName     models.OptionalInput[string] `gel:"last_name" json:"last_name,omitempty" minLength:"2" maxLength:"32"`
	Contact      models.OptionalNull[string]  `gel:"contact" json:"contact,omitempty" `
	Organisation models.OptionalNull[string]  `gel:"organisation" json:"organisation,omitempty" fakesize:"2"` // Organisation code
	Alias        models.OptionalInput[string] `gel:"alias" json:"alias,omitempty"`
	Comment      models.OptionalNull[string]  `gel:"comment" json:"comment,omitempty"`
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
			"first_name":   "<str>item['first_name']",
			"last_name":    "<str>item['last_name']",
			"contact":      "<str>item['contact']",
			"alias":        "<str>item['alias']",
			"comment":      "<str>item['comment']",
			"organisation": "<str>item['organisation']",
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, id, data)
	updated.Meta.Save(e)
	return
}
