package occurrence

import (
	"context"
	"encoding/json"
	"slices"

	_ "embed"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/occurrence/queries"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/taxonomy"
)

type OccurrenceCategory string

//generate:enum skip-gel-unmarshal
const (
	Internal OccurrenceCategory = "Internal"
	External OccurrenceCategory = "External"
)

type TypeStatus string

//generate:enum skip-gel-unmarshal
const (
	Holotype TypeStatus = "Holotype"
	Neotype  TypeStatus = "Neotype"
	Topotype TypeStatus = "Topotype"
)

func (m *OccurrenceCategory) UnmarshalEdgeDBStr(data []byte) error {
	s := string(data)
	switch s {
	case "occurrence::InternalBioMat", "seq::AssembledSequence":
		*m = Internal
	case "seq::ExternalSequence", "occurrence::ExternalOccurrence":
		*m = External
	default:
		*m = OccurrenceCategory(s)
	}
	return nil
}

type WithCategory struct {
	Category OccurrenceCategory `gel:"category" json:"category"`
}

type BaseOccurrence struct {
	ID geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	// Code is a unique identifier for the occurrence within the system.
	Code         string `gel:"code" json:"code"`
	WithCategory `gel:"$inline" json:",inline"`
}

type GenericOccurrence[SamplingType any] struct {
	ID             geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	CodeIdentifier `gel:"$inline" json:",inline"`
	WithCategory   `gel:"$inline" json:",inline"`
	HasSequences   bool                        `gel:"has_sequences" json:"has_sequences"`
	Sampling       SamplingType                `gel:"sampling" json:"sampling"`
	Identification Identification              `gel:"identification" json:"identification"`
	TypeStatus     models.Optional[TypeStatus] `gel:"type_status" json:"type_status,omitzero" nameHint:"TypeStatus"`
	Comments       geltypes.OptionalStr        `gel:"comments" json:"comments,omitempty"`
	Meta           people.Meta                 `gel:"meta" json:"meta"`
}

type Occurrence[SamplingType any] struct {
	GenericOccurrence[SamplingType] `gel:"$inline" json:",inline"`
	Datasets                        []dataset.Dataset                           `gel:"datasets" json:"datasets,omitempty"`
	Internal                        models.Optional[InternalBioMatSpecific]     `gel:"internal" json:"internal,omitempty"`
	External                        models.Optional[ExternalOccurrenceSpecific] `gel:"external" json:"external,omitempty"`
}

// OccurrenceWithCategory represents any occurrence
// with its category (internal, external) and element (biomaterial, sequence).
// Internal sequences are not supposed to be included in this type.
// type OccurrenceWithCategory struct {
// 	Occurrence[SamplingInnerWithSite] `gel:"$inline" json:",inline"`
// 	Category                          OccurrenceCategory `gel:"category" json:"category"`
// }

type SamplingDateWithOccurrences struct {
	ID            geltypes.UUID             `gel:"id" json:"id" format:"uuid"`
	Date          OptionalDateWithPrecision `gel:"date" json:"date,omitzero"`
	Occurrences   []OccurrenceAtSite        `gel:"occurrences" json:"occurrences"`
	OccurringTaxa []taxonomy.Taxon          `gel:"occurring_taxa" json:"occurring_taxa,omitempty"`
}

type SamplingDetailsWithOccurrences struct {
	SamplingInner `gel:"$inline" json:",inline"`
	Occurrences   []OccurrenceAtSite `gel:"occurrences" json:"occurrences"`
	Meta          people.Meta        `gel:"meta" json:"meta,omitempty"`
}

type OccurrenceAtSite struct {
	ID             geltypes.UUID      `gel:"id" json:"id" format:"uuid"`
	Code           string             `gel:"code" json:"code"`
	Identification BaseIdentification `gel:"identification" json:"identification"`
	// SamplingDate      DateWithPrecision   `gel:"sampling_date" json:"sampling_date"`
	Category OccurrenceCategory `gel:"category" json:"category"`
}

type SiteWithOccurrences struct {
	SiteItem  `gel:"$inline" json:",inline"`
	Samplings []SamplingDateWithOccurrences `gel:"samplings" json:"samplings"`
}

