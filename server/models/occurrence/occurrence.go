package occurrence

import (
	"context"
	"encoding/json"
	"slices"

	_ "embed"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/dataset"
	"github.com/lsdch/biome/models/occurrence/queries"
	"github.com/lsdch/biome/models/people"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

type QuantityRange struct {
	Lower int32 `gel:"lower" json:"lower"`
	Upper int32 `gel:"upper" json:"upper"`
}

type TypeStatus string

//generate:enum skip-gel-unmarshal
const (
	Holotype TypeStatus = "Holotype"
	Neotype  TypeStatus = "Neotype"
	Topotype TypeStatus = "Topotype"
)

type BaseOccurrence[SamplingType any] struct {
	ID                     geltypes.UUID `gel:"id" json:"id" format:"uuid"`
	CodeIdentifier         `gel:"$inline" json:",inline"`
	HasSequences           bool                                `gel:"has_sequences" json:"has_sequences"`
	Sampling               SamplingType                        `gel:"sampling" json:"sampling"`
	Identification         Identification                      `gel:"identification" json:"identification"`
	TypeStatus             models.Optional[TypeStatus]         `gel:"type_status" json:"type_status,omitzero" nameHint:"TypeStatus"`
	Comments               geltypes.OptionalStr                `gel:"comments" json:"comments,omitempty"`
	Quantity               models.Optional[QuantityRange]      `gel:"quantity" json:"quantity,omitempty"`
	ContentDescription     geltypes.OptionalStr                `gel:"content_description" json:"content_description,omitempty"`
	VerbatimIdentification geltypes.OptionalStr                `gel:"verbatim_identification" json:"verbatim_identification,omitempty"`
	Collections            []references.CollectionWithVouchers `gel:"collections" json:"collections,omitempty"`
	Meta                   people.Meta                         `gel:"meta" json:"meta"`
}

type OccurrenceListItem BaseOccurrence[SamplingInnerWithSite]

type Occurrence[SamplingType any] struct {
	BaseOccurrence[SamplingType] `gel:"$inline" json:",inline"`
	// ExternalLink       geltypes.OptionalStr           `gel:"external_link" json:"external_link,omitempty"`
	PublishedIn []references.Article    `gel:"published_in" json:"published_in,omitempty"`
	Sources     []references.DataSource `gel:"sources" json:"sources,omitempty"`
	Datasets    []dataset.Dataset       `gel:"datasets" json:"datasets,omitempty"`
	Sequences   []ExternalSequence      `gel:"sequences" json:"sequences,omitempty"`
}

func GetOccurrence(db geltypes.Executor, code string) (occurrence Occurrence[SamplingWithSite], err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence
		select Occurrence {
			*,
			meta: { * },
			datasets: { *, maintainers: { * }, meta: { * } },
			sampling: {
				*,
				performed_by: { * },
				target_taxa: { * },
				fixatives: { * },
				methods: { * },
				habitats: { * },
				occurrences: { *, identification: { **, identified_by: { * } } },
				occurring_taxa: { * },
				site: { *, country: { * } }
			},
			identification: { **, identified_by: { ** } },
			sources: { * },
			sequences: { *, gene: { * }, referenced_in: { * }, meta: { * } },
			# external_link,
			collections: { *, @vouchers },
			quantity,
			published_in: { * },
			content_description
		} filter .code = <str>$0
	`,
		&occurrence, code)
	return occurrence, err
}

type BioMatSortKey string

//generate:enum
const (
	BioMatSortCode           BioMatSortKey = "code"
	BioMatSortSite           BioMatSortKey = "site"
	BioMatSortSamplingDate   BioMatSortKey = "sampling_date"
	BioMatSortIdentifiedOn   BioMatSortKey = "identified_on"
	BioMatSortTaxon          BioMatSortKey = "taxon"
	BioMatSortIdentification BioMatSortKey = "identification"
	BioMatSortIdentifiedBy   BioMatSortKey = "identified_by"
	BioMatSortLastUpdated    BioMatSortKey = "last_updated"
)

var BioMatSortMap = map[BioMatSortKey]string{
	BioMatSortCode:           ".code",
	BioMatSortSite:           "(.site.name ?? .site.code)",
	BioMatSortSamplingDate:   ".sampling.performed_on.date",
	BioMatSortIdentification: ".identification.taxon.name",
	BioMatSortIdentifiedOn:   ".identification.identified_on.date",
	BioMatSortTaxon:          ".identification.taxon.name",
	BioMatSortIdentifiedBy:   ".identification.identified_by.last_name",
	BioMatSortLastUpdated:    ".meta.lastUpdated",
}

type ListOccurrencesOptions struct {
	models.Pagination `json:",inline"`
	models.SortBy[BioMatSortKey]
	models.Filter        `json:",inline"`
	taxonomy.TaxaFilters `json:",inline"`
	Datasets             []string                                   `query:"datasets" json:"datasets,omitzero"`
	HasSequences         models.OptionalInput[bool]                 `query:"has_sequences" json:"has_sequences,omitzero"`
	Confer               models.OptionalInput[bool]                 `query:"confer" json:"confer,omitzero"`
	TypeStatus           models.OptionalInput[TypeStatus]           `query:"type_status" json:"type_status,omitzero"`
	Status               models.OptionalInput[taxonomy.TaxonStatus] `query:"status" json:"status,omitzero"`
	Rank                 []taxonomy.TaxonRank                       `query:"rank" json:"rank,omitzero"`
}

func (o ListOccurrencesOptions) Options() ListOccurrencesOptions {
	return o
}

func (i ListOccurrencesOptions) OrderByString() string {
	if i.SortBy.Key == "" {
		return ""
	}
	if term, ok := BioMatSortMap[i.SortBy.Key]; ok {
		return term + " " + string(i.SortBy.Order)
	} else {
		logrus.Warnf("Unknown sort key: %s", i.SortBy.Key)
		return ""
	}
}

func ListOccurrences(db geltypes.Executor, opts ListOccurrencesOptions) (models.PaginatedList[OccurrenceListItem], error) {
	if err := opts.TaxaFilters.FetchTaxa(db); err != nil {
		return models.PaginatedList[OccurrenceListItem]{}, err
	}
	params, _ := json.Marshal(opts)
	logrus.Debugf("Params: %s", string(params))
	var result = models.PaginatedList[OccurrenceListItem]{
		Items: []OccurrenceListItem{},
	}
	err := db.QuerySingle(context.Background(),
		queries.RenderTemplate("list_occurrences.tmpl.edgeql", opts),
		&result, params, opts.OrderByString())
	return result, err
}

func DeleteOccurrence(db geltypes.Executor, code string) (deleted OccurrenceListItem, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence
			select (
				delete Occurrence filter .code = <str>$0
			) {
        *,
				sampling: { *, site: { *, country: { * } } },
				identification: { **, identified_by: { ** } },
      }
		`,
		&deleted, code)
	return
}

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

