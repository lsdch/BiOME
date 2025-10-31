package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/occurrence/queries"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/specimen"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

type SpecimenVoucher struct {
	Collection geltypes.OptionalStr `gel:"in_collection" json:"collection,omitempty"`
	Item       []string             `gel:"item_vouchers" json:"vouchers,omitempty"`
}

type ExternalOccurrenceSpecific struct {
	Sources            []references.DataSource            `gel:"sources" json:"sources,omitempty"`
	PublishedIn        []references.Article               `gel:"published_in" json:"published_in,omitempty"`
	ExternalLink       geltypes.OptionalStr               `gel:"external_link" json:"external_link,omitempty"`
	OriginalTaxon      geltypes.OptionalStr               `gel:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           models.Optional[specimen.Quantity] `gel:"quantity" json:"quantity,omitempty"`
	ContentDescription geltypes.OptionalStr               `gel:"content_description" json:"content_description,omitempty"`
	Archive            SpecimenVoucher                    `gel:"$inline" json:"archive"`
	Sequences          []ExternalSequence                 `gel:"sequences" json:"sequences,omitempty"`
}

type InternalBioMatSpecific struct {
	HasSequences      bool                            `gel:"has_sequences" json:"has_sequences"`
	IsHomogenous      bool                            `gel:"is_homogenous" json:"is_homogenous"`
	IsCongruent       bool                            `gel:"is_congruent" json:"is_congruent"`
	SequenceConsensus models.Optional[taxonomy.Taxon] `gel:"seq_consensus" json:"seq_consensus,omitempty"`
}

type OccurrenceListItem GenericOccurrence[SamplingInnerWithSite]

type BioMaterialWithDetails Occurrence[SamplingWithSite]

func GetOccurrence(db geltypes.Executor, code string) (occurrence Occurrence[SamplingWithSite], err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence
		select OccurrenceWithType {
			**,
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
			internal: {
				is_homogenous,
				is_congruent,
				seq_consensus
			},
			external: {
				sources: { * },
				sequences: { *, gene: { * }, referenced_in: { * }, meta: { * } },
				external_link,
				in_collection,
				item_vouchers,
				quantity,
				published_in: { * },
				content_description
			}
		} filter .code = <str>$0
	`,
		&occurrence, code)
	return occurrence, err
}

type BioMatSortKey string

//generate:enum
const (
	BioMatSortCode         BioMatSortKey = "code"
	BioMatSortSite         BioMatSortKey = "site"
	BioMatSortSamplingDate BioMatSortKey = "sampling_date"
	BioMatSortIdentifiedOn BioMatSortKey = "identified_on"
	BioMatSortTaxon        BioMatSortKey = "taxon"
	BioMatSortIdentifiedBy BioMatSortKey = "identified_by"
	BioMatSortLastUpdated  BioMatSortKey = "last_updated"
)

var BioMatSortMap = map[BioMatSortKey]string{
	BioMatSortCode:         ".code",
	BioMatSortSite:         "(.site.name ?? .site.code)",
	BioMatSortSamplingDate: ".sampling.performed_on.date",
	BioMatSortIdentifiedOn: ".identification.identified_on.date",
	BioMatSortTaxon:        ".identification.taxon.name",
	BioMatSortIdentifiedBy: ".identification.identified_by.last_name",
	BioMatSortLastUpdated:  ".meta.lastUpdated",
}

