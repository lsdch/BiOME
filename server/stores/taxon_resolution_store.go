package stores

import (
	"context"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
)

type TaxonResolutionStore struct {
}

func NewTaxonResolutionStore() *TaxonResolutionStore {
	return &TaxonResolutionStore{}
}

func (r *TaxonResolutionStore) InitTaxonResolution(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.TaxonResolution, error) {
	resolution, err := q.Queries().InitTaxonResolution(ctx, importID)
	if err != nil {
		return nil, err
	}
	return models.TaxonResolutionFromDBSlice(resolution), nil
}

func (r *TaxonResolutionStore) InsertGBIFBatch(ctx context.Context, q db.Querier, taxa []models.TaxonGBIF) (err error) {
	toInsert := make([]biomedb.InsertGBIFBatchParams, len(taxa))
	for i, taxon := range taxa {
		toInsert[i] = taxon.ToStaging()
	}

	_, err = q.Queries().InsertGBIFBatch(ctx, toInsert)
	return err
}

func (r *TaxonResolutionStore) InsertGBIFCandidatesBatch(ctx context.Context, q db.Querier, candidates map[string][]models.TaxonGBIF) (err error) {
	toInsert := make([]biomedb.InsertTaxonCandidatesBatchParams, 0, len(candidates))
	for inputName, matches := range candidates {
		for _, match := range matches {
			toInsert = append(toInsert, match.ToCandidate(inputName))
		}
	}

	_, err = q.Queries().InsertTaxonCandidatesBatch(ctx, toInsert)
	return err
}

func (r *TaxonResolutionStore) InsertTaxonStaging(ctx context.Context, q db.Querier, importID uuid.UUID, params models.TaxonStagingParams) (err error) {
	return q.Queries().InsertTaxonStaging(ctx, params.ToParams(importID))
}

func (r *TaxonResolutionStore) UpsertTaxonResolution(ctx context.Context, q db.Querier, params biomedb.UpsertTaxonResolutionParams) (err error) {
	return q.Queries().UpsertTaxonResolution(ctx, params)
}

// Build dependency list of GBIF keys that need to be fetched (parents + accepted references for synonyms),
// based on the current set of resolved taxa and their GBIF matches,
// to ensure that all necessary GBIF data is available before materialization.
//
// Returns the list of missing GBIF keys that need to be fetched.
func (r *TaxonResolutionStore) ListMissingGBIFDependencies(ctx context.Context, q db.Querier, importID uuid.UUID) (missingKeys []int32, err error) {
	if err = q.Queries().ExpandGBIFDependencies(ctx, importID); err != nil {
		return nil, err
	}
	missingKeys, err = q.Queries().ListMissingGBIFKeys(ctx, importID)
	return missingKeys, err
}

// Inserts taxa from GBIF staging into the main taxa table, for a given import and rank.
// This is done in two steps: first insert non-synonyms, then insert synonyms (which depend on the accepted taxa being present).
func (r *TaxonResolutionStore) InsertTaxaFromGBIF(ctx context.Context, q db.Querier, importID uuid.UUID, rank models.TaxonRank) (err error) {
	err = q.Queries().InsertTaxaFromGBIF(ctx, biomedb.InsertTaxaFromGBIFParams{ImportID: importID, Rank: string(rank), IsSynonym: false})
	if err != nil {
		return err
	}
	err = q.Queries().InsertTaxaFromGBIF(ctx, biomedb.InsertTaxaFromGBIFParams{ImportID: importID, Rank: string(rank), IsSynonym: true})
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) MarkTaxaNeedingGBIFCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	return q.Queries().MarkTaxaNeedingGBIFCandidates(ctx, importID)
}

func (r *TaxonResolutionStore) MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID, inputNames []string) (err error) {
	return q.Queries().MarkTaxaGBIFImportCompleted(ctx, importID, inputNames)
}

func (r *TaxonResolutionStore) ListTaxnamesToFetchFromGBIF(ctx context.Context, q db.Querier, importID uuid.UUID) (toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow, err error) {
	return q.Queries().ListTaxaToFetchGBIFCandidates(ctx, importID)
}

func (r *TaxonResolutionStore) GenerateLocalTaxonCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {

	err = q.Queries().CreateCandidateTaxaNameExact(ctx, importID)
	if err != nil {
		return err
	}

	err = q.Queries().CreateCandidateTaxaFuzzy(ctx, 0.6, importID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) ListTaxonCandidates(
	ctx context.Context, q db.Querier, importID uuid.UUID,
) (candidatesByName map[string][]models.TaxonCandidate, err error) {

	candidates, err := q.Queries().ListAllTaxonCandidates(ctx, importID)
	if err != nil {
		return nil, err
	}

	candidatesByName = make(map[string][]models.TaxonCandidate)
	for _, candidate := range candidates {
		candidatesByName[candidate.InputName] = append(candidatesByName[candidate.InputName], models.TaxonCandidateFromDB(candidate))
	}
	return candidatesByName, nil
}

func (r *TaxonResolutionStore) GetTaxonResolutionState(ctx context.Context, q db.Querier, importID uuid.UUID) (state *models.TaxonResolutionState, err error) {
	resolutions, err := q.Queries().GetTaxonResolution(ctx, importID)
	if err != nil {
		return nil, err
	}
	candidates, err := r.ListTaxonCandidates(ctx, q, importID)
	if err != nil {
		return nil, err
	}

	return &models.TaxonResolutionState{
		Resolution: models.TaxonResolutionFromDBSlice(resolutions),
		Candidates: candidates,
	}, nil
}
