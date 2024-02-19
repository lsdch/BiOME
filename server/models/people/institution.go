package people

import (
	"context"
	"darco/proto/models"
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

func FindInstitution(db *edgedb.Client, uuid edgedb.UUID) (inst Institution, err error) {
	query := "select people::Institution { *, people:{ * }, meta:{ * } } filter	.id = <uuid>$0;"
	err = db.QuerySingle(context.Background(), query, &inst, uuid)
	return inst, err
}

func ListInstitutions(db *edgedb.Client) (institutions []Institution, err error) {
	query := `select people::Institution { *, people:{ * }, meta:{ * } } order by .code;`
	err = db.Query(context.Background(), query, &institutions)
	return
}

func DeleteInstitution(db *edgedb.Client, code string) (inst Institution, err error) {
	query := `select(
		delete people::Institution filter .code = <str>$0
	) { *, people:{ * }, meta: { * }};`
	err = db.QuerySingle(context.Background(), query, &inst, code)
	return
}

func (inst InstitutionInput) Create(db *edgedb.Client) (created Institution, err error) {
	query := `
		with module people,
			data := <json>$0,
		select ( insert Institution {
			name := <str>data['name'],
			code := <str>data['code'],
			description := <str>json_get(data, 'description'),
			kind := <InstitutionKind>data['kind']
		}) { *, people:{ * }, meta:{ * } };`

	args, _ := json.Marshal(inst)
	err = db.QuerySingle(context.Background(), query, &created, args)
	return
}

type InstitutionUpdate struct {
	Name        *string          `json:"name,omitempty" binding:"omitnil,min=3,max=128" example:"Mos Eisley Laboratory of Environmental Studies"`
	Code        *string          `json:"code,omitempty" binding:"omitnil,alphanum,min=2,max=12" example:"MELES"`
	Description *string          `json:"description,omitempty" binding:"omitnil" example:"The main ecological research lab on Tatooine."`
	Kind        *InstitutionKind `json:"kind,omitempty" example:"Lab" binding:"omitnil,institution_kind"`
} //@name InstitutionUpdate

func (inst InstitutionUpdate) Update(db *edgedb.Client, code string) (id edgedb.UUID, err error) {
	query := `
		with module people,
			data := <json>$1,
		select (
			update Institution filter .code = <str>$0
			set {
				name := <str>json_get(data, 'name') ?? .name,
				code := <str>json_get(data, 'code') ?? .code,
				description := <str>json_get(data, 'description') ??.description,
				kind := <InstitutionKind>json_get(data, 'kind') ?? .kind
			}
		).id;`
	// { *, people:{ * }, meta:{ * } };

	args, _ := json.Marshal(inst)
	logrus.Debugf("Updating institution %s with %+v", code, inst)
	err = db.QuerySingle(context.Background(), query, &id, code, args)
	return
}

func (inst InstitutionUpdate) Validate(db *edgedb.Client, code string) error {
	return nil
}
