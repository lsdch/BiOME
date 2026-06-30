package models

import (
	_ "embed"
	"maps"
	"slices"

	"github.com/lsdch/biome/db/biomedb"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TaxonGBIF struct {
	Key                     int32            `json:"key" gel:"GBIF_ID"`
	Parent                  string           `json:"parent"`
	ParentKey               int32            `json:"parentKey" gel:"parentID"`
	Name                    string           `json:"canonicalName" gel:"name"`
	ScientificName          string           `json:"scientificName" gel:"scientific_name"`
	Status                  string           `json:"taxonomicStatus" gel:"status"`
	Rank                    string           `json:"rank" gel:"rank"`
	NameType                string           `json:"nameType" gel:"name_type"`
	KingdomKey              int32            `json:"kingdomKey"`
	PhylumKey               int32            `json:"phylumKey"`
	ClassKey                int32            `json:"classKey"`
	OrderKey                int32            `json:"orderKey"`
	FamilyKey               int32            `json:"familyKey"`
	GenusKey                int32            `json:"genusKey"`
	SpeciesKey              int32            `json:"speciesKey"`
	HigherClassificationMap map[int32]string `json:"higherClassificationMap"`
	Authorship              Optional[string] `json:"authorship,omitempty" gel:"authorship,omitempty"`
	NumDescendants          int32            `json:"numDescendants" gel:"-"`
	AcceptedKey             int32            `json:"acceptedKey"`
	AcceptedName            string           `json:"accepted"`
}

func (taxon *TaxonGBIF) GetParentKeysList() []int32 {
	parents := []int32{}

	addKeys := func(keys ...int32) {
		for _, key := range keys {
			if key == 0 || key == taxon.Key {
				return
			}
			parents = append(parents, key)
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
	return TaxonRank(cases.Title(language.English).String(taxon.Rank))
}

func (taxon *TaxonGBIF) GetStatus() TaxonStatus {
	switch taxon.Status {
	case "Accepted":
		return biomedb.TaxonStatusACCEPTED
	case "Synonym":
		return biomedb.TaxonStatusSYNONYM
	case "Doubtful":
		return biomedb.TaxonStatusDOUBTFUL
	default:
		return biomedb.TaxonStatusUNCLASSIFIED
	}
}

func (taxon *TaxonGBIF) ToStaging() biomedb.InsertGBIFBatchParams {
	return biomedb.InsertGBIFBatchParams{
		Key:              taxon.Key,
		Parent:           taxon.Parent,
		ParentKey:        taxon.ParentKey,
		CanonicalName:    taxon.Name,
		ScientificName:   taxon.ScientificName,
		Status:           string(taxon.GetStatus()),
		Rank:             string(taxon.GetRank()),
		NameType:         taxon.NameType,
		KingdomKey:       taxon.KingdomKey,
		PhylumKey:        taxon.PhylumKey,
		ClassKey:         taxon.ClassKey,
		OrderKey:         taxon.OrderKey,
		FamilyKey:        taxon.FamilyKey,
		GenusKey:         taxon.GenusKey,
		SpeciesKey:       taxon.SpeciesKey,
		HigherTaxonKeys:  slices.Collect(maps.Keys(taxon.HigherClassificationMap)),
		HigherTaxonNames: slices.Collect(maps.Values(taxon.HigherClassificationMap)),
		Authorship:       taxon.Authorship.GetWithDefault(""),
		NumDescendants:   taxon.NumDescendants,
		AcceptedKey:      taxon.AcceptedKey,
		AcceptedName:     taxon.AcceptedName,
	}
}

func (taxon *TaxonGBIF) ToCandidate(inputName string) biomedb.InsertTaxonCandidatesBatchParams {
	var priority int32
	if taxon.Name == inputName || taxon.ScientificName == inputName {
		if taxon.GetStatus() == biomedb.TaxonStatusACCEPTED {
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

	return biomedb.InsertTaxonCandidatesBatchParams{
		InputName:  inputName,
		Name:       taxon.Name,
		Authorship: taxon.Authorship.ToPtr(),
		Rank:       taxon.GetRank(),
		Status:     taxon.GetStatus(),
		Source:     biomedb.TaxonMatchSourceGBIF,
		MatchType:  biomedb.TaxonMatchTypeExact,
		Priority:   priority,
		GBIFID:     &taxon.Key,
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