type ListOccurrencesOptions struct {
	models.Pagination `json:",inline"`
	models.SortBy[BioMatSortKey]
	models.Filter `json:",inline"`
	Category      models.OptionalInput[OccurrenceCategory]   `query:"category" json:"category,omitzero"`
	Taxon         models.OptionalInput[string]               `query:"taxon" json:"taxon,omitzero"`
	WholeClade    bool                                       `query:"whole_clade" json:"whole_clade"`
	HasSequences  models.OptionalInput[bool]                 `query:"has_sequences" json:"has_sequences,omitzero"`
	Confer        models.OptionalInput[bool]                 `query:"confer" json:"confer,omitzero"`
	IsType        models.OptionalInput[bool]                 `query:"is_type" json:"is_type,omitzero"`
	Status        models.OptionalInput[taxonomy.TaxonStatus] `query:"status" json:"status,omitzero"`
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
	params, _ := json.Marshal(opts)
	logrus.Debugf("Params: %s", string(params))
	var result = models.PaginatedList[OccurrenceListItem]{
		Items: []OccurrenceListItem{},
	}
	err := db.QuerySingle(context.Background(),
		`#edgeql
			with module occurrence,
				params := <json>$0,
				search_term := <str>json_get(params, 'search'),
				category := <OccurrenceCategory>json_get(params, 'category'),
				taxon_name := <str>json_get(params, 'taxon'),
				taxon := (
					(select taxonomy::Taxon filter .name = taxon_name)
					if (exists taxon_name)
					else <taxonomy::Taxon>{}
				),
				whole_clade := <bool>params['whole_clade'],
				with_sequences := <bool>json_get(params, 'has_sequences'),
				confer := <bool>json_get(params, 'confer'),
				status := <taxonomy::TaxonStatus>json_get(params, 'status'),
				is_type := <bool>json_get(params, 'is_type'),
				is_own := <bool>params['owned'],
			items := (
				select OccurrenceWithType { * }
				filter (
					(.code ilike '%%' ++ search_term ++ '%%' if exists search_term else true) and
					(.category = category if exists category else true) and
					(
						(
							taxonomy::is_in_clade(.identification.taxon, taxon) if whole_clade
							else .identification.taxon = taxon
						)
						if exists taxon else true
					) and
					(.identification.taxon.status = status if exists status else true) and
					(.has_sequences = with_sequences if exists with_sequences else true) and
					(.identification.confer = confer if exists confer else true) and
					(.is_type = is_type if exists is_type else true) and
					(.meta.created_by_user = global default::current_user if (is_own and exists global default::current_user) else true)
				)
			),
			select {
				items := (
					select items
					order by <str>$1
					offset <optional int64>json_get(params, 'offset')
					limit <optional int64>json_get(params, 'limit')
				) {
					*,
					sampling: { *, site: { *, country: { * } } },
					identification: { **, identified_by: { ** } },
					meta: { * }
				},
				total_count := count(items),
			};
		`,
		&result, params, opts.OrderByString())
	return result, err
}

func DeleteOccurrence(db geltypes.Executor, code string) (deleted OccurrenceListItem, err error) {
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module occurrence
			select (
				delete BioMaterial filter .code = <str>$0
			) {
        *,
				sampling: { *, site: { *, country: { * } } },
				identification: { **, identified_by: { ** } },
				internal:= [is occurrence::InternalBioMat] {
					has_sequences := exists(.specimens.sequences),
					is_homogenous,
					is_congruent,
					seq_consensus
				}
        external:= [is occurrence::ExternalOccurrence]{
					sources,
          external_link,
          in_collection,
          item_vouchers,
          quantity,
          content_description
        }
      }
		`,
		&deleted, code)
	return
}

type InternalOccurrenceInput struct {
	OccurrenceInput `gel:"$inline" json:",inline"`
	// TODO: Internal-specific fields
}

func (i *InternalOccurrenceInput) WithCreatedMetadata(c *CreatedMetadata) *InternalOccurrenceInput {
	i.OccurrenceInput.WithCreatedMetadata(c)
	return i
}

var internalQuickSaveQuery = queries.InternalBioMatQuery(`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	`#edgeql
		{id, code, category}
	`)

func (i *InternalOccurrenceInput) QuickSave(e geltypes.Executor, samplingNumber int64) (created BaseOccurrence, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		internalQuickSaveQuery,
		&created, samplingNumber, data)
	return
}

var internalSaveQuery = queries.InternalBioMatQuery(
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

func (i *InternalOccurrenceInput) Save(e geltypes.Executor, samplingNumber int64) (created GenericOccurrence[SamplingOutline], err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		internalSaveQuery,
		&created, samplingNumber, data)
	return
}

