package gbif

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lsdch/biome/models"
	"golang.org/x/time/rate"
)

func newTestClient(serverURL string) *GBIFClient {
	return &GBIFClient{
		BaseURL:    serverURL,
		HTTPClient: http.DefaultClient,

		UserAgent: "BiOME-test (mailto: louis.duchemin@univ-lyon1.fr)",

		MaxRetries: 2,
		BaseDelay:  1 * time.Millisecond,
		MaxDelay:   10 * time.Millisecond,

		Limiter: rate.NewLimiter(rate.Inf, 1),

		concurrency: make(chan struct{}, 10),

		BackboneDatasetKey: "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c",
	}
}

func TestSearchParamsEncode(t *testing.T) {
	extinct := true

	p := SearchParams{
		Query:          "Homo sapiens",
		DatasetKey:     "dataset-key",
		HigherTaxonKey: 123,
		Rank:           "SPECIES",
		Status:         "ACCEPTED",
		IsExtinct:      &extinct,
		Limit:          100,
		Offset:         200,
	}

	got := p.Encode()

	expected := url.Values{
		"q":              {"Homo sapiens"},
		"datasetKey":     {"dataset-key"},
		"higherTaxonKey": {"123"},
		"rank":           {"SPECIES"},
		"status":         {"ACCEPTED"},
		"isExtinct":      {"true"},
		"limit":          {"100"},
		"offset":         {"200"},
	}

	if got.Encode() != expected.Encode() {
		t.Fatalf("unexpected query parameters:\ngot:  %s\nwant: %s",
			got.Encode(),
			expected.Encode(),
		)
	}
}

func TestSearchParamsEncodeOmitsZeroValues(t *testing.T) {
	p := SearchParams{}

	got := p.Encode()

	if len(got) != 0 {
		t.Fatalf("expected empty query, got %v", got)
	}
}

func TestSearchParamsApply(t *testing.T) {
	extinct := true

	var p SearchParams

	p.Apply(
		WithHigherTaxonKey(123),
		WithRank("SPECIES"),
		WithStatus("ACCEPTED"),
		WithIsExtinct(&extinct),
		WithLimit(50),
		WithOffset(100),
	)

	if p.HigherTaxonKey != 123 {
		t.Errorf("HigherTaxonKey = %d, want 123", p.HigherTaxonKey)
	}

	if p.Rank != "SPECIES" {
		t.Errorf("Rank = %q, want SPECIES", p.Rank)
	}

	if p.Status != "ACCEPTED" {
		t.Errorf("Status = %q, want ACCEPTED", p.Status)
	}

	if p.IsExtinct == nil || *p.IsExtinct != true {
		t.Errorf("IsExtinct = %v, want true", p.IsExtinct)
	}

	if p.Limit != 50 {
		t.Errorf("Limit = %d, want 50", p.Limit)
	}

	if p.Offset != 100 {
		t.Errorf("Offset = %d, want 100", p.Offset)
	}
}

func TestSearchResponseGetExactMatch(t *testing.T) {
	resp := SearchResponse{
		Results: []models.TaxonGBIF{
			{
				ScientificName: "Homo sapiens",
				Name:           "Homo sapiens",
			},
			{
				ScientificName: "Pan troglodytes",
				Name:           "Pan troglodytes",
			},
		},
	}

	t.Run("scientific name", func(t *testing.T) {
		got, ok := resp.GetExactMatch("Homo sapiens")

		if !ok {
			t.Fatal("expected exact match")
		}

		if got.ScientificName != "Homo sapiens" {
			t.Errorf("ScientificName = %q", got.ScientificName)
		}
	})

	t.Run("name", func(t *testing.T) {
		resp := SearchResponse{
			Results: []models.TaxonGBIF{
				{
					ScientificName: "Foo",
					Name:           "Bar",
				},
			},
		}

		got, ok := resp.GetExactMatch("Bar")

		if !ok {
			t.Fatal("expected exact match")
		}

		if got.Name != "Bar" {
			t.Errorf("Name = %q", got.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		got, ok := resp.GetExactMatch("Unknown species")

		if ok {
			t.Fatal("expected no match")
		}

		if got != nil {
			t.Fatalf("expected nil result, got %+v", got)
		}
	})

	t.Run("multiple matches", func(t *testing.T) {
		resp := SearchResponse{
			Results: []models.TaxonGBIF{
				{
					ScientificName: "Homo sapiens",
					Name:           "Homo sapiens",
				},
				{
					ScientificName: "Homo sapiens",
					Name:           "Homo sapiens",
				},
			},
		}

		got, ok := resp.GetExactMatch("Homo sapiens")

		if ok {
			t.Fatal("expected ambiguous match to fail")
		}

		if got != nil {
			t.Fatalf("expected nil result, got %+v", got)
		}
	})
}

func TestSearchResponseGetFirstResult(t *testing.T) {
	t.Run("returns first result", func(t *testing.T) {
		resp := SearchResponse{
			Results: []models.TaxonGBIF{
				{ScientificName: "First"},
				{ScientificName: "Second"},
			},
		}

		got, err := resp.GetFirstResult()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.ScientificName != "First" {
			t.Errorf("ScientificName = %q, want First", got.ScientificName)
		}
	})

	t.Run("empty results", func(t *testing.T) {
		resp := SearchResponse{}

		got, err := resp.GetFirstResult()

		if err == nil {
			t.Fatal("expected error")
		}

		if got != nil {
			t.Fatalf("expected nil result, got %+v", got)
		}
	})
}

func TestGBIFClientGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}

			if r.URL.Path != "/test" {
				t.Errorf("path = %s, want /test", r.URL.Path)
			}

			if r.Header.Get("User-Agent") != "BiOME-test (mailto: louis.duchemin@univ-lyon1.fr)" {
				t.Errorf("unexpected User-Agent: %q", r.Header.Get("User-Agent"))
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"status": "ok",
			})
		}))
		defer server.Close()

		client := newTestClient(server.URL)

		var result map[string]string

		err := client.get(
			context.Background(),
			"/test",
			url.Values{"foo": {"bar"}},
			&result,
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result["status"] != "ok" {
			t.Errorf("status = %q, want ok", result["status"])
		}
	})

	t.Run("HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		client := newTestClient(server.URL)

		var result map[string]string

		err := client.get(
			context.Background(),
			"/test",
			nil,
			&result,
		)

		if err == nil {
			t.Fatal("expected error")
		}

		if !contains(err.Error(), "500") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`not json`))
		}))
		defer server.Close()

		client := newTestClient(server.URL)

		var result map[string]string

		err := client.get(
			context.Background(),
			"/test",
			nil,
			&result,
		)

		if err == nil {
			t.Fatal("expected JSON decoding error")
		}
	})
}

