package services

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"testing"

	"github.com/lsdch/biome/models"
	csvmodels "github.com/lsdch/biome/models/csv"
	"github.com/stretchr/testify/require"
)

func TestExportDelimited(t *testing.T) {
	data := []models.SamplingWithOccurrences{{
		Occurrences: []models.BaseOccurrence{
			{Code: "comma, tab\t quote\" newline\n café"},
			{Code: "second"},
		},
	}}
	for _, tc := range []struct {
		name      string
		format    models.ExportFormat
		delimiter csvmodels.CSVDelimiter
		want      rune
	}{
		{"csv", models.ExportFormatCSV, "", ','},
		{"semicolon", models.ExportFormatCSV, csvmodels.CSVDelimiterSemicolon, ';'},
		{"tsv", models.ExportFormatTSV, "", '\t'},
		{"tsv overrides delimiter", models.ExportFormatTSV, csvmodels.CSVDelimiterComma, '\t'},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			require.NoError(t, NewExportService().Write(&out, data, tc.format, csvmodels.ExportCSVOptions{Delimiter: tc.delimiter}))
			reader := csv.NewReader(&out)
			reader.Comma = tc.want
			records, err := reader.ReadAll()
			require.NoError(t, err)
			require.Len(t, records, 3)
			index := -1
			for i, name := range records[0] {
				if name == "occurrence_code" {
					index = i
				}
			}
			require.NotEqual(t, -1, index)
			require.Equal(t, data[0].Occurrences[0].Code, records[1][index])
			require.Equal(t, "second", records[2][index])
		})
	}
}

func TestExportJSON(t *testing.T) {
	data := []models.SamplingWithOccurrences{{Occurrences: []models.BaseOccurrence{{Code: "é\t\""}}}}
	var out bytes.Buffer
	require.NoError(t, NewExportService().Write(&out, data, models.ExportFormatJSON, csvmodels.ExportCSVOptions{}))
	expected, err := json.Marshal(data)
	require.NoError(t, err)
	require.JSONEq(t, string(expected), out.String())
}

func TestExportEmpty(t *testing.T) {
	for _, format := range []models.ExportFormat{models.ExportFormatJSON, models.ExportFormatCSV, models.ExportFormatTSV} {
		t.Run(string(format), func(t *testing.T) {
			var out bytes.Buffer
			require.NoError(t, NewExportService().Write(&out, nil, format, csvmodels.ExportCSVOptions{}))
			if format == models.ExportFormatJSON {
				require.JSONEq(t, "[]", out.String())
				return
			}
			reader := csv.NewReader(&out)
			if format == models.ExportFormatTSV {
				reader.Comma = '\t'
			}
			records, err := reader.ReadAll()
			require.NoError(t, err)
			require.Equal(t, [][]string{(csvmodels.SamplingWithOccurrencesCSV{}).Header()}, records)
		})
	}
}

type failingExportWriter struct{ err error }

func (w failingExportWriter) Write(p []byte) (int, error) { return 0, w.err }

func TestExportWriterErrors(t *testing.T) {
	expected := errors.New("write failed")
	for _, format := range []models.ExportFormat{models.ExportFormatJSON, models.ExportFormatCSV, models.ExportFormatTSV} {
		require.ErrorIs(t, NewExportService().Write(failingExportWriter{expected}, nil, format, csvmodels.ExportCSVOptions{}), expected)
	}
}

func TestExportInvalidOptions(t *testing.T) {
	for _, tc := range []struct {
		format  models.ExportFormat
		options csvmodels.ExportCSVOptions
	}{
		{models.ExportFormatDWC, csvmodels.ExportCSVOptions{}},
		{models.ExportFormatCSV, csvmodels.ExportCSVOptions{Delimiter: "invalid"}},
		{models.ExportFormatCSV, csvmodels.ExportCSVOptions{QuoteChar: csvmodels.CSVQuoteCharSingle}},
	} {
		var out bytes.Buffer
		require.Error(t, NewExportService().Write(&out, nil, tc.format, tc.options))
		require.Empty(t, out.String())
	}
}
