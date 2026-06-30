package geoapify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/sirupsen/logrus"
)

const (
	MAX_BATCH_SIZE            = 1000
	REVERSE_GEOCODE_URL       = "https://api.geoapify.com/v1/geocode/reverse"
	BATCH_REVERSE_GEOCODE_URL = "https://api.geoapify.com/v1/batch/geocode/reverse"
)

type GeoapifyUsage biomedb.GeoapifyUsage

type GeoapifyStatus struct {
	APIKey        string `json:"-"`
	Available     bool   `json:"available"`
	HasApiKey     bool   `json:"has_api_key"`
	TodayRequests int32  `json:"requests"`
	Limit         int32  `json:"limit"`
}

type GeoapifyRequestDenial error

var (
	ErrNoAPIKey      = GeoapifyRequestDenial(fmt.Errorf("Geoapify API key is not set"))
	ErrLimitExceeded = GeoapifyRequestDenial(fmt.Errorf("Geoapify usage limit exceeded"))
)

func (s *GeoapifyStatus) AllowRequests(n int32) GeoapifyRequestDenial {
	if !s.HasApiKey {
		return ErrNoAPIKey
	}
	if s.TodayRequests+n > s.Limit {
		return ErrLimitExceeded
	}
	return nil
}

type GeoapifyPendingResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
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

type GeoapifyService struct {
	db     *db.DB
	cfg    config.GeoapifyConfig
	client *http.Client
}

func NewGeoapifyService(db *db.DB, client *http.Client, cfg config.GeoapifyConfig) *GeoapifyService {
	return &GeoapifyService{db: db, cfg: cfg, client: client}
}

func (s *GeoapifyService) ListUsage(ctx context.Context) (usage []GeoapifyUsage, err error) {
	res, err := s.db.Queries().ListGeoapifyUsage(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range res {
		usage = append(usage, GeoapifyUsage(u))
	}

	return usage, nil
}

func (s *GeoapifyService) GetTodayUsage(ctx context.Context) (int32, error) {
	settings, err := s.db.Queries().GetTodayGeoapifyUsage(ctx)
	if err != nil {
		return 0, err
	}
	return settings.RequestsCount, nil
}

func (s *GeoapifyService) IncrementUsage(ctx context.Context, requestsCount int32) (total int32, err error) {
	usage, err := s.db.Queries().IncrementTodayGeoapifyUsage(ctx, requestsCount)
	if err != nil {
		return 0, err
	}
	return usage.RequestsCount, nil
}

func (s *GeoapifyService) GetStatus(ctx context.Context) (GeoapifyStatus, error) {
	hasKey := s.cfg.GeoApifyApiKey != ""
	todayUsage, err := s.GetTodayUsage(ctx)
	if err != nil {
		return GeoapifyStatus{}, err
	}
	limitExceeded := todayUsage >= s.cfg.DailyUsageLimit
	return GeoapifyStatus{
		APIKey:        s.cfg.GeoApifyApiKey,
		Available:     !limitExceeded && hasKey,
		HasApiKey:     hasKey,
		Limit:         s.cfg.DailyUsageLimit,
		TodayRequests: todayUsage,
	}, nil
}

func (s *GeoapifyService) ReverseGeocode(ctx context.Context, lat float32, lon float32) (*GeoapifyResult, error) {
	status, err := s.GetStatus(ctx)
	if err != nil {
		return nil, err
	}
	if err := status.AllowRequests(1); err != nil {
		return nil, err
	}

	// Create URL with query parameters
	u, err := url.Parse(REVERSE_GEOCODE_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lon))
	q.Set("apiKey", status.APIKey)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	logrus.Debugf("Geoapify reverse geocode URL: %s", u.String())

	// Make the request
	resp, err := s.client.Get(u.String())
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

	_, err = s.IncrementUsage(ctx, 1)
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

type LatLonCoords struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (s *GeoapifyService) BatchReverseGeocode(ctx context.Context, locations []LatLonCoords) ([]GeoapifyResult, error) {

	if len(locations) > MAX_BATCH_SIZE {
		return nil, fmt.Errorf("Geoapify batch request exceeds max allowed size (%d/%d)", len(locations), MAX_BATCH_SIZE)
	}

	status, err := s.GetStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Geoapify status: %w", err)
	}

	denial := status.AllowRequests(int32(len(locations)))
	if denial != nil {
		return nil, denial
	}

	// Prepare JSON body
	jsonBody, err := json.Marshal(locations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create URL with query parameters
	u, err := url.Parse(BATCH_REVERSE_GEOCODE_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("apiKey", status.APIKey)
	u.RawQuery = q.Encode()

	// Make the request
	resp, err := s.client.Post(
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

	_, err = s.IncrementUsage(ctx, int32(len(locations)))
	if err != nil {
		return nil, fmt.Errorf("failed to track Geoapify usage: %w", err)
	}

	// Parse the response
	var pending GeoapifyPendingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pending); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.AwaitResult(pending)
}

func (s *GeoapifyService) AwaitResult(p GeoapifyPendingResponse) ([]GeoapifyResult, error) {
	time.Sleep(5 * time.Second)
	for {
		resp, err := s.client.Get(p.URL)
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
