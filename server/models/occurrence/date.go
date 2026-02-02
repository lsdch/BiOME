package occurrence

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
)

type DatePrecision string

//generate:enum
const (
	Day   DatePrecision = "Day"
	Month DatePrecision = "Month"
	Year  DatePrecision = "Year"
)

type YearRange struct {
	Min geltypes.OptionalInt32 `gel:"min" json:"min,omitempty"`
	Max geltypes.OptionalInt32 `gel:"max" json:"max,omitempty"`
}

func GetOccurrencesDateRange(db geltypes.Executor) (dateRange YearRange, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		WITH module occurrence,
			min_date := datetime_get(
				(
					SELECT Occurrence
					filter exists .sampling.performed_on.date
					ORDER BY .sampling.performed_on.date ASC LIMIT 1
				).sampling.performed_on.date,
				"year"),
			max_date := datetime_get(
				(
					SELECT Occurrence
					ORDER BY .sampling.performed_on.date DESC LIMIT 1
				).sampling.performed_on.date,
				"year"),
		SELECT {
			min := <int32>min_date,
			max := <int32>max_date,
		};
	`, &dateRange)
	return
}
