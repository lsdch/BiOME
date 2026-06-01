package stores

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	gbif "github.com/lsdch/biome/models/taxonomy/GBIF"
)

type TaxonResolutionStore struct {
	q *biomedb.Queries
}

func NewTaxonResolutionStore(q *biomedb.Queries) *TaxonResolutionStore {
	return &TaxonResolutionStore{q: q}
}

func (r *TaxonResolutionStore) InitTaxonResolution(ctx context.Context, importHash string) (resolution []biomedb.TaxonResolution, err error) {
	resolution, err = r.q.InitTaxonResolution(ctx, importHash)
	return resolution, err
}

func (r *TaxonResolutionStore) ClaimGBIFImport(ctx context.Context, importHash string) (ok bool, err error) {
	rowsAffected, err := r.q.ClaimGBIFImport(ctx, importHash)
	return rowsAffected > 0, err
}

func (r *TaxonResolutionStore) CompleteGBIFImport(ctx context.Context, importHash string) (status biomedb.GbifImportStatus, err error) {
	return r.q.CompleteGBIFImport(ctx, importHash)
}

func (r *TaxonResolutionStore) FailGBIFImport(ctx context.Context, importHash string) (status biomedb.GbifImportStatus, err error) {
	return r.q.FailGBIFImport(ctx, importHash)
}

func (r *TaxonResolutionStore) InitGBIFCandidatesProgress(ctx context.Context, total int32, importHash string) error {
	return r.q.InitGBIFCandidatesProgress(ctx, total, importHash)
}
func (r *TaxonResolutionStore) IncrementGBIFCandidatesProgress(ctx context.Context, importHash string) (biomedb.IncrementGBIFCandidatesProgressRow, error) {
	return r.q.IncrementGBIFCandidatesProgress(ctx, importHash)
}

type GBIFImportState struct {
	Status    biomedb.GbifImportStatus `json:"status"`
	ClaimedAt pgtype.Timestamptz       `json:"claimed_at"`
}

func (r *TaxonResolutionStore) GetGBIFImportStatus(ctx context.Context, importHash string) (state GBIFImportState, err error) {
	res, err := r.q.GetGBIFImportStatus(ctx, importHash)
	if err != nil {
		return GBIFImportState{}, err
	}
	return GBIFImportState{
		Status:    res.GbifStatus,
		ClaimedAt: res.GbifClaimedAt,
	}, nil
}

func (r *TaxonResolutionStore) InsertGBIFBatch(ctx context.Context, taxa []gbif.TaxonGBIF) (err error) {
	toInsert := make([]biomedb.InsertGBIFBatchParams, len(taxa))
	for i, taxon := range taxa {
		toInsert[i] = taxon.ToStaging()
	}

	_, err = r.q.InsertGBIFBatch(ctx, toInsert)
	return err
}

func (r *TaxonResolutionStore) InsertGBIFCandidatesBatch(ctx context.Context, candidates map[string][]gbif.TaxonGBIF) (err error) {
	toInsert := make([]biomedb.InsertTaxonCandidatesBatchParams, 0, len(candidates))
	for inputName, matches := range candidates {
		for _, match := range matches {
			toInsert = append(toInsert, match.ToCandidate(inputName))
		}
	}

	_, err = r.q.InsertTaxonCandidatesBatch(ctx, toInsert)
	return err
}

func (r *TaxonResolutionStore) InsertTaxonStaging(ctx context.Context, params biomedb.InsertTaxonStagingParams) (err error) {
	return r.q.InsertTaxonStaging(ctx, params)
}

func (r *TaxonResolutionStore) UpsertTaxonResolution(ctx context.Context, params biomedb.UpsertTaxonResolutionParams) (err error) {
	return r.q.UpsertTaxonResolution(ctx, params)
}

// Build dependency list of GBIF keys that need to be fetched (parents + accepted references for synonyms),
// based on the current set of resolved taxa and their GBIF matches,
// to ensure that all necessary GBIF data is available before materialization.
//
// Returns the list of missing GBIF keys that need to be fetched.
func (r *TaxonResolutionStore) ListMissingGBIFDependencies(ctx context.Context, importHash string) (missingKeys []int32, err error) {
	if err = r.q.ExpandGBIFDependencies(ctx, importHash); err != nil {
		return nil, err
	}
	missingKeys, err = r.q.ListMissingGBIFKeys(ctx, importHash)
	return missingKeys, err
}

