package imports

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
	"github.com/lsdch/biome/services/crossref"
	"github.com/sirupsen/logrus"
)

type BibliographyResolver struct {
	crossref *crossref.Client
}

func NewBibliographyResolver(crossref *crossref.Client) *BibliographyResolver {
	return &BibliographyResolver{
		crossref: crossref,
	}
}

// Aggregate rows by DOI or verbatim string to prepare for resolution.
// This function groups rows that have the same DOI or verbatim string into a single resolution input.
func (r *BibliographyResolver) aggregateRows(rows []csvmodels.PublicationImportRow) map[string]*csvmodels.PublicationResolutionInput {
	inputs := make(map[string]*csvmodels.PublicationResolutionInput)

	for _, row := range rows {
		var key string
		if row.DOI != nil {
			key = row.DOI.String()
		} else {
			key = *row.Verbatim
		}

		input, ok := inputs[key]
		if !ok {
			input = &csvmodels.PublicationResolutionInput{
				Authors:  row.Authors,
				Year:     row.Year,
				Verbatim: row.Verbatim,
				DOI:      row.DOI,
				Title:    row.Title,
				Journal:  row.Journal,
			}
			inputs[key] = input
		}

		input.RowNumbers = append(input.RowNumbers, row.RowNumber())
	}

	return inputs
}

func (r *BibliographyResolver) InitBibliographyResolution(ctx context.Context, tx *db.Tx, importID uuid.UUID, rows []csvmodels.PublicationImportRow) error {
	aggregatedRows := r.aggregateRows(rows)
	inputs := make([]biomedb.InitBibliographyResolutionParams, 0, len(aggregatedRows))
	for _, input := range aggregatedRows {
		inputs = append(inputs, input.ToParams(importID))
	}
	batch := tx.Queries().InitBibliographyResolution(ctx, inputs)
	var errs error = nil
	batch.QueryRow(func(i int, missingRows []int32, err error) {
		if err != nil {
			errs = errors.Join(errs, err)
		} else if len(missingRows) > 0 {
			errs = errors.Join(errs,
				fmt.Errorf("missing rows for bibliography resolution: %v at input [DOI: %s; Verbatim: %s]",
					missingRows, models.ValueOrZero(inputs[i].DOI), models.ValueOrZero(inputs[i].Verbatim)),
			)
		}
	})
	if errs != nil {
		return errs
	}

	if err := r.GenerateInternalCandidates(ctx, tx, importID); err != nil {
		return fmt.Errorf("error generating internal candidates: %w", err)
	}

	if err := tx.Queries().GenerateBibliographyManualCandidates(ctx, importID); err != nil {
		return fmt.Errorf("error generating manual candidates: %w", err)
	}

	return r.AutoResolveBibliography(ctx, tx, importID, 0.8, 0.1)
}

func (r *BibliographyResolver) AutoResolveBibliography(ctx context.Context, q db.Querier, importID uuid.UUID, scoreThreshold float32, scoreMargin float32) error {
	err := q.Queries().AutoResolveBibliography(ctx, biomedb.AutoResolveBibliographyParams{
		ImportID:       importID,
		ScoreThreshold: scoreThreshold,
		ScoreMargin:    scoreMargin,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *BibliographyResolver) ResolveBibliographyManualCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	err := q.Queries().ResolveBibliographyManualCandidates(ctx, importID)
	if err != nil {
		return fmt.Errorf("error resolving manual candidates: %w", err)
	}
	return nil
}

func (r *BibliographyResolver) GenerateInternalCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	err := q.Queries().GenerateBibliographyInternalCandidates(ctx, importID)
	if err != nil {
		return fmt.Errorf("error generating internal candidates: %w", err)
	}
	return nil
}

func (r *BibliographyResolver) FetchExternalCandidatesDOI(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.PublicationStagingInput, error) {
	doisToFetch, err := q.Queries().ListDOIsToFetch(ctx, importID)
	if err != nil {
		return nil, fmt.Errorf("error listing DOIs to fetch: %w", err)
	}
	if len(doisToFetch) == 0 {
		logrus.Debugf("No DOIs to fetch for bibliography import %s", importID)
		return nil, nil
	}
	logrus.Debugf("Fetching %d DOIs for bibliography import %s", len(doisToFetch), importID)
	results := r.crossref.WorksBatch(ctx, doisToFetch)
	toStage := make([]models.PublicationStagingInput, 0, len(results))
	for _, res := range results {
		if res.Err != nil {
			logrus.WithError(res.Err).WithField("doi", res.Input).Errorf("error fetching work from CrossRef")
			continue
		}
		if res.Value == nil {
			continue
		}
		toStage = append(toStage, models.PubStagingFromCrossref(res.Value.Message))
	}
	return toStage, nil
}

