package models

import (
	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type TaxonRank = biomedb.TaxonRank
type TaxonStatus = biomedb.TaxonStatus

type Taxon struct {
	ID         uuid.UUID           `json:"id"`
	GBIF_ID    Optional[int32]     `json:"gbif_id"`
	Name       string              `json:"name"`
	Rank       TaxonRank           `json:"rank"`
	Status     TaxonStatus         `json:"status"`
	Authorship Optional[string]    `json:"authorship,omitempty"`
	Comments   Optional[string]    `json:"comments,omitempty"`
	ParentID   Optional[uuid.UUID] `json:"parent_id,omitempty"`
	AcceptedID Optional[uuid.UUID] `json:"accepted_id,omitempty"`
}

func TaxonFromDB(t *biomedb.Taxon) *Taxon {
	if t == nil {
		return nil
	}
	return &Taxon{
		ID:         t.ID,
		GBIF_ID:    NewOptionalFromPtr(t.GBIFID),
		Name:       t.Name,
		Rank:       t.Rank,
		Status:     t.Status,
		Authorship: NewOptionalFromPtr(t.Authorship),
		Comments:   NewOptionalFromPtr(t.Comments),
		ParentID:   NewOptionalFromUUID(t.ParentID),
		AcceptedID: NewOptionalFromUUID(t.AcceptedTaxonID),
	}
}

type TaxonRelations struct {
	ParentTaxon   Optional[Taxon] `json:"parent_taxon,omitempty"`
	AcceptedTaxon Optional[Taxon] `json:"accepted_taxon,omitempty"`
}

type TaxonWithRelations struct {
	Taxon
	TaxonRelations
}

func (t *Taxon) WithRelations(parent *Taxon, accepted *Taxon) TaxonWithRelations {
	return TaxonWithRelations{
		Taxon: *t,
		TaxonRelations: TaxonRelations{
			ParentTaxon:   NewOptionalFromPtr(parent),
			AcceptedTaxon: NewOptionalFromPtr(accepted),
		},
	}
}

type TaxonWithLineage struct {
	Taxon
	TaxonRelations
	Lineage []Taxon `json:"lineage"`
}

func (t *Taxon) WithLineage(parent *Taxon, accepted *Taxon, lineage []Taxon) TaxonWithLineage {
	return TaxonWithLineage{
		Taxon: *t,
		TaxonRelations: TaxonRelations{
			ParentTaxon:   NewOptionalFromPtr(parent),
			AcceptedTaxon: NewOptionalFromPtr(accepted),
		},
		Lineage: lineage,
	}
}
