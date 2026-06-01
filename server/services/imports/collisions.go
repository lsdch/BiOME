package imports

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/db/stores"
)

type OccurrenceCollisionsService struct {
	q *stores.CollisionStore
}

func NewOccurrenceCollisionsService(q *biomedb.Queries) *OccurrenceCollisionsService {
	return &OccurrenceCollisionsService{q: stores.NewCollisionStore(q)}
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type CoordinatesWithPrecision struct {
	Coordinates
	PrecisionM pgtype.Int4 `json:"precision"`
}

type EventDate struct {
	Date      pgtype.Date                    `json:"date"`
	Precision biomedb.NullEventDatePrecision `json:"precision"`
}

type SamplingBaseProps struct {
	Coordinates CoordinatesWithPrecision `json:"coordinates"`
	EventDate   EventDate                `json:"event_date"`
}

type SamplingCollision struct {
	SamplingBaseProps
	DistanceMeters int32 `json:"distance_meters"`
}

type StagingSamplingWithCollisions struct {
	SamplingBaseProps
	OccurrenceRows    []int32 `json:"row_numbers"`
	StagingCollisions []struct {
		RowNumber pgtype.Int4 `json:"row_number"`
		SamplingCollision
	} `json:"staging_collisions" nameHint:"SamplingCollisionStaging"`
	ExistingCollisions []SamplingCollision `json:"existing_collisions"`
}

func (c *StagingSamplingWithCollisions) AddStagingCollision(collision biomedb.DetectBatchSamplingCollisionsRow, rowNumber int32) {
	c.StagingCollisions = append(c.StagingCollisions, struct {
		RowNumber pgtype.Int4 `json:"row_number"`
		SamplingCollision
	}{
		RowNumber:         collision.MatchRowNumber,
		SamplingCollision: samplingCollisionFromRow(collision),
	})
}

func (c *StagingSamplingWithCollisions) AddExistingCollision(collision biomedb.DetectBatchSamplingCollisionsRow) {
	c.ExistingCollisions = append(c.ExistingCollisions, samplingCollisionFromRow(collision))
}

func samplingCollisionFromRow(row biomedb.DetectBatchSamplingCollisionsRow) SamplingCollision {
	return SamplingCollision{
		SamplingBaseProps: SamplingBaseProps{
			Coordinates: CoordinatesWithPrecision{
				Coordinates: Coordinates{
					Latitude:  row.MatchLatitude,
					Longitude: row.MatchLongitude,
				},
				PrecisionM: row.MatchCoordinatesPrecision,
			},
			EventDate: EventDate{
				Date:      row.MatchEventDate,
				Precision: row.MatchEventDatePrecision,
			},
		},
		DistanceMeters: row.DistanceMeters,
	}
}

func (r *OccurrenceCollisionsService) DetectSamplingCollisions(ctx context.Context, importHash string, params stores.CollisionDetectionParams) (collisionsMap map[string]*StagingSamplingWithCollisions, err error) {

	collisionsMap = make(map[string]*StagingSamplingWithCollisions)
	collisions, err := r.q.DetectBatchSamplingCollisions(ctx, importHash, params)
	if err != nil {
		return nil, err
	}
	for _, collision := range collisions {
		hash := collision.SamplingHash
		if _, exists := collisionsMap[hash]; !exists {
			collisionsMap[hash] = &StagingSamplingWithCollisions{
				SamplingBaseProps: SamplingBaseProps{
					Coordinates: CoordinatesWithPrecision{
						Coordinates: Coordinates{
							Latitude:  collision.Latitude,
							Longitude: collision.Longitude,
						},
						PrecisionM: collision.CoordinatesPrecision,
					},
					EventDate: EventDate{
						Date:      collision.EventDate,
						Precision: collision.EventDatePrecision,
					},
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
	SamplingBaseProps
	TaxonName          string                       `json:"taxon_name"`
	TaxonAuthorship    pgtype.Text                  `json:"taxon_authorship"`
	StagingCollisions  []OccurrenceCollisionStaging `json:"staging_collisions"`
	ExistingCollisions []OccurrenceCollision        `json:"existing_collisions"`
}

type OccurrenceCollision struct {
	SamplingCollision
	TaxonName       string      `json:"taxon_name"`
	TaxonAuthorship pgtype.Text `json:"taxon_authorship"`
}

type OccurrenceCollisionStaging struct {
	RowNumber int32 `json:"row_number"`
	OccurrenceCollision
}

func occurrenceCollisionFromRow(row biomedb.DetectBatchOccurrenceCollisionsRow) OccurrenceCollision {
	return OccurrenceCollision{
		SamplingCollision: SamplingCollision{
			SamplingBaseProps: SamplingBaseProps{
				Coordinates: CoordinatesWithPrecision{
					Coordinates: Coordinates{
						Latitude:  row.MatchLatitude,
						Longitude: row.MatchLongitude,
					},
					PrecisionM: row.MatchCoordinatesPrecision,
				},
				EventDate: EventDate{
					Date:      row.MatchEventDate,
					Precision: row.MatchEventDatePrecision,
				},
			},
			DistanceMeters: row.DistanceMeters,
		},
		TaxonName:       row.MatchTaxonName,
		TaxonAuthorship: row.MatchTaxonAuthorship,
	}
}

func (r *OccurrenceCollisionsService) DetectOccurrenceCollisions(ctx context.Context, importHash string, params stores.CollisionDetectionParams) (collisions []OccurrenceCollisionsAtRow, err error) {
	collisions = make([]OccurrenceCollisionsAtRow, 0)

	collisionRows, err := r.q.DetectBatchOccurrencesCollisions(ctx, importHash, params)
	if err != nil {
		return nil, err
	}
	for _, row := range collisionRows {
		collision := OccurrenceCollisionsAtRow{
			RowNumber: row.RowNumber,
			SamplingBaseProps: SamplingBaseProps{
				Coordinates: CoordinatesWithPrecision{
					Coordinates: Coordinates{
						Latitude:  row.Latitude,
						Longitude: row.Longitude,
					},
					PrecisionM: row.CoordinatesPrecision,
				},
				EventDate: EventDate{
					Date:      row.EventDate,
					Precision: row.EventDatePrecision,
				},
			},
			TaxonName:       row.TaxonName,
			TaxonAuthorship: row.TaxonAuthorship,
		}
		if row.DuplicateSource == biomedb.DuplicateSourceExisting {
			collision.ExistingCollisions = append(collision.ExistingCollisions, occurrenceCollisionFromRow(row))
		} else {
			if !row.MatchRowNumber.Valid {
				return nil, fmt.Errorf("staging collision without match row number for row %d", row.RowNumber)
			}
			collision.StagingCollisions = append(collision.StagingCollisions, OccurrenceCollisionStaging{
				RowNumber:           row.MatchRowNumber.Int32,
				OccurrenceCollision: occurrenceCollisionFromRow(row),
			})
		}
		collisions = append(collisions, collision)
	}
	return collisions, nil
}
