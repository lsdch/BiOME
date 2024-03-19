package people

import (
	"context"
	"darco/proto/models"
	_ "embed"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type InstitutionInput struct {
	Name        string             `json:"name" edgedb:"name" example:"Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés" binding:"required,min=10,max=128"`
	Code        string             `json:"code" edgedb:"code" example:"LEHNA" binding:"required,alphanum,min=2,max=12" faker:"word,len=10"`
	Kind        InstitutionKind    `json:"kind" edgedb:"kind" example:"Lab" binding:"required,institution_kind" faker:"oneof: Lab, FundingAgency"`
	Description edgedb.OptionalStr `json:"description,omitempty" edgedb:"description" example:"Where this database was born."`
} // @name InstitutionInput

type Institution struct {
	ID               edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>" binding:"required"`
	InstitutionInput `edgedb:"$inline"`
	People           []Person    `json:"people,omitempty" edgedb:"people"`
	Meta             models.Meta `json:"meta" edgedb:"meta"`
} // @name Institution

func FindInstitution(db edgedb.Executor, uuid edgedb.UUID) (inst Institution, err error) {
	err = db.QuerySingle(context.Background(),
		`select people::Institution { *, people:{ * }, meta:{ * } }
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
	Name        *string          `json:"name,omitempty" binding:"omitnil,min=3,max=128" example:"Laboratoire d'Écologie des Hydrosystèmes Naturels et Anthropisés"`
	Code        *string          `json:"code,omitempty" binding:"omitnil,alphanum,min=2,max=12" example:"LEHNA" faker:"word,len=10"`
	Description *string          `json:"description,omitempty" binding:"omitnil" example:"Where this database was born." faker:"sentence"`
	Kind        *InstitutionKind `json:"kind,omitempty" example:"Lab" binding:"omitnil,institution_kind" faker:"oneof: Lab, FundingAgency"`
} //@name InstitutionUpdate

//go:embed queries/update_institution.edgeql
var institutionUpdateQuery string

func (inst InstitutionUpdate) Update(db edgedb.Executor, code string) (id edgedb.UUID, err error) {
	args, _ := json.Marshal(inst)
	logrus.Debugf("Updating institution %s with %+v", code, inst)
	err = db.QuerySingle(context.Background(), institutionUpdateQuery,
		&id, code, args)
	return
}
