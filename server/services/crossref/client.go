package crossref

import (
	"time"

	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/services"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/sirupsen/logrus"
)

const maxConcurrentRequests = 5

// CrossRefScheduler wraps the CrossRef API client to provide concurrent request handling
// with request queuing and throttling capabilities. It manages both DOI-specific and
// general query requests through separate queues while respecting rate limits.
type CrossRefScheduler struct {
	crossrefapi.CrossRefClient
	DoiQueue         services.Queue[crossrefapi.Works, crossrefapi.CrossRefClient]
	QueryQueue       services.Queue[crossrefapi.WorksQueryResponse, crossrefapi.CrossRefClient]
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
func (c CrossRefScheduler) Start() {
	for {
		if c.ActiveQueries >= c.MaxActiveQueries {
			time.Sleep(time.Millisecond * 300)
			continue
		}
		var item services.ApiRequestItem[crossrefapi.CrossRefClient]
		select {
		case item = <-c.DoiQueue:
		case item = <-c.QueryQueue:
		}
		c.ActiveQueries++
		logrus.Debugf("Sending query ; interval: %d ; limit: %d; active: %d", c.RateLimitInterval, c.RateLimitLimit, c.ActiveQueries)
		go func() {
			item.Execute(&c.CrossRefClient)
			c.ActiveQueries--
		}()
	}
}

var client *CrossRefScheduler

// Initializes a CrossRef API client with mail-to super admin address
// and max concurrent requests throttling
func newClient(maxConcurrent int) *CrossRefScheduler {
	appName := settings.Instance().Name
	mailTo := settings.Get().SuperAdmin.Email

	// Error only occurs if mailTo == "", which is not possible
	crefClient, _ := crossrefapi.NewCrossRefClient(appName, mailTo)

	// Very stringent rate limiting at first, may get relaxed after getting API response
	crefClient.RateLimitInterval = 1
	crefClient.RateLimitLimit = maxConcurrentRequests * 2

	client = &CrossRefScheduler{
		CrossRefClient:   *crefClient,
		DoiQueue:         services.NewQueue[crossrefapi.Works, crossrefapi.CrossRefClient](maxConcurrentRequests),
		QueryQueue:       services.NewQueue[crossrefapi.WorksQueryResponse, crossrefapi.CrossRefClient](maxConcurrentRequests),
		ActiveQueries:    0,
		MaxActiveQueries: maxConcurrent,
	}
	return client
}

func init() {
	client = newClient(maxConcurrentRequests)
	go client.Start()
}

func Client() *CrossRefScheduler {
	return client
}
