package occurrence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/occurrence/queries"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

type OccurrenceBatchMetadataInputs struct {
	Organisations map[string]people.OrganisationInput   `json:"organisations,omitempty"`
	People        map[string]people.PersonInput         `json:"people,omitempty"`
	DataSources   map[string]references.DataSourceInput `json:"data_sources,omitempty"`
	Collections   map[string]references.CollectionInput `json:"collections,omitempty"`
	Taxa          []taxonomy.TaxonInput                 `json:"taxa,omitempty"`
	Bibliography  map[string]references.ArticleInput    `json:"bibliography,omitempty"`
}

func (i OccurrenceBatchMetadataInputs) TaxaByName() map[string]taxonomy.TaxonInput {
	taxa := make(map[string]taxonomy.TaxonInput)
	for _, taxon := range i.Taxa {
		taxa[taxon.Name] = taxon
	}
	return taxa
}

type CreatedMetadata struct {
	Organisations map[string]string `json:"organisations,omitempty"` // input string to code map
	People        map[string]string `json:"people,omitempty"`        // input string to alias map
	DataSources   map[string]string `json:"data_sources,omitempty"`  // input string to code map
	Bibliography  map[string]string `json:"bibliography,omitempty"`  // input string to code map
	Collections   map[string]string `json:"collections,omitempty"`   // input string to code map
}

func (i OccurrenceBatchMetadataInputs) saveNewCollections(tx geltypes.Tx, trackers ...OccurrenceBatchTracker) (map[string]string, error) {

	codes := make(map[string]string)
	for rawColl, coll := range i.Collections {
		created, err := coll.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawColl)
		}
		codes[rawColl] = created.Code
		for _, tracker := range trackers {
			tracker.Progress(1)
		}
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewDataSources(tx geltypes.Tx, trackers ...OccurrenceBatchTracker) (map[string]string, error) {

	codes := make(map[string]string)
	for rawSource, source := range i.DataSources {
		created, err := source.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawSource)
		}
		codes[rawSource] = created.Code
		for _, tracker := range trackers {
			tracker.Progress(1)
		}
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewBibliography(tx geltypes.Tx, trackers ...OccurrenceBatchTracker) (map[string]string, error) {

	codes := make(map[string]string)
	for rawRef, ref := range i.Bibliography {
		created, err := ref.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawRef)
		}
		codes[rawRef] = created.Code
		for _, tracker := range trackers {
			tracker.Progress(1)
		}
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewOrganisations(tx geltypes.Tx, trackers ...OccurrenceBatchTracker) (map[string]string, error) {
	codes := make(map[string]string)
	for rawOrg, org := range i.Organisations {
		logrus.Debugf("Creating organisation '%s' %+v", rawOrg, org)
		created, err := org.Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawOrg)
		}
		codes[rawOrg] = created.Code
		for _, tracker := range trackers {
			tracker.Progress(1)
		}
	}
	return codes, nil
}

func (i OccurrenceBatchMetadataInputs) saveNewPersons(tx geltypes.Tx, orgCodes map[string]string, trackers ...OccurrenceBatchTracker) (map[string]string, error) {
	personAliases := make(map[string]string)
	for rawPerson, person := range i.People {
		created, err := person.WithOrganisationCodes(orgCodes).Save(tx)
		if err != nil {
			return nil, models.WrapErrorPath(err, rawPerson)
		}
		personAliases[rawPerson] = created.Alias
		for _, tracker := range trackers {
			tracker.Progress(1)
		}
	}
	return personAliases, nil
}

