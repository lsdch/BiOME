package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/stores"
)

type TaxonomyService struct {
	gbifClient *gbif.GBIFClient
}

func NewTaxonomyService(client *gbif.GBIFClient) *TaxonomyService {
	return &TaxonomyService{
		gbifClient: client,
	}
}

func (s *TaxonomyService) GetGBIFKingdoms(ctx context.Context, q db.Querier) ([]models.TaxonGBIF, error) {
	kingdoms, err := q.Queries().GetGBIFKingdoms(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.TaxonGBIF, len(kingdoms))
	for i, k := range kingdoms {
		result[i] = models.TaxonGBIF{
			Key:            k.Key,
			ScientificName: k.ScientificName,
			Name:           k.CanonicalName,
			Rank:           k.Rank,
			Status:         k.Status,
		}
	}
	return result, nil
}

func (s *TaxonomyService) FetchGBIFKingdoms(ctx context.Context) ([]models.TaxonGBIF, error) {
	gbifKingdoms, err := s.gbifClient.SearchSpecies(ctx, gbif.SearchParams{
		Status:     "ACCEPTED",
		Rank:       "KINGDOM",
		DatasetKey: "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c",
	})
	if err != nil {
		return nil, err
	}
	return gbifKingdoms.Results, nil
}

func (s *TaxonomyService) ListTaxa(ctx context.Context, q db.Querier, params models.ListTaxaParams) ([]models.Taxon, error) {
	store := stores.NewTaxonomyStore()
	return store.SearchTaxa(ctx, q, params)
}

func (s *TaxonomyService) GetTaxonByID(ctx context.Context, q db.Querier, taxonID uuid.UUID) (*models.Taxon, error) {
	t, err := q.Queries().GetTaxonByID(ctx, taxonID)
	if err != nil {
		return nil, err
	}
	return models.TaxonFromDB(&t), nil
}

func (s *TaxonomyService) GetTaxonByScientificName(ctx context.Context, q db.Querier, scientificName string) (*models.Taxon, error) {
	t, err := q.Queries().GetTaxonByScientificName(ctx, scientificName)
	if err != nil {
		return nil, err
	}
	return models.TaxonFromDB(&t), nil
}

func (s *TaxonomyService) LoadTaxonRelations(ctx context.Context, q db.Querier, taxon *models.Taxon) (*models.TaxonRelations, error) {
	taxonRelations := &models.TaxonRelations{}
	if parentID, ok := taxon.ParentID.Get(); ok {
		parent, err := s.GetTaxonByID(ctx, q, parentID)
		if err != nil {
			return nil, err
		}
		taxonRelations.ParentTaxon = models.NewOptionalFromPtr(parent)
	}
	if acceptedID, ok := taxon.AcceptedID.Get(); ok {
		accepted, err := s.GetTaxonByID(ctx, q, acceptedID)
		if err != nil {
			return nil, err
		}
		taxonRelations.AcceptedTaxon = models.NewOptionalFromPtr(accepted)
	}
	return taxonRelations, nil
}

func (s *TaxonomyService) GetTaxonWithRelations(ctx context.Context, q db.Querier, taxonID uuid.UUID) (*models.TaxonWithRelations, error) {
	taxon, err := s.GetTaxonByID(ctx, q, taxonID)
	if err != nil {
		return nil, err
	}
	taxonRelations, err := s.LoadTaxonRelations(ctx, q, taxon)
	if err != nil {
		return nil, err
	}
	return &models.TaxonWithRelations{
		Taxon:          *taxon,
		TaxonRelations: *taxonRelations,
	}, nil
}

func (s *TaxonomyService) LoadTaxonLineage(ctx context.Context, q db.Querier, taxon *models.Taxon) ([]models.Taxon, error) {
	lineage, err := q.Queries().GetTaxonLineage(ctx, taxon.ID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Taxon, len(lineage))
	for i, t := range lineage {
		result[i] = *models.TaxonFromDB(&t)
	}
	return result, nil
}

func (s *TaxonomyService) GetTaxonWithLineage(ctx context.Context, q db.Querier, taxonID uuid.UUID) (*models.TaxonWithLineage, error) {
	taxon, err := s.GetTaxonByID(ctx, q, taxonID)
	if err != nil {
		return nil, err
	}
	taxonRelations, err := s.LoadTaxonRelations(ctx, q, taxon)
	if err != nil {
		return nil, err
	}
	lineage, err := s.LoadTaxonLineage(ctx, q, taxon)
	if err != nil {
		return nil, err
	}
	return &models.TaxonWithLineage{
		Taxon:          *taxon,
		TaxonRelations: *taxonRelations,
		Lineage:        lineage,
	}, nil
}

func (s *TaxonomyService) GetTaxonWithFullLineage(ctx context.Context, q db.Querier, taxonID uuid.UUID) (*models.TaxonWithFullLineage, error) {
	taxonWithLineage, err := s.GetTaxonWithLineage(ctx, q, taxonID)
	if err != nil {
		return nil, err
	}
	descendants, err := s.GetTaxonDescendants(ctx, q, taxonID)
	if err != nil {
		return nil, err
	}
	return &models.TaxonWithFullLineage{
		TaxonWithLineage: *taxonWithLineage,
		Descendants:      descendants,
	}, nil
}

func (s *TaxonomyService) GetTaxonDescendants(ctx context.Context, q db.Querier, taxonID uuid.UUID) ([]models.Taxon, error) {
	descendants, err := q.Queries().GetTaxonDescendants(ctx, taxonID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Taxon, len(descendants))
	for i, t := range descendants {
		result[i] = *models.TaxonFromDB(&t)
	}
	return result, nil
}

func (s *TaxonomyService) GetTaxaByRank(ctx context.Context, q db.Querier, rank models.TaxonRank) ([]models.Taxon, error) {
	taxa, err := q.Queries().GetTaxaByRank(ctx, (biomedb.TaxonRank)(rank))
	if err != nil {
		return nil, err
	}
	result := make([]models.Taxon, len(taxa))
	for i, t := range taxa {
		result[i] = *models.TaxonFromDB(&t)
	}
	return result, nil
}

func (s *TaxonomyService) CreateTaxon(ctx context.Context, q db.Querier, input *models.CreateTaxonInput) (*models.Taxon, error) {
	taxon, err := q.Queries().InsertTaxon(ctx, *input.ToParams())
	if err != nil {
		return nil, err
	}
	return models.TaxonFromDB(&taxon), nil
}

// DeleteTaxon deletes a taxon by its ID.
//
// It fails if the taxon has any occurrences associated with it.
// Use DeleteTaxonWithOccurrences to delete a taxon and its occurrences.
func (s *TaxonomyService) DeleteTaxon(ctx context.Context, q db.Querier, taxonID uuid.UUID) error {
	err := q.Queries().DeleteTaxonByID(ctx, taxonID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTaxonWithOccurrences deletes a taxon, its descendants, and all associated occurrences from the database.
//
// If orphanSampling is true, it will also delete any sampling events that become orphaned as a result of deleting the occurrences.
func (s *TaxonomyService) DeleteTaxonWithOccurrences(ctx context.Context, tx *db.Tx, taxonID uuid.UUID, orphanSampling bool) error {
	err := tx.Queries().DeleteOccurrencesOfTaxonLineage(ctx, orphanSampling, taxonID)
	if err != nil {
		return err
	}
	err = tx.Queries().DeleteTaxonByID(ctx, taxonID)
	if err != nil {
		return err
	}
	return nil
}
