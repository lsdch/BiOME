package location

import (
	"context"
	"darco/proto/models/people"
	"encoding/json"

	_ "embed"

	"github.com/davecgh/go-spew/spew"
	"github.com/edgedb/edgedb-go"
	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
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
	edgedb.Optional
	HabitatRecord `edgedb:"$inline"`
}

type Habitat struct {
	HabitatRecord `edgedb:"$inline" json:",inline"`
	Meta          people.Meta `edgedb:"meta" json:"meta"`
}

type HabitatInput struct {
	HabitatInner `json:",inline"`
	Description  *string  `json:"description,omitempty" doc:"Optional habitat description"`
	Incompatible []string `json:"incompatible,omitempty" doc:"List of habitat labels this habitat is incompatible with." example:"Lentic"`
}

func ListHabitats(db edgedb.Executor) ([]Habitat, error) {
	var habitats []Habitat
	err := db.Query(context.Background(),
		`select location::Habitat { *, meta: { * }, incompatible: { * } }`,
		&habitats)
	return habitats, err
}

func (i HabitatInput) Create(db edgedb.Executor) (Habitat, error) {
	var created Habitat
	habitat, _ := json.Marshal(i)
	err := db.QuerySingle(context.Background(),
		`with data := <json>$0
 			select (insert location::Habitat {
				label := <str>data['label'],
				description := <str>data['description'],
				depends := (
					select detached location::Habitat
					filter .label in <str>json_array_unpack(data['depends'])
				),
				incompatibleFrom := (
					select detached location::Habitat
					filter .label in <str>json_array_unpack(data['incompatible'])
				)
			}) { *, depends: { * }, meta: { * }, incompatible: { * }`,
		&created, habitat)
	return created, err
}

//go:embed data/habitats.yaml
var habitatsYaml string

func InitialHabitatsSetup(db *edgedb.Client) error {
	var input []HabitatGroupInput
	if err := yaml.Unmarshal([]byte(habitatsYaml), &input); err != nil {
		return err
	}
	spew.Dump("Habitat inputs: %+v", input)
	return db.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		return ImportHabitats(tx, input)
	})
}

func ImportHabitats(tx *edgedb.Tx, habitats []HabitatGroupInput) error {
	items, _ := json.MarshalIndent(habitats, "", "  ")

	logrus.Infof("%s", items)

	err := tx.Execute(context.Background(),
		`with module location,
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

	return tx.Execute(context.Background(),
		`with module location,
		items := json_array_unpack(<json>$0),
		select (for item in items union (
			(update HabitatGroup filter .label = <str>item['label'] set {
				depends := assert_single((
					select Habitat filter .label = <str>json_get(item, 'depends')
				))
			}) union (
				for habitat in json_array_unpack(item['elements']) union (
					update Habitat filter .label = <str>habitat['label'] set {
						incompatible_from := (
							select detached Habitat
							filter .label in <str>json_array_unpack(json_get(habitat, 'incompatible'))
						)
					}
				)
			)
		));`, items)
}
