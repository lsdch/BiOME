package references

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
)

type Collection struct {
	ID          geltypes.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Label       string               `gel:"label" json:"label"`
	Code        string               `gel:"code" json:"code"`
	Personal    bool                 `gel:"personal" json:"personal"`
	Location    geltypes.OptionalStr `gel:"location" json:"location,omitempty"`
	Description geltypes.OptionalStr `gel:"description" json:"description,omitempty"`
	Contact     geltypes.OptionalStr `gel:"contact" json:"contact,omitempty"`
	Meta        people.Meta          `gel:"meta" json:"meta"`
}

type CollectionWithVouchers struct {
	Collection `json:",inline" gel:"$inline"`
	Vouchers   []string `gel:"vouchers" json:"vouchers,omitempty"`
}

type CollectionInput struct {
	Label    string                       `json:"label" validate:"required"`
	Code     string                       `json:"code" validate:"required"`
	Contact  models.OptionalInput[string] `json:"contact,omitempty"`
	Location models.OptionalInput[string] `json:"location,omitempty"`
	Personal models.OptionalInput[bool]   `json:"personal,omitempty"`
	Vouchers []string                     `json:"vouchers,omitempty"`
}

// Used to link to an existing collections, with optional vouchers
type CollectionField struct {
	Name     string   `gel:"name" json:"name"`
	Vouchers []string `gel:"vouchers" json:"vouchers,omitempty"`
}

func (c *CollectionInput) Save(db geltypes.Executor) (created Collection, err error) {
	coll, _ := json.Marshal(c)
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with c := <json>$0
		insert references::Collection {
			label := <str>c["label"],
			code := <str>c["code"],
			personal := <bool>json_get(c, "personal") ?? false,
			location := <str>json_get(c, "location"),
			description := <str>json_get(c, "description"),
			contact := <str>json_get(c, "contact"),
		}`, &created, coll)
	return
}

func ListCollections(db geltypes.Executor) (collections []Collection, err error) {
	err = db.Query(context.Background(),
		`#edgeql
		select references::Collection {
			*, meta: { * }
		}
		order by .label`, &collections)
	return
}

func DeleteCollection(db geltypes.Executor, codeOrName string) (deleted Collection, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with coll := assert_single((
			select references::Collection
			filter .code = <str>$0
			or .label = <str>$0
		))
		select (delete coll) {
			*, meta: { * }
		}`,
		&deleted, codeOrName)
	return
}
