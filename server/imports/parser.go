package imports

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jszwec/csvutil"
	csvmodels "github.com/lsdch/biome/models/csv"
)

type CSVParseError struct {
	RowNumber int32
	Err       error
}

func (e *CSVParseError) Error() string {
	return fmt.Sprintf("error parsing CSV at row %d: %v", e.RowNumber, e.Err)
}

type CSVImportParams struct {
	Reader    io.Reader
	Separator rune
}

type SanitizingReader struct {
	reader   *csv.Reader
	sanitize func([]string) []string
}

func (r *SanitizingReader) Read() ([]string, error) {
	record, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	return r.sanitize(record), nil
}

type CSVParser interface {
	ParseCSV(reader io.Reader, separator rune) ([]csvmodels.OccurrenceImportRow, error)
}

func NewCSVParser() CSVParser {
	return &csvParser{}
}

type csvParser struct{}

func (p *csvParser) ParseCSV(reader io.Reader, separator rune) ([]csvmodels.OccurrenceImportRow, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true
	csvReader.Comma = separator

	sanitizingReader := &SanitizingReader{
		reader: csvReader,
		sanitize: func(record []string) []string {
			for i, v := range record {
				record[i] = strings.TrimSpace(v)

				switch record[i] {
				case "NULL", "N/A", "NA", "-":
					record[i] = ""
				}
			}
			return record
		},
	}

	dec, err := csvutil.NewDecoder(sanitizingReader)
	if err != nil {
		return nil, err
	}

	dec.WithUnmarshalers(
		csvutil.UnmarshalFunc(func(data []byte, f *int32) error {
			if string(data) == "" {
				f = nil
				return nil
			}
			val, err := strconv.ParseInt(string(data), 10, 32)
			if err != nil {
				return err
			}
			*f = int32(val)
			return nil
		}),
	)

	// Sanitize fields
	// dec.Map = func(field, column string, v any) string {
	// 	if !utf8.ValidString(field) {
	// 		logrus.Warnf("Valeur %s : %s contient des bytes invalides", field, column)
	// 	}
	// 	return field
	// }

	var rows []csvmodels.OccurrenceImportRow

	_ = dec.Header()
	row := int32(2) // Start counting from 2 to account for the header row
	for {
		u := csvmodels.OccurrenceImportRow{}
		if err := dec.Decode(&u); err == io.EOF {
			break
		} else if err != nil {
			return nil, &CSVParseError{
				RowNumber: row,
				Err:       err,
			}
		}
		rows = append(rows, u)
		row++
	}
	return rows, nil
}
