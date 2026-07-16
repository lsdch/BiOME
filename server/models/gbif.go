package models

import (
	_ "embed"
	"maps"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type TaxonGBIF struct {
	Key                     int32            `json:"key"`
	Parent                  Optional[string] `json:"parent,omitempty"`
	ParentKey               Optional[int32]  `json:"parentKey,omitempty"`
	Name                    string           `json:"canonicalName"`
	ScientificName          string           `json:"scientificName"`
	Status                  string           `json:"taxonomicStatus"`
	Rank                    string           `json:"rank"`
	NameType                string           `json:"nameType"`
	KingdomKey              Optional[int32]  `json:"kingdomKey,omitempty"`
	PhylumKey               Optional[int32]  `json:"phylumKey,omitempty"`
	ClassKey                Optional[int32]  `json:"classKey,omitempty"`
	OrderKey                Optional[int32]  `json:"orderKey,omitempty"`
	FamilyKey               Optional[int32]  `json:"familyKey,omitempty"`
	GenusKey                Optional[int32]  `json:"genusKey,omitempty"`
	SpeciesKey              Optional[int32]  `json:"speciesKey,omitempty"`
	HigherClassificationMap map[int32]string `json:"higherClassificationMap"`
	Authorship              Optional[string] `json:"authorship,omitempty"`
	NumDescendants          Optional[int32]  `json:"numDescendants,omitempty"`
	AcceptedKey             Optional[int32]  `json:"acceptedKey,omitempty"`
	AcceptedName            Optional[string] `json:"accepted,omitempty"`
}

func (taxon *TaxonGBIF) IsAcceptable() bool {
	return taxon.GetRank().Valid() && taxon.GetStatus().Valid()
}

func (taxon *TaxonGBIF) GetParentKeysList() []int32 {
	parents := []int32{}

	addKeys := func(keys ...Optional[int32]) {
		for _, key := range keys {
			if key, ok := key.Get(); !ok || key == taxon.Key {
				return
			} else {
				parents = append(parents, key)
			}
		}
	}

	addKeys(
		taxon.KingdomKey,
		taxon.PhylumKey,
		taxon.ClassKey,
		taxon.OrderKey,
		taxon.FamilyKey,
		taxon.GenusKey,
		taxon.SpeciesKey,
	)

	return parents
}

func (taxon *TaxonGBIF) GetRank() TaxonRank {
	return TaxonRank(taxon.Rank)
}

func (taxon *TaxonGBIF) GetStatus() TaxonStatus {
	switch strings.ToUpper(taxon.Status) {
	case "ACCEPTED":
		return biomedb.TaxonStatusACCEPTED
	case "SYNONYM":
		return biomedb.TaxonStatusSYNONYM
	case "DOUBTFUL":
		return biomedb.TaxonStatusDOUBTFUL
	default:
		return biomedb.TaxonStatusUNCLASSIFIED
	}
}

func (taxon *TaxonGBIF) ToStaging() biomedb.InsertGBIFBatchParams {
	return biomedb.InsertGBIFBatchParams{
		Key:              taxon.Key,
		Parent:           taxon.Parent.ToPtr(),
		ParentKey:        taxon.ParentKey.ToPtr(),
		CanonicalName:    taxon.Name,
		ScientificName:   taxon.ScientificName,
		Status:           string(taxon.GetStatus()),
		Rank:             string(taxon.GetRank()),
		NameType:         taxon.NameType,
		KingdomKey:       taxon.KingdomKey.ToPtr(),
		PhylumKey:        taxon.PhylumKey.ToPtr(),
		ClassKey:         taxon.ClassKey.ToPtr(),
		OrderKey:         taxon.OrderKey.ToPtr(),
		FamilyKey:        taxon.FamilyKey.ToPtr(),
		GenusKey:         taxon.GenusKey.ToPtr(),
		SpeciesKey:       taxon.SpeciesKey.ToPtr(),
		HigherTaxonKeys:  slices.Collect(maps.Keys(taxon.HigherClassificationMap)),
		HigherTaxonNames: slices.Collect(maps.Values(taxon.HigherClassificationMap)),
		Authorship:       taxon.Authorship.ToPtr(),
		NumDescendants:   taxon.NumDescendants.ToPtr(),
		AcceptedKey:      taxon.AcceptedKey.ToPtr(),
		AcceptedName:     taxon.AcceptedName.ToPtr(),
	}
}

type TaxonGBIFWithPriority struct {
	TaxonGBIF
	Priority int32 `json:"priority"`
}

func (t TaxonGBIF) ComputePriority(res TaxonResolution) int32 {
	var priority int32
	if t.Name == res.InputName || t.ScientificName == res.ScientificName {
		if t.GetStatus() == biomedb.TaxonStatusACCEPTED {
			// Accepted exact matches
			priority = 100
		} else {
			// Synonyms and other non-accepted exact matches
			priority = 80
		}
	} else {
		// Non-exact matches
		priority = 50
	}
	return priority
}

func (t TaxonGBIF) WithPriority(res TaxonResolution) TaxonGBIFWithPriority {
	return TaxonGBIFWithPriority{
		TaxonGBIF: t,
		Priority:  t.ComputePriority(res),
	}
}

func (taxon TaxonGBIFWithPriority) ToCandidate(importID uuid.UUID, resolutionID uuid.UUID) biomedb.InsertTaxonCandidatesBatchParams {
	return biomedb.InsertTaxonCandidatesBatchParams{
		ImportID:     importID,
		ResolutionID: resolutionID,
		Name:         taxon.Name,
		Authorship:   taxon.Authorship.ToPtr(),
		Rank:         taxon.GetRank(),
		Status:       taxon.GetStatus(),
		Source:       biomedb.TaxonMatchSourceGBIF,
		MatchType:    biomedb.TaxonMatchTypeExact,
		Priority:     taxon.Priority,
		GBIFID:       &taxon.Key,
	}
}

func (taxon *TaxonGBIF) Normalize() *TaxonGBIF {
	if authorship, isSet := taxon.Authorship.Get(); isSet && authorship == "" {
		taxon.Authorship.Clear()
	}
	taxon.Status = string(taxon.GetStatus())
	taxon.Rank = string(taxon.GetRank())
	return taxon
}
