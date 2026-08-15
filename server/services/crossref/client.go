package crossref

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lsdch/biome/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

const (
	defaultAPI                      = "https://api.crossref.org"
	DEFAULT_REQUESTS_PER_SECOND     = 3
	DEFAULT_MAX_CONCURRENT_REQUESTS = 10
)

type Work struct {
	Message WorkMessage `json:"message"`
}

type WorkBatchResult struct {
	DOI  types.DOI
	Work *Work
	Err  error
}

type WorkMessage struct {
	DOI       string   `json:"DOI"`
	Title     []string `json:"title"`
	Author    []Author `json:"author"`
	Published struct {
		DateParts [][]int `json:"date-parts"`
	} `json:"published"`
	ContainerTitle []string `json:"container-title"`
	Unstructured   string   `json:"unstructured,omitempty"`
}

func (m *WorkMessage) Year() int32 {
	if len(m.Published.DateParts) == 0 || len(m.Published.DateParts[0]) == 0 {
		return 0
	}

	return int32(m.Published.DateParts[0][0])
}

func (m *WorkMessage) Journal() string {
	if len(m.ContainerTitle) == 0 {
		return ""
	}

	return m.ContainerTitle[0]
}

func (m *WorkMessage) Authors() []string {
	authors := make([]string, 0, len(m.Author))

	for _, author := range m.Author {
		authors = append(authors, author.FullName())
	}

	return authors
}

func (m *WorkMessage) TitleString() string {
	if len(m.Title) == 0 {
		return ""
	}

	return m.Title[0]
}

func (m *WorkMessage) Verbatim() string {
	if m.Unstructured != "" {
		return m.Unstructured
	}
	authorsStr := ""
	for i, author := range m.Author {
		if i > 0 {
			authorsStr += ", "
		}
		authorsStr += author.Family + " " + author.FirstNameInitials()
	}
	return fmt.Sprintf("%s (%d) %s. %s", authorsStr, m.Year(), m.TitleString(), m.Journal())
}

type Author struct {
	Given  string `json:"given"`
	Family string `json:"family"`
}

func (a Author) FullName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", a.Given, a.Family))
}

func (a Author) FirstNameInitials() string {
	parts := strings.Fields(a.Given)
	initials := make([]string, 0, len(parts))

	for _, part := range parts {
		if part != "" {
			initials = append(initials, string(part[0]))
		}
	}

	return strings.Join(initials, "")
}

type WorksQueryResponse struct {
	Message struct {
		TotalResults int `json:"total-results"`

		Items []WorkMessage `json:"items"`
	} `json:"message"`
}

type Client struct {
	httpClient *http.Client

	baseURL string
	appName string
	mailTo  string

	// Limits concurrent HTTP requests.
	semaphore chan struct{}

	mu sync.Mutex

	limiter *rate.Limiter

	// Current CrossRef advertised limits.
	rateLimit int
	interval  time.Duration

	config Config
}

type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	MailTo  string `mapstructure:"MAIL_TO"`

	// Initial fallback when CrossRef has not advertised limits yet.
	RequestsPerSecond float64 `mapstructure:"REQUESTS_PER_SECOND"`

	// Maximum number of HTTP requests running concurrently.
	MaxConcurrentRequests int `mapstructure:"MAX_CONCURRENT_REQUESTS"`

	// Number of workers to use for batch requests. If 0, defaults to MaxConcurrentRequests.
	BatchWorkers int `mapstructure:"BATCH_WORKERS"`

	HTTPClient *http.Client
}

func NewClient(cfg Config) *Client {
	if cfg.RequestsPerSecond == 0 {
		cfg.RequestsPerSecond = DEFAULT_REQUESTS_PER_SECOND
	}

	if cfg.MaxConcurrentRequests <= 0 {
		cfg.MaxConcurrentRequests = DEFAULT_MAX_CONCURRENT_REQUESTS
	}

	if cfg.BatchWorkers <= 0 {
		cfg.BatchWorkers = cfg.MaxConcurrentRequests
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &Client{
		httpClient: httpClient,

		baseURL: defaultAPI,
		appName: cfg.AppName,
		mailTo:  cfg.MailTo,

		rateLimit: 1,
		interval:  time.Second,

		limiter: rate.NewLimiter(
			rate.Limit(cfg.RequestsPerSecond),
			1,
		),
		semaphore: make(chan struct{}, cfg.MaxConcurrentRequests),
		config:    cfg,
	}
}

func (c *Client) Works(
	ctx context.Context,
	doi types.DOI,
) (*Work, error) {

	var result Work

	err := c.getJSON(
		ctx,
		path.Join("works", doi.String()),
		nil,
		&result,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"crossref works %s: %w",
			doi,
			err,
		)
	}

	return &result, nil
}

func (c *Client) QueryWorks(
	ctx context.Context,
	bibliographic string,
	rows int,
) (*WorksQueryResponse, error) {

	query := url.Values{}
	query.Set(
		"query.bibliographic",
		bibliographic,
	)

	if rows > 0 {
		query.Set(
			"rows",
			strconv.Itoa(rows),
		)
	}

	var result WorksQueryResponse

	err := c.getJSON(
		ctx,
		"works",
		&query,
		&result,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"crossref query %q: %w",
			bibliographic,
			err,
		)
	}

	return &result, nil
}

