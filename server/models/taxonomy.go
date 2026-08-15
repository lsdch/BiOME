package models

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/lsdch/biome/db/biomedb"
)

type TaxonRank = biomedb.TaxonRank

func TaxonRankFromString(rank string) (TaxonRank, bool) {
	r := TaxonRank(strings.ToLower(rank))
	return r, r.Valid()
}

func ParentRank(rank TaxonRank) TaxonRank {
	switch rank {
	case biomedb.TaxonRankSpecies:
		return biomedb.TaxonRankGenus
	case biomedb.TaxonRankKingdom:
		return biomedb.TaxonRankKingdom
	default:
		return biomedb.AllTaxonRankValues()[biomedb.TaxonRankHierarchy[rank]+1]
	}
}

func ChildRank(rank TaxonRank) TaxonRank {
	switch rank {
	case biomedb.TaxonRankGenus:
		return biomedb.TaxonRankSpecies
	case biomedb.TaxonRankSubspecies:
		return biomedb.TaxonRankSubspecies
	case biomedb.TaxonRankSubgenus:
		return biomedb.TaxonRankSubgenus
	default:
		return biomedb.AllTaxonRankValues()[biomedb.TaxonRankHierarchy[rank]-1]
	}
}

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

func (t Taxon) Code() string {
	return strings.ReplaceAll(t.Name, " ", "_")
}

func TaxonFromDB(t *biomedb.Taxon) *Taxon {
	if t == nil {
		return nil
	}
	return &Taxon{
		ID:         t.ID,
		GBIF_ID:    NewOptionalFromPtr(t.GBIFID),
		Name:       t.Name,
		Rank:       TaxonRank(t.Rank),
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

// TaxonWithLineage represents a taxon along with its lineage (ancestors), and optionally its parent and accepted taxa.
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

// TaxonWithFullLineage represents a taxon along with its lineage and descendants.
type TaxonWithFullLineage struct {
	TaxonWithLineage
	Descendants []Taxon `json:"descendants"`
}

type ListTaxaParams struct {
	SearchTerm  Optional[string] `query:"search_term"`
	Ranks       []TaxonRank      `query:"ranks"`
	SampledOnly bool             `query:"sampled_only"`
	Pagination
}

type CreateTaxonInput struct {
	Name       string              `json:"name" validate:"required"`
	Rank       TaxonRank           `json:"rank" validate:"required"`
	Status     TaxonStatus         `json:"status" validate:"required"`
	Authorship Optional[string]    `json:"authorship,omitempty"`
	Comments   Optional[string]    `json:"comments,omitempty"`
	ParentID   uuid.UUID           `json:"parent_id"`
	AcceptedID Optional[uuid.UUID] `json:"accepted_id,omitempty"`
}

func (i *CreateTaxonInput) ToParams() *biomedb.InsertTaxonParams {
	return &biomedb.InsertTaxonParams{
		Name:       i.Name,
		Rank:       (biomedb.TaxonRank)(i.Rank),
		Status:     i.Status,
		Authorship: i.Authorship.ToPtr(),
		Comments:   i.Comments.ToPtr(),
		ParentID:   UUIDToPg(i.ParentID),
		AcceptedID: UUIDOpt(i.AcceptedID),
	}
}

// var speciesQualifierRegexp = regexp.MustCompile(
// 	`(?i)\b(?:sp\d?\.(?:\s*)|spp\.(?:\s*)|aff\.(?:\s*)|cf\.(?:\s*))`,
// )

// var genusQualifierRegexp = regexp.MustCompile(
// 	`(?i)\bgen\.\s+(?:n\.|\S+)`,
// )

func InferParentName(canonicalName string) (string, bool) {
	name := strings.TrimSpace(canonicalName)
	if name == "" {
		return "", false
	}

	// Remove parenthetical annotations.
	noParenName := regexp.MustCompile(`\s*\([^)]*\)`).ReplaceAllString(name, " ")
	words := strings.Fields(noParenName)

	// Remove punctuation commonly attached to the final token.
	for i := range words {
		words[i] = strings.TrimRight(words[i], ",;")
	}

	// First look for a species-level qualifier anywhere in the name.
	//
	// Examples:
	//   Candonidae gen. n. sp. I1
	//   Candonidae gen. I1 sp. I1
	//   Schellencandona n. sp. J4
	for i, word := range words {
		lower := strings.ToLower(word)

		if lower != "sp." &&
			lower != "spp." &&
			lower != "aff." &&
			lower != "cf." &&
			!isNumberedSpecies(lower) {
			continue
		}

		parentEnd := i

		// "Schellencandona n. sp. J4"
		//                  ^^
		// The n. belongs to sp., so remove it.
		//
		// "Candonidae gen. n. sp. I1"
		//                ^^^^^^^
		// Here n. belongs to gen. n., so keep it.
		if i >= 1 &&
			strings.EqualFold(words[i-1], "n.") &&
			(i < 2 || !strings.EqualFold(words[i-2], "gen.")) {
			parentEnd = i - 1
		}

		if parentEnd > 0 {
			return strings.Join(words[:parentEnd], " "), true
		}

		return "", false
	}

	// No species-level qualifier: now look for a genus qualifier.
	//
	// Candonidae gen. I1 -> Candonidae
	// Candonidae gen. n.  -> Candonidae
	for i, word := range words {
		if !strings.EqualFold(word, "gen.") {
			continue
		}

		if i > 0 {
			return strings.Join(words[:i], " "), true
		}

		return "", false
	}

	// Generic fallback.
	switch len(words) {
	case 0:
		return "", false
	case 1:
		if strings.Contains(name, "(") {
			return words[0], true
		}
		return "", false

	case 2:
		return words[0], true

	default:
		return strings.Join(words[:len(words)-1], " "), true
	}
}

func isNumberedSpecies(s string) bool {
	if len(s) < 4 || !strings.HasPrefix(s, "sp") {
		return false
	}

	s = strings.TrimPrefix(s, "sp")

	if !strings.HasSuffix(s, ".") {
		return false
	}

	s = strings.TrimSuffix(s, ".")

	if s == "" {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

var unclassifiedQualifiersRegex = regexp.MustCompile(
	`(?i)(?:\baff\.|\(aff\.\)|\bgen\.|\bs?sp\.|\bn\.)`,
)

func InferStatusFromStagingName(name string) TaxonStatus {
	lowerName := strings.ToLower(name)
	if unclassifiedQualifiersRegex.MatchString(lowerName) {
		return biomedb.TaxonStatusUnclassified
	}
	return biomedb.TaxonStatusUnreferenced
}
