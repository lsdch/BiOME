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

type SamplingEventWithOccurrences struct {
	ID            geltypes.UUID      `gel:"id" json:"id" format:"uuid"`
	Date          DateWithPrecision  `gel:"date" json:"date"`
	Target        SamplingTarget     `gel:"$inline" json:"target"`
	Occurrences   []OccurrenceAtSite `gel:"occurrences" json:"occurrences"`
	OccurringTaxa []taxonomy.Taxon   `gel:"occurring_taxa" json:"occurring_taxa,omitempty"`
}

type OccurrenceAtSite struct {
	ID    geltypes.UUID       `gel:"id" json:"id" format:"uuid"`
	Code  string              `gel:"code" json:"code"`
	Taxon taxonomy.TaxonInner `gel:"taxon" json:"taxon"`
	// SamplingDate      DateWithPrecision   `gel:"sampling_date" json:"sampling_date"`
	Category          OccurrenceCategory `gel:"category" json:"category"`
	OccurrenceElement OccurrenceElement  `gel:"element" json:"element"`
}

type SiteWithOccurrences struct {
	SiteItem    `gel:"$inline" json:",inline"`
	Samplings   []SamplingEventWithOccurrences     `gel:"samplings" json:"samplings"`
	LastVisited models.Optional[DateWithPrecision] `gel:"last_visited" json:"last_visited,omitempty"`
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
	Taxa                     []string             `json:"taxa,omitempty" query:"taxa"`
	WholeClade               bool                 `json:"whole_clade" query:"whole_clade"`
	Habitats                 []string             `json:"habitats,omitempty" query:"habitats"`
	SamplingTargetKinds      []SamplingTargetKind `json:"sampling_target_kinds,omitempty" query:"sampling_target_kinds" doc:"List of sampling target names. \"Community\" "`
	SamplingTargetTaxa       []string             `json:"sampling_target_taxa,omitempty" query:"sampling_target_taxa"`
	SamplingTargetWholeClade bool                 `json:"sampling_target_whole_clade" query:"sampling_target_whole_clade"`
	IncludeSites             SiteSamplingStatus   `json:"include_sites,omitempty" query:"include_sites" default:"All" doc:"Include sites with occurrences, sampled sites or all sites. Defaults to sites with at least one occurrence."`
}

func (o OccurrencesBySiteOptions) Options() OccurrencesBySiteOptions {
	return o
}

func OccurrencesBySite(db geltypes.Executor, opts OccurrencesBySiteOptions) ([]SiteWithOccurrences, error) {
	var sites []SiteWithOccurrences
	filters, _ := json.Marshal(opts)
	err := db.Query(context.Background(),
		`#edgeql
			with module occurrence,
				filters := <json>$0,
				country_codes := <str>json_array_unpack(json_get(filters, 'countries')),
				taxa := (
					select taxonomy::Taxon
					filter .code in <str>json_array_unpack(json_get(filters, 'taxa'))
				),
				datasets := (
					select datasets::Dataset
					filter .slug in <str>json_array_unpack(json_get(filters, 'datasets'))
				),
				whole_clade := <bool>json_get(filters, 'whole_clade'),
				habitats := <str>json_array_unpack(json_get(filters, 'habitats')),
				sampling_target_kinds := <events::SamplingTarget>json_array_unpack(json_get(filters, 'sampling_target_kinds')),
				sampling_target_taxa := (
					select taxonomy::Taxon
					filter .name in <str>json_array_unpack(json_get(filters, 'sampling_target_taxa'))
				),
				sampling_target_whole_clade := <bool>json_get(filters, 'sampling_target_whole_clade'),
				sampling_status := <str>json_get(filters, 'include_sites'),
			select location::Site {
				*,
				country: { * },
				last_visited := assert_single((
					select distinct .events.performed_on filter (.date = max(location::Site.events.performed_on.date)) limit 1
				)),
				samplings := (
					select .events.samplings
					filter (
						not exists habitats or all(habitats in .habitats.label)
					)
					and (
						not exists sampling_target_kinds or (
							.sampling_target in sampling_target_kinds
						)
					)
					and (
						not exists sampling_target_taxa or (
							if sampling_target_whole_clade
							then any(.target_taxa in sampling_target_taxa) or any(taxonomy::is_in_clade(.target_taxa, sampling_target_taxa))
							else any(.target_taxa in sampling_target_taxa)
						)
					)
				) {
					id,
					date := .event.performed_on,
					sampling_target,
					target_taxa: { * },
					occurring_taxa: { * },
					occurrences := (
						with occurrences := (
							select .occurrences filter (
								(exists [is BioMaterial].id) or
								(not exists [is seq::ExternalSequence].source_sample)
							)
						),
						select occurrences {
							id,
							code,
							required taxon := (
									[is ExternalBioMat].seq_consensus ??
									[is InternalBioMat].seq_consensus ??
									.identification.taxon
								) { name, status, rank},
							required category := ([is InternalBioMat].category ?? OccurrenceCategory.External),
							required element := (
								if exists [is seq::Sequence].id then 'Sequence'
								else 'BioMaterial'
							)
						}
						filter (
							not exists taxa or (
								if whole_clade
								then (.taxon in taxa) or any(taxonomy::is_in_clade(.taxon, taxa))
								else .taxon in taxa
							)
						)
						and (
							not exists datasets
							or occurrences in datasets[is datasets::OccurrenceDataset].occurrences
						)
					)
				}
		 } filter (
			(
				not exists sampling_status or sampling_status = "All" or (
					sampling_status = "Sampled" and exists .samplings
				) or (
					sampling_status = "Occurrences" and exists .samplings.occurrences
				)
			) and
			(not exists country_codes or .country.code in country_codes) and
			# (not exists habitats or exists .samplings) and
			# (not exists taxa or exists .samplings.occurrences) and
			# (not exists sampling_target_kinds or exists .samplings) and
			(not exists datasets or (
				location::Site in datasets[is datasets::SiteDataset].sites
				?? datasets[is datasets::OccurrenceDataset].sites
			))
		 )
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
