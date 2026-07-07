package imports

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/jszwec/csvutil"
	"github.com/lsdch/biome/models"
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

type CSVParser interface {
	ParseCSV(reader io.Reader, separator rune) ([]models.OccurrenceImportRow, error)
}

func NewCSVParser() CSVParser {
	return &csvParser{}
}

type csvParser struct{}

func (p *csvParser) ParseCSV(reader io.Reader, separator rune) ([]models.OccurrenceImportRow, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = separator
	dec, err := csvutil.NewDecoder(csvReader)
	if err != nil {
		return nil, err
	}

	// Sanitize fields
	dec.Map = func(field, column string, v any) string {
		s := strings.TrimSpace(field)
		if s == "NA" || s == "N/A" {
			return ""
		}
		return s
	}

	var rows []models.OccurrenceImportRow

	_ = dec.Header()
	row := int32(2) // Start counting from 2 to account for the header row
	for {
		u := models.OccurrenceImportRow{}

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
