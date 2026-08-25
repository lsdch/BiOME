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
	InitSamplingTargetsResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (err error)
	GetTaxonResolutions(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error)
	// Materializes taxa from GBIF and local candidates into the main taxon table,
	// ensuring that all necessary GBIF dependencies are fetched and inserted first.
	MaterializeTaxa(ctx context.Context, tx *db.Tx, importID uuid.UUID) (err error)

	// FetchKeysFromGBIF fetches a batch of taxa from GBIF by their keys, returning the corresponding TaxonGBIF objects.
	FetchKeysFromGBIF(ctx context.Context, toFetch []int32, trackers ...progress.ProgressReporter) ([]models.TaxonGBIF, error)

	// ListTaxaToFetch returns a list of taxon resolutions that need to fetch candidates from GBIF.
	ListTaxaToFetch(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.TaxonResolution, error)

	// FetchCandidatesFromGBIF fetches candidate taxa from GBIF for the given taxon resolutions that need candidates.
	// It returns a map of taxon resolution IDs to their corresponding GBIF candidates.
	// All resolutions are included in the result, even if no candidates were found for a particular resolution.
	FetchCandidatesFromGBIF(ctx context.Context, higherTaxonKey int32, toFetch []models.TaxonResolution, trackers ...progress.ProgressReporter) (candidates map[uuid.UUID][]models.TaxonGBIFWithPriority, err error)
	InsertGBIFCandidates(ctx context.Context, q db.Querier, importID uuid.UUID, candidates map[uuid.UUID][]models.TaxonGBIFWithPriority) (err error)
	MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID) error
	AutoResolveUnambiguousCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error)
	ResolveTaxon(ctx context.Context, q db.Querier, importID uuid.UUID, input models.ResolveInput) (err error)
	FillGBIFDependencies(ctx context.Context, q db.Querier, importID uuid.UUID, trackers ...progress.ProgressReporter) error
	AutoCreateManualCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error)
	CreateManualCandidate(ctx context.Context, tx *db.Tx, importID uuid.UUID, input models.TaxonStagingParams) (err error)
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

func (r *taxonResolver) InitSamplingTargetsResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	if err = r.store.InitSamplingTargetResolution(ctx, q, importID); err != nil {
		return fmt.Errorf("error initializing sampling targets resolution: %w", err)
	}
	return nil
}

func (r *taxonResolver) CreateManualCandidate(ctx context.Context, tx *db.Tx, importID uuid.UUID, input models.TaxonStagingParams) (err error) {
	if err = r.store.InsertTaxaStaging(ctx, tx, importID, []models.TaxonStagingParams{input}); err != nil {
		return fmt.Errorf("error inserting manual candidate: %w", err)
	}
	if err = r.AutoResolveUnambiguousCandidates(ctx, tx, importID); err != nil {
		return fmt.Errorf("error auto-resolving unambiguous candidates: %w", err)
	}
	return nil
}