type OccurrenceInput struct {
	Identification         IdentificationInput              `json:"identification" doc:"Occurrence identification"`
	PublishedIn            []string                         `gel:"published_in" json:"published_in,omitempty"`
	Code                   models.OptionalInput[string]     `gel:"code" json:"code,omitzero" doc:"Unique code identifier for the bio material. Generated from taxon and sampling if not provided." example:"Genus_sp[SITE|2001-01]"`
	TypeStatus             models.OptionalInput[TypeStatus] `gel:"type_status" json:"type_status,omitzero" doc:"Flag indicating if the bio material is a type specimen, i.e. the reference specimen used to describe a new species."`
	Sources                []string                         `json:"sources,omitzero"`
	VerbatimIdentification models.OptionalInput[string]     `json:"verbatim_identification,omitzero"`
	// OriginalLink       models.OptionalInput[string]  `json:"external_link,omitzero"`
	Quantity           models.OptionalInput[[]int32] `json:"quantity,omitzero" minItems:"1" maxItems:"2"`
	ContentDescription models.OptionalInput[string]  `json:"content_description,omitzero" doc:"Description of the content of the bio material" example:"2 females, 1 juvenile male"`
	Collections        []references.CollectionField  `json:"collections,omitzero"`
	Comments           models.OptionalInput[string]  `json:"comments,omitzero"`
	Sequences          []ExternalSequenceInput       `json:"sequences,omitzero"`
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
	for i, source := range occ.Sources {
		if s, ok := c.DataSources[source]; ok {
			occ.Sources[i] = s
		}
	}
	for i, col := range occ.Collections {
		if c, ok := c.Collections[col.Name]; ok {
			occ.Collections[i].Name = c
		}
	}
	for i := range occ.Sequences {
		(&occ.Sequences[i]).WithCreatedMetadata(c)
	}
	return occ
}

var OccurrenceSaveExecuteQuery = queries.OccurrenceQuery(
	`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	"")

func (i *OccurrenceInput) SaveExecute(e geltypes.Executor, samplingNumber int64) error {
	data, _ := json.Marshal(i)
	logrus.Debugf("Creating Occurrence with args: %s", string(data))
	return e.Execute(context.Background(), OccurrenceSaveExecuteQuery, samplingNumber, data)
}

var OccurrenceQuickSaveQuery = queries.OccurrenceQuery(
	`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	`#edgeql
		{id, code}
	`)

func (i *OccurrenceInput) QuickSave(e geltypes.Executor, samplingNumber int64) (created CreatedCode, err error) {
	data, _ := json.Marshal(i)
	logrus.Debugf("Creating Occurrence with args: %s", string(data))
	err = e.QuerySingle(context.Background(), OccurrenceQuickSaveQuery, &created, samplingNumber, data)
	return
}

