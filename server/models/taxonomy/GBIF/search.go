package gbif

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/geldata/gel-go/geltypes"
)

const GBIF_BACKBONE_DATASET_KEY = "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c"

type SearchResultTaxon struct {
	TaxonGBIF            `json:",inline"`
	HigherClassification map[int]string `json:"higherClassificationMap"`
}

type SearchResults struct {
	Offset       int                 `json:"offset"`
	EndOfRecords bool                `json:"endOfRecords"`
	Limit        int                 `json:"limit"`
	Count        int                 `json:"count"`
	Results      []SearchResultTaxon `json:"results"`
}

type SearchQuery struct {
	Name       string
	DatasetKey string
	Status     string
	Rank       string
}

func (q SearchQuery) Encode() string {
	v := url.Values{}
	if strings.Trim(q.Name, " ") != "" {
		v.Add("q", q.Name)
	}
	if strings.Trim(q.DatasetKey, " ") != "" {
		v.Add("datasetKey", q.DatasetKey)
	}
	if strings.Trim(q.Status, " ") != "" {
		v.Add("status", q.Status)
	}
	if strings.Trim(q.Rank, " ") != "" {
		v.Add("rank", q.Rank)
	}
	return v.Encode()
}

func (q *SearchQuery) SetRank(rank string) *SearchQuery {
	q.Rank = rank
	return q
}

func (q *SearchQuery) SetStatus(status string) *SearchQuery {
	q.Status = status
	return q
}

func NewSearchQuery(name string) SearchQuery {
	return SearchQuery{
		Name:       name,
		DatasetKey: GBIF_BACKBONE_DATASET_KEY,
		Status:     "ACCEPTED",
	}
}

func Search(query SearchQuery) (results SearchResults, err error) {

	var URL = url.URL{
		Host:     "api.gbif.org",
		Path:     "v1/species/search",
		Scheme:   "https",
		RawQuery: query.Encode(),
	}

	response, err := http.Get(URL.String())
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("request failure: %s", response.Status)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &results)
	return
}

func ImportTaxonByName(e geltypes.Executor, name string) (taxon TaxonGBIF, err error) {
	sq := NewSearchQuery(name)
	results, err := Search(sq)
	if err != nil {
		return
	}
	if results.Count == 0 {
		err = fmt.Errorf("no taxon found for name: %s", name)
		return
	}
	best := results.Results[0]

	// Import parents
	parents, err := fetchParents(best.Key)
	if err != nil {
		return
	}
	insert_count, err := upsertTaxa(e, append(parents, taxon))
	if err != nil {
		return
	}
	return
}

func TaxonExists(e geltypes.Executor, gbifID int) (exists bool, err error) {
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with module taxonomy,
			select exists (Taxon filter .GBIF_ID = <int32>$0)
		`, &exists, gbifID,
	)
	return
}

func CheckLineageExists(e geltypes.Executor, names []string) (missing []string, err error) {
	missing = []string{}
	err = e.QuerySingle(context.Background(),
		`#edgeql
			with module taxonomy,
				lineage_names := <str>array_unpack($0),
			select lineage_names except Taxon.name
		`, &missing, names,
	)
	return
}
