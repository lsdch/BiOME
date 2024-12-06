package occurrence

import "github.com/edgedb/edgedb-go"

type Occurrence struct {
	ID             edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Sampling       SamplingInner      `edgedb:"sampling" json:"sampling"`
	Identification Identification     `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr `edgedb:"comments" json:"comments"`
}
