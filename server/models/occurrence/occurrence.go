package occurrence

import (
	"context"
	"encoding/json"
	"slices"

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

// OccurrenceWithCategory represents any occurrence
// with its category (internal, external) and element (biomaterial, sequence).
// Internal sequences are not supposed to be included in this type.
type OccurrenceWithCategory struct {
	Occurrence        `gel:"$inline" json:",inline"`
	Category          OccurrenceCategory `gel:"category" json:"category"`
	OccurrenceElement OccurrenceElement  `gel:"element" json:"element"`
}

type OccurrenceAtSite struct {
	ID                geltypes.UUID       `gel:"id" json:"id" format:"uuid"`
	Code              string              `gel:"code" json:"code"`
	Taxon             taxonomy.TaxonInner `gel:"taxon" json:"taxon"`
	SamplingDate      DateWithPrecision   `gel:"sampling_date" json:"sampling_date"`
	Category          OccurrenceCategory  `gel:"category" json:"category"`
	OccurrenceElement OccurrenceElement   `gel:"element" json:"element"`
}

type SiteWithOccurrences struct {
	SiteItem    `gel:"$inline" json:",inline"`
	Occurrences []OccurrenceAtSite `gel:"occurrences" json:"occurrences"`
}

type OccurrencesBySiteOptions struct {
	ListSitesOptions
}

func (o OccurrencesBySiteOptions) Options() OccurrencesBySiteOptions {
	return o
}

func OccurrencesBySite(db geltypes.Executor, opts OccurrencesBySiteOptions) ([]SiteWithOccurrences, error) {
	var sites []SiteWithOccurrences
	filters, _ := json.Marshal(opts.ListSitesOptions)
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
			 filters := <json>$0,
			 country_codes := <str>json_get(filters, 'country_codes'),
			select location::Site {
				*,
				country: { * },
				occurrences := (
					(
						select .events.samplings.occurrences
						filter ((exists [is BioMaterial].id) or (not exists [is seq::ExternalSequence].source_sample))
					) {
						id,
						code,
						required sampling_date := .sampling.event.performed_on,
						required taxon := (
										[is ExternalBioMat].seq_consensus ??
										[is InternalBioMat].seq_consensus ??
										.identification.taxon
						) { name, status, rank},
						category := ([is InternalBioMat].category  ?? OccurrenceCategory.External),
						element := (
							if exists [is seq::Sequence].id then 'Sequence'
							else 'BioMaterial'
						)
					}
				)
		 } filter (not exists country_codes or .country.code in country_codes)
		`,
		&sites, filters)
	return sites, err
}

// OccurrenceInnerInput is meant to be embedded in other occurrence input type
type OccurrenceInnerInput struct {
	Identification IdentificationInput                   `json:"identification" doc:"Occurrence identification"`
	Comments       models.OptionalInput[string]          `json:"comments"`
	PublishedIn    []references.OccurrenceReferenceInput `gel:"published_in" json:"published_in,omitempty"`
}

func (occ *OccurrenceInnerInput) WithCreatedMetadata(c CreatedMetadata) OccurrenceInnerInput {
	occ.Identification.WithPersonAliases(c.People)
	for i := range occ.PublishedIn {
		(&occ.PublishedIn[i]).WithArticleCode(c.Bibliography)
	}
	return *occ
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

type occurrenceOverviewQueryResult struct {
	Occurrences   []OccurrenceOverviewItem `gel:"occurrences"`
	NoOccurrences []OccurrenceOverviewItem `gel:"no_occurrences"`
}

func (o occurrenceOverviewQueryResult) toItems() []OccurrenceOverviewItem {
	return slices.Concat(o.Occurrences, o.NoOccurrences)
}

// OccurrenceOverview returns the count of occurrences for each taxon
func OccurrenceOverview(db geltypes.Executor) ([]OccurrenceOverviewItem, error) {
	var items = occurrenceOverviewQueryResult{}
	err := db.QuerySingle(context.Background(),
		`#edgeql
			with module occurrence,
			occ := (
				select Occurrence {
						# use most accurate identification
						taxon:= (
								[is ExternalBioMat].seq_consensus ??
								[is InternalBioMat].seq_consensus ??
								.identification.taxon
						)
				} filter (
						# ignore external bio material that has sequences
						not (Occurrence is ExternalBioMat and exists [is ExternalBioMat].sequences)
						# only account for well identified bio-material
						and [is ExternalBioMat].is_homogenous ?? [is InternalBioMat].is_homogenous ?? true
				)
			),
			groups := (select (group occ by .taxon) { arity := count(.elements)}),
			noOccTaxa := (
				select (taxonomy::Taxon except occ.taxon) filter .rank != taxonomy::Rank.Kingdom
			),

			select {
				occurrences := groups {
					required name:= .key.taxon.name,
					required rank:= .key.taxon.rank,
					required parent_name:= assert_exists(.key.taxon.parent.name),
					required occurrences := <int32>.arity
				},
				no_occurrences := (
					# Rest of taxonomy i.e. taxa having no occurrences
					select noOccTaxa {
					required name := .name,
					required rank := .rank,
					required parent_name:= assert_exists(noOccTaxa.parent.name),
					required occurrences := <int32>0
				})
			}
		`,
		&items)
	return items.toItems(), err
}
