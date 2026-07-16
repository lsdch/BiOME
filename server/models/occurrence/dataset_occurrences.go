package occurrence

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/geldata/gel-go"
// 	"github.com/geldata/gel-go/geltypes"
// 	"github.com/lsdch/biome/db"
// 	"github.com/lsdch/biome/models/dataset"
// 	"github.com/lsdch/biome/models/references"
// 	gbif "github.com/lsdch/biome/models/taxonomy/GBIF"
// 	"github.com/oklog/ulid/v2"
// 	"github.com/sirupsen/logrus"
// )

// type OccurrenceDatasetListItem struct {
// 	dataset.Dataset `gel:"$inline"`
// 	Sites           int64 `gel:"sites_count" json:"sites"`
// 	Occurrences     int64 `gel:"occurrences_count" json:"occurrences"`
// 	IsCongruent     bool  `gel:"is_congruent" json:"is_congruent"`
// }

// func ListOccurrenceDatasets(db geltypes.Executor) ([]OccurrenceDatasetListItem, error) {
// 	datasets := []OccurrenceDatasetListItem{}
// 	err := db.Query(context.Background(),
// 		`#edgeql
// 			select datasets::OccurrenceDataset {
// 				*,
// 				maintainers: { *, user: { * } },
// 				meta: { * },
// 				sites_count := count(.sites),
// 				occurrences_count := count(.occurrences),
// 			}
// 		`,
// 		&datasets,
// 	)
// 	return datasets, err
// }

// type OccurrenceDataset struct {
// 	dataset.Dataset `gel:"$inline"`
// 	Sites           []SiteWithOccurrences `gel:"sites" json:"sites"`
// 	// Occurrences     []OccurrenceWithCategory `gel:"occurrences" json:"occurrences"`
// 	IsCongruent  bool                 `gel:"is_congruent" json:"is_congruent"`
// 	Contributors []string             `gel:"contributors" json:"contributors,omitempty"`
// 	Bibliography []references.Article `gel:"bibliography" json:"bibliography,omitempty"`
// }

// func GetOccurrenceDataset(db geltypes.Executor, slug string) (dataset OccurrenceDataset, err error) {
// 	err = db.QuerySingle(context.Background(),
// 		`#edgeql
// 			with module occurrence,
// 				dataset := (
// 					select datasets::OccurrenceDataset
// 					filter .slug = <str>$0
// 				),
// 			select dataset {
// 				*,
// 				meta: { * },
// 				maintainers: { *, user: { * }, organisations: { * } },
// 				contributors := distinct (.occurrences.identification.identified_by union .occurrences.sampling.performed_by),
// 				bibliography := distinct .occurrences.published_in { *, meta: { * } },
// 				sites: {
// 					*,
// 					country: { * },
// 					samplings: {
// 						id,
// 						date := .performed_on,
// 						occurrences := (
// 							select (dataset.occurrences) {
// 								id,
// 								code,
// 								identification: { confer, addendum, identified_on, taxon: { * } },
// 							} filter.sampling.id = dataset.sites.samplings.id
// 						)
// 					}
// 			}
// 		}
// 		`,
// 		&dataset, slug,
// 	)
// 	return
// }

// type CodeChange struct {
// 	ID          geltypes.UUID `gel:"id" json:"id" format:"uuid"`
// 	Code        string        `gel:"code" json:"code"`
// 	CodeHistory []struct {
// 		Code string    `gel:"code" json:"code"`
// 		Time time.Time `gel:"time" json:"time"`
// 	} `gel:"code_history" json:"code_history"`
// }

// func UpdateOccurrenceCodesInDataset(db geltypes.Executor, slug string) ([]CodeChange, error) {
// 	var changes = []CodeChange{}
// 	err := db.Query(context.Background(),
// 		`#edgeql
// 			with module occurrence,
// 			dataset := (
// 				select datasets::OccurrenceDataset
// 				filter .slug = <str>$0
// 			),
// 			update_pool := (update dataset.occurrences
// 			filter (
// 				with new_code := occurrence::occurrence_code(.identification.taxon, .sampling.code)
// 				select .code != new_code
// 			)
// 			 set {
// 				code_history := (
// 					(.code_history union (code := .code, time := datetime_of_statement()))
// 					if not .code in .code_history.code
// 					else .code_history
// 				)
// 			})
// 			select update_pool {
// 				id,
// 				code,
// 				code_history
// 			}
// 		`,
// 		&changes, slug,
// 	)
// 	return changes, err
// }

