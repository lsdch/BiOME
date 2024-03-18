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
	Name        string             `json:"name" edgedb:"name" example:"Mos Eisley Laboratory of Environmental Studies" binding:"required,min=10,max=128"`
	Code        string             `json:"code" edgedb:"code" example:"MELES" binding:"required,alphanum,min=2,max=12"`
	Description edgedb.OptionalStr `json:"description" edgedb:"description" example:"The main ecological research lab on Tatooine."`
	Kind        InstitutionKind    `json:"kind" edgedb:"kind" example:"Lab" binding:"required,institution_kind"`
} // @name InstitutionInput

type Institution struct {
	ID               edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>" binding:"required"`
	InstitutionInput `edgedb:"$inline"`
	People           []Person    `json:"people,omitempty" edgedb:"people"`
	Meta             models.Meta `json:"meta" edgedb:"meta"`
} // @name Institution

func FindInstitution(db edgedb.Executor, uuid edgedb.UUID) (inst Institution, err error) {
	query := "select people::Institution { *, people:{ * }, meta:{ * } } filter	.id = <uuid>$0;"
	err = db.QuerySingle(context.Background(), query, &inst, uuid)
	return inst, err
}

func ListInstitutions(db edgedb.Executor) (institutions []Institution, err error) {
	query := `select people::Institution { *, people:{ * }, meta:{ * } } order by .code;`
	err = db.Query(context.Background(), query, &institutions)
	return
}

func DeleteInstitution(db edgedb.Executor, code string) (inst Institution, err error) {
	query := `select(
		delete people::Institution filter .code = <str>$0
	) { *, people:{ * }, meta: { * }};`
	err = db.QuerySingle(context.Background(), query, &inst, code)
	return
}

//go:embed queries/create_institution.edgeql
var institutionCreateQuery string

func (inst InstitutionInput) Create(db edgedb.Executor) (created Institution, err error) {
	args, _ := json.Marshal(inst)
	err = db.QuerySingle(context.Background(), institutionCreateQuery, &created, args)
	return
}

type InstitutionUpdate struct {
	Name        *string          `json:"name,omitempty" binding:"omitnil,min=3,max=128" example:"Mos Eisley Laboratory of Environmental Studies"`
	Code        *string          `json:"code,omitempty" binding:"omitnil,alphanum,min=2,max=12" example:"MELES"`
	Description *string          `json:"description,omitempty" binding:"omitnil" example:"The main ecological research lab on Tatooine."`
	Kind        *InstitutionKind `json:"kind,omitempty" example:"Lab" binding:"omitnil,institution_kind"`
} //@name InstitutionUpdate

//go:embed queries/update_institution.edgeql
var institutionUpdateQuery string

func (inst InstitutionUpdate) Update(db edgedb.Executor, code string) (id edgedb.UUID, err error) {
	args, _ := json.Marshal(inst)
	logrus.Debugf("Updating institution %s with %+v", code, inst)
	err = db.QuerySingle(context.Background(), institutionUpdateQuery, &id, code, args)
	return
}

func (inst InstitutionUpdate) Validate(db edgedb.Executor, code string) error {
	return nil
}
