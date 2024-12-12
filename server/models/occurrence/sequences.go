package occurrence

import (
	"context"
	"darco/proto/db"
	"darco/proto/models"
	"darco/proto/models/people"
	"darco/proto/models/references"
	"darco/proto/models/sequences"
	"darco/proto/models/vocabulary"
	"encoding/json"

	"github.com/edgedb/edgedb-go"
)

type LegacySeqID struct {
	ID            int32  `edgedb:"id" json:"id"`
	Code          string `edgedb:"code" json:"code"`
	AlignmentCode string `edgedb:"alignment_code" json:"alignment_code"`
}

type Sequence struct {
	Code     string                       `edgedb:"code" json:"code"`
	Label    edgedb.OptionalStr           `edgedb:"label" json:"label"`
	Sequence edgedb.OptionalStr           `edgedb:"sequence" json:"sequence"`
	Gene     sequences.Gene               `edgedb:"gene" json:"gene"`
	LegacyID models.Optional[LegacySeqID] `edgedb:"legacy" json:"legacy"`
}

type ExternalSeqOrigin struct {
	vocabulary.Vocabulary `edgedb:"$inline" json:",inline"`
	AccessionRequired     bool               `edgedb:"accession_required" json:"accession_required"`
	LinkTemplate          edgedb.OptionalStr `edgedb:"link_template" json:"link_template"`
	Meta                  people.Meta        `edgedb:"meta" json:"meta"`
}

type ExternalSeqOriginInput struct {
	vocabulary.VocabularyInput `edgedb:"$inline" json:",inline"`
	AccessionRequired          bool                         `edgedb:"accession_required" json:"accession_required,omitempty" default:"false"`
	LinkTemplate               models.OptionalInput[string] `edgedb:"link_template" json:"link_template"`
}

func (i ExternalSeqOriginInput) Save(e edgedb.Executor) (created ExternalSeqOrigin, err error) {
	data, _ := json.Marshal(i)
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with data := <json>$0,
			select (insert seq::ExternalSeqOrigin {
				label := <str>data['label'],
				code := <str>data['code'],
				description := <str>json_get(data, 'description'),
				accession_required := <bool>json_get(data, 'accession_required'),
				link_template := <str>json_get(data, 'link_template')
			}) { ** }
		`, &created, data)
	return
}

type ExternalSeqOriginUpdate struct {
	vocabulary.VocabularyUpdate `edgedb:"$inline" json:",inline"`
	AccessionRequired           models.OptionalInput[bool]  `edgedb:"accession_required" json:"accession_required"`
	LinkTemplate                models.OptionalNull[string] `edgedb:"link_template" json:"link_template"`
}

func (u ExternalSeqOriginUpdate) Save(e edgedb.Executor, code string) (updated ExternalSeqOrigin, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: `#edgeql
			with item := <json>$1,
			select (update seq::ExternalSeqOrigin filter .code = <str>$0 set {
				%s
			}) { ** }
		`,
		Mappings: u.FieldMappingsWith("item", map[string]string{
			"accession_required": "<bool>item['accession_required']",
			"link_template":      "<str>item['link_template']",
		}),
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, code, data)
	return
}

type ExternalSequence struct {
	Occurrence `edgedb:"$inline" json:",inline"`
	Sequence   `edgedb:"$inline" json:",inline"`
	References []references.Article `edgedb:"references" json:"references"`
	// SourceSample            `edgedb:"source_sample" json:"source_sample"`
	AccessionNumber    edgedb.OptionalStr `edgedb:"accession_number" json:"accession_number"`
	SpecimenIdentifier string             `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr `edgedb:"original_taxon" json:"original_taxon"`
}
