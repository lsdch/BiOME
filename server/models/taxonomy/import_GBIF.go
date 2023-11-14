package taxonomy

import (
	"context"
	"darco/proto/models"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/edgedb/edgedb-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/thoas/go-funk"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var BASE_URL = "https://api.gbif.org/v1/species/"
var PAGE_SIZE = 1000

type TaxonGBIF struct {
	Key            int                `json:"key" edgedb:"GBIF_ID"`
	Parent         int                `json:"parentKey" edgedb:"parentID"`
	Name           string             `json:"canonicalName" edgedb:"name"`
	Authorship     edgedb.OptionalStr `json:"authorship,omitempty" edgedb:"authorship,omitempty"`
	Status         string             `json:"taxonomicStatus" edgedb:"status"`
	Rank           string             `json:"rank" edgedb:"rank"`
	NumDescendants int                `json:"numDescendants" edgedb:"-"`
	Anchor         bool               `edgedb:"anchor"`
}

func (taxon *TaxonGBIF) normalize() {
	if authorship, isSet := taxon.Authorship.Get(); isSet && authorship == "" {
		taxon.Authorship.Unset()
	}
	if taxon.Status == "ACCEPTED" {
		taxon.Status = "Accepted"
	} else {
		taxon.Status = "Synonym"
	}
	rank := strings.ToLower(taxon.Rank)
	taxon.Rank = cases.Title(language.English, cases.NoLower).String(rank)
}

// ImportProcess represents the progress and status of the taxon import process from GBIF.
type ImportProcess struct {
	Name     string    `json:"name"`
	GBIF_ID  int       `json:"GBIF_ID"`
	Expected int       `json:"expected"`
	Imported int       `json:"imported"`
	Rank     string    `json:"rank"`
	Started  time.Time `json:"started"`
	Done     bool      `json:"done"`
	Error    error     `json:"error"`
}

// ProgressTracker tracks and reports the progress of the import process.
type ProgressTracker struct {
	process ImportProcess
	monitor func(p *ImportProcess)
}

func (p *ProgressTracker) Report() {
	p.monitor(&p.process)
}

func (p *ProgressTracker) Progress(n int) {
	p.process.Imported += n
	p.Report()
}

func (p *ProgressTracker) Errorf(format string, a ...any) error {
	p.process.Error = fmt.Errorf(format, a...)
	p.Report()
	return p.process.Error
}

func (p *ProgressTracker) Terminate() {
	p.process.Done = true
	p.Report()
}

func NewProgressTracker(taxon *TaxonGBIF, f func(p *ImportProcess)) *ProgressTracker {
	process := ImportProcess{
		Name:     taxon.Name,
		GBIF_ID:  taxon.Key,
		Expected: taxon.NumDescendants + 1,
		Imported: 0,
		Rank:     taxon.Rank,
		Started:  time.Now(),
		Done:     false,
		Error:    nil,
	}
	tracker := ProgressTracker{
		process: process,
		monitor: f,
	}
	tracker.Report()
	return &tracker
}

func makeRequest(strURL string, offset int) (body []byte, err error) {
	URL, err := url.ParseRequestURI(strURL)
	params := url.Values{}
	params.Set("limit", fmt.Sprint(PAGE_SIZE))
	params.Set("offset", fmt.Sprint(offset))
	URL.RawQuery = params.Encode()
	strURL = fmt.Sprint(URL)
	response, err := http.Get(strURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("request failure: %s", response.Status)
	}
	return io.ReadAll(response.Body)
}

func requestTaxa(URL string) (taxa []TaxonGBIF, err error) {
	body, err := makeRequest(URL, 0)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &taxa); err != nil {
		return nil, err
	} else {
		return taxa, nil
	}
}

func fetchParents(GBIF_ID int) ([]TaxonGBIF, error) {
	URL, err := url.JoinPath(BASE_URL, fmt.Sprintf("%d/parents", GBIF_ID))
	if err != nil {
		return nil, err
	}
	return requestTaxa(URL)
}

var jsonDB = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	TagKey:                 "edgedb",
}.Froze()

//go:embed upsert_anchor.edgeql
var upsertTaxonCmd string