// type OccurrenceDatasetInput struct {
// 	dataset.DatasetInput
// 	OccurrenceBatchInput
// 	AddOccurrences []string  `json:"add_occurrences,omitempty" default:"[]"`
// 	BatchULID      types.ULID `json:"batch_ulid,omitempty"`
// 	Kingdom        string    `json:"kingdom,omitempty" doc:"This is used to discriminate between homonymous taxa from different kingdoms. For example, if the dataset only contains occurrences of plants, the taxonomic scope can be set to 'Plantae' to avoid fetching homonymous animal taxa from GBIF."`
// }

// func (i *OccurrenceDatasetInput) FetchMissingTaxa(tx geltypes.Executor) (notFound []string, err error) {
// 	missingTaxa, err := i.OccurrenceBatchInput.ListMissingTaxa(tx)
// 	if err != nil || len(missingTaxa) == 0 {
// 		return
// 	}

// 	logrus.Infof("Fetching %d missing GBIF taxa for dataset '%s'", len(missingTaxa), i.Label)
// 	client := gbif.Client()
// 	var kingdomKey int32
// 	if i.Kingdom != "" {
// 		logrus.Warnf("Taxonomic scope is set to '%s' for dataset '%s', fetching kingdom key from GBIF", i.Kingdom, i.Label)
// 		kingdom, err := gbif.FetchByName(client, gbif.SearchParams{
// 			Query: i.Kingdom,
// 			Rank:  "KINGDOM",
// 			Limit: 1,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("Failed to fetch kingdom '%s' from GBIF: %v", i.Kingdom, err)
// 		}
// 		kingdomKey = kingdom.Key
// 	}
// 	gbifTaxa, err := gbif.FetchNames(client, missingTaxa,
// 		// 0 value will be ignored by the search params encoder
// 		gbif.WithHigherTaxonKey(kingdomKey),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logrus.Infof("Fetching parents and synonyms for up to %d GBIF taxa for dataset '%s'", gbifTaxa.Taxa.Count(), i.Label)
// 	if err := gbifTaxa.Taxa.FetchParentsAndSynonyms(tx); err != nil {
// 		return nil, err
// 	}

// 	logrus.Infof("Persisting up to %d GBIF taxa for dataset '%s'", gbifTaxa.Taxa.Count(), i.Label)
// 	if err := gbifTaxa.Taxa.Persist(tx); err != nil {
// 		return nil, err
// 	}

// 	if len(gbifTaxa.NotFound) > 0 {
// 		logrus.Warnf("%d taxa were not found in GBIF for dataset '%s'", len(gbifTaxa.NotFound), i.Label)
// 		return gbifTaxa.NotFound, fmt.Errorf("%d missing taxa were not found in GBIF for dataset '%s'", len(gbifTaxa.NotFound), i.Label)
// 	}

// 	return nil, nil
// }

// func (i *OccurrenceDatasetInput) GenerateBatchULID() *OccurrenceDatasetInput {
// 	i.BatchULID = ulid.Make()
// 	return i
// }

// func (i *OccurrenceDatasetInput) SetTracker(tracker OccurrenceBatchTracker) *OccurrenceDatasetInput {
// 	i.OccurrenceBatchInput.Tracker = tracker
// 	return i
// }

