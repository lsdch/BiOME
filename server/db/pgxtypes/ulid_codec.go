package pgxtypes

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	ulid "github.com/oklog/ulid/v2"
	"github.com/sirupsen/logrus"
)

// =============================================================================
// ULID Codec
// =============================================================================

type ULIDCodec struct{}

func (*ULIDCodec) FormatSupported(format int16) bool {
	return format == pgtype.TextFormatCode
}

func (*ULIDCodec) PreferredFormat() int16 {
	return pgtype.TextFormatCode
}

func (*ULIDCodec) PlanEncode(
	m *pgtype.Map,
	oid uint32,
	format int16,
	value any,
) pgtype.EncodePlan {

	logrus.Printf(
		"ULID PlanEncode oid=%d format=%d type=%T\n",
		oid,
		format,
		value,
	)

	if format != pgtype.TextFormatCode {
		return nil
	}

	switch value.(type) {
	case ulid.ULID, *ulid.ULID:
		return &ulidEncodePlan{}
	default:
		return nil
	}
}

func (*ULIDCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	if format != pgtype.TextFormatCode {
		return nil
	}
	return &ulidScanPlan{}
}

func (*ULIDCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if format != pgtype.TextFormatCode {
		return nil, fmt.Errorf("ulid: unsupported format %d", format)
	}
	if len(src) == 0 {
		return ulid.ULID{}, nil
	}
	return ulid.Parse(string(src))
}

func (*ULIDCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	return (*ULIDCodec)(nil).DecodeDatabaseSQLValue(m, oid, format, src)
}

type ulidEncodePlan struct{}

func (*ulidEncodePlan) Encode(value any, buf []byte) ([]byte, error) {
	logrus.Printf("ULID encode type=%T\n", value)
	switch v := value.(type) {
	case ulid.ULID:
		return append(buf, v.String()...), nil

	case *ulid.ULID:
		if v == nil {
			return nil, nil
		}
		return append(buf, v.String()...), nil

	default:
		return nil, fmt.Errorf("ulid: unsupported type %T", value)
	}
}

type ulidScanPlan struct{}

func (*ulidScanPlan) Scan(src []byte, dst any) error {
	dstULID, ok := dst.(*ulid.ULID)
	if !ok {
		return fmt.Errorf("ulid: expected *ulid.ULID, got %T", dst)
	}
	if len(src) == 0 {
		*dstULID = ulid.ULID{}
		return nil
	}
	id, err := ulid.Parse(string(src))
	if err != nil {
		return fmt.Errorf("ulid: invalid format %q: %w", string(src), err)
	}
	*dstULID = id
	return nil
}

// =============================================================================
// ULID Array Codec
// =============================================================================

type ULIDArrayCodec struct{}

func (*ULIDArrayCodec) FormatSupported(format int16) bool {
	return format == pgtype.TextFormatCode
}

func (*ULIDArrayCodec) PreferredFormat() int16 {
	return pgtype.TextFormatCode
}

func (*ULIDArrayCodec) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	if format != pgtype.TextFormatCode {
		return nil
	}
	return &ulidArrayEncodePlan{}
}

func (*ULIDArrayCodec) PlanScan(m *pgtype.Map, oid uint32, format int16, target any) pgtype.ScanPlan {
	if format != pgtype.TextFormatCode {
		return nil
	}
	return &ulidArrayScanPlan{}
}

func (*ULIDArrayCodec) DecodeDatabaseSQLValue(m *pgtype.Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if format != pgtype.TextFormatCode {
		return nil, fmt.Errorf("ulid[]: unsupported format %d", format)
	}
	if len(src) == 0 {
		return []ulid.ULID{}, nil
	}
	var result []ulid.ULID
	if err := parseULIDArray(string(src), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (*ULIDArrayCodec) DecodeValue(m *pgtype.Map, oid uint32, format int16, src []byte) (any, error) {
	return (*ULIDArrayCodec)(nil).DecodeDatabaseSQLValue(m, oid, format, src)
}

type ulidArrayEncodePlan struct{}

func (*ulidArrayEncodePlan) Encode(value any, buf []byte) ([]byte, error) {
	ulids, ok := value.([]ulid.ULID)
	if !ok {
		return nil, fmt.Errorf("ulid[]: expected []ulid.ULID, got %T", value)
	}
	if len(ulids) == 0 {
		return append(buf, "{}"...), nil
	}
	var b strings.Builder
	b.WriteString("{")
	for i, u := range ulids {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(u.String())
	}
	b.WriteString("}")
	return append(buf, []byte(b.String())...), nil
}

type ulidArrayScanPlan struct{}

func (*ulidArrayScanPlan) Scan(src []byte, dst any) error {
	dstSlice, ok := dst.(*[]ulid.ULID)
	if !ok {
		return fmt.Errorf("ulid[]: expected *[]ulid.ULID, got %T", dst)
	}
	if len(src) == 0 {
		*dstSlice = nil
		return nil
	}
	return parseULIDArray(string(src), dstSlice)
}

func parseULIDArray(s string, dst *[]ulid.ULID) error {
	if s == "{}" || s == "NULL" {
		*dst = []ulid.ULID{}
		return nil
	}
	s = strings.Trim(s, "{}")
	if s == "" {
		*dst = []ulid.ULID{}
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]ulid.ULID, 0, len(parts))
	for i, p := range parts {
		id, err := ulid.Parse(strings.TrimSpace(p))
		if err != nil {
			return fmt.Errorf("ulid[]: invalid at index %d: %w", i, err)
		}
		result = append(result, id)
	}
	*dst = result
	return nil
}
