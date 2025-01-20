package geoapify

import (
	"bytes"
	"darco/proto/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/sirupsen/logrus"
)

const maxBatchSize = 1000

type GeoapifyClient struct {
	apiKey string
	client *http.Client
}

type GeoapifyClientScheduler struct {
	GeoapifyClient
	BatchRequests    services.Queue[GeoapifyResponse, GeoapifyClient]
	ActiveQueries    int
	MaxActiveQueries int
}

type GeoapifyPendingResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

func (c GeoapifyClient) AwaitResult(p GeoapifyPendingResponse) (*GeoapifyResponse, error) {
	time.Sleep(5 * time.Second)
	for {
		resp, err := c.client.Get(p.URL)
		if err != nil {
			return nil, err
		}

		var result GeoapifyResponse
		body, err := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case 200:
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(body, &result)
			return &result, err
		case 202:
			logrus.Infof("Response pending: %+v", string(body))
			time.Sleep(30 * time.Second)
		}
		resp.Body.Close()
	}
}

type GeoapifyResponse []struct {
	Formatted    string  `json:"formatted"`
	City         string  `json:"city"`
	State        string  `json:"state"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	PostalCode   string  `json:"postcode"`
	Street       string  `json:"street"`
	HouseNumber  string  `json:"housenumber"`
	Suburb       string  `json:"suburb"`
	Municipality string  `json:"municipality"`
}

func NewGeoapifyClient(apiKey string) *GeoapifyClient {
	return &GeoapifyClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type LatLongCoords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (g *GeoapifyClient) BatchReverseGeocode(db edgedb.Executor, locations []LatLongCoords) (*GeoapifyResponse, error) {

	if len(locations) > maxBatchSize {
		return nil, fmt.Errorf("Geoapify batch request exceeds max allowed size (%d/%d)", len(locations), maxBatchSize)
	}

	baseURL := "https://api.geoapify.com/v1/batch/geocode/reverse"

	// Prepare JSON body
	jsonBody, err := json.Marshal(locations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("apiKey", g.apiKey)
	u.RawQuery = q.Encode()

	// Update request to POST with JSON body
	// req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonBody))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create request: %w", err)
	// }
	// req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := g.client.Post(
		u.String(),
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	_, err = TrackGeoapifyUsage(db, int32(len(locations)))
	if err != nil {
		return nil, fmt.Errorf("failed to track Geoapify usage: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned non-202 status: %d\n %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var pending GeoapifyPendingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pending); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return g.AwaitResult(pending)
}
