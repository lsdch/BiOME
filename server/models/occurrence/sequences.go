package occurrence

import (
	"darco/proto/models"
	"darco/proto/models/references"
	"darco/proto/models/sequences"

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

type ExternalSeqCategory string

//generate:enum
const (
	NCBI    ExternalSeqCategory = "NCBI"
	PersCom ExternalSeqCategory = "PersCom"
)

type ExternalSequence struct {
	Occurrence `edgedb:"$inline" json:",inline"`
	Sequence   `edgedb:"$inline" json:",inline"`
	References []references.Article `edgedb:"references" json:"references"`
	// SourceSample            `edgedb:"source_sample" json:"source_sample"`
	AccessionNumber    edgedb.OptionalStr `edgedb:"accession_number" json:"accession_number"`
	SpecimenIdentifier string             `edgedb:"specimen_identifier" json:"specimen_identifier"`
	OriginalTaxon      edgedb.OptionalStr `edgedb:"original_taxon" json:"original_taxon"`
}
