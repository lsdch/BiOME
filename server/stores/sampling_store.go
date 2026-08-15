package stores

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/types"
	"github.com/sirupsen/logrus"
	"github.com/uber/h3-go/v4"
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

	logrus.Debugf("ListSamplingsAtProximity: input=%+v", input)
	rows, err := q.Queries().ListSamplingsAtProximity(ctx, input.ToParams())
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingWithDistance, len(rows))
	for i, r := range rows {
		res[i] = models.
			NewSamplingFromDB(r.SamplingsWithCountry).
			WithDistance(r.DistanceMeters)
	}

	return res, nil
}

func (s *SamplingStore) ListSamplingsH3AtProximity(
	ctx context.Context,
	q db.Querier,
	input models.ListSamplingsAtProximityInput,
) ([]models.H3CellWithRichnessAndDistance, error) {

	logrus.Debugf("ListSamplingsH3AtProximity: input=%+v", input)
	rows, err := q.Queries().ListSamplingsH3AtProximity(ctx, input.ToParamsH3())
	if err != nil {
		return nil, err
	}

	res := make([]models.H3CellWithRichnessAndDistance, len(rows))
	for i, r := range rows {
		res[i] = models.CellH3FromDB(h3.Cell(r.H3Index), int32(r.SamplingCount), int32(r.OccurrenceCount)).
			WithDistance(int32(r.DistanceMeters)).
			WithRichness(r.SpeciesRichness, r.GenusRichness, r.FamilyRichness)
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

	sampling := models.NewSamplingFromDB(row)
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

	sampling := models.NewSamplingFromDB(row)
	return &sampling, nil
}

// ------------------------
// Habitats and access points
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
		res[i] = models.NewSamplingFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) GetOccurrencesAtSamplingsBatch(
	ctx context.Context,
	q db.Querier,
	samplingIDs []uuid.UUID,
	occurrenceIDs []types.ULID,
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
	importID uuid.UUID,
) ([]models.SamplingMethodResolution, error) {

	rows, err := q.Queries().InitMethodsResolution(ctx, importID)
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
	importID uuid.UUID,
) ([]models.SamplingMethodResolution, error) {

	rows, err := q.Queries().GetMethodsResolution(ctx, importID)
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
	importID uuid.UUID,
	input models.SamplingMethodResolutionInput,
) (models.SamplingMethodResolution, error) {

	row, err := q.Queries().ResolveMethod(ctx, input.ToParams(importID))
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

// ------------------------
// Fixative Resolution
// ------------------------

func (s *SamplingStore) InitFixativeResolution(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
) ([]models.SamplingFixativeResolution, error) {

	rows, err := q.Queries().InitSamplingFixativesResolution(ctx, importID)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingFixativeResolution, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingFixativeResolutionFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) GetFixativesResolution(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
) ([]models.SamplingFixativeResolution, error) {

	rows, err := q.Queries().GetSamplingFixativesResolution(ctx, importID)
	if err != nil {
		return nil, err
	}

	res := make([]models.SamplingFixativeResolution, len(rows))
	for i, r := range rows {
		res[i] = models.SamplingFixativeResolutionFromDB(r)
	}

	return res, nil
}

func (s *SamplingStore) ResolveFixative(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
	input models.SamplingFixativeResolutionInput,
) (models.SamplingFixativeResolution, error) {

	row, err := q.Queries().ResolveSamplingFixative(ctx, input.ToParams(importID))
	if err != nil {
		return models.SamplingFixativeResolution{}, err
	}

	return models.SamplingFixativeResolutionFromDB(row), nil
}

// ========================
// Materialization
// ========================

func (s *SamplingStore) MaterializeSamplings(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
) error {
	return q.Queries().MaterializeSamplings(ctx, importID)
}

func (s *SamplingStore) MaterializeSamplingMethods(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
) error {
	return q.Queries().MaterializeSamplingMethods(ctx, importID)
}

func (s *SamplingStore) MaterializeSamplingTargets(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
) error {
	return q.Queries().MaterializeSamplingTargets(ctx, importID)
}

func (s *SamplingStore) MaterializeSamplingFixatives(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
) error {
	return q.Queries().MaterializeSamplingFixatives(ctx, importID)
}
