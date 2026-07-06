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

func (d DateWithPrecision) ToCode() string {
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
	Day   int32 `json:"day,omitempty" minimum:"1" maximum:"31" default:"1"`
	Month int32 `json:"month,omitempty" minimum:"1" maximum:"12" default:"1"`
	Year  int32 `json:"year,omitempty" minimum:"1500" maximum:"3000"`
}

func (d CompositeDate) ToTime() time.Time {
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 0, 0, 0, 0, time.UTC)
}

func (d CompositeDate) ToPgDate() pgtype.Date {
	return pgtype.Date{Valid: true, Time: d.ToTime()}
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
		return d.Date.Format(time.RFC3339)
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