func TestGBIFClientGetRetriesOnRateLimit(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := requests.Add(1)

		if n < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.MaxRetries = 2

	var result map[string]string

	err := client.get(
		context.Background(),
		"/test",
		nil,
		&result,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if requests.Load() != 3 {
		t.Errorf("requests = %d, want 3", requests.Load())
	}

	if result["status"] != "ok" {
		t.Errorf("status = %q, want ok", result["status"])
	}
}

func TestGBIFClientGetRateLimitExceeded(t *testing.T) {
	var requests atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.MaxRetries = 2

	var result map[string]string

	err := client.get(
		context.Background(),
		"/test",
		nil,
		&result,
	)

	if err == nil {
		t.Fatal("expected rate limit error")
	}

	if !contains(err.Error(), "rate limit exceeded") {
		t.Errorf("unexpected error: %v", err)
	}

	if requests.Load() != 3 {
		t.Errorf("requests = %d, want 3", requests.Load())
	}
}

func TestGBIFClientGetContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	client.BaseDelay = time.Second

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result map[string]string

	err := client.get(
		ctx,
		"/test",
		nil,
		&result,
	)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}
}

func TestGBIFClientAcquireContextCancellation(t *testing.T) {
	client := newTestClient("https://api.gbif.org")

	// Fill the semaphore.
	client.concurrency <- struct{}{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.acquire(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want context.Canceled", err)
	}

	client.release()
}

func TestGBIFClientSearchSpecies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/species/search" {
			t.Errorf("path = %s, want /v1/species/search", r.URL.Path)
		}

		if r.URL.Query().Get("q") != "Homo sapiens" {
			t.Errorf("q = %q, want Homo sapiens", r.URL.Query().Get("q"))
		}

		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("limit = %q, want 10", r.URL.Query().Get("limit"))
		}

		response := SearchResponse{
			Offset:       0,
			Limit:        10,
			EndOfRecords: true,
			Count:        1,
			Results: []models.TaxonGBIF{
				{
					ScientificName: "Homo sapiens",
					Name:           "Homo sapiens",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	result, err := client.SearchSpecies(
		context.Background(),
		SearchParams{
			Query: "Homo sapiens",
			Limit: 10,
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 1 {
		t.Errorf("Count = %d, want 1", result.Count)
	}

	if len(result.Results) != 1 {
		t.Fatalf("len(Results) = %d, want 1", len(result.Results))
	}

	if result.Results[0].ScientificName != "Homo sapiens" {
		t.Errorf(
			"ScientificName = %q, want Homo sapiens",
			result.Results[0].ScientificName,
		)
	}
}

func TestGBIFClientGetTaxonByKey(t *testing.T) {
	const usageKey int32 = 2436436

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/v1/species/" + strconv.FormatInt(int64(usageKey), 10)

		if r.URL.Path != expectedPath {
			t.Errorf("path = %s, want %s", r.URL.Path, expectedPath)
		}

		taxon := models.TaxonGBIF{
			ScientificName: "Homo sapiens",
			Name:           "Homo sapiens",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(taxon)
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	result, err := client.GetTaxonByKey(
		context.Background(),
		usageKey,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ScientificName != "Homo sapiens" {
		t.Errorf(
			"ScientificName = %q, want Homo sapiens",
			result.ScientificName,
		)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && stringContains(s, substr)
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
