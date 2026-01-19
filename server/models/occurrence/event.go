package occurrence

import (
	"time"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
)

type OptionalDateWithPrecision struct {
	geltypes.Optional `json:"-"`
	DateWithPrecision `gel:"$inline" json:",inline"`
}

type DateWithPrecision struct {
	Date      time.Time     `gel:"date" json:"date"`
	Precision DatePrecision `gel:"precision" json:"precision"`
}

func (d DateWithPrecision) ToCode() string {
	switch d.Precision {
	case Year:
		return d.Date.Format("2006")
	default:
		return d.Date.Format("2006-01")
	}

}

type CompositeDate struct {
	Day   int32 `json:"day,omitempty" minimum:"1" maximum:"31" default:"1"`
	Month int32 `json:"month,omitempty" minimum:"1" maximum:"12" default:"1"`
	Year  int32 `json:"year,omitempty" minimum:"1500" maximum:"3000"`
}

type DateWithPrecisionInput struct {
	Date      CompositeDate `json:"date"`
	Precision DatePrecision `json:"precision"`
}

type Action struct {
	PerformedBy []string                  `json:"performed_by,omitempty" gel:"performed_by"`
	PerformedOn OptionalDateWithPrecision `json:"performed_on,omitzero" gel:"performed_on"`
}

type ActionInput struct {
	PerformedBy []string                                     `json:"performed_by,omitempty"`
	PerformedOn models.OptionalInput[DateWithPrecisionInput] `json:"performed_on,omitzero"`
}

// func (ev *ActionInput) WithPersonAliases(aliases map[string]string) *ActionInput {
// 	for i, alias := range ev.PerformedBy {
// 		if _, ok := aliases[alias]; ok {
// 			ev.PerformedBy[i] = aliases[alias]
// 		}
// 	}
// 	return ev
// }

type ActionUpdate struct {
	PerformedBy models.OptionalNull[[]string]                `gel:"performed_by" json:"performed_by,omitempty"`
	PerformedOn models.OptionalInput[DateWithPrecisionInput] `gel:"performed_on" json:"performed_on,omitempty"`
}