func (i OccurrenceBatchMetadataInputs) Save(tx geltypes.Tx, trackers ...OccurrenceBatchTracker) (*CreatedMetadata, error) {

	total := len(i.Taxa) + len(i.DataSources) + len(i.Bibliography) + len(i.Organisations) + len(i.People)
	for _, tracker := range trackers {
		if total > 0 {
			tracker.Start(total).SetDescription("Saving associated metadata")
		} else {
			tracker.StartUnknownTotal("Saving associated metadata")
		}
		defer tracker.Finish()
	}

	for _, tracker := range trackers {
		tracker.SetDetail(fmt.Sprintf("Saving %d taxa", len(i.Taxa)))
	}

	for j, taxon := range i.Taxa {
		if _, err := taxon.Save(tx); err != nil {
			logrus.Errorf("Failed to save taxon: %+v", taxon)
			return nil, models.WrapErrorIndex(err, j).PrependPath("taxa")
		}
		for _, tracker := range trackers {
			tracker.Progress(1)
		}
	}

	for _, tracker := range trackers {
		tracker.SetDetail(fmt.Sprintf("Saving %d data sources", len(i.DataSources)))
	}

	dataSources, err := i.saveNewDataSources(tx, trackers...)
	if err != nil {
		return nil, models.WrapErrorPath(err, "data_sources")
	}

	for _, tracker := range trackers {
		tracker.SetDetail(fmt.Sprintf("Saving %d data sources", len(i.Collections)))
	}

	collections, err := i.saveNewCollections(tx, trackers...)
	if err != nil {
		return nil, models.WrapErrorPath(err, "collections")
	}

	for _, tracker := range trackers {
		tracker.SetDetail(fmt.Sprintf("Saving %d bibliography entries", len(i.Bibliography)))
	}

	bibliography, err := i.saveNewBibliography(tx, trackers...)
	if err != nil {
		return nil, models.WrapErrorPath(err, "bibliography")
	}

	for _, tracker := range trackers {
		tracker.SetDetail(fmt.Sprintf("Saving %d organisations", len(i.Organisations)))
	}

	organisations, err := i.saveNewOrganisations(tx, trackers...)
	if err != nil {
		return nil, models.WrapErrorPath(err, "organisations")
	}

	for _, tracker := range trackers {
		tracker.SetDetail(fmt.Sprintf("Saving %d people", len(i.People)))
	}

	personAliases, err := i.saveNewPersons(tx, organisations, trackers...)
	if err != nil {
		return nil, models.WrapErrorPath(err, "people")
	}

	for _, tracker := range trackers {
		tracker.Finish()
	}

	return &CreatedMetadata{
		Organisations: organisations,
		People:        personAliases,
		DataSources:   dataSources,
		Bibliography:  bibliography,
		Collections:   collections,
	}, nil
}

// OccurrenceBatchInput is the input type for registering occurrences in bulk,
// including all the necessary upstream data:
// site, events, sampling.
// Occurrences can be registered in bulk, with multiple events and samplings.
type OccurrenceBatchInput struct {
	OccurrenceBatchMetadataInputs `json:",inline"`
	Content                       []SiteOccurrenceInput  `json:"content"`
	BatchSize                     int                    `json:"batch_size,omitempty"`
	Tracker                       OccurrenceBatchTracker `json:"-"`
}

func (i *OccurrenceBatchInput) ensureTracker() *OccurrenceBatchInput {
	if i.Tracker == nil {
		i.Tracker = &NoOpBatchTracker{}
	}
	return i
}

func (batch *OccurrenceBatchInput) WithBatchSize(size int) *OccurrenceBatchInput {
	batch.BatchSize = size
	return batch
}

// TaxaSet returns a set of all taxa names present in the batch
func (batch *OccurrenceBatchInput) TaxaSet() mapset.Set[string] {
	taxa := mapset.NewSet[string]()
	for _, siteWithOccurrences := range batch.Content {
		for _, sampling := range siteWithOccurrences.Samplings {
			for _, t := range sampling.TargetTaxa {
				taxa.Add(t)
			}
			for _, occ := range sampling.Occurrences {
				taxa.Add(occ.Identification.Taxon)
			}
		}
	}
	return taxa
}

func (batch *OccurrenceBatchInput) ListMissingTaxa(tx geltypes.Tx) (missing []string, err error) {
	ocurringTaxa := batch.TaxaSet()
	providedTaxa := batch.OccurrenceBatchMetadataInputs.TaxaByName()
	taxaToCheck := ocurringTaxa.Difference(mapset.NewSetFromMapKeys(providedTaxa))
	for _, taxon := range providedTaxa {
		if _, parentProvided := providedTaxa[taxon.Parent]; !parentProvided {
			taxaToCheck.Add(taxon.Parent)
		}
	}
	if taxaToCheck.Cardinality() == 0 {
		return nil, nil
	}

	return taxonomy.CheckMissingTaxa(tx, taxaToCheck.ToSlice())
}

func (batch OccurrenceBatchInput) SaveSites(tx geltypes.Tx) error {

	batch.ensureTracker()
	batch.Tracker.Start(len(batch.Content)).SetDescription("Saving sites").SetDetail("")
	defer batch.Tracker.Finish()

	wg := &sync.WaitGroup{}
	errorChan := make(chan error, 1)

	for i := 0; i < len(batch.Content); i += batch.BatchSize {
		if len(errorChan) > 0 {
			return <-errorChan
		}
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			endIndex := min(start+batch.BatchSize, len(batch.Content))
			batch.Tracker.SetDetail(fmt.Sprintf("Starting goroutine for site batch [%d, %d]", start, endIndex))

			subBatch := batch.Content[start:endIndex]
			data, _ := json.Marshal(subBatch)
			err := tx.Execute(context.Background(),
				`#edgeql
					with data := <json>$0,
					for site_data in json_array_unpack(data) union (
						location::insert_site(site_data)
					)
				`, data)
			if err != nil {
				if len(errorChan) > 0 {
					return
				}
				errorChan <- models.WrapErrorIndex(err, start).PrependPath("content")
			}
			batch.Tracker.Progress(len(subBatch))
		}(i)
	}
	wg.Wait()
	if len(errorChan) > 0 {
		return <-errorChan
	}
	return nil
}

