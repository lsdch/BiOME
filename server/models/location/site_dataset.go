package location

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"encoding/json"
	"slices"

	"github.com/edgedb/edgedb-go"
)

type SiteDatasetInner struct {
	ID          edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Label       string             `edgedb:"label" json:"label" `
	Description edgedb.OptionalStr `edgedb:"description" json:"description"`
}

type SiteDataset struct {
	SiteDatasetInner `edgedb:"$inline" json:",inline"`
	Sites            []SiteItem          `edgedb:"sites" json:"sites"`
	Maintainers      []people.PersonUser `edgedb:"maintainers" json:"maintainers"`
}

func (d *SiteDataset) AddSites(db edgedb.Executor, site_ids []edgedb.UUID) (*SiteDataset, error) {
	err := db.QuerySingle(context.Background(),
		`with module location,
		select(update <SiteDataset><uuid>$0 set {
			sites := (select .sites union (
				select Site filter .id in array_unpack(<array<uuid>>$1)
			))
		}) { ** }
	`, d, d.ID, site_ids)
	return d, err
}

func (d *SiteDataset) CreateSites(tx *edgedb.Tx, sites []SiteInput) (*SiteDataset, error) {
	for _, s := range sites {
		created, err := s.Create(tx)
		if err != nil {
			return d, err
		} else {
			d.Sites = append(d.Sites, created.SiteItem)
		}
	}
	return d, nil
}

type SiteDatasetInput struct {
	Label       string                       `json:"label" minLength:"4" maxLength:"32"`
	Description models.OptionalInput[string] `json:"description,omitempty"`
	Maintainers []string                     `json:"maintainers" doc:"Dataset maintainers identified by their person alias. Dataset creator is always a maintainer by default."`
	Sites       []string                     `json:"sites,omitempty" doc:"Existing site codes to include in the dataset"`
	NewSites    []SiteInput                  `json:"new_sites,omitempty" doc:"New sites to include in the dataset"`
}

func (i *SiteDatasetInput) ValidateMaintainers(edb edgedb.Executor) ([]edgedb.UUID, []error) {
	checker := db.DBProperty{Object: "people::Person", Property: "alias"}
	maintainers, absents := checker.ExistAll(edb, i.Maintainers)
	if errs := []error{}; absents != nil {
		for _, v := range absents {
			errs = append(errs, v.ErrorDetail("maintainers"))
		}
		return nil, errs
	}
	return maintainers, nil
}

func (i *SiteDatasetInput) ValidateSites(edb edgedb.Executor) ([]edgedb.UUID, []error) {
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

func (i *SiteDatasetInput) Validate(edb edgedb.Executor) (*SiteDatasetInputValidated, []error) {
	maintainers, errsMaintainers := i.ValidateMaintainers(edb)
	sites, errsSites := i.ValidateSites(edb)
	errs := slices.Concat(errsMaintainers, errsSites)
	if errs != nil {
		return nil, errs
	}

	return &SiteDatasetInputValidated{
		Label:       i.Label,
		Description: i.Description,
		Maintainers: maintainers,
		Sites:       sites,
	}, nil
}

type SiteDatasetInputValidated struct {
	Label       string                       `json:"label"`
	Description models.OptionalInput[string] `json:"description"`
	Maintainers []edgedb.UUID                `json:"maintainers"`
	Sites       []edgedb.UUID                `json:"sites"`
	NewSites    []SiteInput                  `json:"new_sites"`
}

func (i *SiteDatasetInputValidated) Create(db edgedb.Client) (*SiteDataset, error) {
	var created SiteDataset
	m, _ := json.Marshal(i)

	err := db.Tx(context.Background(), func(ctx context.Context, tx *edgedb.Tx) error {
		err := tx.QuerySingle(context.Background(),
			`with module location,
				data := <json>$0
			select(insert location::SiteDataset {
				label := <str>data['label'],
				description := <str>json_get(data, 'description'),
				maintainers := (
					select (global current_user).identity
				)
			})`, &created, m)
		if err != nil {
			return err
		}
		if _, err := created.AddSites(tx, i.Sites); err != nil {
			return err
		}
		if _, err := created.CreateSites(tx, i.NewSites); err != nil {
			return err
		}
		return nil
	})
	return &created, err
}