// Inserts taxa from GBIF staging into the main taxa table, for a given import and rank.
// This is done in two steps: first insert non-synonyms, then insert synonyms (which depend on the accepted taxa being present).
func (r *TaxonResolutionStore) InsertTaxaFromGBIF(ctx context.Context, importHash string, rank biomedb.TaxonRank) (err error) {
	err = r.q.InsertTaxaFromGBIF(ctx, biomedb.InsertTaxaFromGBIFParams{ImportHash: importHash, Rank: string(rank), IsSynonym: false})
	if err != nil {
		return err
	}
	err = r.q.InsertTaxaFromGBIF(ctx, biomedb.InsertTaxaFromGBIFParams{ImportHash: importHash, Rank: string(rank), IsSynonym: true})
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) MarkTaxaNeedingGBIFCandidates(ctx context.Context, importHash string) (err error) {
	return r.q.MarkTaxaNeedingGBIFCandidates(ctx, importHash)
}

func (r *TaxonResolutionStore) MarkTaxaGBIFImportCompleted(ctx context.Context, importHash string, inputNames []string) (err error) {
	return r.q.MarkTaxaGBIFImportCompleted(ctx, importHash, inputNames)
}

func (r *TaxonResolutionStore) ListTaxnamesToFetchFromGBIF(ctx context.Context, importHash string) (toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow, err error) {
	return r.q.ListTaxaToFetchGBIFCandidates(ctx, importHash)
}

type TaxonCandidate struct {
	Name       string                   `json:"name"`
	Authorship string                   `json:"authorship"`
	Rank       biomedb.TaxonRank        `json:"rank"`
	Status     biomedb.TaxonStatus      `json:"status"`
	Source     biomedb.TaxonMatchSource `json:"source"`
	MatchType  biomedb.TaxonMatchType   `json:"match_type"`
	Score      pgtype.Float8            `json:"score"`
	TaxonID    pgtype.UUID              `json:"taxon_id,omitempty"`
	GBIF_ID    pgtype.Int4              `json:"gbif_id,omitempty"`
}

func (r *TaxonResolutionStore) GenerateLocalTaxonCandidates(ctx context.Context, importHash string) (err error) {

	err = r.q.CreateCandidateTaxaNameExact(ctx, importHash)
	if err != nil {
		return err
	}

	err = r.q.CreateCandidateTaxaFuzzy(ctx, 0.6, importHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) ListTaxonCandidates(ctx context.Context, importHash string) (candidatesByName map[string][]TaxonCandidate, err error) {

	candidates, err := r.q.ListAllTaxonCandidates(ctx, importHash)
	if err != nil {
		return nil, err
	}

	candidatesByName = make(map[string][]TaxonCandidate)
	for _, candidate := range candidates {
		candidatesByName[candidate.InputName] = append(candidatesByName[candidate.InputName], TaxonCandidate{
			Name:       candidate.TaxonName,
			Authorship: candidate.TaxonAuthorship.String,
			Rank:       candidate.ResolvedRank,
			Status:     candidate.ResolvedStatus,
			Source:     candidate.Source,
			MatchType:  candidate.MatchType,
			Score:      candidate.Score,
			TaxonID:    candidate.ResolvedTaxonID,
			GBIF_ID:    candidate.ResolvedGbifID,
		})
	}
	return candidatesByName, nil
}

type TaxonResolutionState struct {
	Resolution []biomedb.TaxonResolution   `json:"resolution"`
	Candidates map[string][]TaxonCandidate `json:"candidates"`
}

func (r *TaxonResolutionStore) GetTaxonResolutionState(ctx context.Context, importHash string) (state *TaxonResolutionState, err error) {
	resolution, err := r.q.GetTaxonResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}

	candidates, err := r.ListTaxonCandidates(ctx, importHash)
	if err != nil {
		return nil, err
	}

	return &TaxonResolutionState{
		Resolution: resolution,
		Candidates: candidates,
	}, nil
}
