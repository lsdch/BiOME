package people

import (
	"darco/proto/models/validations"
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
)

type InstitutionKind string // @name InstitutionKind

const (
	Lab                InstitutionKind = "Lab"
	FoundingAgency     InstitutionKind = "FundingAgency"
	SequencingPlatform InstitutionKind = "SequencingPlatform"
	Other              InstitutionKind = "Other"
)

var values = []InstitutionKind{Lab, FoundingAgency, SequencingPlatform, Other}

func (m *InstitutionKind) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*m), nil
}

func (m *InstitutionKind) UnmarshalEdgeDBStr(data []byte) error {
	*m = InstitutionKind(string(data))
	return nil
}

func institutionKindValidator(fl validator.FieldLevel) bool {
	return slices.Contains(values, fl.Field().Interface().(InstitutionKind))
}

var InstitutionKindValidator = validations.CustomValidator{
	Tag:     "institution_kind",
	Handler: institutionKindValidator,
	Message: func(fl validator.FieldError) string {
		return fmt.Sprintf("Invalid kind: %s", fl.Value())
	},
}

var _ = validations.RegisterCustomValidator(InstitutionKindValidator)
