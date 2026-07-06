package stores

import (
	"context"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
)

type TaxonResolutionStore struct {
}

func NewTaxonResolutionStore() *TaxonResolutionStore {
	return &TaxonResolutionStore{}
}

func (r *TaxonResolutionStore) InitTaxonResolution(ctx context.Context, q db.Querier, importHash string) ([]models.TaxonResolution, error) {
	resolution, err := q.Queries().InitTaxonResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}
	return models.TaxonResolutionFromDBSlice(resolution), nil
}

func (r *TaxonResolutionStore) ClaimGBIFImport(ctx context.Context, q db.Querier, importHash string) (ok bool, err error) {
	rowsAffected, err := q.Queries().ClaimGBIFImport(ctx, importHash)
	return rowsAffected > 0, err
}

func (r *TaxonResolutionStore) CompleteGBIFImport(ctx context.Context, q db.Querier, importHash string) (status models.GBIFImportStatus, err error) {
	return q.Queries().CompleteGBIFImport(ctx, importHash)
}

func (r *TaxonResolutionStore) FailGBIFImport(ctx context.Context, q db.Querier, importHash string) (status models.GBIFImportStatus, err error) {
	return q.Queries().FailGBIFImport(ctx, importHash)
}

func (r *TaxonResolutionStore) InitGBIFCandidatesProgress(ctx context.Context, q db.Querier, total int32, importHash string) error {
	return q.Queries().InitGBIFCandidatesProgress(ctx, total, importHash)
}
func (r *TaxonResolutionStore) IncrementGBIFCandidatesProgress(ctx context.Context, q db.Querier, importHash string) (biomedb.IncrementGBIFCandidatesProgressRow, error) {
	return q.Queries().IncrementGBIFCandidatesProgress(ctx, importHash)
}

func (r *TaxonResolutionStore) GetGBIFImportStatus(ctx context.Context, q db.Querier, importHash string) (state models.GBIFImportState, err error) {
	res, err := q.Queries().GetGBIFImportStatus(ctx, importHash)
	if err != nil {
		return models.GBIFImportState{}, err
	}
	return models.GBIFImportStateFromDB(res), nil
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

func (r *TaxonResolutionStore) InsertTaxonStaging(ctx context.Context, q db.Querier, importHash string, params models.TaxonStagingParams) (err error) {
	return q.Queries().InsertTaxonStaging(ctx, params.ToParams(importHash))
}

func (r *TaxonResolutionStore) UpsertTaxonResolution(ctx context.Context, q db.Querier, params biomedb.UpsertTaxonResolutionParams) (err error) {
	return q.Queries().UpsertTaxonResolution(ctx, params)
}

// Build dependency list of GBIF keys that need to be fetched (parents + accepted references for synonyms),
// based on the current set of resolved taxa and their GBIF matches,
// to ensure that all necessary GBIF data is available before materialization.
//
// Returns the list of missing GBIF keys that need to be fetched.
func (r *TaxonResolutionStore) ListMissingGBIFDependencies(ctx context.Context, q db.Querier, importHash string) (missingKeys []int32, err error) {
	if err = q.Queries().ExpandGBIFDependencies(ctx, importHash); err != nil {
		return nil, err
	}
	missingKeys, err = q.Queries().ListMissingGBIFKeys(ctx, importHash)
	return missingKeys, err
}

// Inserts taxa from GBIF staging into the main taxa table, for a given import and rank.
// This is done in two steps: first insert non-synonyms, then insert synonyms (which depend on the accepted taxa being present).
func (r *TaxonResolutionStore) InsertTaxaFromGBIF(ctx context.Context, q db.Querier, importHash string, rank models.TaxonRank) (err error) {
	err = q.Queries().InsertTaxaFromGBIF(ctx, biomedb.InsertTaxaFromGBIFParams{ImportHash: importHash, Rank: string(rank), IsSynonym: false})
	if err != nil {
		return err
	}
	err = q.Queries().InsertTaxaFromGBIF(ctx, biomedb.InsertTaxaFromGBIFParams{ImportHash: importHash, Rank: string(rank), IsSynonym: true})
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) MarkTaxaNeedingGBIFCandidates(ctx context.Context, q db.Querier, importHash string) (err error) {
	return q.Queries().MarkTaxaNeedingGBIFCandidates(ctx, importHash)
}

func (r *TaxonResolutionStore) MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importHash string, inputNames []string) (err error) {
	return q.Queries().MarkTaxaGBIFImportCompleted(ctx, importHash, inputNames)
}

func (r *TaxonResolutionStore) ListTaxnamesToFetchFromGBIF(ctx context.Context, q db.Querier, importHash string) (toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow, err error) {
	return q.Queries().ListTaxaToFetchGBIFCandidates(ctx, importHash)
}

func (r *TaxonResolutionStore) GenerateLocalTaxonCandidates(ctx context.Context, q db.Querier, importHash string) (err error) {

	err = q.Queries().CreateCandidateTaxaNameExact(ctx, importHash)
	if err != nil {
		return err
	}

	err = q.Queries().CreateCandidateTaxaFuzzy(ctx, 0.6, importHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) ListTaxonCandidates(
	ctx context.Context, q db.Querier, importHash string,
) (candidatesByName map[string][]models.TaxonCandidate, err error) {

	candidates, err := q.Queries().ListAllTaxonCandidates(ctx, importHash)
	if err != nil {
		return nil, err
	}

	candidatesByName = make(map[string][]models.TaxonCandidate)
	for _, candidate := range candidates {
		candidatesByName[candidate.InputName] = append(candidatesByName[candidate.InputName], models.TaxonCandidateFromDB(candidate))
	}
	return candidatesByName, nil
}

func (r *TaxonResolutionStore) GetTaxonResolutionState(ctx context.Context, q db.Querier, importHash string) (state *models.TaxonResolutionState, err error) {
	resolutions, err := q.Queries().GetTaxonResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}
	candidates, err := r.ListTaxonCandidates(ctx, q, importHash)
	if err != nil {
		return nil, err
	}

	return &models.TaxonResolutionState{
		Resolution: models.TaxonResolutionFromDBSlice(resolutions),
		Candidates: candidates,
	}, nil
}
