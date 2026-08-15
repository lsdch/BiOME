package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/stores"
	"github.com/lsdch/biome/types"
	"github.com/sirupsen/logrus"
	"github.com/uber/h3-go/v4"
	"golang.org/x/sync/errgroup"
)

type OccurrencesService struct {
	samplings   *SamplingService
	datasets    *DatasetsService
	importBatch *ImportBatchService
	taxonomy    *TaxonomyService
	store       *stores.OccurrenceStore
}

func NewOccurrencesService(
	samplings *SamplingService,
	datasets *DatasetsService,
	importBatch *ImportBatchService,
	taxonomy *TaxonomyService,
	store *stores.OccurrenceStore,
) *OccurrencesService {
	return &OccurrencesService{
		samplings:   samplings,
		datasets:    datasets,
		importBatch: importBatch,
		taxonomy:    taxonomy,
		store:       store,
	}
}

func (s *OccurrencesService) GenerateCode(taxon models.Taxon, sampling models.Sampling) string {
	return fmt.Sprintf("%s[%s]", taxon.Code(), sampling.Code())
}

func (s *OccurrencesService) CreateOccurrence(ctx context.Context, tx *db.Tx, input models.FullOccurrenceInput) (*models.OccurrenceWithDetails, error) {
	sampling, err := s.samplings.CreateSampling(ctx, tx, input.Sampling)
	if err != nil {
		return nil, err
	}

	created, err := s.CreateOccurrenceAtSampling(ctx, tx, sampling.ID, input.Occurrence)
	if err != nil {
		return nil, err
	}

	occurrence := created.WithSampling(*sampling)
	return &occurrence, nil
}

func (s *OccurrencesService) CreateOccurrenceAtSampling(ctx context.Context, tx *db.Tx, samplingID uuid.UUID, input models.OccurrenceInput) (*models.OccurrenceWithMetadata, error) {
	taxon, err := s.taxonomy.GetTaxonByID(ctx, tx, input.Identification.TaxonID)
	if err != nil {
		return nil, err
	}

	sampling, err := s.samplings.GetSampling(ctx, tx, samplingID)
	if err != nil {
		return nil, err
	}

	code := s.GenerateCode(*taxon, *sampling)

	created, err := tx.Queries().AddOccurrenceToSampling(ctx, input.ToParams(samplingID, code))
	if err != nil {
		return nil, err
	}

	occurrence := models.BaseOccurrenceFromDB(created.Occurrence, created.Taxon)
	metadata := models.OccurrenceMetadata{}

	for _, collection := range input.Collections {
		col, err := s.AddCollection(ctx, tx, occurrence.ID, collection)
		if err != nil {
			return nil, err
		}
		metadata.Collections = append(metadata.Collections, *col)
	}

	for _, article := range input.PublishedIn {
		err := tx.Queries().AddPublicationToOccurrence(ctx, occurrence.ID, article)
		if err != nil {
			return nil, err
		}
	}

	metadata.References, err = s.loadOccurrenceArticles(ctx, tx, occurrence.ID)
	if err != nil {
		return nil, err
	}

	for _, datasetID := range input.Datasets {
		err := s.datasets.AddOccurrenceToDataset(ctx, tx, datasetID, occurrence.ID)
		if err != nil {
			return nil, err
		}
	}

	metadata.Datasets, err = s.loadOccurrenceDatasets(ctx, tx, occurrence.ID)
	if err != nil {
		return nil, err
	}

	occurrenceWithMetadata := occurrence.WithMetadata(metadata)
	return &occurrenceWithMetadata, nil
}

func (s *OccurrencesService) ListOccurrences(ctx context.Context, q db.Querier, params stores.ListOccurrencesParams) (models.PaginatedList[models.Occurrence], error) {
	occurrences, err := s.store.ListOccurrences(ctx, q, params)
	if err != nil {
		return models.PaginatedList[models.Occurrence]{}, err
	}
	totalCount, err := s.store.ListOccurrencesCount(ctx, q, params)
	if err != nil {
		return models.PaginatedList[models.Occurrence]{}, err
	}

	return models.PaginatedList[models.Occurrence]{Items: occurrences, TotalCount: totalCount}, nil
}