// Fetches missing GBIF dependencies for taxa that have already been resolved to a GBIF taxon,
// to ensure that all necessary GBIF data is available before materialization.
func (r *taxonResolver) FetchKeysFromGBIF(ctx context.Context, toFetch []int32, trackers ...progress.ProgressReporter) ([]models.TaxonGBIF, error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	for _, tracker := range trackers {
		tracker.AddToTotal(int32(len(toFetch)))
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

			for _, tracker := range trackers {
				tracker.Increment(1)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return taxa, nil
}

// FetchCandidatesFromGBIF fetches candidate taxa from GBIF for the given taxon resolutions that need candidates.
// It returns a map of taxon resolution IDs to their corresponding GBIF candidates.
// All resolutions are included in the result, even if no candidates were found for a particular resolution.
func (r *taxonResolver) FetchCandidatesFromGBIF(
	ctx context.Context,
	higherTaxonKey int32,
	toFetch []models.TaxonResolution,
	trackers ...progress.ProgressReporter,
) (candidates map[uuid.UUID][]models.TaxonGBIFWithPriority, err error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	for _, tracker := range trackers {
		tracker.AddToTotal(int32(len(toFetch)))
	}

	var mutex sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	candidatesMap := make(map[uuid.UUID][]models.TaxonGBIFWithPriority, len(toFetch))
	for i := range toFetch {
		resolution := toFetch[i]
		query := resolution.InputName
		g.Go(func() error {
			params := gbif.SearchParams{
				Query:          query,
				Rank:           "",
				HigherTaxonKey: higherTaxonKey,
				Limit:          10,
				DatasetKey:     r.gbif.BackboneDatasetKey,
			}
			if rank, ok := resolution.InputRank.Get(); ok {
				params.Rank = strings.ToUpper(rank)
			} else {
				switch len(strings.Fields(resolution.InputName)) {
				case 1:
					switch strings.ToLower(params.Query) {
					// Help GBIF a little by providing a rank for some known kingdoms
					case "animalia", "plantae", "fungi", "protozoa", "chromista":
						logrus.Infof("Using kingdom rank for taxon '%s'", resolution.InputName)
						params.Rank = "KINGDOM"
						params.HigherTaxonKey = 0
					default:
						params.Rank = "GENUS"
					}
				case 2:
					params.Rank = "SPECIES"
				case 3:
					params.Rank = "SUBSPECIES"
				default:
					params.Rank = ""
				}
			}
			if params.Rank == "SUBSPECIES" {
				// Scientific names in GBIF for subspecies do not include the authorship
				params.Query = resolution.InputName
			}
			resp, err := r.gbif.SearchSpecies(ctx, params)
			if err != nil {
				logrus.Errorf("error fetching GBIF data for taxon '%s': %v", resolution.InputName, err)
				return err
			}

			mutex.Lock()
			defer mutex.Unlock()
			hasExactMatch := false
			for _, match := range resp.Results {
				if match.IsAcceptable() {
					candidate := match.WithPriority(resolution)
					candidatesMap[resolution.ID] = append(candidatesMap[resolution.ID], candidate)
					if candidate.Priority == models.TaxonGBIFPriorityExactAccepted {
						hasExactMatch = true
					}
				} else {
					logrus.Warnf("GBIF taxon '%s' [%d] is not acceptable (rank: %s, status: %s, parent rank: %v)", match.ScientificName, match.Key, match.Rank, match.Status, match.GetParentRank())
				}
			}
			if _, ok := resolution.InputRank.Get(); !ok && !hasExactMatch {
				// If no exact match was found and no rank was specified,
				// broaden the search to include all ranks for this taxon resolution.
				params.Rank = ""
				resp, err := r.gbif.SearchSpecies(ctx, params)
				if err != nil {
					logrus.Errorf("error fetching GBIF data for taxon '%s' with no rank: %v", resolution.InputName, err)
					return err
				}
				for _, match := range resp.Results {
					if match.IsAcceptable() {
						candidate := match.WithPriority(resolution)
						candidatesMap[resolution.ID] = append(candidatesMap[resolution.ID], candidate)
					}
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

func (r *taxonResolver) MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	return r.store.MarkTaxaGBIFImportCompleted(ctx, q, importID)
}

func (r *taxonResolver) GetTaxonResolutions(ctx context.Context, q db.Querier, importID uuid.UUID) (state []models.TaxonResolutionWithCandidates, err error) {
	state, err = r.store.GetTaxonResolutions(ctx, q, importID)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (r *taxonResolver) FillGBIFDependencies(ctx context.Context, q db.Querier, importID uuid.UUID, trackers ...progress.ProgressReporter) error {
	depsKeys, err := r.store.ListMissingGBIFDependencies(ctx, q, importID)
	if err != nil {
		return err
	}

	logrus.Infof("Found %d missing GBIF dependencies for import %s", len(depsKeys), importID)
	deps, err := r.FetchKeysFromGBIF(ctx, depsKeys, trackers...)
	if err != nil {
		return err
	}

	if err = r.store.InsertGBIFBatch(ctx, q, deps); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) MaterializeTaxaGBIF(ctx context.Context, tx *db.Tx, importID uuid.UUID) (err error) {

	for _, rank := range slices.Backward(biomedb.AllTaxonRankValues()) {

		err = r.store.MaterializeTaxaFromGBIF(ctx, tx, importID, models.TaxonRank(rank))
		if err != nil {
			return err
		}
	}

	if err = tx.Queries().UpdateMaterializedGBIFCandidates(ctx, importID); err != nil {
		return err
	}
	return err
}

func (r *taxonResolver) MaterializeTaxaStaging(ctx context.Context, tx *db.Tx, importID uuid.UUID) (err error) {
	for _, rank := range slices.Backward(biomedb.AllTaxonRankValues()) {
		err = r.store.MaterializeTaxaStaging(ctx, tx, importID, rank)
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
	if err = r.MaterializeTaxaStaging(ctx, tx, importID); err != nil {
		return err
	}
	return nil
}

// CreateManualCandidates creates new taxa in the staging table for taxon resolutions that do not have any candidates.
// It also creates a corresponding taxon_candidates record for each new taxon, linking it to the taxon resolution.
// If the parent taxon does not exist, it will create a new taxon resolution for the parent as well.
// After this, AutoResolveUnambiguousCandidates should be run to resolve any unambiguous candidates that may have been created.
func (r *taxonResolver) AutoCreateManualCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	resolutions, err := r.store.ListTaxonResolutionsWithoutCandidates(ctx, q, importID)
	if err != nil {
		return err
	}
	toStage := make([]models.TaxonStagingParams, 0, len(resolutions))
	for _, res := range resolutions {
		r, ok := res.InputRank.Get()
		if !ok {
			continue
		}
		rank, ok := models.TaxonRankFromString(r)
		if !ok {
			logrus.Errorf("invalid rank '%s' for taxon resolution %s [%s]", r, res.InputName, res.ID)
			continue
		}

		parentName, ok := models.InferParentName(res.InputName)
		if !ok {
			logrus.Errorf("failed to infer parent name for taxon resolution %s [%s]", res.InputName, res.ID)
			continue
		}

		toStage = append(toStage, models.TaxonStagingParams{
			Name:         res.InputName,
			Authorship:   res.InputAuthorship,
			Rank:         rank,
			Status:       models.InferStatusFromStagingName(res.InputName),
			ParentName:   parentName,
			ResolutionID: res.ID,
		})

	}
	return r.store.InsertTaxaStaging(ctx, q, importID, toStage)
}

func (r *taxonResolver) AutoResolveUnambiguousCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	thresholds := []int32{100, 80}
	for _, threshold := range thresholds {
		err = r.store.AutoResolveUnambiguousCandidates(ctx, q, importID, threshold)
		if err != nil {
			return err
		}
	}
	if err = r.store.SetNeedsResolution(ctx, q, importID); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) ResolveTaxon(ctx context.Context, q db.Querier, importID uuid.UUID, input models.ResolveInput) (err error) {
	err = r.store.ResolveTaxon(ctx, q, importID, input)
	if err != nil {
		return err
	}
	return nil
}
