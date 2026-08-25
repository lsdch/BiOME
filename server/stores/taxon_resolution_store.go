package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/app_errors"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

type InconsistentTaxon struct {
	Name        string   `json:"name"`
	Authorships []string `json:"authorships"`
	Ranks       []string `json:"ranks"`
}

func InconsistentTaxonFromDB(t biomedb.CheckTaxaConsistencyInImportRow) InconsistentTaxon {
	return InconsistentTaxon{
		Name:        t.InputName,
		Authorships: t.Authorships,
		Ranks:       t.Ranks,
	}
}

func (t InconsistentTaxon) ErrorDetail() *app_errors.AppErrorDetail {
	return &app_errors.AppErrorDetail{
		Message:  fmt.Sprintf("Taxon %s has inconsistent authorships or ranks", t.Name),
		Value:    t,
		Location: t.Name,
	}
}

type ErrInconsistentTaxa struct {
	Taxa []InconsistentTaxon
}

func ErrInconsistentTaxaFromDB(rows []biomedb.CheckTaxaConsistencyInImportRow) ErrInconsistentTaxa {
	taxa := make([]InconsistentTaxon, len(rows))
	for i, row := range rows {
		taxa[i] = InconsistentTaxonFromDB(row)
	}
	return ErrInconsistentTaxa{Taxa: taxa}
}

func (e ErrInconsistentTaxa) Error() string {
	b, _ := json.MarshalIndent(e, "", "  ")
	return fmt.Sprintf("inconsistent taxa: %s", string(b))
}

func (e ErrInconsistentTaxa) AppError() *app_errors.AppError {
	errs := make([]*app_errors.AppErrorDetail, len(e.Taxa))
	for i, taxon := range e.Taxa {
		errs[i] = taxon.ErrorDetail()
	}
	return &app_errors.AppError{
		Code:     app_errors.ErrRegistry[app_errors.ErrorCategoryImport][app_errors.ErrorCodeInconsistentTaxa],
		Category: app_errors.ErrorCategoryImport,
		ErrorModel: huma.ErrorModel{
			Status: http.StatusUnprocessableEntity,
			Title:  "Inconsistent taxa in import",
			Detail: "Some taxa in the import have inconsistent authorship or rank information. Please review the import data and correct any inconsistencies before proceeding.",
			Errors: errs,
		},
	}
}

type TaxonResolutionStore struct {
}

func NewTaxonResolutionStore() *TaxonResolutionStore {
	return &TaxonResolutionStore{}
}

func (r *TaxonResolutionStore) InitTaxonResolution(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.TaxonResolution, error) {
	check, err := q.Queries().CheckTaxaConsistencyInImport(ctx, importID)
	if err != nil {
		return nil, err
	}
	if len(check) > 0 {
		return nil, ErrInconsistentTaxaFromDB(check)
	}
	resolution, err := q.Queries().InitTaxonResolution(ctx, importID)
	if err != nil {
		return nil, err
	}
	return models.TaxonResolutionFromDBSlice(resolution), nil
}
func (r *TaxonResolutionStore) InitSamplingTargetResolution(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	return q.Queries().InitSamplingTargetResolution(ctx, importID)
}

func (r *TaxonResolutionStore) LinkTaxonResolutions(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	return q.Queries().LinkTaxonResolutions(ctx, importID)
}

func (r *TaxonResolutionStore) InsertGBIFBatch(ctx context.Context, q db.Querier, taxa []models.TaxonGBIF) (err error) {
	toInsert := make([]biomedb.InsertGBIFBatchParams, len(taxa))
	for i, taxon := range taxa {
		toInsert[i] = taxon.ToStaging()
	}

	batch := q.Queries().InsertGBIFBatch(ctx, toInsert)
	var errs error = nil
	batch.Exec(func(i int, err error) {
		errs = errors.Join(errs, err)
	})
	return errs
}

func (r *TaxonResolutionStore) InsertGBIFCandidatesBatch(ctx context.Context, q db.Querier, importID uuid.UUID, candidates map[uuid.UUID][]models.TaxonGBIFWithPriority) (err error) {
	toInsert := make([]biomedb.InsertTaxonCandidatesBatchParams, 0, len(candidates))
	for resolutionID, matches := range candidates {
		for _, match := range matches {
			if match.IsAcceptable() {
				toInsert = append(toInsert, match.ToCandidate(importID, resolutionID))
			}
		}
	}

	batch := q.Queries().InsertTaxonCandidatesBatch(ctx, toInsert)
	var errs error = nil
	batch.Exec(func(i int, err error) {
		errs = errors.Join(errs, err)
	})
	return errs
}

func (r *TaxonResolutionStore) ListTaxonResolutionsWithoutCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (resolutions []models.TaxonResolution, err error) {
	res, err := q.Queries().ListTaxonResolutionsWithoutCandidates(ctx, importID)
	if err != nil {
		return nil, err
	}
	resolutions = make([]models.TaxonResolution, len(res))
	for i, res := range res {
		resolutions[i] = models.TaxonResolutionFromDB(res)
	}
	return resolutions, nil
}

