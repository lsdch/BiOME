package services

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/stores"
	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
)

type OccurrencesService struct {
	db          db.Querier
	samplings   *SamplingService
	datasets    *DatasetsService
	importBatch *ImportBatchService
}

func NewOccurrencesService(db db.Querier, samplings *SamplingService, datasets *DatasetsService, importBatch *ImportBatchService) *OccurrencesService {
	return &OccurrencesService{
		db:          db,
		samplings:   samplings,
		datasets:    datasets,
		importBatch: importBatch,
	}
}

func (s *OccurrencesService) ListOccurrences(ctx context.Context, params stores.ListOccurrencesParams) ([]models.Occurrence, error) {
	store := stores.NewOccurrenceStore(s.db)
	return store.ListOccurrences(ctx, params)
}

func (s *OccurrencesService) GetOccurrenceWithDetails(ctx context.Context, occurrenceID ulid.ULID) (*models.OccurrenceWithDetails, error) {
	o, err := s.db.Queries().GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	var (
		samplingMetadata   models.SamplingMetadata
		occurrenceMetadata models.OccurrenceMetadata
	)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (errG error) {
		samplingMetadata, errG = s.samplings.LoadSamplingMetadata(ctx, o.Sampling.ID)
		return errG
	})

	g.Go(func() (errG error) {
		occurrenceMetadata, errG = s.loadOccurrenceMetadata(ctx, occurrenceID)
		return errG
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	occurrenceWithDetails := models.OccurrenceFromDB(
		o.Occurrence, o.Taxon, o.Sampling, o.Country,
	).WithDetails(samplingMetadata, occurrenceMetadata)
	return &occurrenceWithDetails, nil
}

func (s *OccurrencesService) GetOccurrencesGroupsH3(
	ctx context.Context,
	params biomedb.OccurrencesGroupsH3Params,
) ([]biomedb.OccurrencesGroupsH3Row, error) {
	return s.db.Queries().OccurrencesGroupsH3(ctx, params)
}

func (s *OccurrencesService) loadOccurrenceDatasets(ctx context.Context, occurrenceID ulid.ULID) ([]models.Dataset, error) {
	return s.datasets.LoadDatasetsForOccurrence(ctx, occurrenceID)
}

func (s *OccurrencesService) loadOccurrenceArticles(ctx context.Context, occurrenceID ulid.ULID) ([]models.Article, error) {
	articles, err := s.db.Queries().GetOccurrenceArticles(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	references := make([]models.Article, len(articles))
	for i, article := range articles {
		references[i] = models.ArticleFromDB(article)
	}
	return references, nil
}

func (s *OccurrencesService) loadOccurrenceCodeHistory(ctx context.Context, occurrenceID ulid.ULID) ([]models.CodeHistoryEntry, error) {
	codeHistory, err := s.db.Queries().GetOccurrenceCodeHistory(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	codeHistoryEntries := make([]models.CodeHistoryEntry, len(codeHistory))
	for i, entry := range codeHistory {
		codeHistoryEntries[i] = models.CodeHistoryEntryFromDB(entry)
	}
	return codeHistoryEntries, nil
}

func (s *OccurrencesService) loadOccurrenceCollections(ctx context.Context, occurrenceID ulid.ULID) ([]models.Collection, error) {
	collections, err := s.db.Queries().GetOccurrenceCollections(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Collection, len(collections))
	for i, c := range collections {
		result[i] = models.CollectionFromDB(c)
	}
	return result, nil
}

func (s *OccurrencesService) loadOccurrenceMetadata(ctx context.Context, occurrenceID ulid.ULID) (models.OccurrenceMetadata, error) {
	codeHistory, err := s.loadOccurrenceCodeHistory(ctx, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	references, err := s.loadOccurrenceArticles(ctx, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	collections, err := s.loadOccurrenceCollections(ctx, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	datasets, err := s.loadOccurrenceDatasets(ctx, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	importBatch, err := s.importBatch.GetImportBatchForOccurrence(ctx, occurrenceID)
	if err != nil {
		return models.OccurrenceMetadata{}, err
	}

	return models.NewOccurrenceMetadata(codeHistory, references, datasets, collections, importBatch), nil
}
