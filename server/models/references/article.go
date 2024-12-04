package references

import (
	"darco/proto/models/people"

	"github.com/edgedb/edgedb-go"
)

type Article struct {
	ID       edgedb.UUID        `edgedb:"id" json:"id" format:"uuid"`
	Authors  []string           `edgedb:"authors" json:"authors"`
	Year     int32              `edgedb:"year" json:"year"`
	Title    string             `edgedb:"title" json:"title"`
	Journal  edgedb.OptionalStr `edgedb:"journal" json:"journal"`
	Verbatim edgedb.OptionalStr `edgedb:"verbatim" json:"verbatim"`
	DOI      edgedb.OptionalStr `edgedb:"doi" json:"doi"`
	Comments edgedb.OptionalStr `edgedb:"comments" json:"comments"`
	Meta     people.Meta        `edgedb:"meta" json:"meta"`
}
