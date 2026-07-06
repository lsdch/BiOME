package imports

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jszwec/csvutil"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/services"
	gbif "github.com/lsdch/biome/services/gbif"
	"github.com/lsdch/biome/stores"
)

type ImportService struct {
	gbifClient      *gbif.GBIFClient
	samplings       *services.SamplingService
	taxonResolution TaxonResolutionService
	store           *stores.ImportStore
}

func NewImportService(gbifClient *gbif.GBIFClient,
	samplings *services.SamplingService,
	taxonResolution TaxonResolutionService,
) *ImportService {
	return &ImportService{
		gbifClient:      gbifClient,
		samplings:       samplings,
		taxonResolution: taxonResolution,
		store:           stores.NewImportStore(),
	}
}

type CSVImportParams struct {
	Reader    io.Reader
	Separator rune
}

type ImportResolutionState struct {
	Workflow biomedb.ImportWorkflow `json:"workflow"`
	// Taxon resolution state
	Taxonomy *models.TaxonResolutionState `json:"taxonomy"`
	// Sampling methods resolution state
	SamplingMethods []models.SamplingMethodResolution `json:"sampling_methods"`
}

func (r *ImportService) ComputeImportHash(label string, userID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", slug.Make(label), userID.String())
}

func (r *ImportService) InitImportWorkflow(
	ctx context.Context,
	tx *db.Tx,
	label string,
	params CSVImportParams,
) (state *ImportResolutionState, err error) {

	importHash := r.ComputeImportHash(label, ctx.Value(db.CTX_USER_KEY).(uuid.UUID))

	occurrencesData, err := r.ParseCSV(params.Reader, params.Separator)
	if err != nil {
		return nil, fmt.Errorf("CSV parse error: %w", err)
	}

	occurrencesStaging := make([]biomedb.CopyImportStagingParams, len(occurrencesData))
	for i, row := range occurrencesData {
		occurrencesStaging[i] = row.ToStaging(importHash)
	}

	state = &ImportResolutionState{}

	state.Workflow, err = r.store.InitBatchImport(ctx, tx, importHash, label)
	if err != nil {
		return nil, err
	}
	if err = r.store.Bootstrap(ctx, tx, importHash); err != nil {
		return nil, err
	}
	if err = r.store.InsertStaging(ctx, tx, importHash, occurrencesStaging); err != nil {
		return nil, err
	}

	state.Taxonomy, err = r.taxonResolution.InitResolution(ctx, tx, importHash)
	if err != nil {
		return nil, err
	}

	state.SamplingMethods, err = r.samplings.InitMethodResolution(ctx, tx, importHash)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (r *ImportService) EnrichTaxonomyGBIF(ctx context.Context, q db.Querier, importHash string) (err error) {
	_, err = r.taxonResolution.EnrichTaxonResolutionWithGBIF(ctx, q, importHash)
	return err
}

func (s *ImportService) ParseCSV(r io.Reader, sep rune) ([]models.OccurrenceImportRow, error) {
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

	var rows []models.OccurrenceImportRow

	_ = dec.Header()
	row := int32(2) // Start counting from 2 to account for the header row
	for {
		u := models.OccurrenceImportRow{}

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
