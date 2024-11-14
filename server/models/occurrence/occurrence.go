package occurrence

import "github.com/edgedb/edgedb-go"

type Occurrence struct {
	Sampling       Sampling           `edgedb:"sampling" json:"sampling"`
	Identification Identification     `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr `edgedb:"comments" json:"comments"`
}
