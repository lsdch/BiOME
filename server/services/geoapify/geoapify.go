package geoapify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/lsdch/biome/config"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models"
	"github.com/sirupsen/logrus"
)

const (
	MAX_BATCH_SIZE            = 1000
	REVERSE_GEOCODE_URL       = "https://api.geoapify.com/v1/geocode/reverse"
	BATCH_REVERSE_GEOCODE_URL = "https://api.geoapify.com/v1/batch/geocode/reverse"
)

type GeoapifyService struct {
	cfg    config.GeoapifyConfig
	client *http.Client
}

func NewGeoapifyService(client *http.Client, cfg config.GeoapifyConfig) *GeoapifyService {
	return &GeoapifyService{cfg: cfg, client: client}
}

func (s *GeoapifyService) ListUsage(ctx context.Context, q db.Querier) (usage []models.GeoapifyUsage, err error) {
	res, err := q.Queries().ListGeoapifyUsage(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range res {
		usage = append(usage, models.GeoapifyUsageFromDB(u))
	}

	return usage, nil
}

func (s *GeoapifyService) GetTodayUsage(ctx context.Context, q db.Querier) (int32, error) {
	settings, err := q.Queries().GetTodayGeoapifyUsage(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return settings.RequestsCount, nil
}

func (s *GeoapifyService) IncrementUsage(ctx context.Context, q db.Querier, requestsCount int32) (total int32, err error) {
	usage, err := q.Queries().IncrementTodayGeoapifyUsage(ctx, requestsCount)
	if err != nil {
		return 0, err
	}
	return usage.RequestsCount, nil
}

func (s *GeoapifyService) GetStatus(ctx context.Context, q db.Querier) (models.GeoapifyStatus, error) {
	hasKey := s.cfg.GeoApifyApiKey != ""
	todayUsage, err := s.GetTodayUsage(ctx, q)
	if err != nil {
		return models.GeoapifyStatus{}, err
	}
	limitExceeded := todayUsage >= s.cfg.DailyUsageLimit
	return models.GeoapifyStatus{
		APIKey:        s.cfg.GeoApifyApiKey,
		Available:     !limitExceeded && hasKey,
		HasApiKey:     hasKey,
		Limit:         s.cfg.DailyUsageLimit,
		TodayRequests: todayUsage,
	}, nil
}

func (s *GeoapifyService) ReverseGeocode(ctx context.Context, q db.Querier, lat, lon float64) (*models.GeoapifyResult, error) {
	status, err := s.GetStatus(ctx, q)
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

	query := u.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("apiKey", status.APIKey)
	query.Set("format", "json")
	u.RawQuery = query.Encode()

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

	_, err = s.IncrementUsage(ctx, q, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to track Geoapify usage: %w", err)
	}

	// Parse the response
	var result models.ReverseGeoCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Geoapify always returns an array of results, even for single queries
	// so we just take the first one
	// Coordinates in the middle of the ocean also return a result,
	// albeit with some empty fields
	return &result.Results[0], nil
}

func (s *GeoapifyService) BatchReverseGeocode(ctx context.Context, q db.Querier, locations []models.GeoapifyCoords) ([]models.GeoapifyResult, error) {

	if len(locations) > MAX_BATCH_SIZE {
		return nil, fmt.Errorf("Geoapify batch request exceeds max allowed size (%d/%d)", len(locations), MAX_BATCH_SIZE)
	}

	status, err := s.GetStatus(ctx, q)
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

	query := u.Query()
	query.Set("apiKey", status.APIKey)
	u.RawQuery = query.Encode()

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

	_, err = s.IncrementUsage(ctx, q, int32(len(locations)))
	if err != nil {
		return nil, fmt.Errorf("failed to track Geoapify usage: %w", err)
	}

	// Parse the response
	var pending models.GeoapifyPendingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pending); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.AwaitResult(pending)
}

func (s *GeoapifyService) AwaitResult(p models.GeoapifyPendingResponse) ([]models.GeoapifyResult, error) {
	time.Sleep(5 * time.Second)
	for {
		resp, err := s.client.Get(p.URL)
		if err != nil {
			return nil, err
		}

		var result []models.GeoapifyResult
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
