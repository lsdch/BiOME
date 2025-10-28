package dataset

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/sirupsen/logrus"

	"github.com/gosimple/slug"
)

type DatasetInner struct {
	ID          geltypes.UUID        `gel:"id" json:"id" format:"uuid"`
	Label       string               `gel:"label" json:"label"`
	Slug        string               `gel:"slug" json:"slug"`
	Pinned      bool                 `gel:"pinned" json:"pinned"`
	Description geltypes.OptionalStr `gel:"description" json:"description"`
	Category    DatasetCategory      `gel:"category" json:"category"`
}

type Dataset struct {
	DatasetInner `gel:"$inline" json:",inline"`
	Publication  models.Optional[references.Article] `gel:"publication" json:"publication,omitempty"`
	Maintainers  []people.PersonUser                 `gel:"maintainers" json:"maintainers"`
	Meta         people.Meta                         `gel:"meta" json:"meta"`
}

func (d Dataset) RollbackImport(db geltypes.Executor) error {
	return db.Execute(context.Background(),
		`#edgeql
			delete Auditable filter .meta.batch_import_id = <str>$0
		`, d.Meta.BatchID,
	)
}

func (d *Dataset) IsMaintainer(user people.UserInner) bool {
	for _, u := range d.Maintainers {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

type ListDatasetOptions struct {
	Pinned   models.OptionalInput[bool] `query:"pinned"`
	Category DatasetCategory            `query:"category"`
	OrderBy  string                     `query:"orderBy"`
	Limit    int                        `query:"limit" minimum:"1"`
}

func (o ListDatasetOptions) Options() ListDatasetOptions {
	return o
}

func ListDatasets(db geltypes.Executor, options ListDatasetOptions) ([]Dataset, error) {
	logrus.Debugf("Options: %+v", options)
	var datasets []Dataset
	opts, _ := json.Marshal(options)
	query := `#edgeql
			with opts := <json>$0
			select datasets::Dataset { *,
				maintainers: { *, user: { * } },
				meta: { * }
			}
			filter .pinned = (<bool>json_get(opts, 'pinned') ?? .pinned)
			# and .category = <datasets::DatasetCategory>(<str>json_get(opts, 'category') ?? <str>.category)
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

func GetDataset(db geltypes.Executor, slug string) (dataset Dataset, err error) {
	err = db.QuerySingle(context.Background(), `#edgeql
		select datasets::Dataset {
			*,
			maintainers: { *, user: { * } },
			meta: { * },
		} filter .slug = <str>$0
	`, &dataset, slug)
	return dataset, err
}

type DatasetMaintainersInput []string

func (dm DatasetMaintainersInput) Validate(edb geltypes.Executor) ([]geltypes.UUID, []error) {
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
	Publication models.OptionalInput[string] `json:"publication,omitempty"`
	Pinned      models.OptionalInput[bool]   `json:"pinned,omitempty"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
	Maintainers DatasetMaintainersInput      `json:"maintainers" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
}

func (i *DatasetInput) GenerateSlug() {
	i.Slug = slug.Make(i.Label)
}

func (i *DatasetInput) Save(db geltypes.Executor) (created Dataset, err error) {
	i.GenerateSlug()
	data, _ := json.Marshal(i)
	err = db.QuerySingle(context.Background(),
		`#edgeql
      with module occurrence,
        data := <json>$0,
				select (
					insert datasets::OccurrenceDataset {
						label := <str>data['label'],
						slug := <str>data['slug'],
						description := <str>json_get(data, 'description'),
						maintainers := (
							select people::Person
							filter .alias in <str>json_array_unpack(data['maintainers'])
						) ?? (SELECT admin::Settings.superadmin.identity),
					}
				) {
        *,
				maintainers: { * },
				publication: { * },
				meta: { * }
      }
      `, &created, data,
	)
	return
}

type DatasetUpdate struct {
	Label       models.OptionalInput[string]                  `gel:"label" json:"label,omitempty" minLength:"4" maxLength:"32"`
	Description models.OptionalNull[string]                   `gel:"description" json:"description,omitempty"`
	Publication models.OptionalNull[string]                   `gel:"publication" json:"publication,omitempty"`
	Pinned      models.OptionalNull[bool]                     `gel:"pinned" json:"pinned,omitempty"`
	Maintainers models.OptionalInput[DatasetMaintainersInput] `gel:"maintainers" json:"maintainers,omitempty" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
}

func (u DatasetUpdate) Save(e geltypes.Executor, slug string) (updated Dataset, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update datasets::Dataset filter .slug = <str>$0 set {
				%s
			}) { **, sites: { *, country: { * }}}
		`,
		Mappings: map[string]string{
			"label": "<str>item['label']",
			"publication": `#edgeql
				(
					select references::Article
					filter .code = <str>item['publication']
				)`,
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

func TogglePinDataset(db geltypes.Executor, slug string) (dataset Dataset, err error) {
	err = db.QuerySingle(context.Background(), `#edgeql
		select (update datasets::Dataset filter .slug = <str>$0 set {
			pinned := not .pinned
		}) { ** }
	 `, &dataset, slug)
	return
}
