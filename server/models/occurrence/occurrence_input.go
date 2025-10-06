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
			for _, internalBiomat := range sampling.InternalBiomat {
				taxa.Add(internalBiomat.Identification.Taxon)
			}
			for _, externalBiomat := range sampling.ExternalBiomat {
				taxa.Add(externalBiomat.Identification.Taxon)
			}
			for _, sequence := range sampling.Sequences {
				taxa.Add(sequence.Identification.Taxon)
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

// func (i OccurrenceBatchInput) SaveBulk(tx geltypes.Tx) (occurrences []OccurrenceWithCategory, err error) {
// 	replacements, err := i.OccurrenceBatchMetadataInputs.Save(tx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	missingTaxa, err := i.ListMissingTaxa(tx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(missingTaxa) > 0 {
// 		err = os.WriteFile("missing_taxa.txt", []byte(strings.Join(missingTaxa, "\n")), 0644)
// 		if err != nil {
// 			logrus.Errorf("Failed to write missing taxa to file: %v", err)
// 		}
// 		return nil, models.WrapErrorPath(fmt.Errorf("the following taxa are missing: %v.\nPlease add missing taxa definitions in the 'taxa' field of your occurrence batch input", missingTaxa), "occurrences")
// 	}

// 	for _, siteOccurrence := range i.Occurrences {
// 		siteOccurrence.WithCreatedMetadata(replacements)
// 	}

// 	data, _ := json.Marshal(i.Occurrences)
// 	err = tx.Query(context.Background(),
// 		`#edgeql
// 		with data := <json>$0,
// 		for site_data in json_array_unpack(data) union (
// 			with
// 				site := location::insert_site(site_data),
// 				# events
// 			select (
// 				for event_data in json_array_unpack(site_data['events']) union (
// 					with
// 						event := (
// 							insert events::Event {
// 								site := site,
// 								performed_by := (
// 									select people::Person
// 									filter .alias in <str>json_array_unpack(json_get(data, 'performed_by'))
// 								),
// 								performed_by_groups := (
// 									select people::Organisation
// 									filter .code in <str>json_array_unpack(json_get(data,'performed_by_groups'))
// 								),
// 								performed_on := (
// 									select date::from_json_with_precision(data['performed_on'])
// 								),
// 								spottings := (
// 									select taxonomy::Taxon
// 									filter .name in <str>json_array_unpack(json_get(data, 'spottings'))
// 								),
// 							}
// 						),
// 						abiotic_measurements := (
// 							for measurement_data in json_array_unpack(json_get(event_data, 'abiotic_measurements')) union (
// 								with
// 									param := (select events::AbioticParameter filter .code = <str>measurement_data['param']),
// 									value := <float64>measurement_data['value']
// 								insert events::AbioticMeasurement {
// 									event := event,
// 									param := param,
// 									value := value
// 								} unless conflict on ((.event, .param)) else (
// 									update events::AbioticMeasurement set {
// 										param := param,
// 										value := value
// 									}
// 								)
// 							)
// 						),
// 					select (
// 						for sampling_data in json_array_unpack(json_get(event_data, 'samplings')) union (
// 							with
// 								sampling := (
// 										with module events
// 										insert events::Sampling {
// 											event := event,
// 											methods := (
// 												select events::SamplingMethod
// 												filter .code in <str>json_array_unpack(json_get(data, 'methods'))
// 											),
// 											fixatives := (
// 												select samples::Fixative
// 												filter .code in <str>json_array_unpack(json_get(data, 'fixatives'))
// 											),
// 											sampling_target := <events::SamplingTarget>(data['target']['kind']),
// 											target_taxa := (
// 												select taxonomy::Taxon
// 												filter .name in <str>json_array_unpack(json_get(data, 'target', 'taxa'))
// 											),
// 											sampling_duration := <int32>json_get(data, 'duration'),
// 											comments := <str>json_get(data, 'comments'),
// 											habitats := (
// 												select sampling::Habitat
// 												filter .label in <str>json_array_unpack(json_get(data, 'habitats'))
// 											),
// 											access_points := (<str>json_array_unpack(json_get(data, 'access_points')))
// 										}
// 								),
// 								internal_biomats := (
// 									with module occurrence
// 									for biomat in json_array_unpack(json_get(event_data, 'internal_biomats')) union (
// 										with module occurrence,
// 											identification := data['identification'],
// 											taxon := taxonomy::taxonByName(<str>identification['taxon']),
// 											publications := json_array_unpack(json_get(data, 'published_in')),
// 										insert InternalBioMat {
// 											sampling := sampling,
// 											code := <str>json_get(data, 'code') ?? occurrence::biomat_code(taxon, sampling),
// 											identification := (
// 												insert occurrence::Identification {
// 													taxon := taxon,
// 													identified_by := people::personByAlias(<str>identification['identified_by']),
// 													identified_on := date::from_json_with_precision(identification['identified_on']),
// 												}
// 											),
// 											is_type := <bool>json_get(data, 'is_type') ?? false,
// 											comments := <str>json_get(data, 'comments'),
// 											published_in := (select distinct
// 												(for p in publications union (
// 													select references::Article {
// 														@original_source := <bool>json_get(p, 'original')
// 													} filter .code = <str>p['code']
// 												))
// 											)
// 										}
// 									)
// 								),
// 								external_biomats_and_seqs := (
// 									for biomat in json_array_unpack(json_get(event_data, 'external_biomats')) union (
// 										with
// 											ext_biomat := (
// 												with module occurrence,
// 													identification := data['identification'],
// 													taxon := taxonomy::taxonByName(<str>identification['taxon']),
// 													publications := json_array_unpack(json_get(data, 'published_in')),
// 												insert ExternalBioMat {
// 													sampling := sampling,
// 													code := <str>json_get(data, 'code') ?? occurrence::biomat_code(taxon, sampling),
// 													original_source := (
// 														with src := <str>json_get(data, 'original_source')
// 														select (if exists src then default::get_vocabulary(src)[is references::DataSource] else <references::DataSource>{})
// 													),
// 													original_link := <str>json_get(data, 'original_link'),
// 													quantity := <occurrence::QuantityType>json_get(data, 'quantity'),
// 													content_description := <str>json_get(data, 'content_description'),
// 													in_collection := <str>json_get(data, 'collection'),
// 													item_vouchers := <str>json_array_unpack(json_get(data, 'item_vouchers')),
// 													comments := <str>json_get(data, 'comments'),
// 													published_in := (select distinct
// 														(for p in publications union (
// 															select references::Article {
// 																@original_source := <bool>json_get(p, 'original')
// 															} filter .code = <str>p['code']
// 														))
// 													),
// 													identification := (
// 														insert occurrence::Identification {
// 															taxon := taxon,
// 															identified_by := people::personByAlias(<str>identification['identified_by']),
// 															identified_on := date::from_json_with_precision(identification['identified_on']),
// 														}
// 													),
// 													is_type := <bool>json_get(data, 'is_type') ?? false,
// 												}
// 											),
// 											sequences := (
// 												with module seq
// 												for seq in json_array_unpack(json_get(biomat, 'sequences')) union (
// 													insert ExternalSequence {
// 														sampling := sampling,
// 														code := <str>data['code'],
// 														label := <str>json_get(data, 'label'),
// 														sequence := <str>json_get(data, 'sequence'),
// 														gene := seq::geneByCode(<str>data['gene']),
// 														legacy := <tuple<id: int32, code: str, alignment_code: str>>json_get(data, 'legacy'),
// 														origin := <seq::ExtSeqOrigin>json_get(data, 'origin'),
// 														published_in := (
// 															with pubs := json_array_unpack(json_get(data, 'published_in'))
// 															select distinct (
// 																for p in pubs union (
// 																	select references::Article {
// 																		@original_source := <bool>json_get(p, 'original')
// 																	} filter .code = <str>p['code']
// 																)
// 															)
// 														),
// 														identification := (
// 															with identification := data['identification']
// 															insert occurrence::Identification {
// 																identified_by := people::personByAlias(<str>identification['identified_by']),
// 																identified_on := date::from_json_with_precision(identification['identified_on']),
// 																taxon := taxonomy::taxonByName(<str>identification['taxon']),
// 															}
// 														),
// 														referenced_in := (
// 															for ref in json_array_unpack(json_get(data, 'referenced_in'))
// 															insert references::SeqReference {
// 																db := references::dataSourceByCode(<str>ref['db']),
// 																accession := <str>ref['accession'],
// 																is_origin := <bool>json_get(ref, 'is_origin'),
// 															}
// 														),
// 														specimen_identifier := <str>json_get(data, 'specimen_identifier'),
// 														original_taxon := <str>json_get(data, 'original_taxon'),
// 														source_sample := (
// 															with source_sample := <str>json_get(data, 'source_sample')
// 															select if exists source_sample
// 															then occurrence::externalBiomatByCode(source_sample)
// 															else <occurrence::ExternalBioMat>{}
// 														)
// 													}
// 												)
// 											),
// 										select ext_biomat[is occurrence::Occurrence] union sequences[is occurrence::Occurrence]
// 									)
// 								)
// 							select internal_biomats[is occurrence::Occurrence] union external_biomats_and_seqs
// 						)
// 					)
// 				)
// 			)
// 		)
// 	`, &occurrences, data)
// 	return occurrences, err
// }

func (i OccurrenceBatchInput) SaveSites(tx geltypes.Tx) error {
	data, _ := json.Marshal(i.Occurrences)
	logrus.Infof("Marshalling done")
	return tx.Execute(context.Background(),
		`#edgeql
			with data := <json>$0,
			for site_data in json_array_unpack(data) union (
				location::insert_site(site_data)
			)
		`, data)
}

func (i OccurrenceBatchInput) Save(tx geltypes.Tx) (occurrences []OccurrenceWithCategory, err error) {

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

func (i *SiteOccurrenceInput) WithCreatedMetadata(c *CreatedMetadata) *SiteOccurrenceInput {
	for j := range i.Samplings {
		i.Samplings[j].WithCreatedMetadata(c)
	}
	return i
}

func (i SiteOccurrenceInput) Save(tx geltypes.Tx) ([]OccurrenceWithCategory, error) {
	site, err := i.SiteInput.Save(tx)
	if err != nil {
		return nil, err
	}
	occurrences := []OccurrenceWithCategory{}
	for j, sampling := range i.Samplings {
		occ, err := sampling.Save(tx, site.Code)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("samplings")
		} else {
			occurrences = append(occurrences, occ...)
		}
	}

	for j, abioticMeasurement := range i.AbioticMeasurements {
		if err := site.AddAbioticMeasurement(tx, abioticMeasurement); err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("abiotic_measurements")
		}
	}

	return occurrences, nil
}

