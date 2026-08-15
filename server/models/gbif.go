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

// GetParentRank returns the rank of the parent taxon based on the available parent keys.
func (taxon TaxonGBIF) GetParentRank() *TaxonRank {
	rank := new(TaxonRank)
	parentKey, ok := taxon.ParentKey.Get()
	if !ok {
		return rank
	}
	switch parentKey {
	case taxon.KingdomKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankKingdom
	case taxon.PhylumKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankPhylum
	case taxon.ClassKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankClass
	case taxon.OrderKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankOrder
	case taxon.FamilyKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankFamily
	case taxon.GenusKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankGenus
	case taxon.SpeciesKey.GetWithDefault(-1):
		*rank = biomedb.TaxonRankSpecies
	default:
		return nil
	}
	return rank
}

// HasConsistentParent checks if the taxon's parent rank is one level above its own rank.
func (taxon TaxonGBIF) HasConsistentParent() bool {
	if taxon.GetRank() == biomedb.TaxonRankKingdom {
		// Kingdoms have no parent, so they are considered consistent.
		return true
	}
	gbifParentRank := taxon.GetParentRank()
	if gbifParentRank == nil {
		return false
	}
	expectedParentRank := ParentRank(taxon.GetRank())
	return *gbifParentRank == expectedParentRank
}

// IsAcceptable checks if the taxon has a valid rank, status, and consistent parent
// (parent rank must be one level above the current taxon's rank).
func (taxon TaxonGBIF) IsAcceptable() bool {
	return taxon.GetRank().Valid() && taxon.GetStatus().Valid() && taxon.HasConsistentParent()
}

func (taxon TaxonGBIF) GetParentKeysList() []int32 {
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

// GetRank returns the taxon's GBIF rank as a TaxonRank type.
func (taxon TaxonGBIF) GetRank() TaxonRank {
	return TaxonRank(strings.ToLower(taxon.Rank))
}

// GetStatus returns the taxon's GBIF status as a TaxonStatus type.
func (taxon TaxonGBIF) GetStatus() TaxonStatus {
	upperStr := strings.ToUpper(taxon.Status)
	if strings.Contains(upperStr, "SYNONYM") {
		return biomedb.TaxonStatusSynonym
	}
	switch upperStr {
	case "ACCEPTED":
		return biomedb.TaxonStatusAccepted
	case "SYNONYM":
		return biomedb.TaxonStatusSynonym
	case "DOUBTFUL":
		return biomedb.TaxonStatusDoubtful
	default:
		return biomedb.TaxonStatusUnclassified
	}
}

func (taxon TaxonGBIF) ToStaging() biomedb.InsertGBIFBatchParams {
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

type TaxonGBIFPriority int32

//generate:enum
const (
	TaxonGBIFPriorityExactAccepted    TaxonGBIFPriority = 100
	TaxonGBIFPriorityExactNonAccepted TaxonGBIFPriority = 80
	TaxonGBIFPriorityNonExact         TaxonGBIFPriority = 50
)

type TaxonGBIFWithPriority struct {
	TaxonGBIF
	Priority TaxonGBIFPriority `json:"priority"`
}

func (t TaxonGBIF) ComputePriority(res TaxonResolution) TaxonGBIFPriority {
	var priority TaxonGBIFPriority
	if strings.EqualFold(t.Name, res.InputName) || strings.EqualFold(t.ScientificName, res.ScientificName) {
		if t.GetStatus() == biomedb.TaxonStatusAccepted {
			// Accepted exact matches
			priority = TaxonGBIFPriorityExactAccepted
		} else {
			// Synonyms and other non-accepted exact matches
			priority = TaxonGBIFPriorityExactNonAccepted
		}
	} else {
		// Non-exact matches
		priority = TaxonGBIFPriorityNonExact
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
		Priority:     int32(taxon.Priority),
		GBIFID:       &taxon.Key,
	}
}
