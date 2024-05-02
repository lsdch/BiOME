package location

import (
	"context"
	"encoding/json"

	_ "embed"

	"github.com/edgedb/edgedb-go"
	"github.com/goccy/go-yaml"
	"github.com/sirupsen/logrus"
)

type HabitatInner struct {
	Label       string             `edgedb:"label" json:"label" doc:"A short label for the habitat. If the habitat is a specialization of a more general one, it should not repeat the parent label." example:"Lotic"`
	Description edgedb.OptionalStr `edgedb:"description" json:"description,omitempty" doc:"Optional habitat description"`
}

type HabitatRecord struct {
	ID           edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	HabitatInner `edgedb:"$inline" json:",inline"`
}

type Habitat struct {
	ID            edgedb.UUID `edgedb:"id" json:"id" format:"uuid"`
	HabitatRecord `edgedb:"$inline" json:",inline"`
	Incompatible  []HabitatRecord `edgedb:"incompatible" json:"incompatible"`
	Depends       []HabitatRecord `edgedb:"depends" json:"depends,omitempty"`
}

type HabitatInput struct {
	HabitatInner `json:",inline"`
	Depends      []string `json:"depends,omitempty" doc:"List of habitat labels this habitat may specialize." example:"Aquatic, Surface"`
	Incompatible []string `json:"incompatible,omitempty" doc:"List of habitat labels this habitat is incompatible with." example:"Lentic"`
}

func ListHabitats(db edgedb.Executor) ([]Habitat, error) {
	var habitats []Habitat
	err := db.Query(context.Background(),
		`select location::Habitat { *, depends: { * }, incompatible: { * } }`,
		&habitats)
	logrus.Debugf("HABITATS: %+v", habitats)
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
			}) { *, depends: { * }, incompatible: { * }}`,
		&created, habitat)
	return created, err
}

//go:embed data/habitats.yaml
var habitatsYaml string

func InitialHabitatsSetup(db edgedb.Executor) error {
	var input []HabitatInput
	if err := yaml.Unmarshal([]byte(habitatsYaml), &input); err != nil {
		return err
	}
	logrus.Infof("Habitat inputs: %+v", input)
	_, err := ImportHabitats(db, input)
	return err
}

func ImportHabitats(db edgedb.Executor, habitats []HabitatInput) ([]Habitat, error) {
	var created []Habitat
	items, _ := json.Marshal(habitats)

	err := db.Execute(context.Background(),
		`with items := <json>$0,
			for item in json_array_unpack(items) union (
				insert location::Habitat {
					label := <str>item['label'],
					description := <str>item['description'],
				}
			)`, items)
	if err != nil {
		return nil, err
	}
	err = db.Query(context.Background(),
		`with module location,
			items := <json>$0
		select (
			for item in json_array_unpack(items) union (
				update Habitat filter .label = <str>item['label'] set {
					depends := (
						select detached Habitat
						filter .label in <str>json_array_unpack(json_get(item, 'depends'))
					),
					incompatibleFrom := (
						select detached Habitat
						filter .label in <str>json_array_unpack(json_get(item, 'incompatible'))
					)
				}
			)
		) { *, depends: { * }, incompatible: { * }}`,
		&created, items,
	)
	return created, err
}
