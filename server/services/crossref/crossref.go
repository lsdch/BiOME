package crossref

import (
	"darco/proto/models/settings"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

var client *crossrefapi.CrossRefClient

// Client gives a CrossRef API client with mail-to super admin address
func Client() *crossrefapi.CrossRefClient {
	if client == nil {
		appName := settings.Instance().Name
		mailTo := settings.Get().SuperAdmin.Email
		// Error only occurs if mailTo == "", which is not possible
		client, _ = crossrefapi.NewCrossRefClient(appName, mailTo)
	}
	return client
}

// Rename to avoid name collision in generated OpenAPI spec
type CrossRefPerson crossrefapi.Person

func RetrieveDOI(db edgedb.Executor, doi string) (*crossrefapi.Works, error) {
	data, err := Client().Works(doi)
	if data == nil && err == nil {
		return nil, huma.Error404NotFound("No match found")
	}

	return data, err
}

type BibSearchResults struct {
	Total int                   `json:"total"`
	Items []crossrefapi.Message `json:"items"`
}

func BibliographicSearch(bib string) (*BibSearchResults, error) {
	query := crossrefapi.WorksQuery{
		Fields: &crossrefapi.WorksQueryFields{
			Bibliographic: bib,
		},
		Pagination: &crossrefapi.Pagination{
			Rows: 3,
		},
	}
	data, err := Client().QueryWorks(query)
	if err != nil {
		return nil, err
	}
	return &BibSearchResults{
		Total: int(data.Message.TotalResults),
		Items: data.Message.Items,
	}, nil
}
