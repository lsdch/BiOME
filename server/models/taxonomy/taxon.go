package taxonomy

import (
	"context"
	"darco/proto/models"
	_ "embed"
	"fmt"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

// TaxonStatus represents the taxonomic status of a taxon.
type TaxonStatus string // @name TaxonStatus

const (
	Accepted     TaxonStatus = "Accepted"
	Synonym      TaxonStatus = "Synonym"
	Unclassified TaxonStatus = "Unclassified"
)

func (m *TaxonStatus) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*m), nil
}

func (m *TaxonStatus) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonStatus(string(data))
	return nil
}

// TaxonRank represents the taxonomic rank of a taxon.
type TaxonRank string // @name TaxonRank

const (
	Kingdom    TaxonRank = "Kingdom"
	Phylum     TaxonRank = "Phylum"
	Class      TaxonRank = "Class"
	Family     TaxonRank = "Family"
	Genus      TaxonRank = "Genus"
	Species    TaxonRank = "Species"
	Subspecies TaxonRank = "Subspecies"
)

func (m *TaxonRank) MarshalEdgeDBStr() ([]byte, error) {
	return []byte(*m), nil
}

func (m *TaxonRank) UnmarshalEdgeDBStr(data []byte) error {
	*m = TaxonRank(string(data))
	return nil
}

type Taxon struct {
	GBIF_ID    edgedb.OptionalInt32 `edgedb:"GBIF_ID" example:"2206247" validate:"numeric"`
	Name       string               `edgedb:"name" json:"name" example:"Asellus aquaticus" validate:"required, alpha" validatePatch:"alpha"`
	Code       string               `edgedb:"code" json:"code" example:"ASEaquaticus"`
	Status     TaxonStatus          `edgedb:"status" json:"status" example:"Accepted" validate:"required"`
	Authorship edgedb.OptionalStr   `edgedb:"authorship" json:"authorship" example:"(Linnaeus, 1758)"`
	Rank       TaxonRank            `edgedb:"rank" json:"rank" example:"Species" validate:"required"`
} // @name Taxon

// @tags taxonomy
type TaxonDB struct {
	ID     edgedb.UUID `edgedb:"id" json:"id" example:"<UUID>"`
	Taxon  `edgedb:"$inline"`
	Anchor bool        `edgedb:"anchor" json:"anchor"`
	Meta   models.Meta `edgedb:"meta" json:"meta"`
} // @name TaxonDB

type TaxonSelect struct {
	TaxonDB `edgedb:"$inline"`
	Parent  struct {
		edgedb.Optional
		TaxonDB `edgedb:"$inline"`
	} `edgedb:"parent" json:"parent"`
	Children []TaxonDB `edgedb:"children" json:"children,omitempty"`
} // @name TaxonWithRelatives

type TaxonInput struct {
	Taxon  `edgedb:"$inline"`
	Parent string `edgedb:"parent"`
} // @name TaxonInput

func ListTaxa(pattern string, rank TaxonRank, status TaxonStatus) ([]TaxonDB, error) {
	var taxa = make([]TaxonDB, 0)
	query := "select taxonomy::Taxon { *, meta: {**}}"
	expr := &models.FilterGroup{Operator: "AND", Components: make([]models.Expr, 3)}
	if pattern != "" {
		expr.Add(&models.QueryFilter{
			Field: ".name", Op: "ilike", Arg: "<str>$pattern", Value: fmt.Sprintf("%%%s%%", pattern),
		})
	}
	if rank != "" {
		expr.Add(&models.QueryFilter{
			Field: ".rank", Op: "=", Arg: "<taxonomy::Rank>$rank", Value: string(rank),
		})
	}
	if status != "" {
		expr.Add(&models.QueryFilter{
			Field: ".status", Op: "=", Arg: "<taxonomy::TaxonStatus>$status", Value: string(status),
		})
	}
	qb := models.QueryBuilder{Query: query, Expr: expr}
	args := qb.Args()
	logrus.Debugf("Taxonomy list query: %s", qb.String())
	err := models.DB().Query(context.Background(), qb.String(), &taxa, args)
	return taxa, err
}

func GetAnchorTaxa() (taxa []TaxonDB, err error) {
	query := "select taxonomy::Taxon { *, meta: { ** } } filter .anchor"
	err = models.DB().Query(context.Background(), query, &taxa)
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
	err := models.DB().QuerySingle(context.Background(), query, taxon, code)
	return taxon, err
}

func DeleteTaxon(code string) (taxon TaxonDB, err error) {
	query := "delete taxonomy::Taxon filter .code = <str>$0;"
	err = models.DB().QuerySingle(context.Background(), query, &taxon, code)
	return
}

//go:embed update_taxon.edgeql
var updateTaxonCmd string

func UpdateTaxon(code string, taxon TaxonInput) (TaxonSelect, error) {
	res := TaxonSelect{}
	args := models.StructToMap(taxon)
	args["target"] = code
	err := models.DB().QuerySingle(context.Background(), updateTaxonCmd, &res, args)
	return res, err
}