func (batch *OccurrenceBatchInput) Save(client geltypes.Tx) (err error) {

	// occurrencesChan := make(chan []BaseOccurrence)

	batch.ensureTracker()

	batch.Tracker.StartUnknownTotal("Saving associated metadata")
	replacements, err := batch.OccurrenceBatchMetadataInputs.Save(client, batch.Tracker)
	if err != nil {
		return err
	}

	batch.Tracker.SetDetail("Checking for missing taxa")
	missingTaxa, err := batch.ListMissingTaxa(client)
	if err != nil {
		return err
	}
	if len(missingTaxa) > 0 {
		err = os.WriteFile("missing_taxa.txt", []byte(strings.Join(missingTaxa, "\n")), 0644)
		if err != nil {
			logrus.Errorf("Failed to write missing taxa to file: %v", err)
		}
		return models.WrapErrorPath(fmt.Errorf("the following taxa are missing: %v.\nPlease add missing taxa definitions in the 'taxa' field of your occurrence batch input", missingTaxa), "content")
	}
	batch.Tracker.Finish()

	if err := batch.SaveSites(client); err != nil {
		return models.WrapErrorPath(err, "content")
	}

	batch.Tracker.Start(len(batch.Content)).SetDescription("Saving occurrences")
	defer batch.Tracker.Finish()

	// errorChan := make(chan error, 1)
	for j, siteOccurrence := range batch.Content {
		// if len(errorChan) > 0 {
		// 	err = <-errorChan
		// 	break
		// }
		batch.Tracker.SetDetail(fmt.Sprintf("Site %s", siteOccurrence.Code))

		siteOccurrence.WithCreatedMetadata(replacements)
		if err := siteOccurrence.SaveAbiotics(client); err != nil {
			return models.WrapErrorIndex(err, j).PrependPath("content")
		}

		for k, sampling := range siteOccurrence.Samplings {
			err := sampling.Save(client, siteOccurrence.Code)
			if err != nil {
				return models.WrapErrorIndex(err, k).PrependPath("samplings").PrependIndex(j).PrependPath("content")
			}
		}
		batch.Tracker.Progress(1)
	}
	return
}

func (batch *OccurrenceBatchInput) SaveParallel(client *gel.Client, cores int) error {

	batch.ensureTracker()

	replacements, err := batch.OccurrenceBatchMetadataInputs.Save(client, batch.Tracker)
	if err != nil {
		return err
	}

	batch.Tracker.SetDetail("Checking for missing taxa")
	missingTaxa, err := batch.ListMissingTaxa(client)
	if err != nil {
		return err
	}
	if len(missingTaxa) > 0 {
		err = os.WriteFile("missing_taxa.txt", []byte(strings.Join(missingTaxa, "\n")), 0644)
		if err != nil {
			logrus.Errorf("Failed to write missing taxa to file: %v", err)
		}
		return models.WrapErrorPath(fmt.Errorf("the following taxa are missing: %v.\nPlease add missing taxa definitions in the 'taxa' field of your occurrence batch input", missingTaxa), "content")
	}
	batch.Tracker.Finish()

	if err := batch.SaveSites(client); err != nil {
		return models.WrapErrorPath(err, "content")
	}

	logrus.Infof("Saving abiotics")
	for i, siteOccurrence := range batch.Content {
		if err := siteOccurrence.SaveAbiotics(client); err != nil {
			return models.WrapErrorIndex(err, i).PrependPath("content")
		}
	}

	errorChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	mainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	txError := client.
		WithConfig(map[string]any{
			"apply_access_policies": false,
		}).
		// WithTxOptions(gelcfg.NewTxOptions().WithIsolation(gelcfg.RepeatableRead)).
		Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
			batch.Tracker.Start(len(batch.Content)).SetDescription("Saving occurrences").SetDetail("")
			defer batch.Tracker.Finish()
			for i := 0; i*batch.BatchSize < len(batch.Content); i++ {
				wg.Add(1)
				go func(coreID int) {
					defer wg.Done()
					startIndex, endIndex := coreID*batch.BatchSize, min((coreID+1)*batch.BatchSize, len(batch.Content))
					for j, siteOccurrence := range batch.Content[startIndex:endIndex] {
						select {
						case <-mainCtx.Done():
							return
						default:
						}
						siteOccurrence.WithCreatedMetadata(replacements)

						for k, sampling := range siteOccurrence.Samplings {
							if err := sampling.Save(tx, siteOccurrence.Code); err != nil {
								if len(errorChan) > 0 {
									return
								}
								errorChan <- models.WrapErrorIndex(fmt.Errorf("%v [%+v]", err, sampling), k).PrependPath("samplings").PrependIndex(j).PrependPath("content")
								cancel() // stop all other goroutines
								return
							}
						}
						batch.Tracker.Progress(1)
					}
				}(i)
			}
			wg.Wait()
			select {
			case err := <-errorChan:
				return err
			default:
				return nil
			}
		})

	return txError
}

