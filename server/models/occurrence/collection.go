package occurrence

import (
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
)

type Collection struct {
	ID          geltypes.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Label       string               `gel:"label" json:"label"`
	Code        string               `gel:"code" json:"code"`
	Description geltypes.OptionalStr `gel:"description" json:"description,omitempty"`
	Contact     geltypes.OptionalStr `gel:"contact" json:"contact,omitempty"`
	Meta        people.Meta          `gel:"meta" json:"meta"`
}

type CollectionWithVouchers struct {
	Collection `json:",inline" gel:"$inline"`
	Vouchers   []string `gel:"vouchers" json:"vouchers,omitempty"`
}

type CollectionInput struct {
	Label    string                       `json:"label" validate:"required"`
	Code     string                       `json:"code" validate:"required"`
	Contact  models.OptionalInput[string] `json:"contact,omitempty"`
	Vouchers []string                     `json:"vouchers,omitempty"`
}

type CollectionField struct {
	Name     string   `gel:"name" json:"name"`
	Vouchers []string `gel:"vouchers" json:"vouchers,omitempty"`
}