func (c *Client) getJSON(
	ctx context.Context,
	endpoint string,
	query *url.Values,
	dst any,
) error {

	// Limit concurrent HTTP requests.
	if err := c.acquireSemaphore(ctx); err != nil {
		return err
	}
	defer c.releaseSemaphore()

	// Respect CrossRef rate limit.
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}

	u.Path = endpoint

	values := u.Query()
	values.Set("mailto", c.mailTo)

	if query != nil {
		for key, items := range *query {
			for _, item := range items {
				values.Add(key, item)
			}
		}
	}

	u.RawQuery = values.Encode()

	for attempts := 0; attempts < 3; attempts++ {

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			u.String(),
			nil,
		)

		if err != nil {
			return err
		}

		req.Header.Set(
			"User-Agent",
			fmt.Sprintf(
				"%s/1.0 (https://github.com/lsdch/biome; mailto:%s)",
				c.appName,
				c.mailTo,
			),
		)

		logrus.Debugf("Crossref request %s %s", req.Method, req.URL.String())

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{
			"status":   resp.StatusCode,
			"limit":    resp.Header.Get("X-Rate-Limit-Limit"),
			"interval": resp.Header.Get("X-Rate-Limit-Interval"),
			"retry":    resp.Header.Get("Retry-After"),
		}).Debug("CrossRef response")

		c.updateLimits(resp.Header)

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			if err := c.waitRetryAfter(ctx, resp); err != nil {
				return err
			}
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)

			return fmt.Errorf(
				"crossref HTTP %d: %s",
				resp.StatusCode,
				string(body),
			)
		}
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(dst)
	}

	return fmt.Errorf("crossref: too many retries")
}

func (c *Client) updateLimits(
	headers http.Header,
) {
	limit := headers.Get("X-Rate-Limit-Limit")
	interval := headers.Get("X-Rate-Limit-Interval")

	if limit == "" || interval == "" {
		return
	}

	rateLimit, err := strconv.Atoi(limit)
	if err != nil || rateLimit <= 0 {
		return
	}

	duration, err := parseInterval(interval)
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.rateLimit = rateLimit
	c.interval = duration

	c.limiter.SetLimit(
		rate.Limit(
			float64(rateLimit) / duration.Seconds(),
		),
	)

	c.limiter.SetBurst(min(rateLimit, c.config.MaxConcurrentRequests))
}

func parseInterval(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)

	value = strings.TrimSuffix(value, "s")

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return time.Duration(seconds) * time.Second, nil
}

func (c *Client) waitRetryAfter(
	ctx context.Context,
	resp *http.Response,
) error {

	retry := strings.TrimSpace(resp.Header.Get("Retry-After"))
	var delay time.Duration
	if retry == "" {
		delay := time.Second

		logrus.WithField("delay", delay).
			Warn("CrossRef rate limited without Retry-After header")
	} else if seconds, err := strconv.Atoi(retry); err == nil {
		if seconds <= 0 {
			return nil
		}

		delay = time.Duration(seconds) * time.Second

	} else {
		// Retry-After: <http-date>
		t, err := http.ParseTime(retry)
		if err != nil {
			return fmt.Errorf("invalid Retry-After header %q", retry)
		}

		delay = time.Until(t)
		if delay <= 0 {
			return nil
		}
	}

	logrus.WithField("delay", delay).
		Warn("CrossRef rate limited, waiting before retry")

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) acquireSemaphore(ctx context.Context) error {
	select {
	case c.semaphore <- struct{}{}:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) releaseSemaphore() {
	<-c.semaphore
}

type RateLimitInfo struct {
	Limit    int
	Interval time.Duration
}

func (c *Client) CurrentRateLimit() RateLimitInfo {
	c.mu.Lock()
	defer c.mu.Unlock()

	return RateLimitInfo{
		Limit:    c.rateLimit,
		Interval: c.interval,
	}
}

func (c *Client) WorksBatch(
	ctx context.Context,
	dois []types.DOI,
) []BatchResult[types.DOI, *Work] {

	logrus.Debugf("Fetching %d works from CrossRef", len(dois))

	return batch(
		ctx,
		c.config.BatchWorkers,
		dois,
		c.Works,
	)
}

func (c *Client) QueryWorksBatch(
	ctx context.Context,
	queries []string,
	rows int,
) []BatchResult[string, *WorksQueryResponse] {

	logrus.Debugf("Fetching %d bibliographic queries from CrossRef", len(queries))

	return batch(
		ctx,
		c.config.BatchWorkers,
		queries,
		func(ctx context.Context, query string) (*WorksQueryResponse, error) {
			return c.QueryWorks(ctx, query, rows)
		},
	)
}

type BatchResult[TIn any, TOut any] struct {
	Input TIn
	Value TOut
	Err   error
}

func batch[TIn any, TOut any](
	ctx context.Context,
	workers int,
	input []TIn,
	fn func(context.Context, TIn) (TOut, error),
) []BatchResult[TIn, TOut] {

	if len(input) == 0 {
		return nil
	}

	if workers <= 0 {
		workers = 1
	}

	type job struct {
		index int
		item  TIn
	}

	jobs := make(chan job)
	results := make([]BatchResult[TIn, TOut], len(input))

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for job := range jobs {
			result := BatchResult[TIn, TOut]{
				Input: job.item,
			}

			select {
			case <-ctx.Done():
				result.Err = ctx.Err()
			default:
				result.Value, result.Err = fn(ctx, job.item)
			}

			results[job.index] = result
		}
	}

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go worker()
	}

	for i, item := range input {
		select {
		case jobs <- job{
			index: i,
			item:  item,
		}:
		case <-ctx.Done():
			results[i] = BatchResult[TIn, TOut]{
				Input: item,
				Err:   ctx.Err(),
			}
		}
	}

	close(jobs)
	wg.Wait()

	return results
}
