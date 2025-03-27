package occurrence

import (
	"context"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/taxonomy"
)

type OccurrenceCategory string

//generate:enum
const (
	Internal OccurrenceCategory = "Internal"
	External OccurrenceCategory = "External"
)

type GenericOccurrence[SamplingType any] struct {
	Sampling       SamplingType                     `gel:"sampling" json:"sampling"`
	Identification Identification                   `gel:"identification" json:"identification"`
	PublishedIn    []references.OccurrenceReference `gel:"published_in" json:"published_in,omitempty"`
}

type Occurrence struct {
	ID                               geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	GenericOccurrence[SamplingInner] `gel:"$inline" json:",inline"`
	Comments                         geltypes.OptionalStr `gel:"comments" json:"comments"`
}

type OccurrenceElement string

//generate:enum
const (
	BioMaterialElement OccurrenceElement = "BioMaterial"
	SequenceElement    OccurrenceElement = "Sequence"
)

// OccurrenceWithCategory represents any occurrence with its category (internal, external) and element (biomaterial, sequence).
// Internal sequences are not supposed to be included in this type.
type OccurrenceWithCategory struct {
	Occurrence        `gel:"$inline" json:",inline"`
	Category          OccurrenceCategory `gel:"category" json:"category"`
	OccurrenceElement OccurrenceElement  `gel:"element" json:"element"`
}

// OccurrenceInnerInput is meant to be embedded in other occurrence input type
type OccurrenceInnerInput struct {
	Identification IdentificationInput                   `json:"identification" doc:"Occurrence identification"`
	Comments       models.OptionalInput[string]          `json:"comments"`
	PublishedIn    []references.OccurrenceReferenceInput `gel:"published_in" json:"published_in,omitempty"`
}

type OccurrenceUpdate struct {
	SamplingID     models.OptionalInput[geltypes.UUID]        `json:"sampling_id" format:"uuid"`
	Identification models.OptionalInput[IdentificationUpdate] `gel:"identification" json:"identification,omitempty"`
	Comments       models.OptionalNull[string]                `gel:"comments" json:"comments,omitempty"`
}

// OccurrenceOverviewItem is a representation of the occurrences count for one taxon
type OccurrenceOverviewItem struct {
	Name        string             `gel:"name" json:"name"`
	ParentName  string             `gel:"parent_name" json:"parent_name"`
	Occurrences int32              `gel:"occurrences" json:"occurrences"`
	Rank        taxonomy.TaxonRank `gel:"rank" json:"rank"`
}

// OccurrenceOverview returns the count of occurrences for each taxon
func OccurrenceOverview(db geltypes.Executor) ([]OccurrenceOverviewItem, error) {
	var items = []OccurrenceOverviewItem{}
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
			occ := (
				select occurrence::Occurrence {
					taxon:= (
						[is ExternalBioMat].seq_consensus ??
						[is InternalBioMat].seq_consensus ??
						.identification.taxon
					) { *, parent: { * } }
				} filter (
					#  ignore external bio material that has sequences
					not (Occurrence is ExternalBioMat and not exists [is ExternalBioMat].sequences)
					and [is ExternalBioMat].is_homogenous ?? [is InternalBioMat].is_homogenous ?? true
				)
			),
			groups := (select (
				group occ
				using name := .taxon.name, parent := .taxon.parent.name
				by (parent, name)
			) { arity := count(.elements), key : { * } })
			select taxonomy::Taxon {
				name,
				rank,
				required parent_name := assert_exists(.parent.name),
				required occurrences := <int32>assert_single((
					select groups filter .key.name = taxonomy::Taxon.name
				)).arity ?? <int32>0
			} filter .rank != taxonomy::Rank.Kingdom
		`,
		&items)
	return items, err
}
