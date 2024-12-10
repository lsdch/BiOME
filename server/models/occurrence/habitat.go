package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"

	_ "embed"

	"github.com/edgedb/edgedb-go"
	"github.com/goccy/go-yaml"
)

type HabitatInner struct {
	Label string `edgedb:"label" json:"label" doc:"A short label for the habitat." minLength:"3" maxLength:"32" example:"Lotic"`
}

type HabitatRecord struct {
	ID           edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	HabitatInner `edgedb:"$inline" json:",inline"`
	Description  edgedb.OptionalStr `edgedb:"description" json:"description,omitempty" doc:"Optional habitat description"`
	Incompatible []HabitatRecord    `edgedb:"incompatible" json:"incompatible,omitempty"`
}

type OptionalHabitatRecord struct {
	edgedb.Optional `json:"-"`
	_               struct{} `nullable:"true"`
	HabitatRecord   `edgedb:"$inline" json:",omitempty"`
}

type Habitat struct {
	HabitatRecord `edgedb:"$inline" json:",inline"`
	Meta          people.Meta `edgedb:"meta" json:"meta"`
}

type HabitatInput struct {
	HabitatInner `json:",inline"`
	Description  *string `json:"description,omitempty" doc:"Optional habitat description"`
}

func (i HabitatInput) Save(db edgedb.Executor) (Habitat, error) {
	var created Habitat
	habitat, _ := json.Marshal(i)
	err := db.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0
 			select (insert sampling::Habitat {
				label := <str>data['label'],
				description := <str>data['description'],
			}) { *, meta: { * }, incompatible: { * } }
		`,
		&created, habitat)
	return created, err
}

type HabitatUpdate struct {
	Label       models.OptionalInput[string] `edgedb:"label" json:"label,omitempty"`
	Description models.OptionalNull[string]  `edgedb:"description" json:"description,omitempty"`
}

func ListHabitats(db edgedb.Executor) ([]Habitat, error) {
	var habitats []Habitat
	err := db.Query(context.Background(),
		`select sampling::Habitat { *, meta: { * }, incompatible: { * } }`,
		&habitats)
	return habitats, err
}

//go:embed data/habitats.yaml
var habitatsYaml string

func InitialHabitatsSetup(tx *edgedb.Tx) error {
	var input []HabitatGroupInput
	if err := yaml.Unmarshal([]byte(habitatsYaml), &input); err != nil {
		return err
	}
	return ImportHabitats(tx, input)
}

func ImportHabitats(tx *edgedb.Tx, habitats []HabitatGroupInput) error {
	items, _ := json.MarshalIndent(habitats, "", "  ")

	err := tx.Execute(context.Background(),
		`#edgeql
			with module sampling,
			items := json_array_unpack(<json>$0),
			for item in items union (
				with habitatGroup := (insert HabitatGroup {
						label := <str>item['label'],
						exclusive_elements := <bool>json_get(item, 'exclusive_elements') ?? true
					}),
				for habitat in json_array_unpack(item['elements']) union (
					insert Habitat {
						label := <str>habitat['label'],
						description := <str>json_get(habitat, 'description'),
						in_group := habitatGroup,
					}
				)
			);`, items)
	if err != nil {
		return err
	}

	// Link habitat dependencies
	return tx.Execute(context.Background(),
		`#edgeql
		with module sampling,
		items := json_array_unpack(<json>$0),
		select (for item in items union (
			(update HabitatGroup filter .label = <str>item['label'] set {
				depends := assert_single((
					select Habitat filter .label = <str>json_get(item, 'depends')
				))
			})
		));`, items)
}