func (r *TaxonResolutionStore) InsertTaxaStaging(ctx context.Context, q db.Querier, importID uuid.UUID, toStage []models.TaxonStagingParams) (err error) {
	params := make([]biomedb.InsertTaxaStagingParams, len(toStage))
	for i, param := range toStage {
		params[i] = param.ToParams(importID)
	}
	batch := q.Queries().InsertTaxaStaging(ctx, params)
	var errs error = nil
	batch.Exec(func(i int, err error) {
		errs = errors.Join(errs, err)
	})
	return errs
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
func (r *TaxonResolutionStore) MaterializeTaxaFromGBIF(ctx context.Context, tx *db.Tx, importID uuid.UUID, rank models.TaxonRank) (err error) {
	logrus.Infof("Materializing non synonym taxa for import ID %s at rank %s", importID, rank)
	err = tx.Queries().MaterializeTaxaFromGBIF(ctx, biomedb.MaterializeTaxaFromGBIFParams{ImportID: importID, Rank: string(rank), IsSynonym: false})
	if err != nil {
		return err
	}
	logrus.Infof("Materializing synonym taxa for import ID %s at rank %s", importID, rank)
	err = tx.Queries().MaterializeTaxaFromGBIF(ctx, biomedb.MaterializeTaxaFromGBIFParams{ImportID: importID, Rank: string(rank), IsSynonym: true})
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) MaterializeTaxaStaging(ctx context.Context, tx *db.Tx, importID uuid.UUID, rank models.TaxonRank) (err error) {
	logrus.Infof("Materializing taxa from staging for import ID %s at rank %s", importID, rank)
	if err = tx.Queries().MaterializeTaxaStaging(ctx, importID, rank); err != nil {
		return err
	}
	if err := tx.Queries().SyncMaterializedTaxa(ctx, importID, rank); err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) MarkTaxaNeedingGBIFCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	return q.Queries().MarkTaxaNeedingGBIFCandidates(ctx, importID)
}

func (r *TaxonResolutionStore) MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	return q.Queries().MarkTaxaGBIFImportCompleted(ctx, importID)
}

func (r *TaxonResolutionStore) ListTaxnamesToFetchFromGBIF(ctx context.Context, q db.Querier, importID uuid.UUID) (toFetch []models.TaxonResolution, err error) {
	res, err := q.Queries().ListResolutionsToFetchGBIFCandidates(ctx, importID)
	if err != nil {
		return nil, err
	}
	resolutions := make([]models.TaxonResolution, len(res))
	for i, res := range res {
		resolutions[i] = models.TaxonResolutionFromDB(res)
	}
	return resolutions, nil
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
) (candidatesByResolution map[uuid.UUID][]models.TaxonCandidate, err error) {

	candidates, err := q.Queries().ListAllTaxonCandidates(ctx, importID)
	if err != nil {
		return nil, err
	}

	candidatesByResolution = make(map[uuid.UUID][]models.TaxonCandidate)
	for _, candidate := range candidates {
		candidatesByResolution[candidate.ResolutionID] = append(candidatesByResolution[candidate.ResolutionID], models.TaxonCandidateFromDB(candidate))
	}
	return candidatesByResolution, nil
}

func (r *TaxonResolutionStore) GetTaxonResolutions(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error) {
	resolutions, err := q.Queries().GetTaxonResolution(ctx, importID)
	if err != nil {
		return nil, err
	}
	candidates, err := r.ListTaxonCandidates(ctx, q, importID)
	if err != nil {
		return nil, err
	}

	state = make([]models.TaxonResolutionWithCandidates, len(resolutions))
	for i, resolution := range resolutions {
		model := models.TaxonResolutionFromDB(resolution.TaxonResolution)
		model.SetFromResolutionName(resolution.FromResolutionName)
		state[i] = models.TaxonResolutionWithCandidates{
			TaxonResolution: model,
			Candidates:      candidates[resolution.TaxonResolution.ID],
		}
	}
	return state, nil
}

func (r *TaxonResolutionStore) AutoResolveUnambiguousCandidates(ctx context.Context, q db.Querier, importID uuid.UUID, threshold int32) (err error) {
	err = q.Queries().AutoResolveUnambiguousCandidates(ctx, importID, threshold)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) SetNeedsResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	err = q.Queries().SetNeedsResolutionForUnresolvedCandidates(ctx, importID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaxonResolutionStore) ResolveTaxon(ctx context.Context, q db.Querier, importID uuid.UUID, input models.ResolveInput) (err error) {
	err = q.Queries().ResolveTaxon(ctx, biomedb.ResolveTaxonParams{
		ImportID:     importID,
		ResolutionID: input.ResolutionID,
		CandidateID:  input.CandidateID,
	})
	if err != nil {
		return err
	}
	return nil
}
