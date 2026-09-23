package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
)

type ExportService struct {
}

func NewExportService() *ExportService {
	return &ExportService{}
}

// ContentType validates the export options before any response is written.
func (s *ExportService) ContentType(format models.ExportFormat, options csvmodels.ExportCSVOptions) (string, error) {
	switch format {
	case models.ExportFormatJSON:
		return "application/json", nil
	case models.ExportFormatCSV, models.ExportFormatTSV:
		if options.QuoteChar != "" && options.QuoteChar != csvmodels.CSVQuoteCharDouble {
			return "", fmt.Errorf("only double quotes are supported for delimited exports")
		}
		if options.Delimiter != "" && options.Delimiter != csvmodels.CSVDelimiterComma &&
			options.Delimiter != csvmodels.CSVDelimiterSemicolon && options.Delimiter != csvmodels.CSVDelimiterTab {
			return "", fmt.Errorf("unsupported delimiter: %q", options.Delimiter)
		}
		if format == models.ExportFormatTSV {
			return "text/tab-separated-values; charset=utf-8", nil
		}
		return "text/csv; charset=utf-8", nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// Write exports the nested JSON representation or one delimited row per occurrence.
// TSV always uses tabs; CSV defaults to commas and accepts an explicit delimiter.
func (s *ExportService) Write(w io.Writer, data []models.SamplingWithOccurrences, format models.ExportFormat, options csvmodels.ExportCSVOptions) error {
	if _, err := s.ContentType(format, options); err != nil {
		return err
	}
	if format == models.ExportFormatJSON {
		if data == nil {
			data = []models.SamplingWithOccurrences{}
		}
		return json.NewEncoder(w).Encode(data)
	}
	delimiter := ','
	if options.Delimiter != "" {
		delimiter = []rune(options.Delimiter)[0]
	}
	if format == models.ExportFormatTSV {
		delimiter = '\t'
	}
	rows := make([]csvmodels.SamplingWithOccurrencesCSV, len(data))
	for i, row := range data {
		rows[i] = csvmodels.SamplingWithOccurrencesCSV{SamplingWithOccurrences: row}
	}
	return s.writeDelimited(w, rows, delimiter, '|')
}

func (s *ExportService) writeDelimited(
	w io.Writer,
	data []csvmodels.SamplingWithOccurrencesCSV,
	delimiter rune,
	innerDelimiter rune,
) error {
	cw := csv.NewWriter(w)
	cw.Comma = delimiter

	if err := cw.Write((csvmodels.SamplingWithOccurrencesCSV{}).Header()); err != nil {
		return err
	}

	for _, row := range data {
		for record := range row.Records(innerDelimiter) {
			if err := cw.Write(record); err != nil {
				return err
			}
		}
	}

	cw.Flush()
	return cw.Error()
}