var OccurrenceSaveQuery = queries.OccurrenceQuery(
	`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	`#edgeql
		{
			[is occurrence::Occurrence].*,
			sampling: { id, number, performed_on },
			identification: { **, identified_by: { * } },
			meta: { * }
		}
	`)

func (i *OccurrenceInput) Save(e geltypes.Executor, samplingNumber int64) (created BaseOccurrence[SamplingOutline], err error) {
	data, _ := json.Marshal(i)
	logrus.Debugf("Creating Occurrence with args: %s", string(data))
	err = e.QuerySingle(context.Background(),
		OccurrenceSaveQuery,
		&created, samplingNumber, data)
	return
}

type OccurrenceUpdate struct {
	SamplingID     models.OptionalInput[geltypes.UUID]        `json:"sampling_id" format:"uuid"`
	Identification models.OptionalInput[IdentificationUpdate] `gel:"identification" json:"identification,omitempty"`
	Code           models.OptionalInput[string]               `gel:"code" json:"code,omitempty"`
	TypeStatus     models.OptionalNull[TypeStatus]            `gel:"type_status" json:"type_status,omitzero"`
	OriginalSource models.OptionalNull[string]                `gel:"sources" json:"sources,omitempty"`
	// OriginalLink       models.OptionalNull[string]    `gel:"external_link" json:"external_link,omitempty"`
	VerbatimIdentification models.OptionalNull[string]                       `gel:"verbatim_identification" json:"verbatim_identification,omitempty"`
	Quantity               models.OptionalNull[[]int32]                      `gel:"quantity" json:"quantity,omitempty" minItems:"2" maxItems:"2"`
	ContentDescription     models.OptionalNull[string]                       `gel:"content_description" json:"content_description,omitempty"`
	Collections            models.OptionalNull[[]references.CollectionField] `gel:"collections" json:"collections,omitempty"`
	PublishedIn            models.OptionalNull[[]string]                     `gel:"published_in" json:"published_in,omitempty"`
	Comments               models.OptionalNull[string]                       `gel:"comments" json:"comments,omitempty"`
}

func (u OccurrenceUpdate) Save(e geltypes.Executor, code string) (updated BaseOccurrence[SamplingOutline], err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
      with item := <json>$1,
      select (
				update occurrence::ExternalOccurrence
				filter .code = <str>$0
				set { %s }
			 ) {
        [is occurrence::Occurrence].*,
				sampling: { id, number, performed_on },
				identification: { **, identified_by: { * } },
				meta: { * }
      }
    `,
		Mappings: map[string]string{
			"code": "<str>item['code']", // if not explicitly provided, updated code is autogenerated
			"sources": `#edgeql
				(
					select references::DataSource filter .code in <str>json_array_unpack(item['sources'])
				)
			`,
			// "external_link":  "<str>item['external_link']",
			"verbatim_identification": "<str>item['verbatim_identification']",
			"quantity": `#edgeql
				<tuple<lower: int32, upper: int32>>(
					lower := <int32>item['quantity'][0],
					upper := <int32>item['quantity'][1],
					)
			`,
			"content_description": "<str>item['content_description']",
			"collections": `#edgeql
				(
					for col in json_array_unpack(item['collections']) union (
						select references::Collection {
							@vouchers := <array<str>>col['vouchers']
						}
						filter .code = <str>col['name'] or .label = <str>col['name']
					)
				)
			`,
			"item_vouchers":  "<str>json_array_unpack(item['item_vouchers'])",
			"comments":       "<str>item['comments']",
			"type_status":    "<bool>item['type_status']",
			"identification": u.Identification.Value.UpdateQuery(".identification"),
			"published_in": `#edgeql
					(
						select distinct references::Article
						filter .code in <str>json_array_unpack(json_get(item, 'published_in'))
					),
			`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
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
						taxon := (.identification.taxon)
				}
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

type DateRange struct {
	MinDate geltypes.OptionalDateTime `gel:"min_date" json:"min_date,omitzero"`
	MaxDate geltypes.OptionalDateTime `gel:"max_date" json:"max_date,omitzero"`
}

func GetOccurrenceDateRange(db geltypes.Executor) (dateRange DateRange, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence,
		min_date := (
			select min(Occurrence.sampling.performed_on.date)
		)
		select {
			min_date := min(
				select Occurrence
				filter .sampling.performed_on.date != null
				return datetime_get(.sampling.performed_on.date, 'year')
			),
			max_date := max(
				select Occurrence
				filter .sampling.performed_on.date != null
				return datetime_get(.sampling.performed_on.date, 'year')
			)
		}
		`,
		&dateRange)
	return dateRange, err
}
