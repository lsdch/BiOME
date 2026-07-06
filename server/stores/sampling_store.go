package stores

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/oklog/ulid/v2"
)

type SamplingStore struct {
}

func NewSamplingStore() *SamplingStore {
	return &SamplingStore{}
}

// ------------------------
// Samplings
// ------------------------

func (s *SamplingStore) ListSamplingsAtProximity(
	ctx context.Context,
	q db.Querier,
	input models.ListSamplingsAtProximityInput,
) ([]models.SamplingWithDistance, error) {

	rows, err := q.Queries().ListSamplingsAtProximity(ctx, input.ToParams())
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingWithDistance, len(rows))
	for i, r := range rows {
		res[i] = models.
			NewSamplingFromDB(r.Sampling, r.Country).
			WithDistance(r.DistanceMeters)
	}

	return res, nil
}

func (s *SamplingStore) CreateSampling(
	ctx context.Context,
	q db.Querier,
	input models.SamplingInput,
) (*models.Sampling, error) {

	row, err := q.Queries().CreateSampling(ctx, input.ToParams())
	if err != nil {
		return nil, err
	}

	sampling := models.NewSamplingFromDB(row.Sampling, row.Country)
	return &sampling, nil
}

func (s *SamplingStore) GetSampling(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) (*models.Sampling, error) {

	row, err := q.Queries().GetSampling(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	sampling := models.NewSamplingFromDB(row.Sampling, row.Country)
	return &sampling, nil
}

// ------------------------
// Methods
// ------------------------

func (s *SamplingStore) ListUnknownSamplingMethodCodes(
	ctx context.Context,
	q db.Querier,
	codes []string,
) ([]string, error) {
	return q.Queries().ListUnknownSamplingMethodCodes(ctx, codes)
}

func (s *SamplingStore) ReplaceMethodsAtSampling(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
	codes []string,
) error {
	return q.Queries().ReplaceMethodsAtSampling(ctx, samplingID, codes)
}

func (s *SamplingStore) GetSamplingMethodsAtEvent(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) ([]models.SamplingMethod, error) {

	rows, err := q.Queries().GetSamplingMethodsAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingMethod, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingMethodFromDB(r)
	}

	return res, nil
}

// ------------------------
// Fixatives
// ------------------------

func (s *SamplingStore) ListUnknownFixativeCodes(
	ctx context.Context,
	q db.Querier,
	codes []string,
) ([]string, error) {
	return q.Queries().ListUnknownFixativeCodes(ctx, codes)
}

func (s *SamplingStore) ReplaceFixativesAtSampling(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
	codes []string,
) error {
	return q.Queries().ReplaceFixativesAtSampling(ctx, samplingID, codes)
}

func (s *SamplingStore) GetSamplingFixativesAtEvent(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) ([]models.Fixative, error) {

	rows, err := q.Queries().GetSamplingFixativesAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Fixative, len(rows))
	for i, r := range rows {
		res[i] = models.FixativeFromDB(r)
	}

	return res, nil
}

// ------------------------
// Habitats
// ------------------------

