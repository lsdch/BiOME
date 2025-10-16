package occurrence

import (
	"context"
	"encoding/json"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/models/specimen"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

type SpecimenVoucher struct {
	Collection geltypes.OptionalStr `gel:"in_collection" json:"collection,omitempty"`
	Item       []string             `gel:"item_vouchers" json:"vouchers,omitempty"`
}

type ExternalBioMatSpecific struct {
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
			identification: { ** },
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
				published_in,
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

type ListBioMaterialOptions struct {
	models.Pagination `json:",inline"`
	models.SortBy[BioMatSortKey]
	models.Filter `json:",inline"`
	Category      models.OptionalInput[OccurrenceCategory] `query:"category" json:"category,omitzero"`
	Taxon         models.OptionalInput[string]             `query:"taxon" json:"taxon,omitzero"`
	WholeClade    bool                                     `query:"whole_clade" json:"whole_clade"`
	HasSequences  models.OptionalInput[bool]               `query:"has_sequences" json:"has_sequences,omitzero"`
	IsType        models.OptionalInput[bool]               `query:"is_type" json:"is_type,omitzero"`
}

func (o ListBioMaterialOptions) Options() ListBioMaterialOptions {
	return o
}

func (i ListBioMaterialOptions) OrderByString() string {
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

func ListOccurrences(db geltypes.Executor, opts ListBioMaterialOptions) (models.PaginatedList[OccurrenceListItem], error) {
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
					(.has_sequences = with_sequences if exists with_sequences else true) and
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
					identification: { **, identified_by: { * } },
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
				identification: { **, identified_by: { * } },
				internal:= [is occurrence::InternalBioMat] {
					has_sequences := exists(.specimens.sequences),
					is_homogenous,
					is_congruent,
					seq_consensus
				}
        external:= [is occurrence::ExternalBioMat]{
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

func (i InternalOccurrenceInput) Save(e geltypes.Executor, samplingNumber int64) (created GenericOccurrence[SamplingOutline], err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with
				sampling := (
					assert_exists(
						(select events::Sampling filter .number = <int64>$0),
						message := "Failed to find sampling with number: " ++ <str><int64>$0
					)
				),
				data := <json>$1,
				identification := data['identification'],
				taxon := taxonomy::taxonByName(<str>identification['taxon']),
			select (insert occurrence::InternalBioMat {
				code := <str>json_get(data, 'code') ?? occurrence::occurrence_code(taxon, sampling),
				identification := (
					insert occurrence::Identification {
						taxon := taxon,
						identified_by := people::personByAlias(<str>identification['identified_by']),
						identified_on := date::from_json_with_precision(identification['identified_on']),
					}
				),
				sampling := sampling,
				is_type := <bool>json_get(data, 'is_type') ?? false,
			}) {
				[is occurrence::Occurrence].*,
				sampling: { id, number, performed_on },
				identification: { **, identified_by: { * } },
				meta: { * }
			}
		`, &created, samplingNumber, data)
	return
}

type ExternalOccurrenceInput struct {
	OccurrenceInput    `gel:"$inline" json:",inline"`
	OriginalSource     models.OptionalInput[string]            `json:"sources,omitzero"`
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
	if dataSource, ok := bm.OriginalSource.Get(); ok {
		if s, ok := c.DataSources[dataSource]; ok {
			bm.OriginalSource = (&bm.OriginalSource).SetValue(s)
		}
	}
	for i := range bm.Sequences {
		(&bm.Sequences[i]).WithCreatedMetadata(c)
	}
	return bm
}

func (i ExternalOccurrenceInput) Save(e geltypes.Executor, samplingNumber int64) (created GenericOccurrence[SamplingOutline], err error) {
	data, _ := json.Marshal(i)
	logrus.Infof("Creating ExternalBioMat with args: %s", string(data))
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with
				sampling := (assert_exists(
					(select events::Sampling filter .number = <int64>$0),
					message := "Failed to find sampling with number: " ++ <str><int64>$0
				)),
				data := <json>$1,
				biomat := occurrence::insert_external_biomat(sampling, data),
				sequences := (
					for seq_data in json_array_unpack(json_get(data, 'sequences')) union (
						insert seq::ExternalSequence {
							biomat := biomat,
							code := <str>seq_data['code'],
							label := <str>seq_data['label'],
							sequence := <str>seq_data['sequence'],
							gene := seq::geneByCode(<str>seq_data['gene']),
							legacy := <tuple<id: int32, code: str, alignment_code: str>>json_get(seq_data, 'legacy'),
							specimen_identifier := <str>seq_data['specimen_identifier'],
							referenced_in := (
								for ref in json_array_unpack(json_get(seq_data, 'referenced_in'))
								union (
									insert references::SeqReference {
										db := references::dataSourceByCode(<str>ref['db']),
										accession := <str>ref['accession'],
									}
								)
							)
						}
					)
				)
			select biomat {
        [is occurrence::Occurrence].*,
				sampling: { id, number, performed_on },
				identification: { **, identified_by: { * } },
				meta: { * }
      }
		`, &created, samplingNumber, data)
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
				update occurrence::ExternalBioMat
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
