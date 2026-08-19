package csvmodels

import (
	"fmt"

	"github.com/lsdch/biome/models"
)

type CSVDelimiter string

//generate:enum
const (
	CSVDelimiterComma     CSVDelimiter = ","
	CSVDelimiterSemicolon CSVDelimiter = ";"
	CSVDelimiterTab       CSVDelimiter = "\t"
)

type CSVQuoteChar string

//generate:enum
const (
	CSVQuoteCharDouble CSVQuoteChar = `"`
	CSVQuoteCharSingle CSVQuoteChar = `'`
)

type CoordinatesWithPrecisionInput struct {
	models.Coordinates
	PrecisionM *int32 `csv:"coordinates_precision_m,omitempty" validate:"omitempty,gte=0"`
}

func (i CoordinatesWithPrecisionInput) String() string {
	if i.PrecisionM == nil {
		return fmt.Sprintf("%f,%f[NA]", i.Latitude, i.Longitude)
	}
	return fmt.Sprintf("%f,%f[%dm]", i.Latitude, i.Longitude, *i.PrecisionM)
}
