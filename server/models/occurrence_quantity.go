package models

import (
	"fmt"
	"strconv"
	"strings"
)

type OccurrenceQuantity struct {
	Exact int32 `json:"exact,omitzero"`
	Lower int32 `json:"lower,omitzero"`
	Upper int32 `json:"upper,omitzero"`
}

func NewOptionalOccurrenceQuantity(exact, lower, upper *int32) Optional[OccurrenceQuantity] {
	if exact != nil {
		return NewOptional(OccurrenceQuantity{Exact: *exact})
	} else if lower != nil && upper != nil {
		return NewOptional(OccurrenceQuantity{Lower: *lower, Upper: *upper})
	} else {
		return Optional[OccurrenceQuantity]{}
	}
}

func (q OccurrenceQuantity) String() string {
	if q.Exact > 0 {
		return fmt.Sprintf("%d", q.Exact)
	} else if q.Lower > 0 && q.Upper > 0 {
		return fmt.Sprintf("%d-%d", q.Lower, q.Upper)
	} else {
		return ""
	}
}

type QuantityInput struct {
	Exact Optional[int32] `json:"exact,omitzero"`
	Lower Optional[int32] `json:"lower,omitzero"`
	Upper Optional[int32] `json:"upper,omitzero"`
}

func (q *QuantityInput) UnmarshalCSV(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" {
		return nil
	}
	parts := strings.Split(str, "-")
	if len(parts) == 1 {
		exact, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		q.Exact = NewOptional(int32(exact))
	} else if len(parts) == 2 {
		lower, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		upper, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 32)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		q.Lower = NewOptional(int32(lower))
		q.Upper = NewOptional(int32(upper))
	} else {
		return fmt.Errorf("invalid quantity format: %s", str)
	}
	return nil
}
