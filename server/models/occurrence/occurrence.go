package occurrence

import (
	"context"
	"darco/proto/models"
	"darco/proto/models/taxonomy"
	"fmt"

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
	Identification IdentificationInput          `json:"identification"`
	Comments       models.OptionalInput[string] `json:"comments"`
}

func (i OccurrenceInput) GenerateCode(db edgedb.Executor) (string, error) {
	sampling, err := i.GetSampling(db)
	if err != nil {
		return "", fmt.Errorf("Sampling not found")
	}
	return fmt.Sprintf("%s[%s]",
		taxonomy.TaxonCode(i.Identification.Taxon),
		sampling.Code,
	), nil
}

func (i OccurrenceInput) GetSampling(db edgedb.Executor) (sampling SamplingInner, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
			select <events::Sampling><uuid>$0 { * }
		`, i.SamplingID, &sampling)
	return
}

type OccurrenceUpdate struct {
	SamplingID     models.OptionalInput[edgedb.UUID]          `json:"sampling_id" format:"uuid"`
	Identification models.OptionalInput[IdentificationUpdate] `edgedb:"identification" json:"identification,omitempty"`
	Comments       models.OptionalNull[string]                `edgedb:"comments" json:"comments,omitempty"`
}

// OccurrenceOverviewItem is a representation of the occurrences count for one taxon
type OccurrenceOverviewItem struct {
	Name        string             `edgedb:"name" json:"name"`
	ParentName  string             `edgedb:"parent_name" json:"parent_name"`
	Occurrences int32              `edgedb:"occurrences" json:"occurrences"`
	Rank        taxonomy.TaxonRank `edgedb:"rank" json:"rank"`
}

// OccurrenceOverview returns the count of occurrences for each taxon
func OccurrenceOverview(db edgedb.Executor) ([]OccurrenceOverviewItem, error) {
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
					and [is ExternalBioMat].homogenous ?? [is InternalBioMat].homogenous ?? true
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
