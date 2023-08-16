package taxonomy

import (
	"context"
	"darco/proto/models"

	"github.com/edgedb/edgedb-go"
)

type Taxon struct {
	ID         edgedb.UUID `edgedb:"id"`
	GBIF_ID    int32
	Name       string                  `json:"name" edgedb:"name"`
	Code       string                  `json:"code" edgedb:"code"`
	Status     string                  `json:"status" edgedb:"status"`
	Anchor     bool                    `json:"anchor" edgedb:"anchor"`
	Authorship edgedb.OptionalStr      `json:"authorship" edgedb:"authorship"`
	Rank       string                  `json:"rank" edgedb:"rank"`
	Modified   edgedb.OptionalDateTime `json:"modified" edgedb:"modified"`
	Created    edgedb.OptionalDateTime `json:"created" edgedb:"created"`
}

func GetAnchorTaxa() (taxa []Taxon, err error) {
	query := "select taxonomy::Taxon { * } filter .anchor"
	err = models.DB.Query(context.Background(), query, &taxa)
	return
}
