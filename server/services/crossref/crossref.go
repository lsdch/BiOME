package crossref

import (
	"darco/proto/models/settings"
	"fmt"
	"time"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

// Rename to avoid name collision in generated OpenAPI spec
type CrossRefPerson crossrefapi.Person

// ConcurrentClient wraps the CrossRef API client to provide concurrent request handling
// with request queuing and throttling capabilities. It manages both DOI-specific and
// general query requests through separate queues while respecting rate limits.
type ConcurrentClient struct {
	crossrefapi.CrossRefClient
	DoiQueue         Queue[crossrefapi.Works]
	QueryQueue       Queue[crossrefapi.WorksQueryResponse]
	ActiveQueries    int
	MaxActiveQueries int
}

// Start initiates a continuous processing loop for handling API requests.
// It manages concurrent requests by monitoring active queries against a maximum limit.
// The method processes requests from both DOI and general query queues, executing them
// while respecting rate limiting constraints. When the maximum number of active queries
// is reached, the process waits before accepting new requests.
//
// The method runs indefinitely and spawns goroutines for each request execution.
// Each request is processed asynchronously, and the active query count is decremented
// upon completion.
func (c ConcurrentClient) Start() {
	for {
		if client.ActiveQueries >= client.MaxActiveQueries {
			time.Sleep(time.Millisecond * 300)
			continue
		}
		var item ApiRequestItem
		select {
		case item = <-client.DoiQueue:
		case item = <-client.QueryQueue:
		}
		client.ActiveQueries++
		logrus.Debugf("Sending query ; interval: %d ; limit: %d; active: %d", client.RateLimitInterval, client.RateLimitLimit, client.ActiveQueries)
		go func() {
			item.Execute(&client.CrossRefClient)
			client.ActiveQueries--
		}()
	}
}

var client *ConcurrentClient

// Initializes a CrossRef API client with mail-to super admin address
// and max concurrent requests throttling
func newClient(maxConcurrent int) *ConcurrentClient {
	appName := settings.Instance().Name
	mailTo := settings.Get().SuperAdmin.Email
	// Error only occurs if mailTo == "", which is not possible
	crefClient, _ := crossrefapi.NewCrossRefClient(appName, mailTo)
	// Very stringent rate limiting at first, may get relaxed after getting API response
	crefClient.RateLimitInterval = 1
	crefClient.RateLimitLimit = 5
	client = &ConcurrentClient{
		CrossRefClient:   *crefClient,
		DoiQueue:         NewQueue[crossrefapi.Works](5),
		QueryQueue:       NewQueue[crossrefapi.WorksQueryResponse](5),
		ActiveQueries:    0,
		MaxActiveQueries: maxConcurrent,
	}
	return client
}

func init() {
	client = newClient(5)
	go client.Start()
}

func Client() *ConcurrentClient {
	return client
}

// RetrieveDOI queries the CrossRef API to retrieve metadata for a given DOI.
// The function enqueues the request to the CrossRef client queue and waits for the response.
// If no match is found or there's an error, returns a 404 Not Found error.
func RetrieveDOI(db edgedb.Executor, doi string) (*crossrefapi.Works, error) {

	queueItem := QueueItem[crossrefapi.Works]{
		Query: func(client *crossrefapi.CrossRefClient) ApiResponse[crossrefapi.Works] {
			data, err := client.Works(doi)
			return ApiResponse[crossrefapi.Works]{
				Data:  data,
				Error: err,
			}
		},
		Receiver: make(chan ApiResponse[crossrefapi.Works]),
	}

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

	queueItem := QueueItem[BibResponse]{
		Query: func(client *crossrefapi.CrossRefClient) ApiResponse[BibResponse] {
			data, err := client.QueryWorks(query)
			return ApiResponse[BibResponse]{
				Data:  data,
				Error: err,
			}
		},
		Receiver: make(chan ApiResponse[BibResponse]),
	}

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
