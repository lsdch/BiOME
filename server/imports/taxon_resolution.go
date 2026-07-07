package imports

import (
	"context"
	"fmt"
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
	InitResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (state *models.TaxonResolutionState, err error)
	GetTaxonResolutionState(ctx context.Context, q db.Querier, importID uuid.UUID) (state *models.TaxonResolutionState, err error)
	MaterializeTaxa(ctx context.Context, q db.Querier, importID uuid.UUID) (err error)
	FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]models.TaxonGBIF, error)
	ListTaxaToFetch(ctx context.Context, q db.Querier, importID uuid.UUID) ([]biomedb.ListTaxaToFetchGBIFCandidatesRow, error)
	FetchCandidatesFromGBIF(ctx context.Context, q db.Querier, importID uuid.UUID, toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow, trackers ...progress.ProgressReporter) (candidates map[string][]models.TaxonGBIF, err error)
	InsertGBIFCandidates(ctx context.Context, q db.Querier, candidates map[string][]models.TaxonGBIF) (err error)
	MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID, taxa []string) error
}

func NewTaxonResolutionService(gbif *gbif.GBIFClient) TaxonResolver {
	return &taxonResolver{gbif: gbif, store: stores.NewTaxonResolutionStore()}
}

func (r *taxonResolver) InitResolution(ctx context.Context, q db.Querier, importID uuid.UUID) (state *models.TaxonResolutionState, err error) {
	resolution, err := r.store.InitTaxonResolution(ctx, q, importID)
	if err != nil {
		return nil, err
	}
	if err = r.GenerateLocalCandidates(ctx, q, importID); err != nil {
		return nil, err
	}
	candidates, err := r.store.ListTaxonCandidates(ctx, q, importID)
	if err != nil {
		return nil, err
	}
	return &models.TaxonResolutionState{
		Resolution: resolution,
		Candidates: candidates,
	}, nil
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
	toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow,
	trackers ...progress.ProgressReporter,
) (candidates map[string][]models.TaxonGBIF, err error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	var mutex sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	candidatesMap := make(map[string][]models.TaxonGBIF, len(toFetch))
	for i := range toFetch {
		taxon := toFetch[i]
		query := taxon.FullInputName
		if query == "" {
			query = taxon.TaxonName
		}
		g.Go(func() error {
			params := gbif.SearchParams{
				Query:      query,
				Rank:       "",
				Limit:      5,
				DatasetKey: r.gbif.BackboneDatasetKey,
			}
			if taxon.TaxonRank == nil {
				switch len(strings.Split(taxon.TaxonName, " ")) {
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
				logrus.Errorf("error fetching GBIF data for taxon '%s': %v", taxon.TaxonName, err)
				return err
			}

			mutex.Lock()
			defer mutex.Unlock()
			for _, match := range resp.Results {
				candidatesMap[taxon.FullInputName] = append(candidatesMap[taxon.FullInputName], match)
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

func (r *taxonResolver) InsertGBIFCandidates(ctx context.Context, q db.Querier, candidates map[string][]models.TaxonGBIF) (err error) {
	if len(candidates) == 0 {
		return nil
	}
	stagingParams := make([]models.TaxonGBIF, 0, len(candidates))
	for _, matches := range candidates {
		for _, match := range matches {
			stagingParams = append(stagingParams, match)
		}
	}
	if err = r.store.InsertGBIFBatch(ctx, q, stagingParams); err != nil {
		return err
	}
	if err = r.store.InsertGBIFCandidatesBatch(ctx, q, candidates); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolver) GenerateLocalCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	return r.store.GenerateLocalTaxonCandidates(ctx, q, importID)
}

func (r *taxonResolver) ListTaxaToFetch(ctx context.Context, q db.Querier, importID uuid.UUID) ([]biomedb.ListTaxaToFetchGBIFCandidatesRow, error) {
	if err := r.store.MarkTaxaNeedingGBIFCandidates(ctx, q, importID); err != nil {
		return nil, fmt.Errorf("failed to mark taxa needing GBIF candidates for %s", importID)
	}
	return r.store.ListTaxnamesToFetchFromGBIF(ctx, q, importID)
}

func (r *taxonResolver) MarkTaxaGBIFImportCompleted(ctx context.Context, q db.Querier, importID uuid.UUID, taxa []string) error {
	return r.store.MarkTaxaGBIFImportCompleted(ctx, q, importID, taxa)
}

func (r *taxonResolver) GetTaxonResolutionState(ctx context.Context, q db.Querier, importID uuid.UUID) (state *models.TaxonResolutionState, err error) {
	state, err = r.store.GetTaxonResolutionState(ctx, q, importID)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (r *taxonResolver) MaterializeTaxaGBIF(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {

	gbifDepsKeys, err := r.store.ListMissingGBIFDependencies(ctx, q, importID)
	if err != nil {
		return err
	}

	gbifDeps, err := r.FetchKeysFromGBIF(ctx, gbifDepsKeys)
	if err != nil {
		return err
	}

	if err = r.store.InsertGBIFBatch(ctx, q, gbifDeps); err != nil {
		return err
	}

	for _, rank := range biomedb.AllTaxonRankValues() {
		err = r.store.InsertTaxaFromGBIF(ctx, q, importID, rank)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *taxonResolver) MaterializeTaxa(ctx context.Context, q db.Querier, importID uuid.UUID) (err error) {
	if err = r.MaterializeTaxaGBIF(ctx, q, importID); err != nil {
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
