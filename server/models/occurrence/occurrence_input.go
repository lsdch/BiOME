package occurrence

import (
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

func (i OccurrenceBatchInput) Save(tx geltypes.Tx) (occurrences []OccurrenceWithCategory, err error) {

	replacements, err := i.OccurrenceBatchMetadataInputs.Save(tx)
	if err != nil {
		return nil, err
	}

	for i, siteOccurrence := range i.Occurrences {
		occ, err := siteOccurrence.WithCreatedMetadata(*replacements).Save(tx)
		if err != nil {
			return nil, models.WrapErrorIndex(err, i).PrependPath("occurrences")
		} else {
			occurrences = append(occurrences, occ...)
		}
	}
	return
}

/*
SiteOccurrenceInput is the input type for registering a site and its occurrences in bulk.
It includes the site data and a list of events.
Each event can have multiple samplings, spottings, and abiotic measurements.
*/
type SiteOccurrenceInput struct {
	SiteInput `json:",inline"`
	Events    []EventInputWithActions `json:"events"`
}

func (i *SiteOccurrenceInput) WithCreatedMetadata(c CreatedMetadata) SiteOccurrenceInput {
	for j := range i.Events {
		i.Events[j].WithCreatedMetadata(c)
	}
	return *i
}

func (i SiteOccurrenceInput) Save(tx geltypes.Tx) ([]OccurrenceWithCategory, error) {
	site, err := i.SiteInput.Save(tx)
	if err != nil {
		return nil, err
	}
	occurrences := []OccurrenceWithCategory{}
	for j, event := range i.Events {
		occ, err := event.Save(tx, site.Code)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("events")
		} else {
			occurrences = append(occurrences, occ...)
		}
	}
	return occurrences, nil
}

// EventInputWithActions is the input type for registering an event and its occurrences in bulk.
// It includes the event data and a list of samplings.
// Each sampling can have multiple internal and external biomaterials, and sequences.
// It also includes spottings and abiotic measurements.
type EventInputWithActions struct {
	EventInput          `json:",inline"`
	Samplings           []SamplingInputWithOccurrences `json:"samplings"`
	Spottings           SpottingUpdate                 `json:"spottings"`
	AbioticMeasurements []AbioticMeasurementInput      `json:"abiotic_measurements"`
}

func (ev EventInputWithActions) WithCreatedMetadata(c CreatedMetadata) EventInputWithActions {
	ev.EventInput.WithPersonAliases(c.People)
	for i := range ev.Samplings {
		ev.Samplings[i].WithCreatedMetadata(c)
	}
	return ev
}

func (i EventInputWithActions) Save(tx geltypes.Tx, site_code string) ([]OccurrenceWithCategory, error) {
	event, err := i.EventInput.Save(tx, site_code)
	if err != nil {
		return nil, err
	}

	if err := event.AddSpottings(tx, i.Spottings); err != nil {
		return nil, models.WrapErrorPath(err, "spottings")
	}

	for j, abioticMeasurement := range i.AbioticMeasurements {
		if err := event.AddAbioticMeasurement(tx, abioticMeasurement); err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("abiotic_measurements")
		}
	}

	occurrences := []OccurrenceWithCategory{}
	for j, sampling := range i.Samplings {
		occ, err := sampling.Save(tx, event.ID)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("samplings")
		} else {
			occurrences = append(occurrences, occ...)
		}
	}

	return occurrences, nil
}

type SamplingInputWithOccurrences struct {
	SamplingInput  `json:",inline"`
	InternalBiomat []InternalBioMatInput              `json:"internal_biomats"`
	ExternalBiomat []ExternalBioMatInputWithSequences `json:"external_biomats"`
	Sequences      []ExternalSequenceInput            `json:"sequences"`
}

func (s *SamplingInputWithOccurrences) WithCreatedMetadata(c CreatedMetadata) SamplingInputWithOccurrences {
	for i := range s.InternalBiomat {
		(&s.InternalBiomat[i]).WithCreatedMetadata(c)
	}
	for i := range s.ExternalBiomat {
		(&s.ExternalBiomat[i]).WithCreatedMetadata(c)
	}
	for i := range s.Sequences {
		(&s.Sequences[i]).WithCreatedMetadata(c)
	}
	return *s
}

func (i SamplingInputWithOccurrences) Save(tx geltypes.Tx, eventID geltypes.UUID) (occurrences []OccurrenceWithCategory, err error) {

	sampling, err := i.SamplingInput.Save(tx, eventID)
	if err != nil {
		return nil, err
	}

	for j, internalBiomat := range i.InternalBiomat {
		biomat, err := internalBiomat.Save(tx, sampling.ID)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("internal_biomats")
		} else {
			occurrences = append(occurrences, biomat.AsOccurrence())
		}
	}

	for j, externalBiomat := range i.ExternalBiomat {
		occ, err := externalBiomat.Save(tx, sampling.ID)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("external_biomats")
		} else {
			occurrences = append(occurrences, occ...)
		}
	}

	for j, sequence := range i.Sequences {
		sequence.UseSamplingCode(sampling.Code)
		seq, err := sequence.Save(tx, sampling.ID)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("sequences")
		} else {
			occurrences = append(occurrences, seq.AsOccurrence())
		}
	}

	return
}

type ExternalBioMatInputWithSequences struct {
	ExternalBioMatInput `json:",inline"`
	Sequences           []ExternalSequenceInput `json:"sequences"`
}

func (bm *ExternalBioMatInputWithSequences) WithCreatedMetadata(c CreatedMetadata) ExternalBioMatInputWithSequences {
	(&bm.ExternalBioMatInput).WithCreatedMetadata(c)
	for i := range bm.Sequences {
		(&bm.Sequences[i]).WithCreatedMetadata(c)
	}
	return *bm
}

func (i ExternalBioMatInputWithSequences) Save(tx geltypes.Tx, samplingID geltypes.UUID) (occurrences []OccurrenceWithCategory, err error) {
	biomat, err := i.ExternalBioMatInput.Save(tx, samplingID)
	if err != nil {
		return nil, err
	}
	occurrences = append(occurrences, biomat.AsOccurrence())

	for j, sequence := range i.Sequences {
		sequence.SourceSample.SetValue(biomat.Code)
		sequence.UseSamplingCode(biomat.Sampling.Code)
		seq, err := sequence.Save(tx, samplingID)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("sequences")
		} else {
			occurrences = append(occurrences, seq.AsOccurrence())
		}
	}

	return
}
