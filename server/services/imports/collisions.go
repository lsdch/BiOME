package imports

import (
	"context"
	"fmt"

	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	md "github.com/lsdch/biome/models"
	"github.com/lsdch/biome/stores"
)

type OccurrenceCollisionsService struct {
	store *stores.CollisionStore
}

func NewOccurrenceCollisionsService() *OccurrenceCollisionsService {
	return &OccurrenceCollisionsService{store: stores.NewCollisionStore()}
}

type SamplingBaseProps struct {
	Coordinates md.CoordinatesWithPrecision       `json:"coordinates"`
	PerformedOn md.Optional[md.DateWithPrecision] `json:"event_date,omitempty"`
}

type SamplingCollision struct {
	SamplingBaseProps
	DistanceMeters int32 `json:"distance_meters"`
}

type StagingCollision struct {
	RowNumber int32 `json:"row_number"`
	SamplingCollision
}
type StagingSamplingWithCollisions struct {
	SamplingBaseProps
	OccurrenceRows     []int32             `json:"row_numbers"`
	StagingCollisions  []StagingCollision  `json:"staging_collisions"`
	ExistingCollisions []SamplingCollision `json:"existing_collisions"`
}

func (c *StagingSamplingWithCollisions) AddStagingCollision(collision biomedb.DetectBatchSamplingCollisionsRow, rowNumber int32) {
	c.StagingCollisions = append(c.StagingCollisions, StagingCollision{
		RowNumber:         rowNumber,
		SamplingCollision: samplingCollisionFromRow(collision),
	})
}

func (c *StagingSamplingWithCollisions) AddExistingCollision(collision biomedb.DetectBatchSamplingCollisionsRow) {
	c.ExistingCollisions = append(c.ExistingCollisions, samplingCollisionFromRow(collision))
}

func samplingCollisionFromRow(row biomedb.DetectBatchSamplingCollisionsRow) SamplingCollision {
	return SamplingCollision{
		SamplingBaseProps: SamplingBaseProps{
			Coordinates: md.CoordinatesWithPrecisionFromDB(
				row.MatchLatitude,
				row.MatchLongitude,
				row.MatchCoordinatesPrecision,
			),
			PerformedOn: md.MaybeDateWithPrecisionFromDB(row.MatchEventDate, row.MatchEventDatePrecision),
		},
		DistanceMeters: row.DistanceMeters,
	}
}

func (r *OccurrenceCollisionsService) DetectSamplingCollisions(ctx context.Context, q db.Querier, importHash string, params stores.CollisionDetectionParams) (collisionsMap map[string]*StagingSamplingWithCollisions, err error) {

	collisionsMap = make(map[string]*StagingSamplingWithCollisions)
	collisions, err := r.store.DetectBatchSamplingCollisions(ctx, q, importHash, params)
	if err != nil {
		return nil, err
	}
	for _, collision := range collisions {
		hash := collision.SamplingHash
		if _, exists := collisionsMap[hash]; !exists {
			collisionsMap[hash] = &StagingSamplingWithCollisions{
				SamplingBaseProps: SamplingBaseProps{
					Coordinates: md.CoordinatesWithPrecisionFromDB(
						collision.Latitude,
						collision.Longitude,
						collision.CoordinatesPrecision,
					),
					PerformedOn: md.MaybeDateWithPrecisionFromDB(collision.EventDate, collision.EventDatePrecision),
				},
				OccurrenceRows: []int32{collision.RowNumber},
			}
		} else {
			collisionsMap[hash].OccurrenceRows = append(collisionsMap[hash].OccurrenceRows, collision.RowNumber)
		}
		if collision.DuplicateSource == biomedb.DuplicateSourceExisting {
			collisionsMap[hash].AddExistingCollision(collision)
		} else {
			collisionsMap[hash].AddStagingCollision(collision, collision.RowNumber)
		}
	}
	return collisionsMap, nil
}

type OccurrenceCollisionsAtRow struct {
	RowNumber int32 `json:"row_number"`
	OccurrenceCollision
	StagingCollisions  []OccurrenceCollisionStaging `json:"staging_collisions"`
	ExistingCollisions []OccurrenceCollision        `json:"existing_collisions"`
}

type OccurrenceCollision struct {
	SamplingCollision
	TaxonName       string              `json:"taxon_name"`
	TaxonAuthorship md.Optional[string] `json:"taxon_authorship,omitempty"`
}

type OccurrenceCollisionStaging struct {
	RowNumber int32 `json:"row_number"`
	OccurrenceCollision
}

func occurrenceCollisionFromRow(row biomedb.DetectBatchOccurrenceCollisionsRow) OccurrenceCollision {
	return OccurrenceCollision{
		SamplingCollision: SamplingCollision{
			SamplingBaseProps: SamplingBaseProps{
				Coordinates: md.CoordinatesWithPrecisionFromDB(
					row.MatchLatitude,
					row.MatchLongitude,
					row.MatchCoordinatesPrecision,
				),
				PerformedOn: md.MaybeDateWithPrecisionFromDB(row.MatchEventDate, row.MatchEventDatePrecision),
			},
			DistanceMeters: row.DistanceMeters,
		},
		TaxonName:       row.MatchTaxonName,
		TaxonAuthorship: md.NewOptionalFromPtr(row.MatchTaxonAuthorship),
	}
}

func (r *OccurrenceCollisionsService) DetectOccurrenceCollisions(ctx context.Context, q db.Querier, importHash string, params stores.CollisionDetectionParams) (collisions []OccurrenceCollisionsAtRow, err error) {
	collisions = make([]OccurrenceCollisionsAtRow, 0)

	collisionRows, err := stores.NewCollisionStore().DetectBatchOccurrencesCollisions(ctx, q, importHash, params)
	if err != nil {
		return nil, err
	}
	for _, row := range collisionRows {
		collision := OccurrenceCollisionsAtRow{
			RowNumber: row.RowNumber,
			OccurrenceCollision: OccurrenceCollision{
				SamplingCollision: SamplingCollision{
					SamplingBaseProps: SamplingBaseProps{
						Coordinates: md.CoordinatesWithPrecisionFromDB(
							row.Latitude,
							row.Longitude,
							row.CoordinatesPrecision,
						),
						PerformedOn: md.MaybeDateWithPrecisionFromDB(row.EventDate, row.EventDatePrecision),
					},
					DistanceMeters: row.DistanceMeters,
				},
				TaxonName:       row.TaxonName,
				TaxonAuthorship: md.NewOptionalFromPtr(row.TaxonAuthorship),
			},
		}

		if row.DuplicateSource == biomedb.DuplicateSourceExisting {
			collision.ExistingCollisions = append(collision.ExistingCollisions, occurrenceCollisionFromRow(row))
		} else {
			if row.MatchRowNumber == nil {
				return nil, fmt.Errorf("staging collision without match row number for row %d", row.RowNumber)
			}
			collision.StagingCollisions = append(collision.StagingCollisions, OccurrenceCollisionStaging{
				RowNumber:           *row.MatchRowNumber,
				OccurrenceCollision: occurrenceCollisionFromRow(row),
			})
		}
		collisions = append(collisions, collision)
	}
	return collisions, nil
}
