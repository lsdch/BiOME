package taxonomy

import (
	"context"
	"darco/proto/db"
	"darco/proto/db/expr"
	"darco/proto/models"
	_ "embed"
	"encoding/json"
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
	GBIF_ID    edgedb.OptionalInt32 `edgedb:"GBIF_ID" json:"GBIF_ID" example:"2206247" binding:"numeric"`
	Name       string               `edgedb:"name" json:"name" example:"Asellus aquaticus" binding:"required, alpha"`
	Code       string               `edgedb:"code" json:"code" example:"ASEaquaticus" binding:"required"`
	Status     TaxonStatus          `edgedb:"status" json:"status" example:"Accepted" binding:"required"`
	Authorship edgedb.OptionalStr   `edgedb:"authorship" json:"authorship" example:"(Linnaeus, 1758)"`
	Rank       TaxonRank            `edgedb:"rank" json:"rank" example:"Species" binding:"required"`
} // @name Taxon

type TaxonInput struct {
	Taxon  `edgedb:"$inline"`
	Parent edgedb.UUID
} // @name TaxonInput

type TaxonDB struct {
	ID     edgedb.UUID `edgedb:"id" json:"id" example:"<UUID>" binding:"required"`
	Taxon  `edgedb:"$inline"`
	Anchor bool        `edgedb:"anchor" json:"anchor"`
	Meta   models.Meta `edgedb:"meta" json:"meta" binding:"required"`
} // @name TaxonDB

type TaxonSelect struct {
	TaxonDB `edgedb:"$inline"`
	Parent  struct {
		edgedb.Optional
		TaxonDB `edgedb:"$inline"`
	} `edgedb:"parent" json:"parent"`
	Children []TaxonDB `edgedb:"children" json:"children,omitempty"`
} // @name TaxonWithRelatives

func ListTaxa(db *edgedb.Client,
	pattern string, rank TaxonRank, status TaxonStatus,
) ([]TaxonDB, error) {
	var taxa = make([]TaxonDB, 0)
	query := "select taxonomy::Taxon { *, meta: {**}}"
	e := &expr.FilterGroup{Operator: "AND", Components: make([]expr.Expr, 3)}
	if pattern != "" {
		e.Add(&expr.QueryFilter{
			Field: ".name", Op: "ilike", Arg: "<str>$pattern", Value: fmt.Sprintf("%%%s%%", pattern),
		})
	}
	if rank != "" {
		e.Add(&expr.QueryFilter{
			Field: ".rank", Op: "=", Arg: "<taxonomy::Rank>$rank", Value: string(rank),
		})
	}
	if status != "" {
		e.Add(&expr.QueryFilter{
			Field: ".status", Op: "=", Arg: "<taxonomy::TaxonStatus>$status", Value: string(status),
		})
	}
	qb := expr.QueryBuilder{Query: query, Expr: e}
	args := qb.Args()
	logrus.Debugf("Taxonomy list query: %s", qb.String())
	err := db.Query(context.Background(), qb.String(), &taxa, args)
	return taxa, err
}

func ListAnchorTaxa(db *edgedb.Client) (taxa []TaxonDB, err error) {
	query := "select taxonomy::Taxon { *, meta: { ** } } filter .anchor"
	err = db.Query(context.Background(), query, &taxa)
	return
}

func FindByCode(db *edgedb.Client, code string) (taxon TaxonSelect, err error) {
	query := `
		select taxonomy::Taxon { *,
			parent : { * , meta: { ** }},
			children : { * , meta: { ** }}
		}
		filter .code = <str>$0;
	`
	err = db.QuerySingle(context.Background(), query, &taxon, code)
	return taxon, err
}

func Delete(code string) error {
	query := "delete taxonomy::Taxon filter .code = <str>$0;"
	return db.Client().Execute(context.Background(), query, code)
}

//go:embed queries/create_taxon.edgeql
var createTaxonCmd string

func (taxon TaxonInput) Create(db *edgedb.Client) (created TaxonSelect, err error) {
	args, _ := json.Marshal(taxon)
	err = db.QuerySingle(context.Background(), createTaxonCmd, &created, args)
	return created, err
}

type TaxonUpdate struct {
	GBIF_ID    *int32       `json:"GBIF_ID,omitempty" binding:"omitnil,numeric"`
	Name       *string      `json:"name,omitempty" binding:"omitnil,alpha"`
	Code       *string      `json:"code,omitempty" binding:"omitnil,alpha"`
	Status     *TaxonStatus `json:"status,omitempty" binding:"omitnil"`
	Authorship *string      `json:"authorship,omitempty" binding:"omitnil,alphanumunicode"`
	Rank       *TaxonRank   `json:"rank,omitempty" binding:"omitnil"`
	Parent     *string      `json:"parent,omitempty" binding:"omitnil"` // parent code
} // @name TaxonUpdate

//go:embed queries/update_taxon.edgeql
var updateTaxonCmd string

func (taxon TaxonUpdate) Update(db *edgedb.Client) (uuid edgedb.UUID, err error) {
	args, _ := json.Marshal(taxon)
	err = db.QuerySingle(context.Background(), updateTaxonCmd, &uuid, args)
	return uuid, err
}
