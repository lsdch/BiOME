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
	Acronym     string             `json:"acronym" edgedb:"acronym" example:"MELES" binding:"required,alphanum,min=2,max=12"`
	Description edgedb.OptionalStr `json:"description" edgedb:"description" example:"The main ecological research lab on Tatooine."`
} // @name InstitutionInput

type Institution struct {
	ID               edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>" binding:"required"`
	InstitutionInput `edgedb:"$inline"`
	People           []person.Person `json:"people" edgedb:"people"`
	Meta             models.Meta     `json:"meta" edgedb:"meta" binding:"required"`
} // @name Institution

func Find(db *edgedb.Client, acronym string) (inst Institution, err error) {
	query := `select people::Institution {
		id, name, acronym, description, people, meta
	} filter .acronym = <str>$0;`
	err = db.QuerySingle(context.Background(), query, &inst, acronym)
	return
}

func List() (institutions []Institution, err error) {
	query := `select
		people::Institution {
			id, name, acronym, description, people, meta: { * }
		}
		order by .acronym;`
	err = models.DB().Query(context.Background(), query, &institutions)
	return
}

func (inst *InstitutionInput) Create(db *edgedb.Client) (*Institution, error) {
	query := `
	with module people,
	data := <json>$0,
	select (
		insert Institution {
			name := <str>data['name'],
			acronym := <str>data['acronym'],
			description := <str>data['description']
		}
	) { id, name, acronym, description, people, meta: { * } };
	`
	var createdInstitution Institution
	args, _ := json.Marshal(inst)
	if err := db.QuerySingle(context.Background(), query, &createdInstitution, args); err != nil {
		return nil, err
	}
	return &createdInstitution, nil
}

func (inst *Institution) Update(db *edgedb.Client) (*Institution, error) {
	query := `
	with module people,
	data := <json>$0,
	select(
		update Institution filter .id = <uuid>data['id']
		set {
			name := <str>data['name'],
			acronym := <str>data['acronym'],
			description := <str>data['description']
		}
	) { id, name, acronym, description, people, meta: { * } };
	`
	args, _ := json.Marshal(inst)
	var updated Institution
	err := db.QuerySingle(context.Background(), query, &updated, args)
	logrus.Debugf("Updated institution : %+v", updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (inst *Institution) Delete(db *edgedb.Client) error {
	query := `
	with module people,
	delete Institution
	filter .id = <uuid>$0
	`
	if err := db.Execute(context.Background(), query, inst.ID); err != nil {
		return err
	}
	return nil
}
