package location

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

type HabitatGroup struct {
	ID        edgedb.UUID                    `edgedb:"id" json:"id" format:"uuid"`
	Label     string                         `edgedb:"label" json:"label" doc:"Name for the group of habitat tags" example:"Water flow"`
	Exclusive bool                           `edgedb:"exclusive_elements" json:"exclusive_elements"`
	Depends   models.Optional[HabitatRecord] `edgedb:"depends" json:"depends,omitempty"`
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
	Label     string         `json:"label,omitempty"`
	Depends   string         `json:"depends,omitempty"`
	Exclusive *bool          `json:"exclusive_elements,omitempty"`
	Elements  []HabitatInput `json:"elements,omitempty"`
	// Label     string          `json:"label" doc:"Name for the group of habitat tags" example:"Water flow" minLength:"3" maxLength:"32"`
	// Depends   string          `json:"depends,omitempty" doc:"Habitat tag that this group is a refinement of" example:"Aquatic, Surface"`
	// Exclusive *bool           `json:"exclusive_elements,omitempty"`
	// Elements  []HabitatInput `json:"elements,omitempty"`
}

func (t HabitatGroupUpdate) Schema(r huma.Registry) *huma.Schema {
	s := *r.Schema(reflect.TypeFor[HabitatGroupInput](), false, "")
	s.Required = []string{}
	r.Map()["HabitatGroupUpdate"] = &s
	return &s
}

func (t HabitatGroupUpdate) Update(db edgedb.Executor, label string) (id edgedb.UUID, err error) {
	data, _ := json.Marshal(t)
	err = db.QuerySingle(context.Background(),
		`with item := <json>$1,
		update location::HabitatGroup filter .label = $0 set {
			label := <str>json_get(item, 'label') ?? .label,
			depends := (
				if exists json_get(item, 'depends')
				then (select location::Habitat filter .label = <str>item['depends'])
				else .depends
			),
			exclusive := <bool>json_get(item, 'exclusive_elements')
		}`, &id, label, data)
	return
}

func SetHabitatGroupParent(db edgedb.Executor, groupLabel string, parentHabitatLabel string) (*HabitatGroup, error) {
	var updated HabitatGroup
	err := db.QuerySingle(context.Background(),
		`select (update location::HabitatGroup filter .id = <uuid>$0 set {
			depends := assert_exists((select location::Habitat filter .label = <str>$1))
		}) { ** }`,
		&updated, groupLabel, parentHabitatLabel)
	return &updated, err
}