func (batch *OccurrenceBatchInput) SaveOccurrencesBatch(db geltypes.Executor) (err error) {

	batch.ensureTracker()

	batchSize := batch.BatchSize
	if batchSize == 0 {
		batchSize = len(batch.Content)
	}

	insertQuery := fmt.Sprintf(
		`#edgeql
		with data := <json>$0,
		for site_data in json_array_unpack(data) union (
			with
				site := location::insert_site(site_data),
				abiotic_measurements := (
					for abiotic_data in json_array_unpack(json_get(site_data, 'abiotic_measurements'))
					union (
						%s # insert abiotic measurement query
					)
				),
				samplings := (
					for sampling_data in <json>json_array_unpack(json_get(site_data, 'samplings'))
					union (
						with
							samp := (
								%s
							),
							occurrences := (
								for occ_data in json_array_unpack(json_get(<json>sampling_data, 'occurrences'))
								union (
									%s # insert occurrence query
								)
							),
						select samp
					)
				),
				select site
			)
		`,

		queries.AbioticQuery("site", "abiotic_data", ""),
		queries.SamplingQuery("site", "sampling_data", ""),
		queries.OccurrenceQuery("samp", "occ_data", ""),
	)

	// logrus.Infof("Query:  \n%s", insertQuery)
	batch.Tracker.Start(len(batch.Content)).SetDescription("Saving occurrences")
	defer batch.Tracker.Finish()

	batch.Tracker.SetDetail("Processing metadata dependencies")
	replacements, err := batch.OccurrenceBatchMetadataInputs.Save(db)
	if err != nil {
		return err
	}

	for _, siteOccurrence := range batch.Content {
		siteOccurrence.WithCreatedMetadata(replacements)
	}

	missingTaxa, err := batch.ListMissingTaxa(db)
	if err != nil {
		return err
	}
	if len(missingTaxa) > 0 {
		err = os.WriteFile("missing_taxa.txt", []byte(strings.Join(missingTaxa, "\n")), 0644)
		if err != nil {
			logrus.Errorf("Failed to write missing taxa to file: %v", err)
		}
		return fmt.Errorf("the following taxa are missing: %v.\nPlease add missing taxa definitions in the 'taxa' field of your occurrence batch input", missingTaxa)
	}

	for i := 0; i < len(batch.Content); i += batch.BatchSize {
		batch.Tracker.SetDetail(fmt.Sprintf("Inserting batch %d", i))
		endIndex := min(i+batchSize, len(batch.Content))
		actualSize := endIndex - i
		data, _ := json.Marshal(batch.Content[i:endIndex])
		// logrus.Infof("%v", string(data))
		err = db.Execute(context.Background(), insertQuery, data)
		if err != nil {
			return models.WrapErrorPath(err, "content")
		}
		logrus.Infof("Progress: %d", actualSize)
		batch.Tracker.Progress(actualSize)
	}
	return nil
}

/*
SiteOccurrenceInput is the input type for registering a site and its occurrences in bulk.
It includes the site data and a list of samplings and abiotic measurements.
*/
type SiteOccurrenceInput struct {
	SiteInput           `json:",inline"`
	Samplings           []SamplingInputWithOccurrences `json:"samplings,omitempty"`
	AbioticMeasurements []AbioticMeasurementInput      `json:"abiotic_measurements,omitempty"`
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
// Each sampling can have multiple occurrences.
// It also includes spottings and abiotic measurements.

type SamplingInputWithOccurrences struct {
	SamplingInput `json:",inline"`
	Occurrences   []OccurrenceInput `json:"occurrences,omitempty"`
}

func (s *SamplingInputWithOccurrences) WithCreatedMetadata(c *CreatedMetadata) *SamplingInputWithOccurrences {
	for i := range s.Occurrences {
		(&s.Occurrences[i]).WithCreatedMetadata(c)
	}
	return s
}

func (i SamplingInputWithOccurrences) Save(tx geltypes.Tx, siteCode string) (err error) {

	sampling, err := i.SamplingInput.QuickSave(tx, siteCode)
	if err != nil {
		return err
	}

	// Save external occurrences and their sequences
	for j, occurrence := range i.Occurrences {
		err := occurrence.SaveExecute(tx, sampling.Number)
		if err != nil {
			return models.WrapErrorIndex(err, j).PrependPath("occurrences")
		}
	}

	return
}
