package imports

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	gbif "github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/stores"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type taxonResolutionService struct {
	store *stores.TaxonResolutionStore
	gbif  *gbif.GBIFClient
}

type TaxonResolutionService interface {
	InitResolution(ctx context.Context, q db.Querier, importHash string) (state *models.TaxonResolutionState, err error)
	GetTaxonResolutionState(ctx context.Context, q db.Querier, importHash string) (state *models.TaxonResolutionState, err error)
	MaterializeTaxa(ctx context.Context, q db.Querier, importHash string) (err error)
	EnrichTaxonResolutionWithGBIF(ctx context.Context, q db.Querier, importHash string) (status biomedb.GBIFImportStatus, err error)
	FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]models.TaxonGBIF, error)
}

func NewTaxonResolutionService(gbif *gbif.GBIFClient) TaxonResolutionService {
	return &taxonResolutionService{gbif: gbif, store: stores.NewTaxonResolutionStore()}
}

func (r *taxonResolutionService) InitResolution(ctx context.Context, q db.Querier, importHash string) (state *models.TaxonResolutionState, err error) {
	store := stores.NewTaxonResolutionStore()
	resolution, err := store.InitTaxonResolution(ctx, q, importHash)
	if err != nil {
		return nil, err
	}
	if err = r.GenerateLocalCandidates(ctx, q, importHash); err != nil {
		return nil, err
	}
	candidates, err := store.ListTaxonCandidates(ctx, q, importHash)
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
func (r *taxonResolutionService) FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]models.TaxonGBIF, error) {
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

func (r *taxonResolutionService) FetchCandidatesFromGBIF(ctx context.Context, q db.Querier, importHash string, toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow) (candidates map[string][]models.TaxonGBIF, err error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	if err = r.store.InitGBIFCandidatesProgress(ctx, q, int32(len(toFetch)), importHash); err != nil {
		return nil, fmt.Errorf("failed to initialize GBIF candidates fetching progress : %v", err)
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
			if _, err = r.store.IncrementGBIFCandidatesProgress(ctx, q, importHash); err != nil {
				logrus.Errorf("error incrementing GBIF candidates progress for import '%s': %v", importHash, err)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return candidatesMap, nil
}

func (r *taxonResolutionService) InsertGBIFCandidates(ctx context.Context, q db.Querier, candidates map[string][]models.TaxonGBIF) (err error) {
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

func (r *taxonResolutionService) GenerateLocalCandidates(ctx context.Context, q db.Querier, importHash string) error {
	return r.store.GenerateLocalTaxonCandidates(ctx, q, importHash)
}

const (
	ErrGBIFAlreadyRunning string = "GBIF import is already in progress for this import hash"
)

func (r *taxonResolutionService) EnrichTaxonResolutionWithGBIF(ctx context.Context, q db.Querier, importHash string) (status models.GBIFImportStatus, err error) {

	claimed, err := r.store.ClaimGBIFImport(ctx, q, importHash)
	if err != nil {
		return
	}
	if !claimed {
		s, err := r.store.GetGBIFImportStatus(ctx, q, importHash)
		if err != nil {
			return status, fmt.Errorf("failed to get GBIF import status for %s", importHash)
		}
		return s.Status, fmt.Errorf(ErrGBIFAlreadyRunning)
	}

	if err = r.store.MarkTaxaNeedingGBIFCandidates(ctx, q, importHash); err != nil {
		return status, fmt.Errorf("failed to mark taxa needing GBIF candidates for %s", importHash)
	}
	toFetch, err := r.store.ListTaxnamesToFetchFromGBIF(ctx, q, importHash)
	if err != nil {
		return status, fmt.Errorf("failed to list taxa to fetch from GBIF for %s", importHash)
	}
	taxa, err := r.FetchCandidatesFromGBIF(ctx, q, importHash, toFetch)
	if err != nil {
		status, err = r.store.FailGBIFImport(ctx, q, importHash)
		return status, err
	}
	if err = r.InsertGBIFCandidates(ctx, q, taxa); err != nil {
		status, err = r.store.FailGBIFImport(ctx, q, importHash)
		return status, err
	}
	if err = r.store.MarkTaxaGBIFImportCompleted(ctx, q, importHash, slices.Collect(maps.Keys(taxa))); err != nil {
		return status, fmt.Errorf("failed to mark taxa needing GBIF candidates for %s", importHash)
	}

	return r.store.CompleteGBIFImport(ctx, q, importHash)
}

func (r *taxonResolutionService) GetTaxonResolutionState(ctx context.Context, q db.Querier, importHash string) (state *models.TaxonResolutionState, err error) {
	state, err = r.store.GetTaxonResolutionState(ctx, q, importHash)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (r *taxonResolutionService) MaterializeTaxaGBIF(ctx context.Context, q db.Querier, importHash string) (err error) {

	gbifDepsKeys, err := r.store.ListMissingGBIFDependencies(ctx, q, importHash)
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
		err = r.store.InsertTaxaFromGBIF(ctx, q, importHash, rank)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *taxonResolutionService) MaterializeTaxa(ctx context.Context, q db.Querier, importHash string) (err error) {
	if err = r.MaterializeTaxaGBIF(ctx, q, importHash); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolutionService) ResolveToManualTaxon(ctx context.Context, q db.Querier, importHash string, params models.TaxonStagingParams) (err error) {
	if params.ParentSource == biomedb.TaxonMatchSourceManual {
		parentName, ok := params.ParentInputName.Get()
		if !ok || parentName == "" {
			return fmt.Errorf("parent source is %s but parent name is not provided for manual taxon %s", params.ParentSource, params.Name)
		}
		err = r.store.UpsertTaxonResolution(ctx, q, biomedb.UpsertTaxonResolutionParams{
			ImportHash: importHash,
			InputName:  parentName,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert parent taxon resolution for manual taxon %s: %v", params.Name, err)
		}
	}
	return r.store.InsertTaxonStaging(ctx, q, importHash, params)
}
