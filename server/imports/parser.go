package imports

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"github.com/jszwec/csvutil"
	csvmodels "github.com/lsdch/biome/models/csv"
)

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
	ParseBibCSV(reader io.Reader, separator rune) ([]csvmodels.PublicationImportRow, error)
}

func NewCSVParser() CSVParser {
	return &csvParser{
		sanitizer: func(record []string) []string {
			for i, v := range record {
				record[i] = strings.Trim(strings.TrimSpace(v), ",;")

				switch record[i] {
				case "NULL", "N/A", "NA", "-":
					record[i] = ""
				}
			}
			return record
		},
	}
}

type csvParser struct {
	sanitizer func([]string) []string
}

func (p *csvParser) readerWithSanitization(reader io.Reader, separator rune) *SanitizingReader {
	csvReader := csv.NewReader(reader)
	// Disabled to prevent reader from messing without tab-separated files with empty fields
	csvReader.TrimLeadingSpace = false
	csvReader.Comma = separator
	csvReader.FieldsPerRecord = -1

	return &SanitizingReader{
		reader:   csvReader,
		sanitize: p.sanitizer,
	}
}

func (p *csvParser) decoderWithSanitization(reader io.Reader, separator rune) (*csvutil.Decoder, error) {
	sanitizingReader := p.readerWithSanitization(reader, separator)

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

	return dec, nil
}

func (p *csvParser) ParseCSV(reader io.Reader, separator rune) ([]csvmodels.OccurrenceImportRow, error) {
	dec, err := p.decoderWithSanitization(reader, separator)
	if err != nil {
		return nil, err
	}

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
		u.SetRowNumber(row)
		getCSVConfigurator[csvmodels.OccurrenceImportRow]()(&u)
		if err := dec.Decode(&u); err == io.EOF {
			break
		} else if err != nil {
			return nil, &csvmodels.CSVParseError{
				RowNumber: row,
				Err:       err,
			}
		}
		rows = append(rows, u)
		row++
	}
	return rows, nil
}

func (p *csvParser) ParseBibCSV(reader io.Reader, separator rune) ([]csvmodels.PublicationImportRow, error) {
	dec, err := p.decoderWithSanitization(reader, separator)
	if err != nil {
		return nil, err
	}

	var rows []csvmodels.PublicationImportRow

	_ = dec.Header()
	row := int32(2) // Start counting from 2 to account for the header row
	for {
		u := csvmodels.PublicationImportRow{}
		u.SetRowNumber(row)
		if err := dec.Decode(&u); err == io.EOF {
			break
		} else if err != nil {
			return nil, &csvmodels.CSVParseError{
				RowNumber: row,
				Err:       err,
			}
		}
		rows = append(rows, u)
		row++
	}
	return rows, nil
}
