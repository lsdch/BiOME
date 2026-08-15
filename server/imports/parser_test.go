package imports

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestCSVReaderKeepsEmptyFields(t *testing.T) {
	reader := csv.NewReader(strings.NewReader(
		"site\tlocality\tmeasurement\n" +
			"Baradla Cave\t\tstygobiont\n",
	))

	reader.Comma = '\t'

	header, err := reader.Read()
	if err != nil {
		t.Fatalf("failed reading header: %v", err)
	}

	if len(header) != 3 {
		t.Fatalf("expected 3 header fields, got %d: %#v", len(header), header)
	}

	record, err := reader.Read()
	if err != nil {
		t.Fatalf("failed reading record: %v", err)
	}

	if len(record) != 3 {
		t.Fatalf("expected 3 fields, got %d: %#v", len(record), record)
	}

	if record[0] != "Baradla Cave" {
		t.Errorf("unexpected first field: %q", record[0])
	}

	if record[1] != "" {
		t.Errorf("expected empty locality, got %q", record[1])
	}

	if record[2] != "stygobiont" {
		t.Errorf("unexpected measurement: %q", record[2])
	}
}

func TestSanitizingReaderKeepsEmptyFields(t *testing.T) {
	csvReader := csv.NewReader(strings.NewReader(
		"site\tlocality\tmeasurement\n" +
			"Baradla Cave\t\tstygobiont\n",
	))
	csvReader.Comma = '\t'

	reader := &SanitizingReader{
		reader: csvReader,
		sanitize: func(record []string) []string {
			for i, v := range record {
				record[i] = strings.TrimSpace(v)
			}
			return record
		},
	}

	_, err := reader.Read() // header
	if err != nil {
		t.Fatal(err)
	}

	record, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}

	if len(record) != 3 || record[1] != "" {
		t.Fatalf("empty field was not preserved: %#v", record)
	}
}

func TestCSVReaderKeepsEmptyTSVColumnsWithTrim(t *testing.T) {
	reader := csv.NewReader(strings.NewReader(
		"site\tlocality\tmeasurement\n" +
			"Baradla Cave\t\tstygobiont\n",
	))

	reader.Comma = '\t'
	reader.TrimLeadingSpace = false

	_, _ = reader.Read()

	record, err := reader.Read()
	if err != nil {
		t.Fatal(err)
	}

	if len(record) != 3 || record[1] != "" {
		t.Fatalf("expected empty middle field, got %#v", record)
	}
}