func (r *BibliographyResolver) FetchExternalCandidatesQuery(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.PublicationStagingInput, error) {
	queries, err := q.Queries().ListVerbatimToFetch(ctx, importID)
	if err != nil {
		return nil, fmt.Errorf("error listing verbatim queries to fetch: %w", err)
	}
	if len(queries) == 0 {
		logrus.Debugf("No verbatim queries to fetch for bibliography import %s", importID)
		return nil, nil
	}
	logrus.Debugf("Fetching %d verbatim queries for bibliography import %s", len(queries), importID)
	results := r.crossref.QueryWorksBatch(ctx, queries, 10)
	toStage := make([]models.PublicationStagingInput, 0, len(results))
	for _, res := range results {
		if res.Err != nil {
			logrus.WithError(res.Err).WithField("verbatim", res.Input).Errorf("error fetching work from CrossRef")
			continue
		}
		if res.Value == nil {
			continue
		}
		for _, item := range res.Value.Message.Items {
			toStage = append(toStage, models.PubStagingFromCrossref(item))
		}
	}
	return toStage, nil
}

func (r *BibliographyResolver) StageExternalCandidates(ctx context.Context, q db.Querier, importID uuid.UUID, toStage []models.PublicationStagingInput) error {
	if len(toStage) > 0 {
		logrus.Debugf("Staging external candidates for bibliography import %s", importID)
		params := make([]biomedb.StagePublicationsParams, 0, len(toStage))
		for _, input := range toStage {
			params = append(params, input.ToDBParams())
		}
		_, err := q.Queries().StagePublications(ctx, params)
		if err != nil {
			return fmt.Errorf("error staging publications: %w", err)
		}
	}
	if err := q.Queries().GenerateBibliographyExternalCandidates(ctx, importID); err != nil {
		return fmt.Errorf("error generating external candidates: %w", err)
	}
	return nil
}

func (r *BibliographyResolver) ListCandidates(ctx context.Context, q db.Querier, importID uuid.UUID) (map[uuid.UUID][]models.PublicationCandidate, error) {
	candidates, err := q.Queries().ListPublicationCandidates(ctx, importID)
	if err != nil {
		return nil, fmt.Errorf("error listing bibliography candidates: %w", err)
	}

	candidatesByResolution := make(map[uuid.UUID][]models.PublicationCandidate)
	for _, c := range candidates {
		candidate := models.BasePublicationCandidateFromDB(c.PublicationCandidate)
		candidatesByResolution[candidate.ResolutionID] = append(candidatesByResolution[candidate.ResolutionID], models.PublicationCandidate{
			BasePublicationCandidate: candidate,
			DOI:                      models.NewOptionalFromPtr(c.DOI),
			Authors:                  c.Authors,
			Year:                     models.NewOptionalFromPtr(c.Year),
			Title:                    models.NewOptionalFromPtr(c.Title),
			Journal:                  models.NewOptionalFromPtr(c.Journal),
			Verbatim:                 c.Verbatim,
		})
	}
	return candidatesByResolution, nil
}

func (r *BibliographyResolver) ResolvePublication(ctx context.Context, q db.Querier, importID uuid.UUID, input models.ResolveInput) (err error) {
	err = q.Queries().ResolvePublication(ctx, biomedb.ResolvePublicationParams{
		ImportID:            importID,
		ResolutionID:        input.ResolutionID,
		ResolvedCandidateID: input.CandidateID,
	})
	if err != nil {
		return fmt.Errorf("error resolving publication for import %s: %w", importID, err)
	}
	return nil
}

func (r *BibliographyResolver) GetBibliographyResolution(
	ctx context.Context, q db.Querier, importID uuid.UUID,
) ([]models.PublicationResolutionWithCandidates, error) {
	resolution, err := q.Queries().ListPublicationResolutions(ctx, importID)
	if err != nil {
		return nil, fmt.Errorf("error listing bibliography resolution: %w", err)
	}

	candidates, err := r.ListCandidates(ctx, q, importID)
	if err != nil {
		return nil, fmt.Errorf("error listing bibliography candidates: %w", err)
	}

	resolutionState := make([]models.PublicationResolutionWithCandidates, 0, len(resolution))
	for _, res := range resolution {
		resolutionState = append(resolutionState, models.PublicationResolutionWithCandidates{
			PublicationResolution: models.PublicationResolutionFromDB(res),
			Candidates:            candidates[res.ID],
		})
	}
	return resolutionState, nil
}

func (r *BibliographyResolver) MaterializeBibliography(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	err := q.Queries().MaterializeBibliography(ctx, importID)
	if err != nil {
		return fmt.Errorf("error materializing bibliography: %w", err)
	}
	return nil
}
