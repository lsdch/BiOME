package people

import (
	"darco/proto/models/validations"
	"fmt"
	"math/rand"
	"reflect"
	"slices"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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

var _ = faker.AddProvider("institutionKind",
	func(v reflect.Value) (interface{}, error) {
		idx := rand.Intn(len(values))
		x := values[idx]
		logrus.Infof("Called %d: %v", idx, x)
		return InstitutionKind(values[idx]), nil
	})