func (s *OccurrencesService) GetOccurrenceWithDetails(ctx context.Context, q db.Querier, occurrenceID types.ULID) (*models.OccurrenceWithDetails, error) {
	o, err := q.Queries().GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		fmt.Printf("error type: %T\n", err)
		fmt.Printf("error: %v\n", err)

		var pgErr *pgconn.PgError
		fmt.Printf("is PgError: %v\n", errors.As(err, &pgErr))
		fmt.Printf("is pgx.ErrNoRows: %v\n", errors.Is(err, pgx.ErrNoRows))
		fmt.Printf("is sql.ErrNoRows: %v\n", errors.Is(err, sql.ErrNoRows))
		return nil, err
	}

	var (
		samplingMetadata   models.SamplingMetadata
		occurrenceMetadata models.OccurrenceMetadata
	)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (errG error) {
		samplingMetadata, errG = s.samplings.LoadSamplingMetadata(ctx, q, o.SamplingsWithCountry.ID)
		return errG
	})

	g.Go(func() (errG error) {
		occurrenceMetadata, errG = s.loadOccurrenceMetadata(ctx, q, occurrenceID)
		return errG
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	occurrenceWithDetails := models.OccurrenceFromDB(o.Occurrence, o.Taxon, o.SamplingsWithCountry).
		WithDetails(samplingMetadata, occurrenceMetadata)
	return &occurrenceWithDetails, nil
}

func (s *OccurrencesService) GetOccurrencesGroupsH3(ctx context.Context, q db.Querier, params biomedb.OccurrencesGroupsH3Params) ([]biomedb.OccurrencesGroupsH3Row, error) {
	return q.Queries().OccurrencesGroupsH3(ctx, params)
}

func (s *OccurrencesService) loadOccurrenceDatasets(ctx context.Context, q db.Querier, occurrenceID types.ULID) ([]models.Dataset, error) {
	return s.datasets.LoadDatasetsForOccurrence(ctx, q, occurrenceID)
}

func (s *OccurrencesService) loadOccurrenceArticles(ctx context.Context, q db.Querier, occurrenceID types.ULID) ([]models.Publication, error) {
	articles, err := q.Queries().GetOccurrencePublications(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	references := make([]models.Publication, len(articles))
	for i, article := range articles {
		references[i] = models.PublicationFromDB(article)
	}
	return references, nil
}

func (s *OccurrencesService) loadOccurrenceCodeHistory(ctx context.Context, q db.Querier, occurrenceID types.ULID) ([]models.CodeHistoryEntry, error) {
	codeHistory, err := q.Queries().GetOccurrenceCodeHistory(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	codeHistoryEntries := make([]models.CodeHistoryEntry, len(codeHistory))
	for i, entry := range codeHistory {
		codeHistoryEntries[i] = models.CodeHistoryEntryFromDB(entry)
	}
	return codeHistoryEntries, nil
}

func (s *OccurrencesService) loadOccurrenceMetadata(ctx context.Context, q db.Querier, occurrenceID types.ULID) (models.OccurrenceMetadata, error) {
	codeHistory, err := s.loadOccurrenceCodeHistory(ctx, q, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	references, err := s.loadOccurrenceArticles(ctx, q, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	collections, err := s.loadOccurrenceCollections(ctx, q, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	datasets, err := s.loadOccurrenceDatasets(ctx, q, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	importBatch, err := s.importBatch.GetImportBatchForOccurrence(ctx, q, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	return models.NewOccurrenceMetadata(codeHistory, references, datasets, collections, importBatch), nil
}

func (s *OccurrencesService) OccurrencesByTaxaOverview(ctx context.Context, q db.Querier) ([]models.OccurrenceOverviewItem, error) {
	overview, err := q.Queries().OccurrencesByTaxaOverview(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.OccurrenceOverviewItem, len(overview))
	for i, item := range overview {
		result[i] = models.OccurrenceOverviewItemFromDB(item)
	}
	return result, nil
}

func (s *OccurrencesService) ListOccurringTaxaAtCell(ctx context.Context, q db.Querier, cell h3.Cell, resolution int64, params stores.ListOccurrencesParams) ([]models.OccurrenceOverviewItem, error) {
	return s.store.ListOccurringTaxaAtCell(ctx, q, cell, resolution, params)
}

// MaterializeOccurrences triggers the materialization of occurrences for a given import batch.
// This is typically used after an import batch has been processed to ensure that the occurrences are properly stored and indexed.
// Occurrence-related metadata, such as collections and references, are also materialized during this process.
//
// RefreshOccurrenceCodes should be called at some point after materialization to ensure that occurrence codes are up-to-date.
func (s *OccurrencesService) MaterializeOccurrences(ctx context.Context, tx *db.Tx, importID uuid.UUID) error {
	logrus.Infof("Materializing occurrences for import batch ID %s", importID)
	if err := tx.Queries().GenerateCodesStaging(ctx, importID); err != nil {
		return err
	}
	missingCodes, err := tx.Queries().CheckStagingCodesGenerated(ctx, importID)
	if err != nil {
		return fmt.Errorf("failed to check occurrence codes generation: %v", err)
	}
	if len(missingCodes) > 0 {
		return fmt.Errorf("found %d occurrences without generated codes", len(missingCodes))
	}
	if err := tx.Queries().MaterializeOccurrences(ctx, importID); err != nil {
		return err
	}
	if err := tx.Queries().MaterializeCollections(ctx, importID); err != nil {
		return err
	}
	return nil
}

func (s *OccurrencesService) RefreshOccurrenceCodes(ctx context.Context, q db.Querier) error {
	return q.Queries().RefreshOccurrenceCodes(ctx)
}

func (s *OccurrencesService) ListOccurrencesAtProximity(ctx context.Context, q db.Querier, params models.ListSamplingsAtProximityInput) ([]*models.SamplingWithOccurrencesAndDistance, error) {
	samplings, err := s.samplings.ListSamplingsAtProximity(ctx, q, params)
	if err != nil {
		return nil, err
	}
	samplingIDs := make([]uuid.UUID, len(samplings))
	for i, sampling := range samplings {
		samplingIDs[i] = sampling.ID
	}
	occurrences, err := q.Queries().LoadOccurrencesForSamplings(ctx, samplingIDs)
	if err != nil {
		return nil, err
	}
	// Then bind occurrences to samplings
	samplingMap := make(map[uuid.UUID]*models.SamplingWithOccurrencesAndDistance)
	for _, row := range samplings {
		so := row.WithOccurrences([]models.BaseOccurrence{})
		samplingMap[row.ID] = &so
	}
	for _, occurrence := range occurrences {
		occ := models.BaseOccurrenceFromDB(occurrence.Occurrence, occurrence.Taxon)
		if sampling, ok := samplingMap[occurrence.Occurrence.SamplingID]; ok {
			sampling.Occurrences = append(sampling.Occurrences, occ)
		}
	}
	return slices.Collect(maps.Values(samplingMap)), nil

}

func (s *OccurrencesService) ListOccurrencesH3(ctx context.Context, q db.Querier, resolution int64, params stores.ListOccurrencesParams) ([]models.H3CellWithRichness, error) {
	return s.store.ListOccurrencesH3(ctx, q, resolution, params)
}

func (s *OccurrencesService) ListSamplingsH3(ctx context.Context, q db.Querier, resolution int64, params stores.ListSamplingsParams) ([]models.H3CellWithRichness, error) {
	return s.store.ListSamplingsH3(ctx, q, resolution, params)
}

func (s *OccurrencesService) ListSamplingsWithOccurrences(ctx context.Context, q db.Querier, params stores.ListOccurrencesParams) ([]models.SamplingWithOccurrences, error) {
	samplings, err := s.store.ListSamplings(ctx, q, params)
	if err != nil {
		return nil, err
	}
	occurrences, err := s.store.ListBaseOccurrences(ctx, q, params)
	if err != nil {
		return nil, err
	}
	return s.assembleSamplingsWithOccurrences(samplings, occurrences), nil
}

func (s *OccurrencesService) ListSamplingsWithOccurrencesAtCell(ctx context.Context, q db.Querier, cell h3.Cell, resolution int64, params stores.ListOccurrencesParams) ([]models.SamplingWithOccurrences, error) {
	samplings, err := s.store.ListSamplingsAtCell(ctx, q, cell, resolution, params)
	if err != nil {
		return nil, err
	}
	occurrences, err := s.store.ListOccurrencesAtCell(ctx, q, cell, resolution, params)
	if err != nil {
		return nil, err
	}
	return s.assembleSamplingsWithOccurrences(samplings, occurrences), nil
}

func (s *OccurrencesService) assembleSamplingsWithOccurrences(
	samplings []models.Sampling,
	occurrences []models.BaseOccurrenceWithSamplingID,
) []models.SamplingWithOccurrences {
	byID := make(map[uuid.UUID]*models.SamplingWithOccurrences, len(samplings))
	result := make([]models.SamplingWithOccurrences, len(samplings))

	for i, sampling := range samplings {
		result[i] = models.SamplingWithOccurrences{
			Sampling:    sampling,
			Occurrences: nil,
		}
		byID[sampling.ID] = &result[i]
	}

	for _, occurrence := range occurrences {
		if sampling, ok := byID[occurrence.SamplingID]; ok {
			sampling.Occurrences = append(sampling.Occurrences, occurrence.BaseOccurrence)
		}
	}

	return result
}
