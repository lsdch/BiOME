package imports

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/pgutils"
	"github.com/lsdch/biome/db/stores"
	"github.com/lsdch/biome/models"
	gbif "github.com/lsdch/biome/models/taxonomy/GBIF"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type taxonResolutionService struct {
	store *stores.TaxonResolutionStore
	gbif  *gbif.GBIFClient
}

type TaxonResolutionService interface {
	InitResolution(ctx context.Context, importHash string) (state *stores.TaxonResolutionState, err error)
	GetTaxonResolutionState(ctx context.Context, importHash string) (state *stores.TaxonResolutionState, err error)
	MaterializeTaxa(ctx context.Context, importHash string) (err error)
	EnrichTaxonResolutionWithGBIF(ctx context.Context, importHash string) (status biomedb.GbifImportStatus, err error)
	// FetchCandidatesFromGBIF(ctx context.Context, importHash string, toFetch []biomedb.ListTaxaWithoutExactCandidateRow) (candidates map[string][]gbif.TaxonGBIF, err error)
	// InsertGBIFCandidates(ctx context.Context, candidates map[string][]gbif.TaxonGBIF) (err error)
	FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]gbif.TaxonGBIF, error)
}

func NewTaxonResolutionService(q *biomedb.Queries, gbif *gbif.GBIFClient) TaxonResolutionService {
	return &taxonResolutionService{
		store: stores.NewTaxonResolutionStore(q),
		gbif:  gbif,
	}
}

// Initializes the taxon resolution state for a given import hash. This is typically called when starting a new import workflow.
// It fetches the current resolution state, generates local candidates, and identifies any missing GBIF dependencies that need to be fetched.
func (r *taxonResolutionService) InitResolution(ctx context.Context, importHash string) (state *stores.TaxonResolutionState, err error) {
	resolution, err := r.store.InitTaxonResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}
	err = r.GenerateLocalCandidates(ctx, importHash)
	if err != nil {
		return nil, err
	}
	candidates, err := r.store.ListTaxonCandidates(ctx, importHash)
	if err != nil {
		return nil, err
	}
	return &stores.TaxonResolutionState{
		Resolution: resolution,
		Candidates: candidates,
	}, nil
}

