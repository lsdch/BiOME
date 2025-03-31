package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"

	_ "embed"

	"github.com/goccy/go-yaml"
)

type HabitatLabel struct {
	Label string `gel:"label" json:"label" doc:"A short label for the habitat." minLength:"3" maxLength:"32" example:"Lotic"`
}

type HabitatInner struct {
	ID           geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	HabitatLabel `gel:"$inline" json:",inline"`
	Description  geltypes.OptionalStr `gel:"description" json:"description,omitempty" doc:"Optional habitat description"`
}

type HabitatRecord struct {
	HabitatInner `gel:"$inline" json:",inline"`
	Incompatible []HabitatInner `gel:"incompatible" json:"incompatible"`
}

type OptionalHabitatRecord struct {
	geltypes.Optional `json:"-"`
	_                 struct{} `nullable:"true"`
	HabitatRecord     `gel:"$inline" json:",omitempty"`
}

type Habitat struct {
	HabitatRecord `gel:"$inline" json:",inline"`
	Meta          people.Meta `gel:"meta" json:"meta"`
}

type HabitatInput struct {
	HabitatLabel `json:",inline"`
	Description  *string `json:"description,omitempty" doc:"Optional habitat description"`
}

func (i HabitatInput) Save(db geltypes.Executor) (Habitat, error) {
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
	Label       models.OptionalInput[string] `gel:"label" json:"label,omitempty"`
	Description models.OptionalNull[string]  `gel:"description" json:"description,omitempty"`
}

func ListHabitats(db geltypes.Executor) ([]Habitat, error) {
	var habitats []Habitat
	err := db.Query(context.Background(),
		`select sampling::Habitat { *, meta: { * }, incompatible: { * } }`,
		&habitats)
	return habitats, err
}

//go:embed data/habitats.yaml
var habitatsYaml string

func InitialHabitatsSetup(tx geltypes.Tx) error {
	var input []HabitatGroupInput
	if err := yaml.Unmarshal([]byte(habitatsYaml), &input); err != nil {
		return err
	}
	return ImportHabitats(tx, input)
}

func ImportHabitats(tx geltypes.Tx, habitats []HabitatGroupInput) error {
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
