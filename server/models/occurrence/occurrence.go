package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/taxonomy"

	"github.com/edgedb/edgedb-go"
)

type OccurrenceCategory string

//generate:enum
const (
	Internal OccurrenceCategory = "Internal"
	External OccurrenceCategory = "External"
)

type GenericOccurrence[SamplingType any] struct {
	ID             edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Sampling       SamplingType       `edgedb:"sampling" json:"sampling"`
	Identification Identification     `edgedb:"identification" json:"identification"`
	Comments       edgedb.OptionalStr `edgedb:"comments" json:"comments"`
}

type Occurrence GenericOccurrence[SamplingInner]

type OccurrenceInput struct {
	SamplingID     edgedb.UUID                  `json:"sampling_id" format:"uuid"`
	Identification IdentificationInput          `edgedb:"identification" json:"identification"`
	Comments       models.OptionalInput[string] `edgedb:"comments" json:"comments"`
}

type OccurrenceOverviewItem struct {
	Name        string             `edgedb:"name" json:"name"`
	ParentName  string             `edgedb:"parent_name" json:"parent_name"`
	Occurrences int32              `edgedb:"occurrences" json:"occurrences"`
	Rank        taxonomy.TaxonRank `edgedb:"rank" json:"rank"`
}

func OccurrenceOverview(db edgedb.Executor) ([]OccurrenceOverviewItem, error) {
	var items = []OccurrenceOverviewItem{}
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
			occ := (
				select occurrence::Occurrence {
					taxon:= (
						[is ExternalBioMat].sequence_consensus ??
						[is InternalBioMat].sequence_consensus ??
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
