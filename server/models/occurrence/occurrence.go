package occurrence

import (
	_ "embed"
)

// type QuantityRange struct {
// 	Lower int32 `gel:"lower" json:"lower"`
// 	Upper int32 `gel:"upper" json:"upper"`
// }

// type TypeStatus string

// //generate:enum skip-gel-unmarshal
// const (
// 	Holotype TypeStatus = "Holotype"
// 	Neotype  TypeStatus = "Neotype"
// 	Topotype TypeStatus = "Topotype"
// )

// type CollectionField struct {
// 	Name     string   `gel:"name" json:"name"`
// 	Vouchers []string `gel:"vouchers" json:"vouchers,omitempty"`
// }

// type BaseOccurrence[SamplingType any] struct {
// 	ID                     geltypes.UUID `gel:"id" json:"id" format:"uuid"`
// 	CodeIdentifier         `gel:"$inline"`
// 	HasSequences           bool                           `gel:"has_sequences" json:"has_sequences"`
// 	Sampling               SamplingType                   `gel:"sampling" json:"sampling"`
// 	Identification         Identification                 `gel:"identification" json:"identification"`
// 	TypeStatus             models.Optional[TypeStatus]    `gel:"type_status" json:"type_status,omitzero" nameHint:"TypeStatus"`
// 	Comments               geltypes.OptionalStr           `gel:"comments" json:"comments,omitempty"`
// 	Quantity               models.Optional[QuantityRange] `gel:"quantity" json:"quantity,omitempty"`
// 	ContentDescription     geltypes.OptionalStr           `gel:"content_description" json:"content_description,omitempty"`
// 	VerbatimIdentification geltypes.OptionalStr           `gel:"verbatim_identification" json:"verbatim_identification,omitempty"`
// 	Collections            []CollectionField              `gel:"collections" json:"collections,omitempty"`
// 	Meta                   people.Meta                    `gel:"meta" json:"meta"`
// }

// type OccurrenceListItem BaseOccurrence[SamplingInnerWithSite]

// type Occurrence[SamplingType any] struct {
// 	BaseOccurrence[SamplingType] `gel:"$inline"`
// 	// ExternalLink       geltypes.OptionalStr           `gel:"external_link" json:"external_link,omitempty"`
// 	PublishedIn []references.Article    `gel:"published_in" json:"published_in,omitempty"`
// 	Sources     []references.DataSource `gel:"sources" json:"sources,omitempty"`
// 	Datasets    []dataset.Dataset       `gel:"datasets" json:"datasets,omitempty"`
// 	Sequences   []ExternalSequence      `gel:"sequences" json:"sequences,omitempty"`
// }

// func GetOccurrence(db geltypes.Executor, code string) (occurrence Occurrence[SamplingWithSite], err error) {
// 	err = db.QuerySingle(context.Background(),
// 		`#edgeql
// 		with module occurrence
// 		select Occurrence {
// 			*,
// 			meta: { * },
// 			datasets: { *, maintainers: { * }, meta: { * } },
// 			sampling: {
// 				*,
// 				target_taxa: { * },
// 				fixatives: { * },
// 				methods: { * },
// 				habitats: { * },
// 				occurrences: { *, identification: { ** } },
// 				site: { *, country: { * } }
// 			},
// 			identification: { ** },
// 			sources: { * },
// 			sequences: { *, gene: { * }, referenced_in: { * }, meta: { * } },
// 			# external_link,
// 			collections,
// 			quantity,
// 			published_in: { * },
// 			content_description
// 		} filter .code = <str>$0
// 	`,
// 		&occurrence, code)
// 	return occurrence, err
// }

// type BioMatSortKey string

// //generate:enum
// const (
// 	BioMatSortCode           BioMatSortKey = "code"
// 	BioMatSortSite           BioMatSortKey = "site"
// 	BioMatSortSamplingDate   BioMatSortKey = "sampling_date"
// 	BioMatSortIdentifiedOn   BioMatSortKey = "identified_on"
// 	BioMatSortTaxon          BioMatSortKey = "taxon"
// 	BioMatSortIdentification BioMatSortKey = "identification"
// 	BioMatSortLastUpdated    BioMatSortKey = "last_updated"
// )

