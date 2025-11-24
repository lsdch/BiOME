package occurrence

import (
	"context"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/dataset"
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
	IsCongruent bool `gel:"is_congruent" json:"is_congruent"`
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

func (i *OccurrenceDatasetInput) GatherBatch(db geltypes.Executor) error {
	if i.BatchULID.IsZero() {
		return nil
	}
	return db.Execute(context.Background(),
		`#edgeql
			with
			dataset := (
				assert_single((
					select datasets::OccurrenceDataset
					filter .meta.batch_import_id = <str>$0
				),
				message := "Multiple datasets found with the same batch_import_id")
			),
			update occurrence::Occurrence
			filter .meta.batch_import_id = <str>$0
			set {
				datasets := dataset
			}
		`, i.BatchULID.String(),
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
		return
	}

	if err = dataset.OccurrenceBatchInput.WithBatchSize(batchSize).SaveParallel(client, cores); err != nil {
		logrus.Infof("Rolling back dataset")
		created.RollbackImport(client)
		return
	}
	logrus.Infof("Gathering batch in dataset")
	err = dataset.GatherBatch(client)
	if err != nil {
		logrus.Fatalf("Failed to gather batch: %v", err)
	}
	logrus.Infof("Done")
	return

}
