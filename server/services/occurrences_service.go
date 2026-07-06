package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/stores"
	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
)

type OccurrencesService struct {
	samplings   *SamplingService
	datasets    *DatasetsService
	importBatch *ImportBatchService
	taxonomy    *TaxonomyService
}

func NewOccurrencesService(
	samplings *SamplingService,
	datasets *DatasetsService,
	importBatch *ImportBatchService,
	taxonomy *TaxonomyService,
) *OccurrencesService {
	return &OccurrencesService{
		samplings:   samplings,
		datasets:    datasets,
		importBatch: importBatch,
		taxonomy:    taxonomy,
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
		err := tx.Queries().AddArticleToOccurrence(ctx, occurrence.ID, article)
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
	store := stores.NewOccurrenceStore()
	occurrences, err := store.ListOccurrences(ctx, q, params)
	if err != nil {
		return models.PaginatedList[models.Occurrence]{}, err
	}
	totalCount, err := store.ListOccurrencesCount(ctx, q, params)
	if err != nil {
		return models.PaginatedList[models.Occurrence]{}, err
	}

	return models.PaginatedList[models.Occurrence]{Items: occurrences, TotalCount: totalCount}, nil
}

func (s *OccurrencesService) GetOccurrenceWithDetails(ctx context.Context, q db.Querier, occurrenceID ulid.ULID) (*models.OccurrenceWithDetails, error) {
	o, err := q.Queries().GetOccurrenceByID(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}

	var (
		samplingMetadata   models.SamplingMetadata
		occurrenceMetadata models.OccurrenceMetadata
	)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() (errG error) {
		samplingMetadata, errG = s.samplings.LoadSamplingMetadata(ctx, q, o.Sampling.ID)
		return errG
	})

	g.Go(func() (errG error) {
		occurrenceMetadata, errG = s.loadOccurrenceMetadata(ctx, q, occurrenceID)
		return errG
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	occurrenceWithDetails := models.OccurrenceFromDB(o.Occurrence, o.Taxon, o.Sampling, o.Country).WithDetails(samplingMetadata, occurrenceMetadata)
	return &occurrenceWithDetails, nil
}

func (s *OccurrencesService) GetOccurrencesGroupsH3(ctx context.Context, q db.Querier, params biomedb.OccurrencesGroupsH3Params) ([]biomedb.OccurrencesGroupsH3Row, error) {
	return q.Queries().OccurrencesGroupsH3(ctx, params)
}

func (s *OccurrencesService) loadOccurrenceDatasets(ctx context.Context, q db.Querier, occurrenceID ulid.ULID) ([]models.Dataset, error) {
	return s.datasets.LoadDatasetsForOccurrence(ctx, q, occurrenceID)
}

func (s *OccurrencesService) loadOccurrenceArticles(ctx context.Context, q db.Querier, occurrenceID ulid.ULID) ([]models.Article, error) {
	articles, err := q.Queries().GetOccurrenceArticles(ctx, occurrenceID)
	if err != nil {
		return nil, err
	}
	references := make([]models.Article, len(articles))
	for i, article := range articles {
		references[i] = models.ArticleFromDB(article)
	}
	return references, nil
}

func (s *OccurrencesService) loadOccurrenceCodeHistory(ctx context.Context, q db.Querier, occurrenceID ulid.ULID) ([]models.CodeHistoryEntry, error) {
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

func (s *OccurrencesService) loadOccurrenceMetadata(ctx context.Context, q db.Querier, occurrenceID ulid.ULID) (models.OccurrenceMetadata, error) {
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
		result[i] = models.OccurrenceOverviewItem{
			Name:        item.Name,
			ParentName:  models.NewOptionalFromPtr(item.ParentName),
			Authorship:  models.NewOptionalFromPtr(item.Authorship),
			Occurrences: item.OccurrencesCount,
			Rank:        item.Rank,
		}
	}
	return result, nil
}
