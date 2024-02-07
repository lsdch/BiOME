package institution

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/person"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type InstitutionInput struct {
	Name        string             `json:"name" edgedb:"name" example:"Mos Eisley Laboratory of Environmental Studies" binding:"required,min=10,max=128"`
	Code        string             `json:"code" edgedb:"code" example:"MELES" binding:"required,alphanum,min=2,max=12"`
	Description edgedb.OptionalStr `json:"description" edgedb:"description" example:"The main ecological research lab on Tatooine."`
} // @name InstitutionInput

type InstitutionUpdate struct {
	Name        string `json:"name" example:"Mos Eisley Laboratory of Environmental Studies" binding:"omitempty,min=10,max=128"`
	Code        string `json:"code" example:"MELES" binding:"omitempty,alphanum,min=2,max=12"`
	Description string `json:"description" example:"The main ecological research lab on Tatooine." binding:"omitempty"`
} //@name InstitutionUpdate

type Institution struct {
	ID               edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>" binding:"required"`
	InstitutionInput `edgedb:"$inline"`
	People           []person.Person `json:"people" edgedb:"people"`
	Meta             models.Meta     `json:"meta" edgedb:"meta" binding:"required"`
} // @name Institution

func Find(db *edgedb.Client, code string) (inst Institution, err error) {
	logrus.Debugf("Institution look up code: %s", code)
	query := "select people::Institution { *, meta:{ * } } filter .code = <str>$0;"
	err = db.QuerySingle(context.Background(), query, &inst, code)
	return
}

func List(db *edgedb.Client) (institutions []Institution, err error) {
	query := `select people::Institution { *, people:{ * }, meta:{ * } } order by .code;`
	err = db.Query(context.Background(), query, &institutions)
	return
}

func (inst InstitutionInput) Create(db *edgedb.Client) (created Institution, err error) {
	query := `
		with module people,
			data := <json>$0,
		select ( insert Institution {
			name := <str>data['name'],
			code := <str>data['code'],
			description := <str>data['description']
		}) { *, people:{ * }, meta:{ * } };`

	args, _ := json.Marshal(inst)
	err = db.QuerySingle(context.Background(), query, &created, args)
	return
}

func (inst Institution) Update(db *edgedb.Client) (updated Institution, err error) {
	query := `
	with module people,
		data := <json>$0,
	select(
		update Institution filter .id = <uuid>data['id']
		set {
			name := <str>data['name'],
			code := <str>data['code'],
			description := <str>data['description']
		}
	) { *, people:{ * }, meta:{ * } };`

	args, _ := json.Marshal(inst)
	logrus.Debugf("Updating institution : %+v", inst)
	err = db.QuerySingle(context.Background(), query, &updated, args)
	return
}

func Delete(db *edgedb.Client, code string) (inst Institution, err error) {
	query := `select(
		delete people::Institution filter .code = <str>$0
	) { *, people:{ * }, meta: { * }};`
	err = db.QuerySingle(context.Background(), query, &inst, code)
	return
}