var internalSaveExecuteQuery = queries.InternalBioMatQuery(
	`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	"")

func (i *InternalOccurrenceInput) SaveExecute(e geltypes.Executor, samplingNumber int64) error {
	data, _ := json.Marshal(i)
	return e.Execute(context.Background(),
		internalSaveExecuteQuery, samplingNumber, data)
}

type ExternalOccurrenceInput struct {
	OccurrenceInput    `gel:"$inline" json:",inline"`
	Sources            []string                                `json:"sources,omitzero"`
	OriginalTaxon      models.OptionalInput[string]            `json:"original_taxon,omitzero"`
	OriginalLink       models.OptionalInput[string]            `json:"external_link,omitzero"`
	Quantity           models.OptionalInput[specimen.Quantity] `json:"quantity,omitzero"`
	ContentDescription models.OptionalInput[string]            `json:"content_description,omitzero" doc:"Description of the content of the bio material" example:"2 females, 1 juvenile male"`
	Collection         models.OptionalInput[string]            `json:"collection,omitzero"`
	ItemVouchers       []string                                `json:"vouchers,omitzero"`
	Comments           models.OptionalInput[string]            `json:"comments,omitzero"`
	Sequences          []ExternalSequenceInput                 `json:"sequences,omitzero"`
}

func (bm *ExternalOccurrenceInput) WithCreatedMetadata(c *CreatedMetadata) *ExternalOccurrenceInput {
	bm.OccurrenceInput.WithCreatedMetadata(c)
	for i, source := range bm.Sources {
		if s, ok := c.DataSources[source]; ok {
			bm.Sources[i] = s
		}
	}
	for i := range bm.Sequences {
		(&bm.Sequences[i]).WithCreatedMetadata(c)
	}
	return bm
}

var externalSaveExecuteQuery = queries.ExternalOccurrenceQuery(
	`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	"")

func (i *ExternalOccurrenceInput) SaveExecute(e geltypes.Executor, samplingNumber int64) error {
	data, _ := json.Marshal(i)
	logrus.Debugf("Creating ExternalOccurrence with args: %s", string(data))
	return e.Execute(context.Background(), externalSaveExecuteQuery, samplingNumber, data)
}

var externalQuickSaveQuery = queries.ExternalOccurrenceQuery(
	`#edgeql
		select events::Sampling filter .number = <int64>$0
	`,
	`<json>$1`,
	`#edgeql
		{id, code, category}
	`)

func (i *ExternalOccurrenceInput) QuickSave(e geltypes.Executor, samplingNumber int64) (created BaseOccurrence, err error) {
	data, _ := json.Marshal(i)
	logrus.Debugf("Creating ExternalOccurrence with args: %s", string(data))
	err = e.QuerySingle(context.Background(), externalQuickSaveQuery, &created, samplingNumber, data)
	return
}

var externalSaveQuery = queries.ExternalOccurrenceQuery(
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

func (i ExternalOccurrenceInput) Save(e geltypes.Executor, samplingNumber int64) (created GenericOccurrence[SamplingOutline], err error) {
	data, _ := json.Marshal(i)
	logrus.Debugf("Creating ExternalOccurrence with args: %s", string(data))
	err = e.QuerySingle(context.Background(),
		externalSaveQuery,
		&created, samplingNumber, data)
	return
}

type ExternalOccurrenceUpdate struct {
	OccurrenceUpdate   `gel:"$inline" json:",inline"`
	OriginalSource     models.OptionalNull[string]            `gel:"sources" json:"sources,omitempty"`
	OriginalLink       models.OptionalNull[string]            `gel:"external_link" json:"external_link,omitempty"`
	OriginalTaxon      models.OptionalNull[string]            `gel:"original_taxon" json:"original_taxon,omitempty"`
	Quantity           models.OptionalNull[specimen.Quantity] `gel:"quantity" json:"quantity,omitempty"`
	ContentDescription models.OptionalNull[string]            `gel:"content_description" json:"content_description,omitempty"`
	Collection         models.OptionalNull[string]            `gel:"in_collection" json:"collection,omitempty"`
	ItemVouchers       models.OptionalInput[[]string]         `gel:"item_vouchers" json:"vouchers,omitempty"`
	Comments           models.OptionalNull[string]            `gel:"comments" json:"comments,omitempty"`
	PublishedIn        models.OptionalNull[[]string]          `gel:"published_in" json:"published_in,omitempty"`
}

func (u ExternalOccurrenceUpdate) Save(e geltypes.Executor, code string) (updated GenericOccurrence[SamplingOutline], err error) {
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
			"external_link":       "<str>item['external_link']",
			"original_taxon":      "<str>item['original_taxon']",
			"quantity":            "<occurrence::QuantityType>item['quantity']",
			"content_description": "<str>item['content_description']",
			"in_collection":       "<str>item['collection']",
			"item_vouchers":       "<str>json_array_unpack(item['item_vouchers'])",
			"comments":            "<str>item['comments']",
			"is_type":             "<bool>item['is_type']",
			"identification":      u.Identification.Value.UpdateQuery(".identification"),
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