func upsertTaxa(tx *edgedb.Tx, taxa []TaxonGBIF) (n int, err error) {
	taxa = funk.Map(taxa, func(taxon TaxonGBIF) TaxonGBIF {
		taxon.normalize()
		return taxon
	}).([]TaxonGBIF)

	ctx := context.Background()
	err = models.DB().Tx(ctx, func(ctx context.Context, tx *edgedb.Tx) (err error) {
		for _, taxon := range taxa {
			log.Debugf("Inserting taxon from GBIF %+v", &taxon)
			args, _ := jsonDB.Marshal(&taxon)
			if err = tx.Execute(ctx, upsertTaxonCmd, args); err != nil {
				return
			} else {
				n++
			}
		}
		return
	})
	return
}

type ChildrenGBIF struct {
	Offset       int         `json:"offset"`
	Limit        int         `json:"limit"`
	EndOfRecords bool        `json:"endOfRecords"`
	Results      []TaxonGBIF `json:"results"`
}

func fetchChildren(GBIF_ID int, offset int) (children ChildrenGBIF, err error) {
	URL, err := url.JoinPath(BASE_URL, fmt.Sprintf("%d/children", GBIF_ID))
	if err != nil {
		return
	}
	body, err := makeRequest(URL, offset)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &children); err != nil {
		return
	}
	return
}

func importChildren(tx *edgedb.Tx, GBIF_ID int, tracker *ProgressTracker) error {
	var taxa []TaxonGBIF
	endReached := false
	offset := 0
	for !endReached {
		children, err := fetchChildren(GBIF_ID, offset)
		if err != nil {
			return tracker.Errorf("failed to fetch data from GBIF\n%w", err)
		}
		taxa = append(taxa, children.Results...)
		endReached = children.EndOfRecords
		offset += PAGE_SIZE
	}

	taxa = funk.Filter(taxa, func(taxon TaxonGBIF) bool {
		return taxon.Rank != "UNRANKED" && taxon.Rank != "VARIETY" && taxon.Status != "DOUBTFUL"
	}).([]TaxonGBIF)

	if len(taxa) > 0 {
		inserted, err := upsertTaxa(tx, taxa)
		if err != nil {
			return tracker.Errorf("failed to insert taxon imported from GBIF\n%w", err)
		}
		tracker.Progress(inserted)
	}

	for _, taxon := range taxa {
		if taxon.NumDescendants > 0 {
			if err := importChildren(tx, taxon.Key, tracker); err != nil {
				return err
			}
		}
	}
	return nil
}

func fetchTaxon(GBIF_ID int) (taxon TaxonGBIF, err error) {
	taxonURL, _ := url.JoinPath(BASE_URL, fmt.Sprint(GBIF_ID))
	body, err := makeRequest(taxonURL, 0)
	if err != nil {
		err = fmt.Errorf("failed to fetch GBIF record of anchor taxon #%d \n %w", GBIF_ID, err)
		return
	}
	if err = json.Unmarshal(body, &taxon); err != nil {
		err = fmt.Errorf("failed to parse JSON response from GBIF \n %w", err)
		log.WithFields(log.Fields{"body": body}).Errorf("%s", err)
		return
	}
	taxon.Anchor = true
	return
}

func ImportTaxon(GBIF_ID int, monitor func(p *ImportProcess)) (err error) {

	taxon, err := fetchTaxon(GBIF_ID)
	if err != nil {
		return
	}
	log.Infof("Started import of taxon : %s [GBIF %d]", taxon.Name, taxon.Key)

	tracker := NewProgressTracker(&taxon, monitor)

	go models.DB().Tx(context.Background(),
		func(ctx context.Context, tx *edgedb.Tx) error {
			parents, err := fetchParents(GBIF_ID)
			if err != nil {
				return tracker.Errorf("failed to fetch parent taxa of %s[%d] from GBIF\n%w", taxon.Name, taxon.Key, err)
			}

			insert_count, err := upsertTaxa(tx, append(parents, taxon))
			if err != nil {
				return tracker.Errorf("failed to insert a parent of taxon %s[%d] \n%w",
					taxon.Name, taxon.Key, err)
			}
			tracker.Progress(insert_count)

			importChildren(tx, GBIF_ID, tracker)
			tracker.Terminate()
			return nil
		})

	return err
}
