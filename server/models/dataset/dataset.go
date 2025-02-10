package dataset

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"

	"github.com/edgedb/edgedb-go"
	"github.com/gosimple/slug"
)

type DatasetInner struct {
	ID          edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Label       string             `edgedb:"label" json:"label"`
	Slug        string             `edgedb:"slug" json:"slug"`
	Pinned      bool               `edgedb:"pinned" json:"pinned"`
	Description edgedb.OptionalStr `edgedb:"description" json:"description"`
}

type AbstractDataset struct {
	DatasetInner `edgedb:"$inline" json:",inline"`
	Maintainers  []people.PersonUser `edgedb:"maintainers" json:"maintainers"`
	Meta         people.Meta         `edgedb:"meta" json:"meta"`
}

func (d *AbstractDataset) IsMaintainer(user people.UserInner) bool {
	for _, u := range d.Maintainers {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

type PolymorphicDataset struct {
	AbstractDataset `edgedb:"$inline" json:",inline"`
	Category        DatasetCategory `edgedb:"category" json:"category"`
}

type ListDatasetOptions struct {
	Pinned  bool   `json:"pinned" query:"pinned"`
	OrderBy string `json:"orderBy,omitempty" query:"orderBy"`
	Limit   int    `json:"limit,omitempty" query:"limit" minimum:"1"`
}

func (o ListDatasetOptions) Options() ListDatasetOptions {
	return o
}

func ListDatasets(db edgedb.Executor, options ListDatasetOptions) ([]PolymorphicDataset, error) {
	var datasets []PolymorphicDataset
	opts, _ := json.Marshal(options)
	query := `#edgeql
			with opts := <json>$0
			select datasets::AnyDataset { ** }
			filter .pinned = <bool>json_get(opts, 'pinned') ?? .pinned
		`
	if options.OrderBy != "" {
		query += fmt.Sprintf(` order by .%s asc`, options.OrderBy)
	}
	if options.Limit != 0 {
		query += fmt.Sprintf(` limit <int32>%d`, options.Limit)
	}

	err := db.Query(context.Background(), query, &datasets, opts)
	return datasets, err
}

type DatasetMaintainersInput []string

func (dm DatasetMaintainersInput) Validate(edb edgedb.Executor) ([]edgedb.UUID, []error) {
	checker := db.DBProperty{Object: "people::Person", Property: "alias"}
	maintainers, absents := checker.ExistAll(edb, dm)
	if errs := []error{}; absents != nil {
		for _, v := range absents {
			errs = append(errs, v.ErrorDetail("maintainers"))
		}
		return nil, errs
	}
	return maintainers, nil
}

type DatasetInput struct {
	Label       string                       `json:"label" minLength:"4" maxLength:"32"`
	Slug        string                       `json:"slug"`
	Pinned      models.OptionalInput[bool]   `json:"pinned,omitempty"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
	Maintainers DatasetMaintainersInput      `json:"maintainers" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
}

func (i *DatasetInput) GenerateSlug() {
	i.Slug = slug.Make(i.Label)
}

type DatasetUpdate struct {
	Label       models.OptionalInput[string]                  `edgedb:"label" json:"label,omitempty" minLength:"4" maxLength:"32"`
	Description models.OptionalNull[string]                   `edgedb:"description" json:"description,omitempty"`
	Pinned      models.OptionalNull[bool]                     `edgedb:"pinned" json:"pinned,omitempty"`
	Maintainers models.OptionalInput[DatasetMaintainersInput] `edgedb:"maintainers" json:"maintainers,omitempty" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
}

func (u DatasetUpdate) Save(e edgedb.Executor, slug string) (updated PolymorphicDataset, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update datasets::AnyDataset filter .slug = <str>$0 set {
				%s
			}) { **, sites: { *, country: { * }}}
		`,
		Mappings: map[string]string{
			"label":       "<str>item['label']",
			"description": "<str>item['description']",
			"pinned":      "<bool>item['pinned']",
			"maintainers": `#edgeql
				(
					select people::Person
					filter .alias in <str>json_array_unpack(item['maintainers'])
				)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, slug, data)
	updated.Meta.Save(e)
	return
}

func TogglePinDataset(db edgedb.Executor, slug string) (dataset PolymorphicDataset, err error) {
	err = db.QuerySingle(context.Background(), `#edgeql
		select (update datasets::AnyDataset filter .slug = <str>$0 set {
			pinned := not .pinned
		}) { ** }
	 `, &dataset, slug)
	return
}
