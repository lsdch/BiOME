package people

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type OrganisationInfos struct {
	Name string  `json:"name" edgedb:"name" example:"Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés" minLength:"10" maxLength:"128" fake:"{company}"`
	Code string  `json:"code" edgedb:"code" example:"LEHNA" minLength:"2" maxLength:"12" fake:"{word}"`
	Kind OrgKind `json:"kind" edgedb:"kind" example:"Lab"`
}

type OrganisationInput struct {
	OrganisationInfos `edgedb:"$inline"`
	Description       models.OptionalInput[string] `json:"description,omitempty" edgedb:"description" example:"Where this database was born."`
}

type OrganisationInner struct {
	ID                edgedb.UUID `json:"id" edgedb:"id" format:"uuid" binding:"required"`
	OrganisationInfos `edgedb:"$inline" json:",inline"`
	Description       edgedb.OptionalStr `json:"description,omitempty" edgedb:"description" example:"Where this database was born."`
}

type Organisation struct {
	OrganisationInner `edgedb:"$inline" json:",inline"`
	People            []PersonUser `json:"people,omitempty" edgedb:"people" doc:"Known members of this organisation"`
	Meta              Meta         `json:"meta" edgedb:"meta"`
}

func FindOrganisation(db edgedb.Executor, uuid edgedb.UUID) (org Organisation, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select people::Organisation { *, people:{ *, user: { * } }, meta:{ * } }
			filter .id = <uuid>$0;
		`, &org, uuid)
	return org, err
}

func ListOrganisations(db edgedb.Executor) (organisations []Organisation, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select people::Organisation { *, people:{ * }, meta:{ * } } order by .code;
		`, &organisations)
	return
}

func DeleteOrganisation(db edgedb.Executor, code string) (org Organisation, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select(
				delete people::Organisation filter .code = <str>$0 limit 1
			) { *, people:{ * }, meta: { * }};
		`, &org, code)
	return
}

func (org Organisation) Delete(db edgedb.Executor) (Organisation, error) {
	return DeleteOrganisation(db, org.Code)
}

func (inst OrganisationInput) Save(db edgedb.Executor) (created Organisation, err error) {
	args, _ := json.Marshal(inst)
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module people,
  	data := <json>$0,
		select (insert_organisation(data)) { *, people:{ * }, meta:{ * } };
	`, &created, args)
	return
}

type OrganisationUpdate struct {
	Name        models.OptionalInput[string] `edgedb:"name" json:"name,omitempty" example:"Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés"`
	Code        models.OptionalInput[string] `edgedb:"code" json:"code,omitempty" example:"LEHNA"`
	Description models.OptionalNull[string]  `edgedb:"description" json:"description,omitempty" example:"Where this database was born."`
	Kind        models.OptionalNull[OrgKind] `edgedb:"kind" json:"kind,omitempty" example:"Lab"`
}

func (org OrganisationUpdate) Save(e edgedb.Executor, code string) (updated Organisation, err error) {

	query := db.UpdateQuery{
		Frame: `#edgeql
			with module people, data := <json>$1
			select(
				update Organisation filter .code = <str>$0 set { %s }
			) { *, people:{ *, user: { * } }, meta:{ * } }
		`,
		Mappings: map[string]string{
			"name":        "<str>data['name']",
			"code":        "<str>data['code']",
			"description": "<str>data['description']",
			"kind":        "<OrgKind>data['kind']",
		},
	}
	args, _ := json.Marshal(org)
	logrus.Debugf("Updating organisation %s with %+v", code, org)
	err = e.QuerySingle(context.Background(), query.Query(org), &updated, code, args)
	updated.Meta.Save(e)
	return
}
