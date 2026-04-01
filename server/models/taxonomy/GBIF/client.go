package gbif

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/lsdch/biome/models/settings"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	GBIF_BASE_URL           = "https://api.gbif.org/"
	MAX_CONCURRENT_REQUESTS = 20
)

type GBIFClient struct {
	BaseURL    string
	HTTPClient *http.Client

	UserAgent string

	// retry behaviour
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration

	// client-side limiter
	Limiter *rate.Limiter

	// concurrent request limiter
	concurrency chan struct{}
}

var client *GBIFClient

func Client() *GBIFClient {
	if client == nil {
		client = NewClient()
	}
	return client
}

func NewClient() *GBIFClient {
	superadmin := settings.Get().SuperAdmin

	httpClient := http.DefaultClient

	userAgent := fmt.Sprintf("biome-client (%s)", superadmin.Email)

	return &GBIFClient{
		BaseURL:    GBIF_BASE_URL,
		HTTPClient: httpClient,
		UserAgent:  userAgent,

		MaxRetries: 5,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   10 * time.Second,

		Limiter: rate.NewLimiter(rate.Limit(20), 50),

		concurrency: make(chan struct{}, MAX_CONCURRENT_REQUESTS),
	}
}

func (c *GBIFClient) acquire(ctx context.Context) error {

	select {
	case c.concurrency <- struct{}{}:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *GBIFClient) release() {
	<-c.concurrency
}

func (c *GBIFClient) get(
	ctx context.Context,
	path string,
	query url.Values,
	dest interface{},
) error {

	if err := c.acquire(ctx); err != nil {
		return err
	}
	defer c.release()

	u, err := url.JoinPath(c.BaseURL, path)
	if err != nil {
		return err
	}
	if query != nil {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	// client-side rate limiting
	if c.Limiter != nil {
		if err := c.Limiter.Wait(ctx); err != nil {
			return err
		}
	}

	delay := c.BaseDelay

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return err
		}

		// retry on rate limit
		if resp.StatusCode == http.StatusTooManyRequests {

			resp.Body.Close()

			if attempt == c.MaxRetries {
				return errors.New("gbif rate limit exceeded after retries")
			}

			jitter := time.Duration(rand.Int63n(int64(delay / 2)))
			sleep := delay + jitter

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sleep):
			}

			delay *= 2
			if delay > c.MaxDelay {
				delay = c.MaxDelay
			}

			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return fmt.Errorf("gbif api error: %s", resp.Status)
		}

		return json.NewDecoder(resp.Body).Decode(dest)
	}

	return errors.New("unreachable")
}

type SearchParams struct {
	Query          string
	DatasetKey     string
	HigherTaxonKey int32
	Rank           string
	Status         string
	IsExtinct      *bool
	Limit          int
	Offset         int
}

type searchParamOption func(*SearchParams)

func newSearchParamOption[T any](value T, setter func(*SearchParams, T)) searchParamOption {
	return func(p *SearchParams) {
		setter(p, value)
	}
}

var (
	WithHigherTaxonKey = func(v int32) searchParamOption {
		return newSearchParamOption(v, func(p *SearchParams, v int32) { p.HigherTaxonKey = v })
	}
	WithRank = func(v string) searchParamOption {
		return newSearchParamOption(v, func(p *SearchParams, v string) { p.Rank = v })
	}
	WithStatus = func(v string) searchParamOption {
		return newSearchParamOption(v, func(p *SearchParams, v string) { p.Status = v })
	}
	WithIsExtinct = func(v *bool) searchParamOption {
		return newSearchParamOption(v, func(p *SearchParams, v *bool) { p.IsExtinct = v })
	}
	WithLimit = func(v int) searchParamOption {
		return newSearchParamOption(v, func(p *SearchParams, v int) { p.Limit = v })
	}
	WithOffset = func(v int) searchParamOption {
		return newSearchParamOption(v, func(p *SearchParams, v int) { p.Offset = v })
	}
)

// Apply applies options to the SearchParams
func (p *SearchParams) Apply(opts ...searchParamOption) {
	for _, opt := range opts {
		opt(p)
	}
}

func (p SearchParams) Encode() url.Values {
	q := url.Values{}

	if p.Query != "" {
		q.Set("q", p.Query)
	}

	if p.DatasetKey != "" {
		q.Set("datasetKey", p.DatasetKey)
	}

	if p.HigherTaxonKey != 0 {
		q.Set("higherTaxonKey", fmt.Sprintf("%d", p.HigherTaxonKey))
	}

	if p.Rank != "" {
		q.Set("rank", p.Rank)
	}

	if p.Status != "" {
		q.Set("status", p.Status)
	}

	if p.IsExtinct != nil {
		q.Set("isExtinct", fmt.Sprintf("%t", *p.IsExtinct))
	}

	if p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}

	if p.Offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", p.Offset))
	}

	return q
}

type SearchResponse struct {
	Offset       int         `json:"offset"`
	Limit        int         `json:"limit"`
	EndOfRecords bool        `json:"endOfRecords"`
	Count        int         `json:"count"`
	Results      []TaxonGBIF `json:"results"`
}

func (resp SearchResponse) GetExactMatch(name string) (*TaxonGBIF, bool) {
	matches := []TaxonGBIF{}
	for _, taxon := range resp.Results {
		if taxon.ScientificName == name || taxon.Name == name {
			matches = append(matches, taxon)
		}
	}

	if len(matches) == 0 {
		logrus.Warnf("exact match not found for name: %s", name)
		return nil, false
	} else if len(matches) > 1 {
		logrus.Warnf("multiple exact matches found for name: %s", name)
		return nil, false
	}
	return &matches[0], true
}

func (resp SearchResponse) GetFirstResult() (*TaxonGBIF, error) {
	if len(resp.Results) == 0 {
		return nil, errors.New("no results found")
	}
	return &resp.Results[0], nil
}

func (c *GBIFClient) SearchSpecies(ctx context.Context, p SearchParams) (*SearchResponse, error) {

	q := p.Encode()

	var result SearchResponse

	err := c.get(ctx, "/v1/species/search", q, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *GBIFClient) GetTaxonByKey(ctx context.Context, usageKey int32) (*TaxonGBIF, error) {

	path := "/v1/species/" + strconv.FormatInt(int64(usageKey), 10)

	var result TaxonGBIF

	err := c.get(ctx, path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