func ListSamplingsAtSite(db geltypes.Executor, siteCode string) ([]SamplingDetailsWithOccurrences, error) {
	var samplings = []SamplingDetailsWithOccurrences{}
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
				site := (select location::Site filter .code = <str>$0),
			select site.samplings {
				id,
				number,
				performed_on,
				target_taxa: { * },
				habitats: { * },
				fixatives: { * },
				methods: { * },
				meta: { * },
				occurrences := (
					select .occurrences {
						id,
						code,
						category,
						identification: { identified_on, confer, addendum, taxon: { * } },
					}
				)
			}
		`,
		&samplings, siteCode)
	return samplings, err
}

type SiteSamplingStatus string

//generate:enum
const (
	IncludeAllSites        SiteSamplingStatus = "All"
	IncludeSampled         SiteSamplingStatus = "Sampled"
	IncludeWithOccurrences SiteSamplingStatus = "Occurrences"
)

type OccurrencesBySiteOptions struct {
	ListSitesOptions
	taxonomy.TaxaFilters     `json:",inline"`
	Habitats                 []string           `json:"habitats,omitempty" query:"habitats"`
	SamplingTargetTaxa       []string           `json:"sampling_target_taxa,omitempty" query:"sampling_target_taxa"`
	SamplingTargetWholeClade bool               `json:"sampling_target_whole_clade" query:"sampling_target_whole_clade"`
	IncludeSites             SiteSamplingStatus `json:"include_sites,omitempty" query:"include_sites" default:"All" doc:"Include sites with occurrences, sampled sites or all sites. Defaults to sites with at least one occurrence."`
}

func (o OccurrencesBySiteOptions) Options() OccurrencesBySiteOptions {
	return o
}

func OccurrencesBySite(db geltypes.Executor, opts OccurrencesBySiteOptions) ([]SiteWithOccurrences, error) {
	if err := opts.FetchTaxa(db); err != nil {
		return nil, err
	}
	var sites []SiteWithOccurrences
	filters, _ := json.Marshal(opts)
	err := db.Query(context.Background(),
		queries.RenderTemplate("occurrences_by_site.tmpl.edgeql", opts),
		&sites, filters)
	return sites, err
}

// OccurrenceInput is meant to be embedded in other occurrence input type
type OccurrenceInput struct {
	Identification IdentificationInput          `json:"identification" doc:"Occurrence identification"`
	Comments       models.OptionalInput[string] `json:"comments,omitzero"`
	PublishedIn    []string                     `gel:"published_in" json:"published_in,omitempty"`
	Code           models.OptionalInput[string] `gel:"code" json:"code,omitzero" doc:"Unique code identifier for the bio material. Generated from taxon and sampling if not provided." example:"Genus_sp[SITE|2001-01]"`
	TypeStatus     models.OptionalInput[bool]   `gel:"type_status" json:"type_status,omitzero" doc:"Flag indicating if the bio material is a type specimen, i.e. the reference specimen used to describe a new species."`
}

func (i *OccurrenceInput) SetCode(code string) {
	i.Code.SetValue(code)
}

func (occ *OccurrenceInput) WithCreatedMetadata(c *CreatedMetadata) *OccurrenceInput {
	occ.Identification.WithPersonAliases(c.People)
	for i, code := range occ.PublishedIn {
		if c, ok := c.Bibliography[code]; ok {
			occ.PublishedIn[i] = c
		}
	}
	return occ
}

type OccurrenceUpdate struct {
	SamplingID     models.OptionalInput[geltypes.UUID]        `json:"sampling_id" format:"uuid"`
	Identification models.OptionalInput[IdentificationUpdate] `gel:"identification" json:"identification,omitempty"`
	Code           models.OptionalInput[string]               `gel:"code" json:"code,omitempty"`
	TypeStatus     models.OptionalNull[TypeStatus]            `gel:"type_status" json:"type_status,omitzero"`
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
						taxon := (
								[is InternalBioMat].seq_consensus ??
								.identification.taxon
						)
				} filter (
						# only account for well identified bio-material
						[is InternalBioMat].is_homogenous ?? true
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
