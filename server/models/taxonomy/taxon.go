package taxonomy

import (
	"context"
	"darco/proto/models"
	_ "embed"

	"github.com/edgedb/edgedb-go"
)

type TaxonStatus string

const (
	Accepted     TaxonStatus = "Accepted"
	Synonym      TaxonStatus = "Synonym"
	Unclassified TaxonStatus = "Unclassified"
)

func (m *TaxonStatus) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonStatus(string(data))
	return nil
}

type TaxonRank string

const (
	Kingdom    TaxonRank = "Kingdom"
	Phylum     TaxonRank = "Phylum"
	Class      TaxonRank = "Class"
	Family     TaxonRank = "Family"
	Genus      TaxonRank = "Genus"
	Species    TaxonRank = "Species"
	Subspecies TaxonRank = "Subspecies"
)

func (m *TaxonRank) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonRank(string(data))
	return nil
}

type Taxon struct {
	GBIF_ID    edgedb.OptionalInt32 `example:"2206247" validate:"numeric"`
	Name       string               `edgedb:"name" json:"name" example:"Asellus aquaticus" validate:"required, alpha" validatePatch:"alpha"`
	Code       string               `edgedb:"code" json:"code" example:"ASEaquaticus"`
	Status     TaxonStatus          `edgedb:"status" json:"status" example:"Accepted" validate:"required,oneof=TaxonStatus"`
	Authorship edgedb.OptionalStr   `edgedb:"authorship" json:"authorship" example:"(Linnaeus, 1758)"`
	Rank       TaxonRank            `edgedb:"rank" json:"rank" example:"Species" validate:"required"`
}

type TaxonDB struct {
	ID     edgedb.UUID `edgedb:"id" json:"id" example:"<UUID>"`
	Taxon  `edgedb:"$inline"`
	Slug   string      `edgedb:"slug" json:"slug" example:"asellus-aquaticus"`
	Anchor bool        `edgedb:"anchor" json:"anchor"`
	Meta   models.Meta `edgedb:"meta" json:"meta"`
}

type TaxonSelect struct {
	TaxonDB `edgedb:"$inline"`
	Parent  struct {
		edgedb.Optional
		TaxonDB `edgedb:"$inline"`
	} `edgedb:"parent" json:"parent"`
	Children []TaxonDB `edgedb:"children" json:"children,omitempty"`
}

type TaxonInput struct {
	Taxon  `edgedb:"$inline"`
	Parent string `edgedb:"parent"`
}

func GetAnchorTaxa() (taxa []TaxonDB, err error) {
	query := "select taxonomy::Taxon { *, meta: { ** } } filter .anchor"
	err = models.DB.Query(context.Background(), query, &taxa)
	return
}

func GetTaxon(code string) (*TaxonSelect, error) {
	taxon := &TaxonSelect{}
	query := `
		select taxonomy::Taxon { *,
			parent : { * , meta: { ** }},
			children : { * , meta: { ** }}
		}
		filter .code = <str>$0;
	`
	err := models.DB.QuerySingle(context.Background(), query, taxon, code)
	return taxon, err
}

func DeleteTaxon(code string) (taxon TaxonDB, err error) {
	query := "delete taxonomy::Taxon filter .code = <str>$0;"
	err = models.DB.QuerySingle(context.Background(), query, &taxon, code)
	return
}

//go:embed update_taxon.edgeql
var updateTaxonCmd string

func UpdateTaxon(code string, taxon TaxonInput) (TaxonSelect, error) {
	res := TaxonSelect{}
	args := models.StructToMap(taxon)
	args["target"] = code
	err := models.DB.QuerySingle(context.Background(), updateTaxonCmd, &res, args)
	return res, err
}
