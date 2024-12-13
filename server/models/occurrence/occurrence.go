package occurrence

import (
	"darco/proto/models"

	"github.com/edgedb/edgedb-go"
)

type OccurrenceCategory string

//generate:enum
const (
	Internal OccurrenceCategory = "Internal"
	External OccurrenceCategory = "External"
)

type GenericOccurrence[SamplingType any] struct {
	ID             edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Sampling       SamplingType       `edgedb:"sampling" json:"sampling"`
	Identification Identification     `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr `edgedb:"comments" json:"comments"`
}

type Occurrence GenericOccurrence[SamplingInner]

type OccurrenceInput struct {
	SamplingID     edgedb.UUID                  `json:"sampling_id" format:"uuid"`
	Identification IdentificationInput          `edgedb:"identification" json:"identification"`
	Comments       models.OptionalInput[string] `edgedb:"comments" json:"comments"`
}
