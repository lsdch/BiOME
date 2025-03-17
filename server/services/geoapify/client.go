package geoapify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/settings"
	"github.com/lsdch/biome/services"

	"github.com/sirupsen/logrus"
)

const maxBatchSize = 1000

type GeoapifyClient struct {
	apiKey string
	client *http.Client
}

type GeoapifyClientScheduler struct {
	GeoapifyClient
	BatchRequests    services.Queue[[]GeoapifyResult, GeoapifyClient]
	ActiveQueries    int
	MaxActiveQueries int
}

type GeoapifyPendingResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

func (c GeoapifyClient) AwaitResult(p GeoapifyPendingResponse) ([]GeoapifyResult, error) {
	time.Sleep(5 * time.Second)
	for {
		resp, err := c.client.Get(p.URL)
		if err != nil {
			return nil, err
		}

		var result []GeoapifyResult
		body, err := io.ReadAll(resp.Body)
		switch resp.StatusCode {
		case 200:
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(body, &result)
			return result, err
		case 202:
			logrus.Infof("Response pending: %+v", string(body))
			time.Sleep(30 * time.Second)
		}
		resp.Body.Close()
	}
}

type GeoapifyResult struct {
	Formatted    string  `json:"formatted"`
	Municipality string  `json:"municipality"`
	City         string  `json:"city"`
	County       string  `json:"county"`
	State        string  `json:"state"`
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	PostalCode   string  `json:"postcode"`
	Street       string  `json:"street"`
	HouseNumber  string  `json:"housenumber"`
	Suburb       string  `json:"suburb"`
}

type ReverseGeoCodeResponse struct {
	Results []GeoapifyResult       `json:"results"`
	Query   map[string]interface{} `json:"query"`
}

type clientOption func(*GeoapifyClient) error

func WithApiKey(apiKey string) clientOption {
	return func(c *GeoapifyClient) error {
		c.apiKey = apiKey
		return nil
	}
}

func NewClient(opts ...clientOption) (*GeoapifyClient, error) {
	client := &GeoapifyClient{
		client: &http.Client{},
	}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, fmt.Errorf("Failed to initialize Geoapify client: %w", err)
		}
	}
	if client.apiKey == "" {
		key, ok := settings.Get().GeoapifyApiKey.Get()
		if !ok {
			return nil, fmt.Errorf("Geoapify API key not set")
		}
		client.apiKey = key
	}
	return client, nil
}

type LatLongCoords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (g *GeoapifyClient) ReverseGeocode(db geltypes.Executor, coords LatLongCoords) (*GeoapifyResult, error) {
	todayUsage, err := TodayGeoapifyUsage(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's Geoapify usage: %w", err)
	}
	if todayUsage.Requests >= CREDIT_LIMIT {
		return nil, fmt.Errorf("Geoapify usage limit exceeded")
	}

	baseURL := "https://api.geoapify.com/v1/geocode/reverse"

	// Create URL with query parameters
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", coords.Lat))
	q.Set("lon", fmt.Sprintf("%f", coords.Lon))
	q.Set("apiKey", g.apiKey)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	logrus.Debugf("Geoapify reverse geocode URL: %s", u.String())

	// Make the request
	resp, err := g.client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, huma.NewError(
			resp.StatusCode,
			"Geoapify API returned non-200 status",
			fmt.Errorf("%s", string(body)),
		)
	}

	_, err = TrackGeoapifyUsage(db, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to track Geoapify usage: %w", err)
	}

	// Parse the response
	var result ReverseGeoCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Geoapify always returns an array of results, even for single queries
	// so we just take the first one
	// Coordinates in the middle of the ocean also return a result,
	// albeit with some empty fields
	return &result.Results[0], nil
}

func (g *GeoapifyClient) BatchReverseGeocode(db geltypes.Executor, locations []LatLongCoords) ([]GeoapifyResult, error) {

	if len(locations) > maxBatchSize {
		return nil, fmt.Errorf("Geoapify batch request exceeds max allowed size (%d/%d)", len(locations), maxBatchSize)
	}

	todayUsage, err := TodayGeoapifyUsage(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's Geoapify usage: %w", err)
	}
	if int(todayUsage.Requests)+len(locations) > CREDIT_LIMIT {
		return nil, fmt.Errorf("Geoapify usage limit exceeded")
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

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return nil, huma.NewError(
			resp.StatusCode,
			"Geoapify API returned non-200 status",
			fmt.Errorf("%s", string(body)),
		)
	}

	_, err = TrackGeoapifyUsage(db, int32(len(locations)))
	if err != nil {
		return nil, fmt.Errorf("failed to track Geoapify usage: %w", err)
	}

	// Parse the response
	var pending GeoapifyPendingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pending); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return g.AwaitResult(pending)
}

type Status struct {
	Available     bool `json:"available"`
	HasApiKey     bool `json:"has_api_key"`
	TodayRequests int  `json:"requests"`
	Limit         int  `json:"limit"`
}

func GetStatus(db geltypes.Executor) (*Status, error) {
	_, hasKey := settings.Get().GeoapifyApiKey.Get()
	todayUsage, err := TodayGeoapifyUsage(db)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's Geoapify usage: %w", err)
	}
	limitExceeded := todayUsage.Requests >= CREDIT_LIMIT
	return &Status{
		Available:     !limitExceeded && hasKey,
		HasApiKey:     hasKey,
		Limit:         CREDIT_LIMIT,
		TodayRequests: int(todayUsage.Requests),
	}, nil
}
