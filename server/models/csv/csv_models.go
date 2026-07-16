package csvmodels

import "github.com/lsdch/biome/models"

type CoordinatesWithPrecisionInput struct {
	models.Coordinates
	PrecisionM *int32 `csv:"coordinates_precision_m,omitempty" validate:"omitempty,gte=0"`
}
