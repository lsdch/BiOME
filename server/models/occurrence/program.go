package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type ProgramInner struct {
	ID          edgedb.UUID          `edgedb:"id" json:"id" format:"uuid"`
	Label       string               `edgedb:"label" json:"label"`
	Code        string               `edgedb:"code" json:"code"`
	StartYear   edgedb.OptionalInt32 `edgedb:"start_year" json:"start_year,omitempty" minimum:"1900" example:"2019"`
	EndYear     edgedb.OptionalInt32 `edgedb:"end_year" json:"end_year,omitempty" example:"2025"`
	Description edgedb.OptionalStr   `edgedb:"description" json:"description,omitempty"`
}

type Program struct {
	ProgramInner    `edgedb:"$inline" json:",inline"`
	Managers        []people.PersonInner      `edgedb:"managers" json:"managers" minItems:"1"`
	FundingAgencies []people.InstitutionInner `edgedb:"funding_agencies" json:"funding_agencies"`
	Meta            people.Meta               `edgedb:"meta" json:"meta"`
}

func ListPrograms(db edgedb.Executor) ([]Program, error) {
	var programs = []Program{}
	err := db.Query(context.Background(), `select events::Program { ** }`, &programs)
	return programs, err
}

func FindProgram(db edgedb.Executor, code string) (program Program, err error) {
	err = db.QuerySingle(context.Background(),
		`select events::Program { ** } filter .code = <str>$0`,
		&program, code)
	return
}

type ProgramInput struct {
	Label           string                         `json:"label" example:"Alice's PhD"`
	Code            string                         `json:"code" example:"PHD_ALICE"`
	Managers        []string                       `json:"managers" minItems:"1" example:"[\"adoe\",\"fmalard\"]"`
	FundingAgencies models.OptionalInput[[]string] `json:"funding_agencies,omitempty" example:"[\"CNRS\"]"`
	StartYear       models.OptionalInput[int32]    `json:"start_year,omitempty" minimum:"1900" example:"2022"`
	EndYear         models.OptionalInput[int32]    `json:"end_year,omitempty" example:"2025"`
	Description     models.OptionalInput[string]   `json:"description,omitempty"`
}

func (i ProgramInput) Create(db edgedb.Executor) (created Program, err error) {
	args, _ := json.Marshal(i)
	err = db.QuerySingle(context.Background(),
		`#edgeql
			data := <json>$0,
			managers := (
				select people::Person
				filter .alias in array_unpack(<array<str>>json_get(data, 'managers'))
			),
			institutions := (
				select people::Institution
				filter .code in array_unpack(<array<str>>json_get(data, 'funding_agencies'))
			),
		select(insert events::Program {
			label := <str>data['label'],
			code := <str>data['code'],
			managers := managers,
			funding_agencies := institutions,
			start_year := <int32>json_get(data, 'start_year'),
			end_year := <int32>json_get(data, 'end_year'),
			description := <str>json_get(data, 'description')
		}) { ** };`,
		&created, args)
	return created, err
}

type ProgramUpdate struct {
	Label           models.OptionalInput[string]  `json:"label,omitempty" example:"Alice's PhD"`
	Code            models.OptionalInput[string]  `json:"code,omitempty" example:"PHD_ALICE"`
	Managers        models.OptionalNull[[]string] `json:"managers,omitempty" minItems:"1" example:"[\"adoe\",\"fmalard\"]"`
	FundingAgencies models.OptionalNull[[]string] `json:"funding_agencies,omitempty" example:"[\"CNRS\"]"`
	StartYear       models.OptionalNull[int32]    `json:"start_year,omitempty" minimum:"1900" example:"2022"`
	EndYear         models.OptionalNull[int32]    `json:"end_year,omitempty" example:"2025"`
	Description     models.OptionalNull[string]   `json:"description,omitempty"`
}

func (u ProgramUpdate) Save(e edgedb.Executor, code string) (updated Program, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `with item := <json>$1,
		select (update events::Program filter .code = <str>$0 set {
			%s
		}) { ** }`,
		Mappings: map[string]string{
			"label":       "<str>item['label']",
			"code":        "<str>item['code']",
			"start_year":  "<int32>item['start_year']",
			"end_year":    "<int32>item['end_year']",
			"description": "<str>item['description']",
			"managers": `(
				select people::Person
				filter .alias in array_unpack(<array<str>>json_get(item, 'managers'))
			)`,
			"funding_agencies": `(
				select people::Institution
				filter .code in array_unpack(<array<str>>json_get(data, 'funding_agencies'))
			)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	updated.Meta.Save(e)
	return
}

func DeleteProgram(db edgedb.Executor, code string) (deleted Program, err error) {
	err = db.QuerySingle(context.Background(),
		`select (
			 delete events::Program filter .code = <str>$0
		 ) { ** };`,
		&deleted, code)
	return
}
