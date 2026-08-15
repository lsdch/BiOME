//go:build integration

package crossref_test

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lsdch/biome/services/crossref"
	"github.com/lsdch/biome/types"
)

func loadTestDOIs(t *testing.T) []types.DOI {
	t.Helper()

	data, err := os.ReadFile("data/test_DOIs.txt")
	if err != nil {
		t.Fatalf("read DOI test file: %v", err)
	}

	var dois []types.DOI

	for _, line := range strings.Split(string(data), "\n") {
		doi := strings.TrimSpace(line)

		if doi == "" || strings.HasPrefix(doi, "#") {
			continue
		}

		dois = append(
			dois,
			types.ParseDOI(doi),
		)
	}

	if len(dois) == 0 {
		t.Fatal("no DOI configured")
	}

	return dois
}

func newIntegrationClient(t *testing.T) *crossref.Client {
	t.Helper()

	mailTo := os.Getenv("CROSSREF_MAILTO")
	appName := os.Getenv("CROSSREF_APP_NAME")

	if mailTo == "" {
		t.Skip("CROSSREF_MAILTO is required")
	}

	if appName == "" {
		t.Skip("CROSSREF_APP_NAME is required")
	}

	return crossref.NewClient(
		crossref.Config{
			AppName: appName,
			MailTo:  mailTo,

			// Conservative startup value.
			// It will be replaced by CrossRef headers.
			RequestsPerSecond: 1,

			MaxConcurrentRequests: 3,
		},
	)
}

func TestWorks_Integration(t *testing.T) {
	client := newIntegrationClient(t)

	dois := loadTestDOIs(t)

	work, err := client.Works(
		context.Background(),
		dois[0],
	)

	if err != nil {
		t.Fatal(err)
	}

	if work.Message.DOI == "" {
		t.Fatal("missing DOI in response")
	}
}

func TestWorks_UpdatesRateLimit(t *testing.T) {
	client := newIntegrationClient(t)

	dois := loadTestDOIs(t)

	before := client.CurrentRateLimit()

	_, err := client.Works(
		context.Background(),
		dois[0],
	)

	if err != nil {
		t.Fatal(err)
	}

	after := client.CurrentRateLimit()

	if before == after {
		t.Fatalf(
			"rate limit was not updated: before=%+v after=%+v",
			before,
			after,
		)
	}

	if after.Limit <= 0 {
		t.Fatalf(
			"invalid updated rate limit: %+v",
			after,
		)
	}
}

func TestWorks_ConcurrentRequests(t *testing.T) {
	client := newIntegrationClient(t)

	dois := loadTestDOIs(t)

	if len(dois) < 3 {
		t.Skip(
			"need at least 3 DOIs",
		)
	}

	ctx := context.Background()

	start := time.Now()

	var wg sync.WaitGroup

	errs := make(chan error, len(dois[:3]))

	for _, doi := range dois[:3] {
		wg.Add(1)

		go func(doi types.DOI) {
			defer wg.Done()

			_, err := client.Works(ctx, doi)

			if err != nil {
				errs <- err
			}

		}(doi)
	}

	wg.Wait()

	close(errs)

	for err := range errs {
		t.Fatal(err)
	}

	elapsed := time.Since(start)

	t.Logf(
		"completed 3 concurrent DOI requests in %s",
		elapsed,
	)

	if elapsed <= 0 {
		t.Fatal("invalid elapsed time")
	}
}

func TestWorks_RespectsRateLimit(t *testing.T) {
	client := newIntegrationClient(t)

	dois := loadTestDOIs(t)

	if len(dois) < 3 {
		t.Skip("need at least 3 DOIs")
	}

	ctx := context.Background()

	start := time.Now()

	for _, doi := range dois[:3] {
		_, err := client.Works(ctx, doi)

		if err != nil {
			t.Fatal(err)
		}
	}

	elapsed := time.Since(start)

	limit := client.CurrentRateLimit()

	minimum := time.Duration(
		float64(limit.Interval) /
			float64(limit.Limit) *
			2,
	)

	t.Logf(
		"elapsed=%s expected-minimum=%s rate=%+v",
		elapsed,
		minimum,
		limit,
	)

	if elapsed < minimum {
		t.Fatalf(
			"rate limiting not respected: elapsed=%s expected >=%s",
			elapsed,
			minimum,
		)
	}
}
