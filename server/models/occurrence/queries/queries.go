package queries

import (
	"bytes"
	_ "embed"
	"text/template"
)

type shapeT string

func (q shapeT) Shape() string {
	return string(q)
}

type withShape interface {
	Shape() string
}

func parseTemplateOrDie(name string, tmplStr string) *template.Template {
	return template.Must(template.New(name).Parse(tmplStr))
}

func buildQueryWithShape(q withShape, templ *template.Template) string {
	buf := new(bytes.Buffer)
	_ = templ.Execute(buf, q)
	return buf.String() + " " + q.Shape() + "\n"
}

/****************************
 * Sampling
 ****************************/

//go:embed insert_sampling.tmpl.edgeql
var insertSamplingQueryFile string
var insertSamplingTemplate = parseTemplateOrDie(
	"insert_sampling",
	insertSamplingQueryFile,
)

type insertSamplingQuery struct {
	Site string
	JSON string
	shapeT
}

func SamplingQuery(site string, JSON string, shape string) (res string) {
	q := insertSamplingQuery{
		Site:   site,
		JSON:   JSON,
		shapeT: shapeT(shape),
	}
	return buildQueryWithShape(q, insertSamplingTemplate)
}

/****************************
 * Abiotics
 ****************************/

//go:embed insert_abiotic.tmpl.edgeql
var insertAbioticQueryFile string
var insertAbioticTemplate = parseTemplateOrDie(
	"insert_abiotic",
	insertAbioticQueryFile,
)

type insertAbioticQuery struct {
	Site string
	JSON string
	shapeT
}

func AbioticQuery(site string, JSON string, shape string) (res string) {
	q := insertAbioticQuery{
		Site:   site,
		JSON:   JSON,
		shapeT: shapeT(shape),
	}
	return buildQueryWithShape(q, insertAbioticTemplate)
}

/****************************
 * Internal Biomaterial
 ****************************/

//go:embed insert_internal_biomat.tmpl.edgeql
var insertInternalBioMatQueryFile string
var insertInternalBioMatTemplate = parseTemplateOrDie(
	"insert_internal_biomat",
	insertInternalBioMatQueryFile,
)

type insertInternalBioMatQuery struct {
	Sampling string
	JSON     string
	shapeT
}

func InternalBioMatQuery(sampling string, JSON string, shape string) (res string) {
	q := insertInternalBioMatQuery{
		Sampling: sampling,
		JSON:     JSON,
		shapeT:   shapeT(shape),
	}
	return buildQueryWithShape(q, insertInternalBioMatTemplate)
}

/****************************
 * External Occurrence
 ****************************/

//go:embed insert_external_occurrence.tmpl.edgeql
var insertExternalOccurrenceQueryFile string
var insertExternalOccurrenceTemplate = parseTemplateOrDie(
	"insert_external_occurrence",
	insertExternalOccurrenceQueryFile,
)

type insertExternalOccurrenceQuery struct {
	Sampling            string
	JSON                string
	InsertSequenceQuery string
	shapeT
}

func ExternalOccurrenceQuery(sampling string, JSON string, shape string) (res string) {
	q := insertExternalOccurrenceQuery{
		Sampling:            sampling,
		JSON:                JSON,
		InsertSequenceQuery: ExternalSequenceQuery("occurrence", "seq_data", ""),
		shapeT:              shapeT(shape),
	}
	return buildQueryWithShape(q, insertExternalOccurrenceTemplate)
}

/****************************
 * External Sequence
 ****************************/

//go:embed insert_external_sequence.tmpl.edgeql
var insertExternalSequenceQueryFile string
var insertExternalSequenceTemplate = parseTemplateOrDie(
	"insert_external_sequence",
	insertExternalSequenceQueryFile,
)

type insertExternalSequenceQuery struct {
	Occurrence string
	JSON       string
	shapeT
}

func ExternalSequenceQuery(occurrence string, JSON string, shape string) (res string) {
	q := insertExternalSequenceQuery{
		occurrence, JSON, shapeT(shape),
	}
	return buildQueryWithShape(q, insertExternalSequenceTemplate)
}