// func (i *OccurrenceDatasetInput) CheckAdditionalOccurrenceCodes(db geltypes.Executor) error {
// 	if len(i.AddOccurrences) == 0 {
// 		return nil
// 	}
// 	var missing []string
// 	err := db.Query(context.Background(),
// 		`#edgeql
// 			with module occurrence,
// 			codes := <str>array_unpack(<array<str>>$0),
// 			existing := (select Occurrence filter .code in codes),
// 			missing := codes except existing.code
// 			select missing`,
// 		&missing, i.AddOccurrences)
// 	if err != nil {
// 		return err
// 	}
// 	if len(missing) > 0 {
// 		return fmt.Errorf("the following occurrence codes are not in the database: %v", missing)
// 	}
// 	logrus.Infof("Adding %d existing occurrences to dataset %s", len(i.AddOccurrences), i.Label)
// 	return nil
// }

// func (i *OccurrenceDatasetInput) GatherBatch(db geltypes.Executor) error {
// 	if err := i.CheckAdditionalOccurrenceCodes(db); err != nil {
// 		return err
// 	}
// 	if i.BatchULID.IsZero() {
// 		return nil
// 	}
// 	return db.Execute(context.Background(),
// 		`#edgeql
// 			with
// 			import_id := <str>$0,
// 			occurrences := (
// 				select occurrence::Occurrence
// 				filter .meta.batch_import_id = import_id
// 				or .code in <str>array_unpack(<optional array<str>>$1)
// 			),
// 			dataset := (
// 				assert_single((
// 					select datasets::OccurrenceDataset
// 					filter .meta.batch_import_id = import_id
// 				),
// 				message := "Multiple datasets found with the same batch_import_id")
// 			),
// 			update occurrences
// 			set {
// 				datasets := distinct (.datasets union dataset)
// 			}
// 		`, i.BatchULID.String(), i.AddOccurrences,
// 	)
// }

// func (dataset *OccurrenceDatasetInput) Save(client *gel.Client, batchSize int) (created dataset.Dataset, err error) {
// 	err = db.WithBatchMode(client, dataset.BatchULID).Tx(
// 		context.Background(),
// 		func(ctx context.Context, tx geltypes.Tx) error {

// 			created, err = dataset.DatasetInput.Save(tx)
// 			if err != nil {
// 				return err
// 			}

// 			err = dataset.OccurrenceBatchInput.WithBatchSize(batchSize).Save(tx)
// 			if err != nil {
// 				return err
// 			}

// 			return dataset.GatherBatch(tx)
// 		},
// 	)
// 	return
// }

// func (dataset OccurrenceDatasetInput) SaveParallel(client *gel.Client, batchSize int, cores int) (created dataset.Dataset, err error) {

// 	notFound, err := dataset.FetchMissingTaxa(client)
// 	if err != nil {
// 		if len(notFound) > 0 {
// 			return created, fmt.Errorf("The following %d taxa were not found in GBIF for dataset '%s':\n%v", len(notFound), dataset.Label, notFound)
// 		} else {
// 			return created, fmt.Errorf("Failed to fetch missing taxa for dataset '%s': %v", dataset.Label, err)
// 		}
// 	}

// 	dataset.GenerateBatchULID()
// 	client = db.WithBatchMode(client, dataset.BatchULID)
// 	created, err = dataset.DatasetInput.Save(client)
// 	if err != nil {
// 		return created, fmt.Errorf("failed to initialize dataset '%s': %v", dataset.Label, err)
// 	}
// 	logrus.Debugf("Initialized dataset '%s'", dataset.Label)

// 	logrus.Infof("Fetching DOIs for dataset '%s'", dataset.Label)
// 	err = dataset.FetchMissingDOIs(client)
// 	if err != nil {
// 		logrus.Errorf("failed to fetch missing DOIs for dataset '%s': %v", dataset.Label, err)
// 		created.RollbackImport(client)
// 		return
// 	}

// 	if err = dataset.OccurrenceBatchInput.WithBatchSize(batchSize).SaveParallel(client, cores); err != nil {
// 		logrus.Errorf("%v", err)
// 		created.RollbackImport(client)
// 		return
// 	}
// 	logrus.Infof("Gathering batch in dataset")
// 	err = dataset.GatherBatch(client)
// 	if err != nil {
// 		logrus.Errorf("Failed to gather batch: %v", err)
// 		created.RollbackImport(client)
// 		return
// 	}
// 	logrus.Infof("Done")
// 	return

// }
