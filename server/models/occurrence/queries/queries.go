package queries

import (
	"bytes"
	"embed"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

//go:embed *.tmpl.edgeql
var templateFS embed.FS
var queryTemplates = template.Must(template.New("").Funcs(template.FuncMap{
	"ToLower": strings.ToLower,
}).ParseFS(templateFS, "*.tmpl.edgeql"))

type shapeT string

func (q shapeT) Shape() string {
	return string(q)
}

type withShape interface {
	Shape() string
}

func ParseTemplateOrDie(name string, tmplStr string) *template.Template {
	return template.Must(template.New(name).Parse(tmplStr))
}

func renderWithShape(name string, q withShape) string {
	return RenderTemplate(name, q) + " " + q.Shape() + "\n"
}

func CompileQuery(tmpl *template.Template, data any) string {
	buf := new(bytes.Buffer)
	_ = tmpl.Execute(buf, data)
	return buf.String()
}

func RenderTemplate(name string, data any) string {
	var buf = bytes.Buffer{}
	err := queryTemplates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		logrus.Error(err)
	}
	return buf.String()
}

/****************************
 * Sampling
 ****************************/

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
	return renderWithShape("insert_sampling.tmpl.edgeql", q)
}

/****************************
 * Abiotics
 ****************************/

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
	return renderWithShape("insert_abiotic.tmpl.edgeql", q)
}

/****************************
 * Occurrence
 ****************************/

type insertOccurrenceQuery struct {
	Sampling            string
	JSON                string
	InsertSequenceQuery insertExternalSequenceQuery
	shapeT
}

func OccurrenceQuery(sampling string, JSON string, shape string) (res string) {
	q := insertOccurrenceQuery{
		Sampling:            sampling,
		JSON:                JSON,
		InsertSequenceQuery: insertExternalSequenceQuery{"occurrence", "seq_data", ""},
		shapeT:              shapeT(shape),
	}
	return renderWithShape("insert_occurrence.tmpl.edgeql", q)
}

/****************************
 * External Sequence
 ****************************/

type insertExternalSequenceQuery struct {
	Occurrence string
	JSON       string
	shapeT
}

func ExternalSequenceQuery(occurrence string, JSON string, shape string) (res string) {
	q := insertExternalSequenceQuery{
		occurrence, JSON, shapeT(shape),
	}
	return renderWithShape("insert_external_sequence.tmpl.edgeql", q)
}
