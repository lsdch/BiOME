package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/people"
)

type ProgramInner struct {
	ID          geltypes.UUID          `gel:"id" json:"id" format:"uuid"`
	Label       string                 `gel:"label" json:"label"`
	Code        string                 `gel:"code" json:"code"`
	StartYear   geltypes.OptionalInt32 `gel:"start_year" json:"start_year,omitempty" minimum:"1900" example:"2019"`
	EndYear     geltypes.OptionalInt32 `gel:"end_year" json:"end_year,omitempty" example:"2025"`
	Description geltypes.OptionalStr   `gel:"description" json:"description,omitempty"`
}

type Program struct {
	ProgramInner    `gel:"$inline" json:",inline"`
	Managers        []people.PersonInner       `gel:"managers" json:"managers" minItems:"1"`
	FundingAgencies []people.OrganisationInner `gel:"funding_agencies" json:"funding_agencies"`
	Datasets        []dataset.DatasetInner     `gel:"datasets" json:"datasets"`
	Meta            people.Meta                `gel:"meta" json:"meta"`
}

func ListPrograms(db geltypes.Executor) ([]Program, error) {
	var programs = []Program{}
	err := db.Query(context.Background(),
		`#edgeql
			select datasets::ResearchProgram { ** }
		`, &programs)
	return programs, err
}

func FindProgram(db geltypes.Executor, code string) (program Program, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select datasets::ResearchProgram { ** } filter .code = <str>$0
		`,
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
	Datasets        models.OptionalInput[[]string] `json:"datasets,omitempty" example:"[\"dataset1\"]"`
}

func (i ProgramInput) Save(db geltypes.Executor) (created Program, err error) {
	args, _ := json.Marshal(i)
	err = db.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			managers := (
				select people::Person
				filter .alias in array_unpack(<array<str>>json_get(data, 'managers'))
			),
			organisations := (
				select people::Organisation
				filter .code in array_unpack(<array<str>>json_get(data, 'funding_agencies'))
			),
		select(insert events::Program {
			label := <str>data['label'],
			code := <str>data['code'],
			managers := managers,
			funding_agencies := organisations,
			start_year := <int32>json_get(data, 'start_year'),
			end_year := <int32>json_get(data, 'end_year'),
			description := <str>json_get(data, 'description'),
			datasets := (
				select datasets::Dataset
				filter .slug in <str>json_array_unpack(json_get(data, 'datasets'))
			)
		}) { ** };`,
		&created, args)
	return created, err
}

type ProgramUpdate struct {
	Label           models.OptionalInput[string]  `gel:"label" json:"label,omitempty" example:"Alice's PhD"`
	Code            models.OptionalInput[string]  `gel:"code" json:"code,omitempty" example:"PHD_ALICE"`
	Managers        models.OptionalNull[[]string] `gel:"managers"  json:"managers,omitempty" minItems:"1" example:"[\"adoe\",\"fmalard\"]"`
	FundingAgencies models.OptionalNull[[]string] `gel:"funding_agencies" json:"funding_agencies,omitempty" example:"[\"CNRS\"]"`
	StartYear       models.OptionalNull[int32]    `gel:"start_year" json:"start_year,omitempty" minimum:"1900" example:"2022"`
	EndYear         models.OptionalNull[int32]    `gel:"end_year" json:"end_year,omitempty" example:"2025"`
	Description     models.OptionalNull[string]   `gel:"description" json:"description,omitempty"`
	Datasets        models.OptionalNull[[]string] `gel:"datasets" json:"datasets,omitempty" example:"[\"dataset1\"]"`
}

func (u ProgramUpdate) Save(e geltypes.Executor, code string) (updated Program, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update events::Program filter .code = <str>$0 set { %s }) { ** }
		`,
		Mappings: map[string]string{
			"label":       "<str>item['label']",
			"code":        "<str>item['code']",
			"start_year":  "<int32>item['start_year']",
			"end_year":    "<int32>item['end_year']",
			"description": "<str>item['description']",
			"managers": `#edgeql
				(
					select people::Person
					filter .alias in array_unpack(<array<str>>json_get(item, 'managers'))
				)`,
			"funding_agencies": `#edgeql
				(
					select people::Organisation
					filter .code in array_unpack(<array<str>>json_get(item, 'funding_agencies'))
				)`,
			"datasets": `#edgeql
				(
					select datasets::Dataset
					filter .slug in <str>json_array_unpack(json_get(item, 'datasets'))
				)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	updated.Meta.Save(e)
	return
}

func DeleteProgram(db geltypes.Executor, code string) (deleted Program, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		select (
			 delete datasets::ResearchProgram filter .code = <str>$0
		 ) { ** };`,
		&deleted, code)
	return
}
