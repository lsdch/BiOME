package occurrence

import (
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
)

// OccurrenceBatchInput is the input type for registering occurrences in bulk,
// including all the necessary upstream data:
// site, events, sampling.
// Occurrences can be registered in bulk, with multiple events and samplings.
// Occurrences types include: BioMaterial (internal/external) and external sequences.
type OccurrenceBatchInput []SiteOccurrenceInput

func (i OccurrenceBatchInput) Save(tx geltypes.Tx) (occurrences []OccurrenceWithCategory, err error) {
	for i, siteOccurrence := range i {
		occ, err := siteOccurrence.Save(tx)
		if err != nil {
			return nil, models.WrapErrorIndex(err, i)
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
