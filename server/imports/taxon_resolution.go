package imports

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/lib/progress"
	"github.com/lsdch/biome/models"
	gbif "github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type taxonResolver struct {
	store *stores.TaxonResolutionStore
	gbif  *gbif.GBIFClient
}

type TaxonResolver interface {
	InitResolution(ctx context.Context, tx *db.Tx, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error)
	GetTaxonResolutions(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error)
	// Materializes taxa from GBIF and local candidates into the main taxon table,
	// ensuring that all necessary GBIF dependencies are fetched and inserted first.
	MaterializeTaxa(ctx context.Context, tx *db.Tx, importID uuid.UUID) (err error)
	FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]models.TaxonGBIF, error)
	ListTaxaToFetch(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.TaxonResolution, error)
	FetchCandidatesFromGBIF(ctx context.Context, q db.Querier, importID uuid.UUID, toFetch []models.TaxonResolution, trackers ...progress.ProgressReporter) (candidates map[uuid.UUID][]models.TaxonGBIFWithPriority, err error)
	InsertGBIFCandidates(ctx context.Context, q db.Querier, importID uuid.UUID, candidates map[uuid.UUID][]models.TaxonGBIFWithPriority) (err error)
	MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID, resolutionIDs []uuid.UUID) error
	AutoResolveUnambiguousCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error)
	ResolveTaxon(ctx context.Context, q db.Querier, importID uuid.UUID, input models.ResolveTaxonInput) (err error)
	FillGBIFDependencies(ctx context.Context, q db.Querier, importID uuid.UUID) error
}

func NewTaxonResolutionService(gbif *gbif.GBIFClient) TaxonResolver {
	return &taxonResolver{gbif: gbif, store: stores.NewTaxonResolutionStore()}
}

func (r *taxonResolver) InitResolution(ctx context.Context, tx *db.Tx, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error) {
	resolution, err := r.store.InitTaxonResolution(ctx, tx, importID)
	if err != nil {
		return nil, fmt.Errorf("error initializing taxon resolution: %w", err)
	}

	if err = r.store.LinkTaxonResolutions(ctx, tx, importID); err != nil {
		return nil, fmt.Errorf("error linking taxon resolutions: %w", err)
	}

	if err = r.GenerateLocalCandidates(ctx, tx, importID); err != nil {
		return nil, fmt.Errorf("error generating local candidates: %w", err)
	}

	if err = r.AutoResolveUnambiguousCandidates(ctx, tx, importID); err != nil {
		return nil, fmt.Errorf("error auto-resolving unambiguous candidates: %w", err)
	}

	candidates, err := r.store.ListTaxonCandidates(ctx, tx, importID)
	if err != nil {
		return nil, fmt.Errorf("error listing taxon candidates: %w", err)
	}

	resolutionState := make([]models.TaxonResolutionWithCandidates, 0, len(resolution))
	for _, res := range resolution {
		resolutionState = append(resolutionState, models.TaxonResolutionWithCandidates{
			TaxonResolution: res,
			Candidates:      candidates[res.ID],
		})
	}
	return resolutionState, nil
}