// Fetches missing GBIF dependencies for taxa that have already been resolved to a GBIF taxon,
// to ensure that all necessary GBIF data is available before materialization.
func (r *taxonResolutionService) FetchKeysFromGBIF(ctx context.Context, toFetch []int32) ([]gbif.TaxonGBIF, error) {

	if len(toFetch) == 0 {
		return nil, nil
	}

	g, ctx := errgroup.WithContext(ctx)
	var (
		mu   sync.Mutex
		taxa = make([]gbif.TaxonGBIF, 0, len(toFetch))
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

func (r *taxonResolutionService) FetchCandidatesFromGBIF(ctx context.Context, importHash string, toFetch []biomedb.ListTaxaToFetchGBIFCandidatesRow) (candidates map[string][]gbif.TaxonGBIF, err error) {
	if len(toFetch) == 0 {
		return nil, nil
	}

	if err = r.store.InitGBIFCandidatesProgress(ctx, int32(len(toFetch)), importHash); err != nil {
		return nil, fmt.Errorf("failed to initialize GBIF candidates fetching progress : %v", err)
	}

	var mutex sync.Mutex
	g, ctx := errgroup.WithContext(ctx)
	candidatesMap := make(map[string][]gbif.TaxonGBIF, len(toFetch))
	for i := range toFetch {
		taxon := toFetch[i]
		query := taxon.FullInputName
		if query == "" {
			query = taxon.TaxonName
		}
		g.Go(func() error {
			params := gbif.SearchParams{
				Query:      query,
				Rank:       taxon.TaxonRank.String,
				Limit:      5,
				DatasetKey: gbif.GBIF_BACKBONE_DATASET_KEY,
			}
			if !taxon.TaxonRank.Valid {
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
			if _, err = r.store.IncrementGBIFCandidatesProgress(ctx, importHash); err != nil {
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

func (r *taxonResolutionService) InsertGBIFCandidates(ctx context.Context, candidates map[string][]gbif.TaxonGBIF) (err error) {
	if len(candidates) == 0 {
		return nil
	}
	stagingParams := make([]gbif.TaxonGBIF, 0, len(candidates))
	for _, matches := range candidates {
		for _, match := range matches {
			stagingParams = append(stagingParams, match)
		}
	}
	if err = r.store.InsertGBIFBatch(ctx, stagingParams); err != nil {
		return err
	}
	if err = r.store.InsertGBIFCandidatesBatch(ctx, candidates); err != nil {
		return err
	}
	return nil
}

func (r *taxonResolutionService) GenerateLocalCandidates(ctx context.Context, importHash string) error {
	return r.store.GenerateLocalTaxonCandidates(ctx, importHash)
}

const (
	ErrGBIFAlreadyRunning string = "GBIF import is already in progress for this import hash"
)

func (r *taxonResolutionService) EnrichTaxonResolutionWithGBIF(ctx context.Context, importHash string) (status biomedb.GbifImportStatus, err error) {

	claimed, err := r.store.ClaimGBIFImport(ctx, importHash)
	if err != nil {
		return
	}
	if !claimed {
		s, err := r.store.GetGBIFImportStatus(ctx, importHash)
		if err != nil {
			return status, fmt.Errorf("failed to get GBIF import status for %s", importHash)
		}
		return s.Status, fmt.Errorf(ErrGBIFAlreadyRunning)
	}

	if err = r.store.MarkTaxaNeedingGBIFCandidates(ctx, importHash); err != nil {
		return status, fmt.Errorf("failed to mark taxa needing GBIF candidates for %s", importHash)
	}
	toFetch, err := r.store.ListTaxnamesToFetchFromGBIF(ctx, importHash)
	if err != nil {
		return status, fmt.Errorf("failed to list taxa to fetch from GBIF for %s", importHash)
	}
	taxa, err := r.FetchCandidatesFromGBIF(ctx, importHash, toFetch)
	if err != nil {
		status, err = r.store.FailGBIFImport(ctx, importHash)
		return status, err
	}
	if err = r.InsertGBIFCandidates(ctx, taxa); err != nil {
		status, err = r.store.FailGBIFImport(ctx, importHash)
		return status, err
	}
	if err = r.store.MarkTaxaGBIFImportCompleted(ctx, importHash, slices.Collect(maps.Keys(taxa))); err != nil {
		return status, fmt.Errorf("failed to mark taxa needing GBIF candidates for %s", importHash)
	}

	return r.store.CompleteGBIFImport(ctx, importHash)
}

func (r *taxonResolutionService) GetTaxonResolutionState(ctx context.Context, importHash string) (state *stores.TaxonResolutionState, err error) {
	state, err = r.store.GetTaxonResolutionState(ctx, importHash)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (r *taxonResolutionService) MaterializeTaxaGBIF(ctx context.Context, importHash string) (err error) {

	gbifDepsKeys, err := r.store.ListMissingGBIFDependencies(ctx, importHash)
	if err != nil {
		return err
	}

	gbifDeps, err := r.FetchKeysFromGBIF(ctx, gbifDepsKeys)
	if err != nil {
		return err
	}

	if err = r.store.InsertGBIFBatch(ctx, gbifDeps); err != nil {
		return err
	}

	for _, rank := range biomedb.AllTaxonRankValues() {
		err = r.store.InsertTaxaFromGBIF(ctx, importHash, rank)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *taxonResolutionService) MaterializeTaxa(ctx context.Context, importHash string) (err error) {
	// For now we only materialize GBIF-resolved taxa, but in the future we may want to also materialize manually resolved taxa that are not in GBIF
	if err = r.MaterializeTaxaGBIF(ctx, importHash); err != nil {
		return err
	}

	// TODO : materialize manually resolved taxa that are not in GBIF
	return nil
}

type TaxonStagingParams struct {
	Name            string                          `json:"name"`
	Authorship      models.OptionalInput[string]    `json:"authorship,omitempty"`
	Rank            biomedb.TaxonRank               `json:"rank"`
	Status          biomedb.TaxonStatus             `json:"status"`
	ParentSource    biomedb.TaxonMatchSource        `json:"parent_source"`
	ParentID        models.OptionalInput[uuid.UUID] `json:"parent_taxa_id,omitempty"`
	ParentGbifID    models.OptionalInput[int32]     `json:"parent_gbif_id,omitempty"`
	ParentInputName models.OptionalInput[string]    `json:"parent_input_name,omitempty"`
}

func (r *taxonResolutionService) ResolveToManualTaxon(ctx context.Context, importHash string, params TaxonStagingParams) (err error) {
	if params.ParentSource == biomedb.TaxonMatchSourceManual {

		if parentName, ok := params.ParentInputName.Get(); !ok || parentName == "" {
			return fmt.Errorf("parent source is %s but parent name is not provided for manual taxon %s", params.ParentSource, params.Name)
		} else {
			r.store.UpsertTaxonResolution(ctx, biomedb.UpsertTaxonResolutionParams{
				ImportHash: importHash,
				InputName:  parentName,
			})
		}

	}
	err = r.store.InsertTaxonStaging(ctx, biomedb.InsertTaxonStagingParams{
		ImportHash:      importHash,
		Name:            params.Name,
		Authorship:      pgutils.TextOpt(params.Authorship),
		Rank:            params.Rank,
		Status:          params.Status,
		ParentSource:    params.ParentSource,
		ParentTaxaID:    pgutils.UUIDOpt(params.ParentID),
		ParentGbifID:    pgutils.Int4Opt(params.ParentGbifID),
		ParentInputName: pgutils.TextOpt(params.ParentInputName),
	})
	return err
}
