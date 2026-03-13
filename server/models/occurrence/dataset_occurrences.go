package occurrence

import (
	"context"
	"fmt"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/references"
	"github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

type OccurrenceDatasetListItem struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           int64 `gel:"sites_count" json:"sites"`
	Occurrences     int64 `gel:"occurrences_count" json:"occurrences"`
	IsCongruent     bool  `gel:"is_congruent" json:"is_congruent"`
}

func ListOccurrenceDatasets(db geltypes.Executor) ([]OccurrenceDatasetListItem, error) {
	datasets := []OccurrenceDatasetListItem{}
	err := db.Query(context.Background(),
		`#edgeql
			select datasets::OccurrenceDataset {
				*,
				maintainers: { *, user: { * } },
				meta: { * },
				sites_count := count(.sites),
				occurrences_count := count(.occurrences),
			}
		`,
		&datasets,
	)
	return datasets, err
}

type OccurrenceDataset struct {
	dataset.Dataset `gel:"$inline" json:",inline"`
	Sites           []SiteWithOccurrences `gel:"sites" json:"sites"`
	// Occurrences     []OccurrenceWithCategory `gel:"occurrences" json:"occurrences"`
	IsCongruent  bool                 `gel:"is_congruent" json:"is_congruent"`
	Contributors []string             `gel:"contributors" json:"contributors,omitempty"`
	Bibliography []references.Article `gel:"bibliography" json:"bibliography,omitempty"`
}

func GetOccurrenceDataset(db geltypes.Executor, slug string) (dataset OccurrenceDataset, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			with module occurrence,
				dataset := (
					select datasets::OccurrenceDataset
					filter .slug = <str>$0
				),
			select dataset {
				*,
				meta: { * },
				maintainers: { *, user: { * }, organisations: { * } },
				contributors := distinct (.occurrences.identification.identified_by union .occurrences.sampling.performed_by),
				bibliography := distinct .occurrences.published_in { *, meta: { * } },
				sites: {
					*,
					country: { * },
					samplings: {
						id,
						date := .performed_on,
						occurring_taxa: { * },
						occurrences := (
							select (dataset.occurrences) {
								id,
								code,
								identification: { confer, addendum, identified_on, taxon: { * } },
							} filter.sampling.id = dataset.sites.samplings.id
						)
					}
			}
		}
		`,
		&dataset, slug,
	)
	return
}

type OccurrenceDatasetInput struct {
	dataset.DatasetInput `json:",inline"`
	OccurrenceBatchInput `json:",inline"`
	AddOccurrences       []string  `json:"add_occurrences,omitempty"`
	BatchULID            ulid.ULID `json:"batch_ulid,omitempty"`
}

func (i *OccurrenceDatasetInput) GenerateBatchULID() *OccurrenceDatasetInput {
	i.BatchULID = ulid.Make()
	return i
}

func (i *OccurrenceDatasetInput) SetTracker(tracker OccurrenceBatchTracker) *OccurrenceDatasetInput {
	i.OccurrenceBatchInput.Tracker = tracker
	return i
}

func (i *OccurrenceDatasetInput) CheckAdditionalOccurrenceCodes(db geltypes.Executor) error {
	if len(i.AddOccurrences) == 0 {
		return nil
	}
	var missing []string
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
			codes := <str>array_unpack(<array<str>>$0),
			existing := (select Occurrence filter .code in codes),
			missing := codes except existing.code
			select missing`,
		&missing, i.AddOccurrences)
	if err != nil {
		return err
	}
	if len(missing) > 0 {
		return fmt.Errorf("the following occurrence codes are not in the database: %v", missing)
	}
	logrus.Infof("Adding %d existing occurrences to dataset %s", len(i.AddOccurrences), i.Label)
	return nil
}

func (i *OccurrenceDatasetInput) GatherBatch(db geltypes.Executor) error {
	if err := i.CheckAdditionalOccurrenceCodes(db); err != nil {
		return err
	}
	if i.BatchULID.IsZero() {
		return nil
	}
	return db.Execute(context.Background(),
		`#edgeql
			with
			import_id := <str>$0,
			occurrences := (
				select occurrence::Occurrence
				filter .meta.batch_import_id = import_id
				or .code in <str>array_unpack(<optional array<str>>$1)
			),
			dataset := (
				assert_single((
					select datasets::OccurrenceDataset
					filter .meta.batch_import_id = import_id
				),
				message := "Multiple datasets found with the same batch_import_id")
			),
			update occurrences
			set {
				datasets := distinct (.datasets union dataset)
			}
		`, i.BatchULID.String(), i.AddOccurrences,
	)
}

func (dataset *OccurrenceDatasetInput) SaveBulk(client *gel.Client, batchSize int) (created dataset.Dataset, err error) {
	err = db.WithBatchMode(client, dataset.BatchULID).Tx(
		context.Background(),
		func(ctx context.Context, tx geltypes.Tx) error {

			created, err = dataset.DatasetInput.Save(tx)
			if err != nil {
				return err
			}

			err = dataset.OccurrenceBatchInput.WithBatchSize(batchSize).Save(tx)
			if err != nil {
				return err
			}

			return dataset.GatherBatch(tx)
		},
	)
	return
}

func (dataset OccurrenceDatasetInput) SaveParallel(client *gel.Client, batchSize int, cores int) (created dataset.Dataset, err error) {
	dataset.GenerateBatchULID()
	client = db.WithBatchMode(client, dataset.BatchULID)
	created, err = dataset.DatasetInput.Save(client)
	if err != nil {
		return created, fmt.Errorf("failed to initialize dataset '%s': %v", dataset.Label, err)
	}
	logrus.Debugf("Initialized dataset '%s'", dataset.Label)

	if err = dataset.OccurrenceBatchInput.WithBatchSize(batchSize).SaveParallel(client, cores); err != nil {
		logrus.Errorf("%v", err)
		created.RollbackImport(client)
		return
	}
	logrus.Infof("Gathering batch in dataset")
	err = dataset.GatherBatch(client)
	if err != nil {
		logrus.Errorf("Failed to gather batch: %v", err)
		created.RollbackImport(client)
		return
	}
	logrus.Infof("Done")
	return

}
