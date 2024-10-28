package events

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type Program struct {
	ID              edgedb.UUID               `edgedb:"id" json:"id" format:"uuid"`
	Label           string                    `edgedb:"label" json:"label"`
	Code            string                    `edgedb:"code" json:"code"`
	Managers        []people.PersonInner      `edgedb:"managers" json:"managers" minItems:"1"`
	FundingAgencies []people.InstitutionInner `edgedb:"funding_agencies" json:"funding_agencies"`
	StartYear       edgedb.OptionalInt32      `edgedb:"start_year" json:"start_year" minimum:"1900" example:"2019"`
	EndYear         edgedb.OptionalInt32      `edgedb:"end_year" json:"end_year" example:"2025"`
	Description     edgedb.OptionalStr        `edgedb:"description" json:"description"`
	Meta            people.Meta               `edgedb:"meta" json:"meta"`
}

func ListPrograms(db edgedb.Executor) ([]Program, error) {
	var programs = []Program{}
	err := db.Query(context.Background(), `select events::Program { ** }`, &programs)
	return programs, err
}

type ProgramInput struct {
	Label           string                       `edgedb:"label" json:"label" example:"Alice's PhD"`
	Code            string                       `edgedb:"code" json:"code" example:"PHD_ALICE"`
	Managers        []string                     `edgedb:"managers" json:"managers" minItems:"1" example:"[\"adoe\",\"fmalard\"]"`
	FundingAgencies []string                     `edgedb:"funding_agencies" json:"funding_agencies" example:"[\"CNRS\"]"`
	StartYear       models.OptionalInput[int32]  `edgedb:"start_year" json:"start_year" minimum:"1900" example:"2022"`
	EndYear         models.OptionalInput[int32]  `edgedb:"end_year" json:"end_year" example:"2025"`
	Description     models.OptionalInput[string] `edgedb:"description" json:"description"`
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
				and .kind = people::FundingAgency
			),
		select(insert events::Program {
			label := data['label'],
			code := data['code'],
			managers := managers,
			funding_agencies := institutions,
			start_year := <datetime>json_get(data, 'start_year'),
			end_year := <datetime>json_get(data, 'end_year'),
			description := <str>json_get(data, 'description')
		}) { ** };`,
		&created, args)
	return created, err
}