// Fetches missing GBIF dependencies for taxa that have already been resolved to a GBIF taxon,
// to ensure that all necessary GBIF data is available before materialization.
func (r *taxonResolver) FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]models.TaxonGBIF, error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	g, ctx := errgroup.WithContext(ctx)
	var (
		mu   sync.Mutex
		taxa = make([]models.TaxonGBIF, 0, len(toFetch))
	)

	for i := range toFetch {
		key := toFetch[i]

		g.Go(func() error {
			taxon, err := r.gbif.GetTaxonByKey(ctx, key)
			if err != nil {
				logrus.Errorf("error fetching GBIF taxon by key %d: %v", key, err)
				return err
			}

			mu.Lock()
			taxa = append(taxa, *taxon)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return taxa, nil
}

func (r *taxonResolver) FetchCandidatesFromGBIF(
	ctx context.Context,
	q db.Querier,
	importID uuid.UUID,
	toFetch []models.TaxonResolution,
	trackers ...progress.ProgressReporter,
) (candidates map[uuid.UUID][]models.TaxonGBIFWithPriority, err error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	var mutex sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	candidatesMap := make(map[uuid.UUID][]models.TaxonGBIFWithPriority, len(toFetch))
	for i := range toFetch {
		resolution := toFetch[i]
		query := resolution.ScientificName
		if query == "" {
			query = resolution.InputName
		}
		g.Go(func() error {
			params := gbif.SearchParams{
				Query:      query,
				Rank:       "",
				Limit:      5,
				DatasetKey: r.gbif.BackboneDatasetKey,
			}
			if rank, ok := resolution.InputRank.Get(); ok {
				params.Rank = rank
			} else {
				switch len(strings.Split(resolution.InputName, " ")) {
				case 1:
					params.Rank = "GENUS"
				case 2:
					params.Rank = "SPECIES"
				default:
					// Chance of an exact match is very low for names with more than 2 parts
					params.Rank = ""
				}
			}
			resp, err := r.gbif.SearchSpecies(ctx, params)
			if err != nil {
				logrus.Errorf("error fetching GBIF data for taxon '%s': %v", resolution.InputName, err)
				return err
			}

			mutex.Lock()
			defer mutex.Unlock()
			for _, match := range resp.Results {
				if match.IsAcceptable() {
					candidatesMap[resolution.ID] = append(candidatesMap[resolution.ID], match.WithPriority(resolution))
				}
			}
			for _, tracker := range trackers {
				tracker.Increment(1)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return candidatesMap, nil
}

func (r *taxonResolver) InsertGBIFCandidates(ctx context.Context, q db.Querier, importID uuid.UUID, candidates map[uuid.UUID][]models.TaxonGBIFWithPriority) (err error) {
	if len(candidates) == 0 {
		return nil
	}
	stagingParams := make(map[int32]models.TaxonGBIF, len(candidates))
	for _, matches := range candidates {
		for _, match := range matches {
			stagingParams[match.Key] = match.TaxonGBIF
		}
	}
	if err = r.store.InsertGBIFBatch(ctx, q, slices.Collect(maps.Values(stagingParams))); err != nil {
		return err
	}
	if err = r.store.InsertGBIFCandidatesBatch(ctx, q, importID, candidates); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) GenerateLocalCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	return r.store.GenerateLocalTaxonCandidates(ctx, q, importID)
}

func (r *taxonResolver) ListTaxaToFetch(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.TaxonResolution, error) {
	if err := r.store.MarkTaxaNeedingGBIFCandidates(ctx, q, importID); err != nil {
		return nil, fmt.Errorf("failed to mark taxa needing GBIF candidates for %s", importID)
	}
	return r.store.ListTaxnamesToFetchFromGBIF(ctx, q, importID)
}

func (r *taxonResolver) MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID, resolutionIDS []uuid.UUID) error {
	return r.store.MarkTaxaGBIFImportCompleted(ctx, q, importID, resolutionIDS)
}

func (r *taxonResolver) GetTaxonResolutions(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error) {
	state, err = r.store.GetTaxonResolutions(ctx, q, importID)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (r *taxonResolver) FillGBIFDependencies(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	depsKeys, err := r.store.ListMissingGBIFDependencies(ctx, q, importID)
	if err != nil {
		return err
	}

	deps, err := r.FetchKeysFromGBIF(ctx, depsKeys)
	if err != nil {
		return err
	}

	if err = r.store.InsertGBIFBatch(ctx, q, deps); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) MaterializeTaxaGBIF(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {

	for _, rank := range slices.Backward(biomedb.AllTaxonRankValues()) {

		err = r.store.MaterializeTaxaFromGBIF(ctx, q, importID, rank)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *taxonResolver) MaterializeTaxa(ctx context.Context, tx *db.Tx, importID uuid.UUID) (err error) {
	if err = r.MaterializeTaxaGBIF(ctx, tx, importID); err != nil {
		return err
	}
	if err = tx.Queries().UpdateMaterializedTaxonCandidates(ctx, importID); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) ResolveToManualTaxon(ctx context.Context, q db.Querier, importID uuid.UUID, params models.TaxonStagingParams) (err error) {
	if params.ParentSource == biomedb.TaxonMatchSourceManual {
		parentName, ok := params.ParentInputName.Get()
		if !ok || parentName == "" {
			return fmt.Errorf("parent source is %s but parent name is not provided for manual taxon %s", params.ParentSource, params.Name)
		}
		err = r.store.UpsertTaxonResolution(ctx, q, biomedb.UpsertTaxonResolutionParams{
			ImportID:  importID,
			InputName: parentName,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert parent taxon resolution for manual taxon %s: %v", params.Name, err)
		}
	}
	return r.store.InsertTaxonStaging(ctx, q, importID, params)
}

func (r *taxonResolver) AutoResolveUnambiguousCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	err = r.store.AutoResolveUnambiguousCandidates(ctx, q, importID)
	if err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) ResolveTaxon(ctx context.Context, q db.Querier, importID uuid.UUID, input models.ResolveTaxonInput) (err error) {
	err = r.store.ResolveTaxon(ctx, q, importID, input)
	if err != nil {
		return err
	}
	return nil
}
