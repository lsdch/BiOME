package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
)

type TaxonomyService struct {
	db db.Querier
}

func NewTaxonomyService(db db.Querier) *TaxonomyService {
	return &TaxonomyService{
		db: db,
	}
}

func (s *TaxonomyService) GetTaxonByID(ctx context.Context, taxonID uuid.UUID) (*models.Taxon, error) {
	t, err := s.db.Queries().GetTaxonByID(ctx, taxonID)
	if err != nil {
		return nil, err
	}
	return models.TaxonFromDB(&t), nil
}

func (s *TaxonomyService) GetTaxonByScientificName(ctx context.Context, scientificName string) (*models.Taxon, error) {
	t, err := s.db.Queries().GetTaxonByScientificName(ctx, scientificName)
	if err != nil {
		return nil, err
	}
	return models.TaxonFromDB(&t), nil
}

func (s *TaxonomyService) LoadTaxonRelations(ctx context.Context, taxon *models.Taxon) (*models.TaxonRelations, error) {
	taxonRelations := &models.TaxonRelations{}
	if parentID, ok := taxon.ParentID.Get(); ok {
		parent, err := s.GetTaxonByID(ctx, parentID)
		if err != nil {
			return nil, err
		}
		taxonRelations.ParentTaxon = models.NewOptionalFromPtr(parent)
	}
	if acceptedID, ok := taxon.AcceptedID.Get(); ok {
		accepted, err := s.GetTaxonByID(ctx, acceptedID)
		if err != nil {
			return nil, err
		}
		taxonRelations.AcceptedTaxon = models.NewOptionalFromPtr(accepted)
	}
	return taxonRelations, nil
}

func (s *TaxonomyService) GetTaxonWithRelations(ctx context.Context, taxonID uuid.UUID) (*models.TaxonWithRelations, error) {
	taxon, err := s.GetTaxonByID(ctx, taxonID)
	if err != nil {
		return nil, err
	}
	taxonRelations, err := s.LoadTaxonRelations(ctx, taxon)
	if err != nil {
		return nil, err
	}
	return &models.TaxonWithRelations{
		Taxon:          *taxon,
		TaxonRelations: *taxonRelations,
	}, nil
}

func (s *TaxonomyService) LoadTaxonLineage(ctx context.Context, taxon *models.Taxon) ([]models.Taxon, error) {
	lineage, err := s.db.Queries().GetTaxonLineage(ctx, taxon.ID)
	if err != nil {
		return nil, err
	}
	result := make([]models.Taxon, len(lineage))
	for i, t := range lineage {
		result[i] = *models.TaxonFromDB(&t)
	}
	return result, nil
}

func (s *TaxonomyService) GetTaxonWithLineage(ctx context.Context, taxonID uuid.UUID) (*models.TaxonWithLineage, error) {
	taxon, err := s.GetTaxonByID(ctx, taxonID)
	if err != nil {
		return nil, err
	}
	taxonRelations, err := s.LoadTaxonRelations(ctx, taxon)
	if err != nil {
		return nil, err
	}
	lineage, err := s.LoadTaxonLineage(ctx, taxon)
	if err != nil {
		return nil, err
	}
	return &models.TaxonWithLineage{
		Taxon:          *taxon,
		TaxonRelations: *taxonRelations,
		Lineage:        lineage,
	}, nil
}

func (s *TaxonomyService) GetTaxaByRank(ctx context.Context, rank models.TaxonRank) ([]models.Taxon, error) {
	taxa, err := s.db.Queries().GetTaxaByRank(ctx, rank)
	if err != nil {
		return nil, err
	}
	result := make([]models.Taxon, len(taxa))
	for i, t := range taxa {
		result[i] = *models.TaxonFromDB(&t)
	}
	return result, nil
}

func (s *TaxonomyService) DeleteTaxon(ctx context.Context, taxonID uuid.UUID) error {
	err := s.db.Queries().DeleteTaxonByID(ctx, taxonID)
	if err != nil {
		return err
	}
	return nil
}