// var BioMatSortMap = map[BioMatSortKey]string{
// 	BioMatSortCode:           ".code",
// 	BioMatSortSite:           "(.sampling.site.name ?? .sampling.site.code)",
// 	BioMatSortSamplingDate:   ".sampling.performed_on.date",
// 	BioMatSortIdentification: ".identification.taxon.name",
// 	BioMatSortIdentifiedOn:   ".identification.identified_on.date",
// 	BioMatSortTaxon:          ".identification.taxon.name",
// 	BioMatSortLastUpdated:    ".meta.lastUpdated",
// }

// type ListOccurrencesOptions struct {
// 	models.Pagination
// 	models.SortBy[BioMatSortKey]
// 	models.Filter
// 	taxonomy.TaxaFilters
// 	Year         models.Optional[int32]                `query:"year" json:"year,omitzero"`
// 	YearEnd      models.OptionalNull[int32]            `query:"year_end" json:"year_end,omitzero"`
// 	Datasets     []string                              `query:"datasets" json:"datasets,omitzero"`
// 	HasSequences models.Optional[bool]                 `query:"has_sequences" json:"has_sequences,omitzero"`
// 	Confer       models.Optional[bool]                 `query:"confer" json:"confer,omitzero"`
// 	TypeStatus   models.Optional[[]TypeStatus]         `query:"type_status" json:"type_status,omitzero"`
// 	Status       models.Optional[taxonomy.TaxonStatus] `query:"status" json:"status,omitzero"`
// 	Rank         []taxonomy.TaxonRank                  `query:"rank" json:"rank,omitzero"`
// }

// func (o ListOccurrencesOptions) Options() ListOccurrencesOptions {
// 	return o
// }

// func (i ListOccurrencesOptions) OrderByString() string {
// 	if i.SortBy.Key == "" {
// 		return ""
// 	}
// 	if term, ok := BioMatSortMap[i.SortBy.Key]; ok {
// 		return term + " " + string(i.SortBy.Order)
// 	} else {
// 		logrus.Warnf("Unknown sort key: %s", i.SortBy.Key)
// 		return ""
// 	}
// }

// func ListOccurrences(db geltypes.Executor, opts ListOccurrencesOptions) (models.PaginatedList[OccurrenceListItem], error) {
// 	if err := opts.TaxaFilters.FetchTaxa(db); err != nil {
// 		return models.PaginatedList[OccurrenceListItem]{}, err
// 	}
// 	params, _ := json.Marshal(opts)
// 	logrus.Debugf("Params: %s", string(params))

// 	var result = models.PaginatedList[OccurrenceListItem]{
// 		Items: []OccurrenceListItem{},
// 	}
// 	err := db.QuerySingle(context.Background(),
// 		queries.RenderTemplate("list_occurrences.tmpl.edgeql", opts),
// 		&result, params)
// 	return result, err
// }

// func DeleteOccurrence(db geltypes.Executor, code string) (deleted OccurrenceListItem, err error) {
// 	err = db.QuerySingle(context.Background(),
// 		`#edgeql
// 		with module occurrence
// 			select (
// 				delete Occurrence filter .code = <str>$0
// 			) {
//         *,
// 				sampling: { *, site: { *, country: { * } } },
// 				identification: { ** },
//       }
// 		`,
// 		&deleted, code)
// 	return
// }

// type SamplingDateWithOccurrences struct {
// 	ID          geltypes.UUID             `gel:"id" json:"id" format:"uuid"`
// 	Date        OptionalDateWithPrecision `gel:"performed_on" json:"date,omitzero"`
// 	Occurrences []OccurrenceAtSite        `gel:"occurrences" json:"occurrences"`
// }

// type SamplingDetailsWithOccurrences struct {
// 	SamplingInner `gel:"$inline"`
// 	Occurrences   []OccurrenceAtSite `gel:"occurrences" json:"occurrences"`
// 	Meta          people.Meta        `gel:"meta" json:"meta,omitempty"`
// }

// type OccurrenceAtSite struct {
// 	ID             geltypes.UUID                  `gel:"id" json:"id" format:"uuid"`
// 	Code           string                         `gel:"code" json:"code"`
// 	Identification IdentificationWithLineageNames `gel:"identification" json:"identification"`
// 	// SamplingDate      DateWithPrecision   `gel:"sampling_date" json:"sampling_date"`
// }

// type SiteWithOccurrences struct {
// 	SiteItem  `gel:"$inline"`
// 	Samplings []SamplingDateWithOccurrences `gel:"samplings" json:"samplings"`
// }

