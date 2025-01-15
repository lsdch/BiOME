package crossref

import (
	"darco/proto/services"
	"fmt"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

// RetrieveDOI queries the CrossRef API to retrieve metadata for a given DOI.
// The function enqueues the request to the CrossRef client queue and waits for the response.
// If no match is found or there's an error, returns a 404 Not Found error.
func RetrieveDOI(db edgedb.Executor, doi string) (*crossrefapi.Works, error) {

	queueItem := services.NewQueueItem(
		func(client *crossrefapi.CrossRefClient) services.ApiResponse[crossrefapi.Works] {
			data, err := client.Works(doi)
			return services.ApiResponse[crossrefapi.Works]{
				Data:  data,
				Error: err,
			}
		},
	)

	Client().DoiQueue <- queueItem

	// Wait for response
	resp := <-queueItem.Receiver
	data, err := resp.Data, resp.Error
	if data == nil || err == nil {
		return nil, huma.Error404NotFound("No match found")
	}

	return data, err
}

type BibSearchResults struct {
	Total int                   `json:"total"`
	Items []crossrefapi.Message `json:"items"`
}

type BibResponse = crossrefapi.WorksQueryResponse

func BibliographicSearch(bib string) (*BibSearchResults, error) {
	query := crossrefapi.WorksQuery{
		Fields: &crossrefapi.WorksQueryFields{
			Bibliographic: bib,
		},
		Pagination: &crossrefapi.Pagination{
			Rows: 3,
		},
	}

	queueItem := services.NewQueueItem(
		func(client *crossrefapi.CrossRefClient) services.ApiResponse[BibResponse] {
			data, err := client.QueryWorks(query)
			return services.ApiResponse[BibResponse]{
				Data:  data,
				Error: err,
			}
		},
	)

	// Add query to queue and wait for result
	Client().QueryQueue <- queueItem

	// Wait for response
	resp := <-queueItem.Receiver
	if resp.Error != nil {
		return nil, resp.Error
	}
	data := resp.Data

	if data == nil || data.Message == nil {
		return nil, fmt.Errorf("Empty response: %+v", data)
	}
	return &BibSearchResults{
		Total: int(data.Message.TotalResults),
		Items: data.Message.Items,
	}, nil
}
