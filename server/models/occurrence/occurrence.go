package occurrence

import "github.com/edgedb/edgedb-go"

type GenericOccurrence[SamplingType any] struct {
	ID             edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Sampling       SamplingType       `edgedb:"sampling" json:"sampling"`
	Identification Identification     `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr `edgedb:"comments" json:"comments"`
}

type Occurrence GenericOccurrence[SamplingInner]
