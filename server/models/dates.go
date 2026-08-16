package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
)

type EventDatePrecision = biomedb.EventDatePrecision

type DateWithPrecision struct {
	Date      time.Time          `json:"date"`
	Precision EventDatePrecision `json:"precision"`
}

func (d DateWithPrecision) UpperBound() DateWithPrecision {
	switch d.Precision {
	case biomedb.EventDatePrecisionYear:
		return DateWithPrecision{
			Date:      time.Date(d.Date.Year(), 12, 31, 0, 0, 0, 0, time.UTC),
			Precision: biomedb.EventDatePrecisionYear,
		}
	case biomedb.EventDatePrecisionMonth:
		lastDay := time.Date(d.Date.Year(), d.Date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		return DateWithPrecision{
			Date:      time.Date(d.Date.Year(), d.Date.Month(), lastDay, 0, 0, 0, 0, time.UTC),
			Precision: biomedb.EventDatePrecisionMonth,
		}
	default:
		return d
	}
}

func ParseDateWithPrecision(dateStr string) (*DateWithPrecision, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return nil, nil
	}
	var d DateWithPrecision
	switch len(dateStr) {
	case 4:
		// Year precision
		parsedDate, err := time.Parse("2006", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid year format: %s", dateStr)
		}
		d.Date = parsedDate
		d.Precision = biomedb.EventDatePrecisionYear
	case 7:
		// Month precision
		parsedDate, err := time.Parse("2006-01", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid month format: %s", dateStr)
		}
		d.Date = parsedDate
		d.Precision = biomedb.EventDatePrecisionMonth
	case 10:
		// Day precision
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid day format: %s", dateStr)
		}
		d.Date = parsedDate
		d.Precision = biomedb.EventDatePrecisionDay
	default:
		return nil, fmt.Errorf("invalid date format: %s", dateStr)
	}
	return &d, nil
}

func (d DateWithPrecision) String() string {
	switch d.Precision {
	case biomedb.EventDatePrecisionYear:
		return d.Date.Format("2006")
	case biomedb.EventDatePrecisionMonth:
		return d.Date.Format("2006-01")
	default:
		return d.Date.Format("2006-01-02")
	}
}

func MaybeDateWithPrecisionFromDB(date pgtype.Date, precision *biomedb.EventDatePrecision) Optional[DateWithPrecision] {
	d := Optional[DateWithPrecision]{}
	if date.Valid && precision.Valid() {
		d = NewOptional(DateWithPrecision{Date: date.Time, Precision: *precision})
	}
	return d
}

type CompositeDate struct {
	Day   int32 `json:"day,omitempty" minimum:"1" maximum:"31" default:"1" query:"day"`
	Month int32 `json:"month,omitempty" minimum:"1" maximum:"12" default:"1" query:"month"`
	Year  int32 `json:"year,omitempty" minimum:"1500" maximum:"3000" query:"year"`
}

func (d CompositeDate) ToTime() time.Time {
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
}

func (d CompositeDate) ToPgDate() pgtype.Date {
	return pgtype.Date{Valid: true, Time: d.ToTime()}
}

// Add adds the specified value to the CompositeDate based on the given unit (day, month, or year).
//
// A negative value will subtract from the date, while a positive value will add to it.
// The unit parameter determines whether the value is added to the day, month, or year component of the CompositeDate.
func (d *CompositeDate) Add(v int32, unit EventDatePrecision) {
	t := d.ToTime()
	switch unit {
	case biomedb.EventDatePrecisionDay:
		t = t.AddDate(0, 0, int(v))
	case biomedb.EventDatePrecisionMonth:
		t = t.AddDate(0, int(v), 0)
	case biomedb.EventDatePrecisionYear:
		t = t.AddDate(int(v), 0, 0)
	}

	d.Day = int32(t.Day())
	d.Month = int32(t.Month())
	d.Year = int32(t.Year())
}

type DateWithPrecisionInput struct {
	Date      CompositeDate      `json:"date"`
	Precision EventDatePrecision `json:"precision"`
}

type EventDateInput struct {
	Date      time.Time
	Precision EventDatePrecision
	IsSet     bool
}

func (d *EventDateInput) PrecisionPtr() *EventDatePrecision {
	if !d.IsSet {
		return nil
	}
	return &d.Precision
}

func (d *EventDateInput) DatePg() pgtype.Date {
	if !d.IsSet {
		return pgtype.Date{}
	}
	return pgtype.Date{Valid: true, Time: d.Date}
}

func (d *EventDateInput) String() string {
	if !d.IsSet {
		return ""
	}
	switch d.Precision {
	case biomedb.EventDatePrecisionDay:
		return d.Date.Format("2006-01-02")
	case biomedb.EventDatePrecisionMonth:
		return d.Date.Format("2006-01")
	case biomedb.EventDatePrecisionYear:
		return d.Date.Format("2006")
	default:
		return d.Date.Format("2006-01-02")
	}
}

func (d *EventDateInput) UnmarshalCSV(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" {
		d.IsSet = false
		return nil
	}
	date, err := time.Parse("2006-01-02", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionDay
		d.IsSet = true
		return nil
	}
	date, err = time.Parse("2006-01-2", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionDay
		d.IsSet = true
		return nil
	}
	date, err = time.Parse("2006-01", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionMonth
		d.IsSet = true
		return nil
	}
	date, err = time.Parse("2006", str)
	if err == nil {
		d.Date = date
		d.Precision = biomedb.EventDatePrecisionYear
		d.IsSet = true
		return nil
	}
	return fmt.Errorf("invalid date format: %s", str)
}