func (s *SamplingStore) GetHabitatsAtEvent(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) ([]models.HabitatWithGroupName, error) {

	rows, err := q.Queries().GetHabitatsAtEvent(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	res := make([]models.HabitatWithGroupName, len(rows))
	for i, r := range rows {
		res[i] = models.HabitatWithGroupNameFromDB(
			r.Habitat,
			r.HabitatGroup.Label,
		)
	}

	return res, nil
}

// ------------------------
// Taxa
// ------------------------

func (s *SamplingStore) GetSamplingTargetTaxa(
	ctx context.Context,
	q db.Querier,
	samplingID uuid.UUID,
) ([]models.Taxon, error) {

	rows, err := q.Queries().GetSamplingTargetTaxa(ctx, samplingID)
	if err != nil {
		return nil, err
	}

	res := make([]models.Taxon, len(rows))
	for i, r := range rows {
		res[i] = *models.TaxonFromDB(&r)
	}

	return res, nil
}

// ------------------------
// Batch operations
// ------------------------

func (s *SamplingStore) GetSamplingBatch(
	ctx context.Context,
	q db.Querier,
	ids []uuid.UUID,
) ([]models.Sampling, error) {

	rows, err := q.Queries().GetSamplingBatch(ctx, ids)
	if err != nil {
		return nil, err
	}

	res := make([]models.Sampling, len(rows))
	for i, r := range rows {
		res[i] = models.NewSamplingFromDB(r.Sampling, r.Country)
	}

	return res, nil
}

func (s *SamplingStore) GetOccurrencesAtSamplingsBatch(
	ctx context.Context,
	q db.Querier,
	samplingIDs []uuid.UUID,
	occurrenceIDs []ulid.ULID,
) (map[uuid.UUID][]models.BaseOccurrence, error) {

	rows, err := q.Queries().GetOccurrencesAtSamplingsBatch(ctx, samplingIDs, occurrenceIDs)
	if err != nil {
		return nil, err
	}

	res := make(map[uuid.UUID][]models.BaseOccurrence, len(samplingIDs))
	for _, r := range rows {
		occurrence := models.BaseOccurrenceFromDB(r.Occurrence, r.Taxon)
		res[r.Occurrence.SamplingID] = append(res[r.Occurrence.SamplingID], occurrence)
	}

	return res, nil
}

// ------------------------
// Sampling Methods (CRUD)
// ------------------------

func (s *SamplingStore) ListSamplingMethods(
	ctx context.Context,
	q db.Querier,
) ([]models.SamplingMethod, error) {

	rows, err := q.Queries().ListSamplingMethods(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingMethod, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingMethodFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) CreateSamplingMethod(
	ctx context.Context,
	q db.Querier,
	input models.SamplingMethodInput,
) (models.SamplingMethod, error) {

	row, err := q.Queries().CreateSamplingMethod(ctx, input.ToDBParams())
	if err != nil {
		return models.SamplingMethod{}, err
	}

	return models.SamplingMethodFromDB(row), nil
}

func (s *SamplingStore) UpdateSamplingMethod(
	ctx context.Context,
	q db.Querier,
	code string,
	input models.SamplingMethodUpdateParams,
) (models.SamplingMethod, error) {

	row, err := q.Queries().UpdateSamplingMethod(ctx, input.ToParams(code))
	if err != nil {
		return models.SamplingMethod{}, err
	}

	return models.SamplingMethodFromDB(row), nil
}

func (s *SamplingStore) DeleteSamplingMethod(
	ctx context.Context,
	q db.Querier,
	code string,
) error {

	return q.Queries().DeleteSamplingMethod(ctx, code)
}

// ------------------------
// Sampling Method Resolution
// ------------------------

func (s *SamplingStore) InitMethodResolution(
	ctx context.Context,
	q db.Querier,
	importHash string,
) ([]models.SamplingMethodResolution, error) {

	rows, err := q.Queries().InitMethodsResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingMethodResolution, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingMethodResolutionFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) GetMethodsResolution(
	ctx context.Context,
	q db.Querier,
	importHash string,
) ([]models.SamplingMethodResolution, error) {

	rows, err := q.Queries().GetMethodsResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingMethodResolution, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingMethodResolutionFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) ResolveMethod(
	ctx context.Context,
	q db.Querier,
	importHash string,
	input models.SamplingMethodResolutionInput,
) (models.SamplingMethodResolution, error) {

	row, err := q.Queries().ResolveMethod(ctx, input.ToParams(importHash))
	if err != nil {
		return models.SamplingMethodResolution{}, err
	}

	return models.SamplingMethodResolutionFromDB(row), nil
}

// ------------------------
// Fixatives
// ------------------------

func (s *SamplingStore) ListSamplingFixatives(
	ctx context.Context,
	q db.Querier,
) ([]models.Fixative, error) {

	rows, err := q.Queries().ListFixatives(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.Fixative, len(rows))
	for i, r := range rows {
		res[i] = models.FixativeFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) CreateSamplingFixative(
	ctx context.Context,
	q db.Querier,
	input models.FixativeInput,
) (models.Fixative, error) {

	row, err := q.Queries().CreateFixative(ctx, input.ToDBParams())
	if err != nil {
		return models.Fixative{}, err
	}

	return models.FixativeFromDB(row), nil
}

func (s *SamplingStore) UpdateSamplingFixative(
	ctx context.Context,
	q db.Querier,
	code string,
	input models.FixativeUpdateParams,
) (models.Fixative, error) {

	row, err := q.Queries().UpdateFixative(ctx, input.ToParams(code))
	if err != nil {
		return models.Fixative{}, err
	}

	return models.FixativeFromDB(row), nil
}

func (s *SamplingStore) DeleteSamplingFixative(
	ctx context.Context,
	q db.Querier,
	code string,
) error {

	tag, err := q.Queries().DeleteFixative(ctx, code)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fixative with code %s not found", code)
	}

	return nil
}
