package models

import (
	"fmt"
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
		var exact int32
		_, err := fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &exact)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		q.Exact = NewOptional(exact)
	} else if len(parts) == 2 {
		var lower, upper int32
		_, err := fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &lower)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		_, err = fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &upper)
		if err != nil {
			return fmt.Errorf("invalid quantity format: %s", str)
		}
		q.Lower = NewOptional(lower)
		q.Upper = NewOptional(upper)
	} else {
		return fmt.Errorf("invalid quantity format: %s", str)
	}
	return nil
}
