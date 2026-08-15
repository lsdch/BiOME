package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type DOI string

// Value implémente driver.Valuer.
// Elle convertit le type DOI en valeur SQL.
func (d DOI) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}

	return string(d), nil
}

// Scan implémente sql.Scanner.
// Elle convertit une valeur SQL en DOI.
func (d *DOI) Scan(src any) error {
	if src == nil {
		*d = ""
		return nil
	}

	switch v := src.(type) {
	case string:
		*d = DOI(v)
		return nil

	case []byte:
		*d = DOI(string(v))
		return nil

	default:
		return fmt.Errorf("doi: cannot scan type %T", src)
	}
}

func normalizeDOI(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "https://doi.org/")
	s = strings.TrimPrefix(s, "http://doi.org/")
	s = strings.TrimPrefix(s, "doi:")
	s = strings.TrimPrefix(s, "doi/")
	s = strings.TrimPrefix(s, "doi.org/")
	return s
}

func (d *DOI) UnmarshalCSV(text []byte) error {
	s := normalizeDOI(strings.TrimSpace(string(text)))
	*d = DOI(s)
	return nil
}

func (d *DOI) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = normalizeDOI(s)
	*d = DOI(s)
	return nil
}

func ParseDOI(s string) DOI {
	s = normalizeDOI(s)
	return DOI(s)
}

func ParseDOIPtr(s *string) DOI {
	if s == nil {
		return ""
	}
	return ParseDOI(*s)
}

func (d DOI) String() string {
	return string(d)
}
