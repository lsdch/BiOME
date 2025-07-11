package occurrence

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/location"
	"github.com/lsdch/biome/services/geoapify"

	"github.com/gosimple/slug"
)

type OccurrencePerSite map[string]EventInputWithActions

// SiteDataset represents a dataset of sites.
type SiteDataset struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           []SiteItem `gel:"sites" json:"sites"`
}

func (d *SiteDataset) ToOccurrenceDataset() *OccurrenceDataset {
	return &OccurrenceDataset{
		Dataset: dataset.Dataset{
			DatasetInner: dataset.DatasetInner{
				ID:          d.ID,
				Label:       d.Label,
				Slug:        d.Slug,
				Description: d.Description,
			},
			Maintainers: d.Maintainers,
			Meta:        d.Meta,
		},
	}
}

func (d *SiteDataset) AddSites(db geltypes.Executor, site_ids []geltypes.UUID) (*SiteDataset, error) {
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

func (d *SiteDataset) CreateSites(tx geltypes.Tx, sites []SiteInput) error {
	sitesData, _ := json.Marshal(sites)
	query := `#edgeql
		with
			dataset := <datasets::SiteDataset><uuid>$0,
			sites := <json>$1,
			created_sites := (
				for site in json_array_unpack(sites) union (
					location::insert_site(site)
				)
			)
		select (update dataset set {
			sites := (select distinct (
				.sites union created_sites
			))
		}) { **, sites: { *, country: { * } } }
	`

	return tx.QuerySingle(context.Background(), query, d, d.ID, sitesData)
}

func ListSiteDatasets(db geltypes.Executor) (datasets []SiteDataset, err error) {
	err = db.Query(context.Background(),
		`#edgeql
			select datasets::SiteDataset { ** }
		`,
		&datasets,
	)
	return
}

// GetSiteDataset retrieves a dataset of sites by its slug.
func GetSiteDataset(db geltypes.Executor, slug string) (dataset SiteDataset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select datasets::SiteDataset {
				**,
				sites: { *, country: { * } }
			} filter .slug = <str>$0
		`,
		&dataset, slug,
	)
	return
}

type SiteInputList []SiteInput

// RequestBody returns the list of coordinates to be used in a batch reverse geocoding request.
func (d SiteInputList) RequestBody() []geoapify.LatLongCoords {
	body := make([]geoapify.LatLongCoords, len(d))
	for i, v := range d {
		body[i] = geoapify.LatLongCoords{
			Lat: v.Coordinates.Latitude,
			Lon: v.Coordinates.Longitude,
		}
	}
	return body
}

// FillPlaces fills the locality and country code of each site in the list,
// based on their coordinates using the Geoapify API.
func (d SiteInputList) FillPlaces(db geltypes.Executor, apiKey string) error {
	client, err := geoapify.NewClient(
		geoapify.WithApiKey(apiKey),
	)
	if err != nil {
		return err
	}
	response, err := client.BatchReverseGeocode(db, d.RequestBody())
	if err != nil {
		return err
	}

	for i, v := range response {
		d[i].CountryCode.SetValue(strings.ToUpper(v.CountryCode))
		switch d[i].Coordinates.Precision {
		case location.M100, location.KM1, location.KM10:
			locality := v.City
			if v.State != "" {
				locality = locality + ", " + v.State
			}
			d[i].Locality.SetValue(locality)
		case location.KM100:
			d[i].Locality.SetValue(v.State)
		case location.Unknown:
			// skip
			break
		}
	}
	return nil
}

func (i SiteInputList) Save(e geltypes.Executor) (created []Site, err error) {
	data, _ := json.Marshal(i)
	err = e.Query(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (
				for site in json_array_unpack(data) union (
				 location::insert_site(site)
				)
			) { *, country: { * }, meta: { * }, datasets: { * } }
		`,
		&created, data)
	return
}

// SiteDatasetInput represents the input for creating a dataset of sites.
// Dataset is populated with existing sites using their codes and new sites are created from the input.
type SiteDatasetInput struct {
	dataset.DatasetInput `json:",inline"`
	Sites                []string      `json:"sites,omitempty" doc:"Existing site codes to include in the dataset"`
	NewSites             SiteInputList `json:"new_sites,omitempty" doc:"New sites to include in the dataset"`
}

// ValidateExistingSites checks if the sites in the input exist in the database.
func (i *SiteDatasetInput) ValidateExistingSites(edb geltypes.Executor) ([]geltypes.UUID, []error) {
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

// ValidateNewSites checks if the new sites in the input are valid.
func (s *SiteDatasetInput) ValidateNewSites(edb geltypes.Executor) []error {
	var errors []error
	for i, site := range s.NewSites {
		if errs := site.Validate(edb); errs != nil {
			errors = slices.Concat(errors, errs.WithLocation(fmt.Sprintf("new_sites[%d].", i)))
		}
	}
	return errors
}

// Validate checks if the input is valid and returns a validated version of it.
func (i *SiteDatasetInput) Validate(edb geltypes.Executor) (*SiteDatasetInputValidated, []error) {
	maintainers, errsMaintainers := i.Maintainers.Validate(edb)
	sites, errsSites := i.ValidateExistingSites(edb)
	errsNewSites := i.ValidateNewSites(edb)
	errs := slices.Concat(errsMaintainers, errsSites, errsNewSites)
	if errs != nil {
		return nil, errs
	}

	return &SiteDatasetInputValidated{
		Label:       i.Label,
		Slug:        slug.Make(i.Label),
		Description: i.Description,
		Maintainers: maintainers,
		Sites:       sites,
		NewSites:    i.NewSites,
	}, nil
}

// SiteDatasetInputValidated represents a validated input for creating a dataset of sites.
type SiteDatasetInputValidated struct {
	Label       string                       `json:"label"`
	Slug        string                       `json:"slug"`
	Description models.OptionalInput[string] `json:"description"`
	Maintainers []geltypes.UUID              `json:"maintainers"`
	Sites       []geltypes.UUID              `json:"sites"`
	NewSites    SiteInputList                `json:"new_sites"`
}

func (i *SiteDatasetInputValidated) SaveTx(tx geltypes.Tx) (*SiteDataset, error) {
	var created SiteDataset
	m, _ := json.Marshal(i)

	err := tx.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0
			select(insert datasets::SiteDataset {
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
					) ?? (select admin::Settings.superadmin.identity)
				)
			}) { ** }
		`, &created, m)
	if err != nil {
		return nil, fmt.Errorf("Failed to create dataset %s: %v", i.Label, err)
	}

	// Add existing sites to the dataset
	if len(i.Sites) > 0 {
		if _, err := created.AddSites(tx, i.Sites); err != nil {
			return nil, fmt.Errorf(
				"Failed to add existing sites into dataset %s: %v",
				i.Label, err,
			)
		}
	}

	// Create new sites and add them to the dataset
	if len(i.NewSites) > 0 {
		if err := created.CreateSites(tx, i.NewSites); err != nil {
			return nil, fmt.Errorf(
				"Failed to save new sites into dataset %s: %v",
				i.Label, err,
			)
		}
	}

	return &created, nil
}

func (i *SiteDatasetInputValidated) Save(db *gel.Client) (*SiteDataset, error) {
	var created = new(SiteDataset)
	err := db.Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) (err error) {
		created, err = i.SaveTx(tx)
		if err != nil {
			created = nil
		}
		return err
	})
	return created, err
}
