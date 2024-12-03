package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

type HabitatGroup struct {
	ID        edgedb.UUID                    `edgedb:"id" json:"id" format:"uuid"`
	Label     string                         `edgedb:"label" json:"label" doc:"Name for the group of habitat tags" example:"Water flow"`
	Exclusive bool                           `edgedb:"exclusive_elements" json:"exclusive_elements"`
	Depends   models.Optional[HabitatRecord] `edgedb:"depends" json:"depends"`
	Elements  []HabitatRecord                `edgedb:"elements" json:"elements"`
	Meta      people.Meta                    `edgedb:"meta" json:"meta"`
}

func (g *HabitatGroup) AddHabitat(e edgedb.Executor, h HabitatInput) error {
	data, _ := json.Marshal(h)
	var new_habitat HabitatRecord
	err := e.QuerySingle(context.Background(),
		`#edgeql
      with data := <json>$1,
			select (insert sampling::Habitat {
				label := <str>data['label'],
				description := <str>json_get(data, 'description'),
				in_group := (assert_single(assert_exists(
					(select <sampling::HabitatGroup><uuid>$0)

				)))
			}) { id, label, description, incompatible: { * } }
		`, &new_habitat, g.ID, data)
	if err != nil {
		return err
	}
	g.Elements = append(g.Elements, new_habitat)
	return nil
}

func (g *HabitatGroup) DeleteHabitat(e edgedb.Executor, label string) error {
	err := e.Execute(context.Background(),
		`#edgeql
			delete assert_exists(sampling::Habitat filter .label = <str>$0)
		`, label,
	)
	if err != nil {
		return err
	}
	g.Elements = slices.DeleteFunc(g.Elements, func(elt HabitatRecord) bool {
		return elt.Label == label
	})
	return nil
}

func (g *HabitatGroup) UpdateHabitat(e edgedb.Executor, label string, h HabitatUpdate) error {
	data, _ := json.Marshal(h)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with data := <json>$1,
      select (update assert_exists(sampling::Habitat filter .label = <str>$0) set {
        %s
      }) { *, incompatible: { * } }
		`,
		Mappings: map[string]string{
			"label":       "<str>data['label']",
			"description": "<str>data['description']",
		},
	}
	var updated HabitatRecord
	err := e.QuerySingle(context.Background(), query.Query(h), &updated, label, data)
	if err != nil {
		return err
	}
	indexToReplace := slices.IndexFunc(g.Elements, func(elt HabitatRecord) bool {
		return elt.Label == label
	})
	g.Elements[indexToReplace] = updated
	return nil
}

type HabitatGroupInput struct {
	Label     string         `json:"label" doc:"Name for the group of habitat tags" example:"Water flow" minLength:"3" maxLength:"32"`
	Depends   string         `json:"depends,omitempty" doc:"Habitat tag that this group is a refinement of" example:"Aquatic, Surface"`
	Exclusive *bool          `json:"exclusive_elements,omitempty"`
	Elements  []HabitatInput `json:"elements" minItems:"1"`
}

func (g HabitatGroupInput) Save(db edgedb.Executor) (created HabitatGroup, err error) {
	data, _ := json.Marshal(g)
	err = db.QuerySingle(context.Background(),
		`#edgeql

			with module sampling,
			data := <json>$0,
			newGroup := (insert HabitatGroup {
				label := <str>data['label'],
				exclusive_elements := <bool>json_get(data, 'exclusive_elements') ?? true
			}),
			habitats := (for habitat in json_array_unpack(json_get(data, 'elements')) union (insert Habitat {
					label := <str>habitat['label'],
					description := <str>json_get(habitat, 'description'),
					in_group := newGroup
				})
			),
			select newGroup { **, elements := habitats {*} }

		`, &created, data,
	)
	return
}

func FindHabitatGroup[ID string | edgedb.UUID](db edgedb.Executor, id ID) (*HabitatGroup, error) {
	query, err := func() (string, error) {
		switch any(id).(type) {
		case edgedb.UUID:
			return `#edgeql
				select sampling::HabitatGroup { ** } filter .id = <uuid>$0`, nil
		case string:
			return `#edgeql
				select sampling::HabitatGroup { ** } filter .label = <str>$0`, nil
		}
		return ``, fmt.Errorf("Invalid identifier type for find argument")
	}()
	if err != nil {
		return nil, err
	}
	var group HabitatGroup
	err = db.QuerySingle(context.Background(), query, &group, id)
	return &group, err
}

func ListHabitatGroups(db edgedb.Executor) ([]HabitatGroup, error) {
	var groups []HabitatGroup
	err := db.Query(context.Background(),
		`#edgeql
		select sampling::HabitatGroup {
			*, depends: { * }, elements: { *, incompatible : { * } }
		}`,
		&groups)
	return groups, err
}

func DeleteHabitatGroup(db edgedb.Executor, label string) (deleted HabitatGroup, err error) {
	query := `#edgeql
		select(
			delete sampling::HabitatGroup filter .label = <str>$0
		){ ** };`
	err = db.QuerySingle(context.Background(), query, &deleted, label)
	return
}

type HabitatGroupUpdate struct {
	Label      models.OptionalInput[string]                   `edgedb:"label" json:"label,omitempty"`
	Depends    models.OptionalNull[string]                    `edgedb:"depends" json:"depends,omitempty"`
	Exclusive  models.OptionalInput[bool]                     `edgedb:"exclusive_elements" json:"exclusive_elements,omitempty"`
	CreateTags models.OptionalInput[[]HabitatInput]           `json:"create_tags,omitempty"`
	UpdateTags models.OptionalInput[map[string]HabitatUpdate] `json:"update_tags,omitempty"`
	DeleteTags models.OptionalInput[[]string]                 `json:"delete_tags,omitempty"`
}

func (u HabitatGroupUpdate) Save(e edgedb.Executor, label string) (updated HabitatGroup, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (
				update sampling::HabitatGroup filter .label = <str>$0 set {
					%s # update clauses
				}
			) { ** }
		`,
		Mappings: map[string]string{
			"depends":            "(select sampling::Habitat filter .label = <str>item['depends'])",
			"label":              "<str>item['label']",
			"exclusive_elements": "<bool>item['exclusive_elements']",
		},
	}
	logrus.Infof("Value: %+v", u)
	logrus.Infof("Query: %v", query.Query(u))
	return updated, e.(*edgedb.Client).Tx(
		context.Background(),
		func(ctx context.Context, tx *edgedb.Tx) (err error) {
			if err = tx.QuerySingle(
				context.Background(),
				query.Query(u),
				&updated, label, data,
			); err != nil {
				return err
			}
			if inputs, create := u.CreateTags.Get(); create {
				for _, input := range inputs {
					logrus.Debugf("Adding habitat: %+v", input)
					if err = updated.AddHabitat(tx, input); err != nil {
						return err
					}
				}
			}
			if labels, delete := u.DeleteTags.Get(); delete {
				for _, label := range labels {
					logrus.Debugf("Deleting habitat %s", label)
					if err = updated.DeleteHabitat(tx, label); err != nil {
						return err
					}
				}
			}
			if inputs, update := u.UpdateTags.Get(); update {
				for label, input := range inputs {
					logrus.Debugf("Updating habitat %s : %+v", label, input)
					if err = updated.UpdateHabitat(tx, label, input); err != nil {
						return err
					}
				}
			}

			updated.Meta.Save(tx)
			return nil
		})
}
