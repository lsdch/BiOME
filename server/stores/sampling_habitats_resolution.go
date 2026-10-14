package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
)

var (
	ErrHabitatBatchNotEditable     = errors.New("habitat resolutions require a staged batch")
	ErrConflictingHabitatCreations = errors.New("the same habitat name has conflicting creation definitions in this batch")
	ErrHabitatsNotReady            = errors.New("habitats are not ready for materialization")
)

func (s *SamplingStore) InitHabitatResolution(ctx context.Context, q db.Querier, importID uuid.UUID) error {
	if err := q.Queries().InitSamplingHabitatsResolution(ctx, importID); err != nil {
		return err
	}
	return q.Queries().InitSamplingHabitatTargets(ctx, importID)
}

func (s *SamplingStore) GetHabitatsResolution(ctx context.Context, q db.Querier, importID uuid.UUID) ([]models.SamplingHabitatResolution, error) {
	var result []models.SamplingHabitatResolution
	err := q.WithTx(ctx, func(tx *db.Tx) error {
		// Keep the resolution headers and their targets consistent across the two reads.
		if _, err := tx.Queries().ReadImportBatchForHabitats(ctx, importID); err != nil {
			return err
		}
		resolutions, err := tx.Queries().GetSamplingHabitatsResolution(ctx, importID)
		if err != nil {
			return err
		}
		targets, err := tx.Queries().GetSamplingHabitatTargets(ctx, importID)
		if err != nil {
			return err
		}
		result = make([]models.SamplingHabitatResolution, len(resolutions))
		byInput := make(map[string]*models.SamplingHabitatResolution, len(resolutions))
		for i, resolution := range resolutions {
			result[i] = models.SamplingHabitatResolution{
				ImportID: importID,
				SamplingHabitatResolutionInput: models.SamplingHabitatResolutionInput{
					InputText: resolution.InputText, Status: resolution.Status,
					HabitatIDs: []uuid.UUID{}, CreateHabitats: []models.HabitatInput{},
				},
			}
			byInput[resolution.InputText] = &result[i]
		}
		for _, target := range targets {
			resolution := byInput[target.InputText]
			if target.HabitatID.Valid {
				resolution.HabitatIDs = append(resolution.HabitatIDs, target.HabitatID.Bytes)
			} else if target.CreationName != nil {
				resolution.CreateHabitats = append(resolution.CreateHabitats, models.HabitatInput{
					Name:        *target.CreationName,
					Facet:       models.NewOptionalFromPtr(target.CreationFacet),
					Description: models.NewOptionalFromPtr(target.CreationDescription),
				})
			}
		}
		return nil
	})
	return result, err
}

func (s *SamplingStore) ResolveHabitats(ctx context.Context, q db.Querier, importID uuid.UUID, input models.SamplingHabitatResolutionInput) (models.SamplingHabitatResolution, error) {
	err := q.WithTx(ctx, func(tx *db.Tx) error {
		batch, err := tx.Queries().LockImportBatchForHabitats(ctx, importID)
		if err != nil {
			return err
		}
		if batch.Status != biomedb.ImportBatchStatusStaged {
			return ErrHabitatBatchNotEditable
		}
		changed, err := tx.Queries().UpdateSamplingHabitatsResolution(ctx, biomedb.UpdateSamplingHabitatsResolutionParams{
			ImportID: importID, InputText: input.InputText, Status: input.Status,
		})
		if err != nil {
			return err
		}
		if changed == 0 {
			return pgx.ErrNoRows
		}
		if err := tx.Queries().DeleteSamplingHabitatTargets(ctx, importID, input.InputText); err != nil {
			return err
		}
		targets := make([]biomedb.SamplingHabitatTarget, 0, len(input.HabitatIDs)+len(input.CreateHabitats))
		for _, id := range input.HabitatIDs {
			targets = append(targets, biomedb.SamplingHabitatTarget{HabitatID: pgtype.UUID{Bytes: id, Valid: true}})
		}
		for _, creation := range input.CreateHabitats {
			targets = append(targets, biomedb.SamplingHabitatTarget{
				CreationName: &creation.Name, CreationFacet: creation.Facet.ToPtr(), CreationDescription: creation.Description.ToPtr(),
			})
		}
		if len(targets) > 0 {
			encoded, err := json.Marshal(targets)
			if err != nil {
				return fmt.Errorf("encode habitat targets: %w", err)
			}
			if err := tx.Queries().InsertSamplingHabitatTargets(ctx, biomedb.InsertSamplingHabitatTargetsParams{
				ImportID: importID, InputText: input.InputText, Targets: encoded,
			}); err != nil {
				return err
			}
		}
		conflict, err := tx.Queries().HasConflictingHabitatCreations(ctx, importID)
		if err != nil {
			return err
		}
		if conflict {
			return ErrConflictingHabitatCreations
		}
		return nil
	})
	return models.SamplingHabitatResolution{ImportID: importID, SamplingHabitatResolutionInput: input}, err
}

func (s *SamplingStore) MaterializeSamplingHabitats(ctx context.Context, tx *db.Tx, importID uuid.UUID) error {
	batch, err := tx.Queries().LockImportBatchForHabitats(ctx, importID)
	if err != nil {
		return err
	}
	if batch.Status != biomedb.ImportBatchStatusStaged {
		return ErrHabitatBatchNotEditable
	}
	// A selected habitat must not disappear between readiness checking and linking.
	if _, err := tx.Queries().LockSamplingHabitatSelections(ctx, importID); err != nil {
		return err
	}
	ready, err := tx.Queries().CheckSamplingHabitatsReady(ctx, importID)
	if err != nil {
		return err
	}
	if !ready {
		return ErrHabitatsNotReady
	}
	conflict, err := tx.Queries().HasConflictingHabitatCreations(ctx, importID)
	if err != nil {
		return err
	}
	if conflict {
		return ErrConflictingHabitatCreations
	}
	// Separate statements see habitats committed by another batch after an ON CONFLICT wait.
	if err := tx.Queries().MaterializeRequestedHabitats(ctx, importID); err != nil {
		return err
	}
	if err := tx.Queries().DeleteRedundantCreatedHabitatTargets(ctx, importID); err != nil {
		return err
	}
	if err := tx.Queries().ResolveCreatedHabitatTargets(ctx, importID); err != nil {
		return err
	}
	if err := tx.Queries().MarkHabitatCreationsResolved(ctx, importID); err != nil {
		return err
	}
	pending, err := tx.Queries().HasUnmaterializedHabitatCreations(ctx, importID)
	if err != nil {
		return err
	}
	if pending {
		return ErrHabitatsNotReady
	}
	return tx.Queries().MaterializeSamplingHabitats(ctx, importID)
}