// func ListSamplingsAtSite(db geltypes.Executor, siteCode string) ([]SamplingDetailsWithOccurrences, error) {
// 	var samplings = []SamplingDetailsWithOccurrences{}
// 	err := db.Query(context.Background(),
// 		`#edgeql
// 			with module occurrence,
// 				site := (select location::Site filter .code = <str>$0),
// 			select site.samplings {
// 				id,
// 				number,
// 				performed_on,
// 				target_taxa: { * },
// 				habitats: { * },
// 				fixatives: { * },
// 				methods: { * },
// 				meta: { * },
// 				occurrences := (
// 					select .occurrences {
// 						id,
// 						code,
// 						identification: { identified_on, confer, addendum, taxon: { * } },
// 					}
// 				)
// 			}
// 		`,
// 		&samplings, siteCode)
// 	return samplings, err
// }

// type SiteSamplingStatus string

// //generate:enum
// const (
// 	IncludeAllSites        SiteSamplingStatus = "All"
// 	IncludeSampled         SiteSamplingStatus = "Sampled"
// 	IncludeWithOccurrences SiteSamplingStatus = "Occurrences"
// )

// type OccurrencesBySiteOptions struct {
// 	ListSitesOptions
// 	taxonomy.TaxaFilters
// 	Habitats       []string             `json:"habitats,omitempty" query:"habitats"`
// 	SamplingTarget taxonomy.TaxaFilters `json:"sampling_target" query:"sampling_target"`
// 	IncludeSites   SiteSamplingStatus   `json:"include_sites,omitempty" query:"include_sites" default:"Occurrences" doc:"Include sites with occurrences, sampled sites or all sites. Defaults to sites with at least one occurrence."`
// }

// func (o OccurrencesBySiteOptions) Options() OccurrencesBySiteOptions {
// 	return o
// }

// func OccurrencesBySite(db geltypes.Executor, opts OccurrencesBySiteOptions) ([]SiteWithOccurrences, error) {
// 	if err := opts.FetchTaxa(db); err != nil {
// 		return nil, err
// 	}
// 	var sites []SiteWithOccurrences
// 	filters, _ := json.Marshal(opts)
// 	err := db.Query(context.Background(),
// 		queries.RenderTemplate("occurrences_by_site.tmpl.edgeql", opts),
// 		&sites, filters)
// 	return sites, err
// }

// OccurrenceOverviewItem is a representation of the occurrences count for one taxon
// type OccurrenceOverviewItem struct {
// 	Name        string               `gel:"name" json:"name"`
// 	ParentName  geltypes.OptionalStr `gel:"parent_name" json:"parent_name,omitempty"`
// 	Occurrences int32                `gel:"occurrences" json:"occurrences"`
// 	Rank        taxonomy.TaxonRank   `gel:"rank" json:"rank"`
// }

// type occurrenceOverviewQueryResult struct {
// 	Occurrences   []OccurrenceOverviewItem `gel:"occurrences"`
// 	NoOccurrences []OccurrenceOverviewItem `gel:"no_occurrences"`
// }

// func (o occurrenceOverviewQueryResult) toItems() []OccurrenceOverviewItem {
// 	return slices.Concat(o.Occurrences, o.NoOccurrences)
// }

// // OccurrenceOverview returns the count of occurrences for each taxon
// func OccurrenceOverview(db geltypes.Executor) ([]OccurrenceOverviewItem, error) {
// 	var items = occurrenceOverviewQueryResult{}
// 	err := db.QuerySingle(context.Background(),
// 		`#edgeql
// 			with module occurrence,
// 			occ := (
// 				select Occurrence {
// 						taxon := (.identification.taxon)
// 				}
// 			),
// 			groups := (select (group occ by .taxon) { arity := count(.elements)}),
// 			noOccTaxa := (
// 				select (taxonomy::Taxon except occ.taxon)
// 			),

// 			select {
// 				occurrences := groups {
// 					required name:= .key.taxon.name,
// 					required rank:= .key.taxon.rank,
// 					parent_name:= .key.taxon.parent.name,
// 					required occurrences := <int32>.arity
// 				},
// 				no_occurrences := (
// 					# Rest of taxonomy i.e. taxa having no occurrences
// 					select noOccTaxa {
// 					required name := .name,
// 					required rank := .rank,
// 					parent_name:= noOccTaxa.parent.name,
// 					required occurrences := <int32>0
// 				})
// 			}
// 		`,
// 		&items)
// 	return items.toItems(), err
// }
