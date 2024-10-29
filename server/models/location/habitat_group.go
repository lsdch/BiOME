package location

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"fmt"

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

type HabitatGroupInput struct {
	Label     string         `json:"label" doc:"Name for the group of habitat tags" example:"Water flow" minLength:"3" maxLength:"32"`
	Depends   string         `json:"depends,omitempty" doc:"Habitat tag that this group is a refinement of" example:"Aquatic, Surface"`
	Exclusive *bool          `json:"exclusive_elements,omitempty"`
	Elements  []HabitatInput `json:"elements,omitempty"`
}

func (g HabitatGroupInput) Create(db edgedb.Executor) (created HabitatGroup, err error) {
	data, _ := json.Marshal(g)
	err = db.QuerySingle(context.Background(),
		`with module location,
		data := <json>$0,
		newGroup := (insert HabitatGroup {
			label := <str>data['label'],
			exclusive_elements := <bool>json_get(data, 'exclusive_elements') ?? true,
		}),
		habitats := (for habitat in json_array_unpack(json_get(data, 'elements')) union (insert Habitat {
				label := <str>habitat['label'],
				description := <str>json_get(habitat, 'description'),
				in_group := newGroup
			})
		),
		select newGroup { **, elements := habitats {*} }`, &created, data,
	)
	return
}

func FindHabitatGroup[ID string | edgedb.UUID](db edgedb.Executor, id ID) (*HabitatGroup, error) {
	query, err := func() (string, error) {
		switch any(id).(type) {
		case edgedb.UUID:
			return `select location::HabitatGroup { ** } filter .id = <uuid>$0`, nil
		case string:
			return `select location::HabitatGroup { ** } filter .label = <str>$0`, nil
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
		`select location::HabitatGroup { *, depends: { * }, elements: { *, incompatible : { * } } }`,
		&groups)
	return groups, err
}

func DeleteHabitatGroup(db edgedb.Executor, label string) (deleted HabitatGroup, err error) {
	query := `select(
			delete location::HabitatGroup filter .label = <str>$0
		){ ** };`
	err = db.QuerySingle(context.Background(), query, &deleted, label)
	return
}

type HabitatGroupUpdate struct {
	Label     models.OptionalInput[string] `json:"label,omitempty"`
	Depends   models.OptionalNull[string]  `json:"depends"`
	Exclusive models.OptionalInput[bool]   `json:"exclusive_elements,omitempty"`
}

// func (t HabitatGroupUpdate) Schema(r huma.Registry) *huma.Schema {
// 	s := *r.Schema(reflect.TypeFor[HabitatGroupInput](), false, "")
// 	delete(s.Properties, "elements")
// 	s.Required = []string{}
// 	r.Map()["HabitatGroupUpdate"] = &s
// 	return &s
// }

func (u HabitatGroupUpdate) Update(e edgedb.Executor, label string) (updated HabitatGroup, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `with item := <json>$1,
		select (update location::HabitatGroup filter .label = <str>$0 set {
			%s
		}) { **, elements := habitats {*} }`,
		Mappings: map[string]string{
			"depends":   "(select location::Habitat filter .label = <str>item['depends'])",
			"label":     "<str>item['label']",
			"exclusive": "<bool>item['exclusive_elements']",
		},
	}
	logrus.Infof("Value: %+v", u)
	logrus.Infof("Query: %v", query.Query(u))
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, label, data)
	updated.Meta.Update(e)
	return
}
