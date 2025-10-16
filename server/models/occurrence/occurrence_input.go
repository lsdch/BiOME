package occurrence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

type OccurrenceBatchMetadataInputs struct {
	Organisations map[string]people.OrganisationInput   `json:"organisations,omitempty"`
	People        map[string]people.PersonInput         `json:"people,omitempty"`
	DataSources   map[string]references.DataSourceInput `json:"data_sources,omitempty"`
	Taxa          []taxonomy.TaxonInput                 `json:"taxa,omitempty"`
	Bibliography  map[string]references.ArticleInput    `json:"bibliography,omitempty"`
}

type CreatedMetadata struct {
	Organisations map[string]string `json:"organisations,omitempty"` // input string to code map
	People        map[string]string `json:"people,omitempty"`        // input string to alias map
	DataSources   map[string]string `json:"data_sources,omitempty"`  // input string to code map
	Bibliography  map[string]string `json:"bibliography,omitempty"`  // input string to code map
}

func (i OccurrenceBatchMetadataInputs) saveNewDataSources(tx geltypes.Tx) (map[string]string, error) {
	codes := make(map[string]string)
	for rawSource, source := range i.DataSources {
		created, err := source.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawSource)
		}
		codes[rawSource] = created.Code
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewBibliography(tx geltypes.Tx) (map[string]string, error) {
	codes := make(map[string]string)
	for rawRef, ref := range i.Bibliography {
		created, err := ref.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawRef)
		}
		codes[rawRef] = created.Code
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewOrganisations(tx geltypes.Tx) (map[string]string, error) {
	codes := make(map[string]string)
	for rawOrg, org := range i.Organisations {
		logrus.Infof("Creating organisation '%s' %+v", rawOrg, org)
		created, err := org.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawOrg)
		}
		codes[rawOrg] = created.Code
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewPersons(tx geltypes.Tx, orgCodes map[string]string) (map[string]string, error) {
	personAliases := make(map[string]string)
	for rawPerson, person := range i.People {
		created, err := person.WithOrganisationCodes(orgCodes).Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawPerson)
		}
		personAliases[rawPerson] = created.Alias
	}
	return personAliases, nil
}

func (i OccurrenceBatchMetadataInputs) Save(tx geltypes.Tx) (*CreatedMetadata, error) {

	for j, taxon := range i.Taxa {
		if _, err := taxon.Save(tx); err != nil {
			logrus.Errorf("Failed to save taxon: %+v", taxon)
			return nil, models.WrapErrorIndex(err, j).PrependPath("taxa")
		}
	}

	dataSources, err := i.saveNewDataSources(tx)
	if err != nil {
		return nil, models.WrapErrorPath(err, "data_sources")
	}

	bibliography, err := i.saveNewBibliography(tx)
	if err != nil {
		return nil, models.WrapErrorPath(err, "bibliography")
	}

	organisations, err := i.saveNewOrganisations(tx)
	if err != nil {
		return nil, models.WrapErrorPath(err, "organisations")
	}

	personAliases, err := i.saveNewPersons(tx, organisations)
	if err != nil {
		return nil, models.WrapErrorPath(err, "people")
	}
	return &CreatedMetadata{
		Organisations: organisations,
		People:        personAliases,
		DataSources:   dataSources,
		Bibliography:  bibliography,
	}, nil
}

// OccurrenceBatchInput is the input type for registering occurrences in bulk,
// including all the necessary upstream data:
// site, events, sampling.
// Occurrences can be registered in bulk, with multiple events and samplings.
// Occurrences types include: BioMaterial (internal/external) and external sequences.
type OccurrenceBatchInput struct {
	OccurrenceBatchMetadataInputs `json:",inline"`
	Occurrences                   []SiteOccurrenceInput `json:"occurrences"`
}

func (i OccurrenceBatchInput) ListMissingTaxa(tx geltypes.Tx) (missing []string, err error) {
	taxa := mapset.NewSet[string]()
	for _, siteWithOccurrences := range i.Occurrences {
		for _, sampling := range siteWithOccurrences.Samplings {
			for _, internalBiomat := range sampling.Internal {
				taxa.Add(internalBiomat.Identification.Taxon)
			}
			for _, externalBiomat := range sampling.External {
				taxa.Add(externalBiomat.Identification.Taxon)
			}
		}
	}
	taxaList := taxa.ToSlice()
	missingTaxa := []string{}
	err = tx.Query(context.Background(),
		`#edgeql
			with module taxonomy,
			existing := (select Taxon.name)
			select array_unpack(<array<str>>$0) except existing
		`,
		&missingTaxa, taxaList)
	return missingTaxa, err
}