// EventInputWithActions is the input type for registering an event and its occurrences in bulk.
// It includes the event data and a list of samplings.
// Each sampling can have multiple internal and external biomaterials, and sequences.
// It also includes spottings and abiotic measurements.

type SamplingInputWithOccurrences struct {
	SamplingInput  `json:",inline"`
	InternalBiomat []InternalBioMatInput              `json:"internal_biomats"`
	ExternalBiomat []ExternalBioMatInputWithSequences `json:"external_biomats"`
	Sequences      []ExternalSequenceInput            `json:"sequences"`
}

func (s *SamplingInputWithOccurrences) WithCreatedMetadata(c *CreatedMetadata) *SamplingInputWithOccurrences {
	s.ActionInput.WithPersonAliases(c.People)
	for i := range s.InternalBiomat {
		(&s.InternalBiomat[i]).WithCreatedMetadata(c)
	}
	for i := range s.ExternalBiomat {
		(&s.ExternalBiomat[i]).WithCreatedMetadata(c)
	}
	for i := range s.Sequences {
		(&s.Sequences[i]).WithCreatedMetadata(c)
	}
	return s
}

func (i SamplingInputWithOccurrences) Save(tx geltypes.Tx, siteCode string) (occurrences []OccurrenceWithCategory, err error) {

	sampling, err := i.SamplingInput.Save(tx, siteCode)
	if err != nil {
		return nil, err
	}

	for j, internalBiomat := range i.InternalBiomat {
		biomat, err := internalBiomat.Save(tx, sampling.Number)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("internal_biomats")
		} else {
			occurrences = append(occurrences, biomat.AsOccurrence())
		}
	}

	for j, externalBiomat := range i.ExternalBiomat {
		occ, err := externalBiomat.Save(tx, sampling.Number)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("external_biomats")
		} else {
			occurrences = append(occurrences, occ...)
		}
	}

	for j, sequence := range i.Sequences {
		sequence.UseSamplingCode(sampling.Code(siteCode))
		seq, err := sequence.Save(tx, sampling.Number)
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

func (bm *ExternalBioMatInputWithSequences) WithCreatedMetadata(c *CreatedMetadata) *ExternalBioMatInputWithSequences {
	(&bm.ExternalBioMatInput).WithCreatedMetadata(c)
	for i := range bm.Sequences {
		(&bm.Sequences[i]).WithCreatedMetadata(c)
	}
	return bm
}

func (i ExternalBioMatInputWithSequences) Save(tx geltypes.Tx, samplingNumber int64) (occurrences []OccurrenceWithCategory, err error) {
	biomat, err := i.ExternalBioMatInput.Save(tx, samplingNumber)
	if err != nil {
		return nil, err
	}
	occurrences = append(occurrences, biomat.AsOccurrence())

	for j, sequence := range i.Sequences {
		sequence.SourceSample.SetValue(biomat.Code)
		sequence.UseSamplingCode(biomat.Sampling.Code(biomat.Sampling.Site.Code))
		seq, err := sequence.Save(tx, samplingNumber)
		if err != nil {
			return nil, models.WrapErrorIndex(err, j).PrependPath("sequences")
		} else {
			occurrences = append(occurrences, seq.AsOccurrence())
		}
	}

	return
}
