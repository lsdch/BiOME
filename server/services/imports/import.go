package imports

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jszwec/csvutil"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/stores"
	gbif "github.com/lsdch/biome/models/taxonomy/GBIF"
	"github.com/lsdch/biome/services"
)

type ImportService struct {
	db         *pgx.Conn
	gbifClient *gbif.GBIFClient
}

func NewImportService(db *pgx.Conn, gbifClient *gbif.GBIFClient) *ImportService {
	return &ImportService{
		db:         db,
		gbifClient: gbifClient,
	}
}

type CSVImportParams struct {
	Reader    io.Reader
	Separator rune
}

type ImportResolutionState struct {
	Workflow biomedb.ImportWorkflow `json:"workflow"`
	// Taxon resolution state
	Taxonomy *stores.TaxonResolutionState `json:"taxonomy"`
	// Sampling methods resolution state
	SamplingMethods []biomedb.SamplingMethodsResolution `json:"sampling_methods"`
}

func (r *ImportService) InitImportWorkflow(
	ctx context.Context,
	importHash,
	label string,
	params CSVImportParams,
) (state *ImportResolutionState, err error) {

	occurrencesData, err := r.ParseCSV(params.Reader, params.Separator)
	if err != nil {
		return nil, fmt.Errorf("CSV parse error: %w", err)
	}

	occurrencesStaging := make([]biomedb.CopyImportStagingParams, len(occurrencesData))
	for i, row := range occurrencesData {
		occurrencesStaging[i] = row.ToStaging(importHash)
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := biomedb.New(tx)

	storeTx := stores.NewImportStore(qtx)

	workflow, err := storeTx.InitBatchImport(ctx, importHash, label)
	if err != nil {
		return nil, err
	}
	if err = storeTx.Bootstrap(ctx, importHash); err != nil {
		return nil, err
	}
	if err = storeTx.InsertStaging(ctx, importHash, occurrencesStaging); err != nil {
		return nil, err
	}

	resolver := NewTaxonResolutionService(qtx, r.gbifClient)
	resolutionState, err := resolver.InitResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}

	samplingService := services.NewSamplingService(qtx)
	methodsResolution, err := samplingService.InitMethodResolution(ctx, importHash)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	go func(importHash string) {
		resolver := NewTaxonResolutionService(biomedb.New(r.db), r.gbifClient)
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		resolver.EnrichTaxonResolutionWithGBIF(bgCtx, importHash)
	}(importHash)

	return &ImportResolutionState{
		Workflow:        workflow,
		Taxonomy:        resolutionState,
		SamplingMethods: methodsResolution,
	}, nil
}

func (r *ImportService) EnrichTaxonomyGBIF(ctx context.Context, importHash string) (err error) {
	resolver := NewTaxonResolutionService(biomedb.New(r.db), r.gbifClient)
	_, err = resolver.EnrichTaxonResolutionWithGBIF(ctx, importHash)
	return err
}

func (r *ImportService) DetectSamplingCollisions(ctx context.Context, importHash string, params stores.CollisionDetectionParams) (collisionsMap map[string]*StagingSamplingWithCollisions, err error) {
	collisionService := NewOccurrenceCollisionsService(biomedb.New(r.db))
	return collisionService.DetectSamplingCollisions(ctx, importHash, params)
}

func (r *ImportService) DetectOccurrenceCollisions(ctx context.Context, importHash string, params stores.CollisionDetectionParams) (collisions []OccurrenceCollisionsAtRow, err error) {
	collisionService := NewOccurrenceCollisionsService(biomedb.New(r.db))
	return collisionService.DetectOccurrenceCollisions(ctx, importHash, params)
}

func (s *ImportService) ParseCSV(r io.Reader, sep rune) ([]OccurrenceImportRow, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comma = sep
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return nil, err
	}

	// Sanitize fields
	dec.Map = func(field, column string, v any) string {
		s := strings.TrimSpace(field)
		if s == "NA" || s == "N/A" {
			return ""
		}
		return s
	}

	var rows []OccurrenceImportRow

	_ = dec.Header()
	row := int32(2) // Start counting from 2 to account for the header row
	for {
		u := OccurrenceImportRow{}

		if err := dec.Decode(&u); err == io.EOF {
			break
		} else if err != nil {
			return nil, &CSVParseError{
				RowNumber: row,
				Err:       err,
			}
		}
		rows = append(rows, u)
		row++
	}
	return rows, nil
}