func (i OccurrenceBatchInput) SaveSites(tx geltypes.Tx) error {
	data, _ := json.Marshal(i.Occurrences)
	return tx.Execute(context.Background(),
		`#edgeql
			with data := <json>$0,
			for site_data in json_array_unpack(data) union (
				location::insert_site(site_data)
			)
		`, data)
}

func (i OccurrenceBatchInput) Save(tx geltypes.Tx) (occurrences []GenericOccurrence[SamplingOutline], err error) {

	replacements, err := i.OccurrenceBatchMetadataInputs.Save(tx)
	if err != nil {
		return nil, err
	}

	missingTaxa, err := i.ListMissingTaxa(tx)
	if err != nil {
		return nil, err
	}
	if len(missingTaxa) > 0 {
		err = os.WriteFile("missing_taxa.txt", []byte(strings.Join(missingTaxa, "\n")), 0644)
		if err != nil {
			logrus.Errorf("Failed to write missing taxa to file: %v", err)
		}
		return nil, models.WrapErrorPath(fmt.Errorf("the following taxa are missing: %v.\nPlease add missing taxa definitions in the 'taxa' field of your occurrence batch input", missingTaxa), "occurrences")
	}

	logrus.Infof("Saving %d sites", len(i.Occurrences))
	if err := i.SaveSites(tx); err != nil {
		return nil, models.WrapErrorPath(err, "occurrences")
	}
	logrus.Infof("Saving occurrences")

	for i, siteOccurrence := range i.Occurrences {
		siteOccurrence.WithCreatedMetadata(replacements)

		if err := siteOccurrence.SaveAbiotics(tx); err != nil {
			return nil, models.WrapErrorIndex(err, i).PrependPath("occurrences")
		}
		for j, sampling := range siteOccurrence.Samplings {
			occ, err := sampling.Save(tx, siteOccurrence.Code)
			if err != nil {
				return nil, models.WrapErrorIndex(err, j).PrependPath("samplings").PrependIndex(i).PrependPath("occurrences")
			} else {
				occurrences = append(occurrences, occ...)
			}
		}
	}
	return
}

/*
SiteOccurrenceInput is the input type for registering a site and its occurrences in bulk.
It includes the site data and a list of samplings and abiotic measurements.
*/
type SiteOccurrenceInput struct {
	SiteInput           `json:",inline"`
	Samplings           []SamplingInputWithOccurrences `json:"samplings"`
	AbioticMeasurements []AbioticMeasurementInput      `json:"abiotic_measurements"`
}

func (site *SiteOccurrenceInput) WithCreatedMetadata(c *CreatedMetadata) *SiteOccurrenceInput {
	for i := range site.Samplings {
		site.Samplings[i].WithCreatedMetadata(c)
	}
	return site
}

func (site *SiteOccurrenceInput) SaveAbiotics(tx geltypes.Tx) error {
	for i, abiotic := range site.AbioticMeasurements {
		_, err := abiotic.Save(tx, site.Code)
		if err != nil {
			return models.WrapErrorIndex(err, i).PrependPath("abiotic_measurements")
		}
	}
	return nil
}

// EventInputWithActions is the input type for registering an event and its occurrences in bulk.
// It includes the event data and a list of samplings.
// Each sampling can have multiple internal and external biomaterials, and sequences.
// It also includes spottings and abiotic measurements.

type SamplingInputWithOccurrences struct {
	SamplingInput `json:",inline"`
	Internal      []InternalOccurrenceInput `json:"internal_occurrences"`
	External      []ExternalOccurrenceInput `json:"external_occurrences"`
}

func (s *SamplingInputWithOccurrences) WithCreatedMetadata(c *CreatedMetadata) *SamplingInputWithOccurrences {
	s.ActionInput.WithPersonAliases(c.People)
	for i := range s.Internal {
		(&s.Internal[i]).WithCreatedMetadata(c)
	}
	for i := range s.External {
		(&s.External[i]).WithCreatedMetadata(c)
	}
	return s
}

func (i SamplingInputWithOccurrences) Save(tx geltypes.Tx, siteCode string) (occurrences []GenericOccurrence[SamplingOutline], err error) {

	sampling, err := i.SamplingInput.Save(tx, siteCode)
	if err != nil {
		return nil, err
	}

	// Save internal occurrences
	for j, internalBiomat := range i.Internal {
		biomat, err := internalBiomat.Save(tx, sampling.Number)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("internal_biomats")
		} else {
			occurrences = append(occurrences, biomat)
		}
	}

	// Save external occurrences and their sequences
	for j, externalBiomat := range i.External {
		occ, err := externalBiomat.Save(tx, sampling.Number)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("external_biomats")
		} else {
			occurrences = append(occurrences, occ)
		}
	}

	return
}
