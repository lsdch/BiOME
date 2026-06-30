package crossref

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/lib/queue"
	"github.com/lsdch/biome/models"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/sirupsen/logrus"
)

// Rename to avoid name collision in generated OpenAPI spec
type CrossRefPerson crossrefapi.Person
type CrossRefDateRange crossrefapi.DateRange

type BibSearchResults struct {
	Total int                   `json:"total"`
	Items []crossrefapi.Message `json:"items"`
}

type BibResponse = crossrefapi.WorksQueryResponse

// CrossRefScheduler wraps the CrossRef API client to provide concurrent request handling
// with request queuing and throttling capabilities. It manages both DOI-specific and
// general query requests through separate queues while respecting rate limits.
type CrossRefScheduler struct {
	crossrefapi.CrossRefClient
	DoiQueue         queue.Queue[crossrefapi.Works, crossrefapi.CrossRefClient]
	QueryQueue       queue.Queue[crossrefapi.WorksQueryResponse, crossrefapi.CrossRefClient]
	MaxActiveQueries int
	semaphor         chan struct{} // handles throttling by limiting concurrent access to the API
}

// Initializes a CrossRef API client with mail-to super admin address
// and max concurrent requests throttling
func NewClient(cfg config.Config) *CrossRefScheduler {
	appName := cfg.API.ContactName
	mailTo := cfg.API.ContactEmail

	// Error only occurs if mailTo == "", which is not possible
	crefClient, _ := crossrefapi.NewCrossRefClient(appName, mailTo)

	// Very stringent rate limiting at first, may get relaxed after getting API response
	crefClient.RateLimitInterval = 1
	crefClient.RateLimitLimit = cfg.CrossRefMaxConcurrent * 2

	return &CrossRefScheduler{
		CrossRefClient:   *crefClient,
		DoiQueue:         queue.NewQueue[crossrefapi.Works, crossrefapi.CrossRefClient](cfg.CrossRefMaxConcurrent),
		QueryQueue:       queue.NewQueue[crossrefapi.WorksQueryResponse, crossrefapi.CrossRefClient](cfg.CrossRefMaxConcurrent),
		MaxActiveQueries: cfg.CrossRefMaxConcurrent,
	}
}

// RetrieveDOI queries the CrossRef API to retrieve metadata for a given DOI.
// The function enqueues the request to the CrossRef client queue and waits for the response.
// If no match is found or there's an error, returns a 404 Not Found error.
func (c *CrossRefScheduler) RetrieveDOI(doi models.DOI) (*crossrefapi.Works, error) {
	doiString := doi.String()
	queueItem := queue.NewQueueItem(
		func(client *crossrefapi.CrossRefClient) queue.ApiResponse[crossrefapi.Works] {
			logrus.Debugf("Querying crossref for DOI : %s", doiString)
			data, err := client.Works(doiString)
			return queue.ApiResponse[crossrefapi.Works]{
				Data:  data,
				Error: err,
			}
		},
	)

	logrus.Debugf("Enqueuing DOI query for %s", doi)
	c.DoiQueue <- queueItem

	// Wait for response
	resp := <-queueItem.Receiver
	data, err := resp.Data, resp.Error
	if data == nil || err != nil {
		return nil, huma.Error404NotFound(fmt.Sprintf("No match found for DOI %s", doi))
	}

	logrus.Debugf("Received response for DOI %s: %+v", doi, data)

	return data, err
}

func (c *CrossRefScheduler) BibliographicSearch(bib string) (*BibSearchResults, error) {
	query := crossrefapi.WorksQuery{
		Fields: &crossrefapi.WorksQueryFields{
			Bibliographic: bib,
		},
		Pagination: &crossrefapi.Pagination{
			Rows: 3,
		},
	}

	queueItem := queue.NewQueueItem(
		func(client *crossrefapi.CrossRefClient) queue.ApiResponse[BibResponse] {
			data, err := client.QueryWorks(query)
			return queue.ApiResponse[BibResponse]{
				Data:  data,
				Error: err,
			}
		},
	)

	// Add query to queue and wait for result
	c.QueryQueue <- queueItem

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

// Start initiates a continuous processing loop for handling API requests.
// It manages concurrent requests by monitoring active queries against a maximum limit.
// The method processes requests from both DOI and general query queues, executing them
// while respecting rate limiting constraints. When the maximum number of active queries
// is reached, the process waits before accepting new requests.
//
// The method runs indefinitely and spawns goroutines for each request execution.
// Each request is processed asynchronously, and the active query count is decremented
// upon completion.
func (c *CrossRefScheduler) Start() {
	logrus.Debugf("[crossref] Starting CrossRef Scheduler with max concurrent requests: %d", c.MaxActiveQueries)
	c.semaphor = make(chan struct{}, c.MaxActiveQueries)
	for {
		// if c.ActiveQueries >= c.MaxActiveQueries {
		// 	time.Sleep(time.Millisecond * 300)
		// 	continue
		// }
		var item queue.ApiRequestItem[crossrefapi.CrossRefClient]
		select {
		case item = <-c.DoiQueue:
			logrus.Infof("[crossref] Received DOI query request; active queries: %d", len(c.semaphor))
			go c.Execute(item)
		case item = <-c.QueryQueue:
			logrus.Infof("[crossref] Received general query request; active queries: %d", len(c.semaphor))
			go c.Execute(item)
		}
	}
}

func (c *CrossRefScheduler) Execute(item queue.ApiRequestItem[crossrefapi.CrossRefClient]) {
	c.semaphor <- struct{}{} // acquire a slot
	defer func() {
		<-c.semaphor // release the slot after execution
	}()
	logrus.Infof("[crossref] Sending query ; interval: %d ; limit: %d; active: %d", c.RateLimitInterval, c.RateLimitLimit, len(c.semaphor))
	item.Execute(&c.CrossRefClient)
}
