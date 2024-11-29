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
	"github.com/gosimple/slug"
)

type DatasetInner struct {
	ID          edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Label       string             `edgedb:"label" json:"label"`
	Slug        string             `edgedb:"slug" json:"slug"`
	Description edgedb.OptionalStr `edgedb:"description" json:"description"`
}

type Dataset struct {
	DatasetInner `edgedb:"$inline" json:",inline"`
	Sites        []SiteItem          `edgedb:"sites" json:"sites"`
	Maintainers  []people.PersonUser `edgedb:"maintainers" json:"maintainers"`
	Meta         people.Meta         `edgedb:"meta" json:"meta"`
}

func FindDataset(db edgedb.Executor, slug string) (*Dataset, error) {
	var dataset Dataset
	err := db.QuerySingle(context.Background(),
		`select datasets::Dataset { **, sites: {*, country: { * }}} filter .slug = <str>$0`,
		&dataset, slug,
	)
	return &dataset, err
}

func (d *Dataset) IsMaintainer(user people.UserInner) bool {
	for _, u := range d.Maintainers {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}

func (d *Dataset) AddSites(db edgedb.Executor, site_ids []edgedb.UUID) (*Dataset, error) {
	err := db.QuerySingle(context.Background(),
		`#edgeql
		select(update <datasets::Dataset><uuid>$0 set {
			sites := (select distinct (
				.sites union (
					select Site filter .id in array_unpack(<array<uuid>>$1)
				)
			))
		}) { **, sites: { *, country: { * } } }
	`, d, d.ID, site_ids)
	return d, err
}

func (d *Dataset) CreateSites(tx *edgedb.Tx, sites []SiteInput) (*Dataset, error) {
	sitesData, _ := json.Marshal(sites)
	query := fmt.Sprintf(`#edgeql
		dataset := <datasets::Dataset><uuid>$0,
		sites := <json>$1,
		created_sites := (
			for site in json_array_unpack(sites) union (
				%s
			)
		)
		select (update dataset set {
			sites := (select distinct (
				.sites union created_sites
			))
		}) { **, sites: { *, country: { * } } }
	`, sites[0].InsertQuery("site"))
	err := tx.QuerySingle(context.Background(), query, d, d.ID, sitesData)
	return d, err
}

func ListDatasets(db edgedb.Executor) ([]Dataset, error) {
	var datasets []Dataset
	err := db.Query(context.Background(),
		`select datasets::Dataset { **, meta: { * } }`,
		&datasets,
	)
	return datasets, err
}

type DatasetMaintainers []string

func (dm DatasetMaintainers) Validate(edb edgedb.Executor) ([]edgedb.UUID, []error) {
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
	Description models.OptionalInput[string] `json:"description,omitempty"`
	Maintainers DatasetMaintainers           `json:"maintainers" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
	Sites       []string                     `json:"sites,omitempty" doc:"Existing site codes to include in the dataset"`
	NewSites    []SiteInput                  `json:"new_sites,omitempty" doc:"New sites to include in the dataset"`
}

func (i *DatasetInput) ValidateExistingSites(edb edgedb.Executor) ([]edgedb.UUID, []error) {
	sites, absents := db.DBProperty{
		Object:   "location::Site",
		Property: "code",
	}.ExistAll(edb, i.Sites)
	if errs := []error{}; absents != nil {
		for _, v := range absents {
			errs = append(errs, v.ErrorDetail("sites"))
		}
		return nil, errs
	}
	return sites, nil
}

func (s *DatasetInput) ValidateNewSites(edb edgedb.Executor) []error {
	var errors []error
	for i, site := range s.NewSites {
		if errs := site.Validate(edb); errs != nil {
			errors = slices.Concat(errors, errs.WithLocation(fmt.Sprintf("new_sites[%d].", i)))
		}
	}
	return errors
}

func (i *DatasetInput) Validate(edb edgedb.Executor) (*DatasetInputValidated, []error) {
	maintainers, errsMaintainers := i.Maintainers.Validate(edb)
	sites, errsSites := i.ValidateExistingSites(edb)
	errsNewSites := i.ValidateNewSites(edb)
	errs := slices.Concat(errsMaintainers, errsSites, errsNewSites)
	if errs != nil {
		return nil, errs
	}

	return &DatasetInputValidated{
		Label:       i.Label,
		Slug:        slug.Make(i.Label),
		Description: i.Description,
		Maintainers: maintainers,
		Sites:       sites,
		NewSites:    i.NewSites,
	}, nil
}

type DatasetInputValidated struct {
	Label       string                       `json:"label"`
	Slug        string                       `json:"slug"`
	Description models.OptionalInput[string] `json:"description"`
	Maintainers []edgedb.UUID                `json:"maintainers"`
	Sites       []edgedb.UUID                `json:"sites"`
	NewSites    []SiteInput                  `json:"new_sites"`
}

func (i *DatasetInputValidated) Save(db *edgedb.Client) (*Dataset, error) {
	var created Dataset
	m, _ := json.Marshal(i)

	err := db.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		err := tx.QuerySingle(context.Background(),
			`#edgeql
			with data := <json>$0
			select(insert datasets::Dataset {
				label := <str>data['label'],
				slug := <str>data['slug'],
				description := <str>json_get(data, 'description'),
				maintainers := (
					select distinct (
						(
							(global default::current_user).identity
						) union (
							select distinct people::Person filter .alias in array_unpack(<array<str>>json_get(data, 'maintainers'))
						)
					)
				)
			}) { ** }
		`, &created, m)
		if err != nil {
			return err
		}

		if len(i.Sites) > 0 {
			if _, err := created.AddSites(tx, i.Sites); err != nil {
				return err
			}
		}

		if len(i.NewSites) > 0 {
			if _, err := created.CreateSites(tx, i.NewSites); err != nil {
				return err
			}
		}

		return nil
	})
	return &created, err
}

type DatasetUpdate struct {
	Label       models.OptionalInput[string]             `json:"label,omitempty" minLength:"4" maxLength:"32"`
	Description models.OptionalNull[string]              `json:"description,omitempty"`
	Maintainers models.OptionalInput[DatasetMaintainers] `json:"maintainers,omitempty" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
}

func (u DatasetUpdate) Save(e edgedb.Executor, slug string) (updated Dataset, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update datasets::Dataset filter .slug = <str>$0 set {
				%s
			}) { **, sites: { *, country: { * }}}
		`,
		Mappings: map[string]string{
			"label":       "<str>item['label']",
			"description": "<str>item['description']",
			"maintainers": `(
				select people::Person
				filter .alias in <str>json_array_unpack(item['maintainers']))
			`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, slug, data)
	updated.Meta.Save(e)
	return
}
