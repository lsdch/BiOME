package csvmodels

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

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

type ExportCSVOptions struct {
	Delimiter CSVDelimiter `query:"delimiter"`
	QuoteChar CSVQuoteChar `query:"quoteChar" default:"\""`
}

type SamplingCSV struct {
	models.Sampling
}

func (s SamplingCSV) Header() []string {
	h := [13]string{
		"sampling_id",
		"site_name",
		"site_code",
		"site_locality",
		"country",
		"country_code",
		"latitude",
		"longitude",
		"coordinates_precision_m",
		"altitude_m",
		"sampling_date",
		"sampling_performed_by",
		"sampling_comments",
	}
	return h[:]
}

func (s SamplingCSV) Record() []string {
	record := [13]string{
		s.ID.String(),
		s.Site.Name.GetWithDefault(""),
		s.Site.Code.GetWithDefault(""),
		s.Site.Locality.GetWithDefault(""),
		s.Site.Country.MapString(func(c models.Country) string { return c.Name }),
		s.Site.Country.MapString(func(c models.Country) string { return c.Code }),
		strconv.FormatFloat(s.Coordinates.Latitude, 'f', -1, 64),
		strconv.FormatFloat(s.Coordinates.Longitude, 'f', -1, 64),
		s.Coordinates.Precision.MapString(func(p int32) string { return strconv.FormatInt(int64(p), 10) }),
		s.Altitude.MapString(func(i int32) string { return strconv.FormatInt(int64(i), 10) }),
		s.PerformedOn.MapString(func(dwp models.DateWithPrecision) string { return dwp.String() }),
		strings.Join(s.PerformedBy, ","),
		s.Comments.GetWithDefault(""),
	}
	return record[:]
}

type BaseOccurrenceCSV struct {
	models.BaseOccurrence
}

func (o BaseOccurrenceCSV) Header() []string {
	h := [17]string{
		"occurrence_id",
		"occurrence_code",
		"identification",
		"taxon_id",
		"taxon_name",
		"taxon_authorship",
		"taxon_rank",
		"taxon_status",
		"taxon_gbif_id",
		"identification_date",
		"identified_by",
		"identification_confer",
		"identification_addendum",
		"specimen_quantity",
		"specimen_type_status",
		"content_description",
		"occurrence_comments",
	}
	return h[:]
}

func (occ BaseOccurrenceCSV) Record(innerDelimiter rune) []string {
	record := [17]string{
		occ.ID.String(),
		occ.Code,
		occ.Identification.String(),
		occ.Identification.Taxon.ID.String(),
		occ.Identification.Taxon.Name,
		occ.Identification.Taxon.Authorship.GetWithDefault(""),
		string(occ.Identification.Taxon.Rank),
		string(occ.Identification.Taxon.Status),
		occ.Identification.Taxon.GBIF_ID.MapString(func(id int32) string { return strconv.FormatInt(int64(id), 10) }),
		occ.Identification.IdentifiedOn.MapString(func(dwp models.DateWithPrecision) string { return dwp.String() }),
		strings.Join(occ.Identification.IdentifiedBy, string(innerDelimiter)),
		strconv.FormatBool(occ.Identification.Confer),
		occ.Identification.Addendum.GetWithDefault(""),
		occ.Quantity.MapString(func(q models.OccurrenceQuantity) string { return q.String() }),
		occ.TypeStatus.MapString(func(ts models.OccurrenceTypeStatus) string { return string(ts) }),
		occ.ContentDescription.GetWithDefault(""),
		occ.Comments.GetWithDefault(""),
	}
	return record[:]
}

type SamplingWithOccurrencesCSV struct {
	models.SamplingWithOccurrences
}

func (s SamplingWithOccurrencesCSV) Header() []string {
	return slices.Concat(
		SamplingCSV{}.Header()[:],
		BaseOccurrenceCSV{}.Header()[:],
	)
}

func (s SamplingWithOccurrencesCSV) Records(innerDelimiter rune) iter.Seq[[]string] {
	samplingRecord := SamplingCSV{Sampling: s.SamplingWithOccurrences.Sampling}.Record()
	return func(yield func([]string) bool) {
		for _, occ := range s.SamplingWithOccurrences.Occurrences {
			occRecord := BaseOccurrenceCSV{BaseOccurrence: occ}.Record(innerDelimiter)
			record := slices.Concat(samplingRecord, occRecord)
			if !yield(record) {
				return
			}
		}
	}
}
