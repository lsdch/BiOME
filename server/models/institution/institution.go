package institution

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/users"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type InstitutionInput struct {
	Name        string             `json:"name" edgedb:"name" example:"Mos Eisley Laboratory of Environmental Studies" binding:"required,min=10,max=128"`
	Acronym     string             `json:"acronym" edgedb:"acronym" example:"MELES" binding:"required,alphanum,min=2,max=12"`
	Description edgedb.OptionalStr `json:"description" edgedb:"description" example:"The main ecological research lab on Tatooine."`
} // @name InstitutionInput

type Institution struct {
	ID               edgedb.UUID `json:"id" edgedb:"id" example:"<UUID>" binding:"required"`
	InstitutionInput `edgedb:"$inline"`
	People           []users.Person `json:"people" edgedb:"people"`
} // @name Institution

func Find(db *edgedb.Client, acronym string) (inst Institution, err error) {
	query := `select people::Institution {
		id, name, acronym, description, people
	} filter .acronym = <str>$0;`
	err = db.QuerySingle(context.Background(), query, &inst, acronym)
	return
}

func List() (institutions []Institution, err error) {
	query := `select
		people::Institution {
			id, name, acronym, description, people
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
	) { id, name, acronym, description };
	`
	var createdInstitution Institution
	args, _ := json.Marshal(inst)
	if err := db.QuerySingle(context.Background(), query, &createdInstitution, args); err != nil {
		return nil, err
	}
	return &createdInstitution, nil
}

func (inst *Institution) Update(db *edgedb.Client) error {
	query := `
	with module people,
	data := <json>$0,
	update Institution filter .id = <uuid>data['id']
	set {
		name := <str>data['name'],
		acronym := <str>data['acronym'],
		description := <str>data['description']
	};
	`
	args, _ := json.Marshal(inst)
	if err := db.Execute(context.Background(), query, args); err != nil {
		return err
	}
	return nil
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
