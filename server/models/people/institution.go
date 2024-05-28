package people

import (
	"context"
	"darco/proto/models"
	_ "embed"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type InstitutionInfos struct {
	Name string          `json:"name" edgedb:"name" example:"Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés" minLength:"10" maxLength:"128" fake:"{company}"`
	Code string          `json:"code" edgedb:"code" example:"LEHNA" minLength:"2" maxLength:"12" fake:"{word}"`
	Kind InstitutionKind `json:"kind" edgedb:"kind" example:"Lab"`
}

type InstitutionInput struct {
	InstitutionInfos `edgedb:"$inline"`
	Description      models.OptionalInput[string] `json:"description,omitempty" edgedb:"description" example:"Where this database was born."`
}

type InstitutionInner struct {
	ID               edgedb.UUID `json:"id" edgedb:"id" format:"uuid" binding:"required"`
	InstitutionInfos `edgedb:"$inline" json:",inline"`
	Description      edgedb.OptionalStr `json:"description,omitempty" edgedb:"description" example:"Where this database was born."`
}

type Institution struct {
	InstitutionInner `edgedb:"$inline" json:",inline"`
	People           []PersonUser `json:"people,omitempty" edgedb:"people" doc:"Known members of this institution"`
	Meta             Meta         `json:"meta" edgedb:"meta"`
}

func FindInstitution(db edgedb.Executor, uuid edgedb.UUID) (inst Institution, err error) {
	err = db.QuerySingle(context.Background(),
		`select people::Institution { *, people:{ *, user: { * } }, meta:{ * } }
			filter .id = <uuid>$0;`,
		&inst, uuid)
	return inst, err
}

func ListInstitutions(db edgedb.Executor) (institutions []Institution, err error) {
	err = db.Query(context.Background(),
		`select people::Institution { *, people:{ * }, meta:{ * } } order by .code;`,
		&institutions)
	return
}

func DeleteInstitution(db edgedb.Executor, code string) (inst Institution, err error) {
	err = db.QuerySingle(context.Background(),
		`select(
			delete people::Institution filter .code = <str>$0 limit 1
		) { *, people:{ * }, meta: { * }};`, &inst, code)
	return
}

func (inst Institution) Delete(db edgedb.Executor) (Institution, error) {
	return DeleteInstitution(db, inst.Code)
}

//go:embed queries/create_institution.edgeql
var institutionCreateQuery string

func (inst InstitutionInput) Create(db edgedb.Executor) (created Institution, err error) {
	args, _ := json.Marshal(inst)
	err = db.QuerySingle(context.Background(), institutionCreateQuery, &created, args)
	return
}

type InstitutionUpdate struct {
	Name        models.OptionalInput[string]         `json:"name,omitempty" example:"Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés"`
	Code        models.OptionalInput[string]         `json:"code,omitempty" example:"LEHNA"`
	Description models.OptionalNull[string]          `json:"description,omitempty" example:"Where this database was born."`
	Kind        models.OptionalNull[InstitutionKind] `json:"kind,omitempty" example:"Lab"`
}

//go:embed queries/update_institution.edgeql
var institutionUpdateQuery string

func (inst InstitutionUpdate) Update(db edgedb.Executor, code string) (id edgedb.UUID, err error) {

	query := models.UpdateQuery{
		Frame: `with module people, data := <json>$1
			select(update Institution filter .code = <str>$0
			set {
				%s
			}).id`,
		Set: map[models.OptionalNullable]models.FieldMapping{
			inst.Name:        {Field: "name", Value: "<str>data['name']"},
			inst.Code:        {Field: "code", Value: "<str>data['code']"},
			inst.Description: {Field: "description", Value: "<str>data['description']"},
			inst.Kind:        {Field: "kind", Value: "<InstitutionKind>data['kind']"},
		},
	}
	args, _ := json.Marshal(inst)
	logrus.Debugf("Updating institution %s with %+v", code, inst)
	err = db.QuerySingle(context.Background(), query.Query(), &id, code, args)
	return
}
