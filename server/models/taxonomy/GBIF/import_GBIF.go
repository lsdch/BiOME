package gbif

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/geldata/gel-go"
	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/taxonomy"

	"github.com/sirupsen/logrus"

	"github.com/thoas/go-funk"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var BASE_URL = "https://api.gbif.org/v1/species/"
var PAGE_SIZE = 1000

type TaxonInnerGBIF struct {
	Key    int    `json:"key" gel:"GBIF_ID"`
	Parent int    `json:"parentKey" gel:"parentID"`
	Name   string `json:"canonicalName" gel:"name"`
	Status string `json:"taxonomicStatus" gel:"status"`
	Rank   string `json:"rank" gel:"rank"`
}

type TaxonGBIF struct {
	TaxonInnerGBIF `json:",inline" gel:"$inline"`
	Authorship     geltypes.OptionalStr `json:"authorship,omitempty" gel:"authorship,omitempty"`
	NumDescendants int                  `json:"numDescendants" gel:"-"`
	Anchor         bool                 `json:"anchor" gel:"anchor"`
}

func (taxon *TaxonGBIF) normalize() {
	if authorship, isSet := taxon.Authorship.Get(); isSet && authorship == "" {
		taxon.Authorship.Unset()
	}
	if taxon.Status == "ACCEPTED" {
		taxon.Status = "Accepted"
	} else {
		taxon.Status = "Unclassified"
	}
	rank := strings.ToLower(taxon.Rank)
	taxon.Rank = cases.Title(language.English, cases.NoLower).String(rank)
}

// ImportProcess represents the progress and status of the taxon import process from GBIF.
type ImportProcess struct {
	Name     string             `json:"name"`
	GBIF_ID  int                `json:"GBIF_ID"`
	Expected int                `json:"expected"`
	Imported int                `json:"imported"`
	Rank     taxonomy.TaxonRank `json:"rank"`
	Started  time.Time          `json:"started"`
	Done     bool               `json:"done"`
	Error    error              `json:"error"`
}

func makeRequest(strURL string, offset int) (body []byte, err error) {
	URL, err := url.ParseRequestURI(strURL)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Set("limit", fmt.Sprint(PAGE_SIZE))
	params.Set("offset", fmt.Sprint(offset))
	URL.RawQuery = params.Encode()
	strURL = fmt.Sprint(URL)
	response, err := http.Get(strURL)
	if err != nil {
		return
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

func upsertTaxa(tx geltypes.Tx, taxa []TaxonGBIF) (n int, err error) {
	taxa = funk.Map(taxa, func(taxon TaxonGBIF) TaxonGBIF {
		taxon.normalize()
		return taxon
	}).([]TaxonGBIF)

	ctx := context.Background()

	for _, taxon := range taxa {
		logrus.Debugf("Inserting taxon from GBIF %+v", &taxon)
		args, _ := json.Marshal(&taxon)
		if err = tx.Execute(ctx,
			`#edgeql
				with module taxonomy,
					data := <json>$0,
					anchor := <bool>data['anchor'],
				insert Taxon {
					name := <str>data['canonicalName'],
					GBIF_ID := <int32>data['key'],
					status := <TaxonStatus>data['taxonomicStatus'],
					parent := (
						select detached Taxon filter .GBIF_ID = <int32>data['parentKey']
					),
					rank := <Rank>data['rank'],
					authorship := <str>data['authorship'],
					anchor := anchor
				}
				unless conflict on .GBIF_ID else (
					update Taxon set {
						anchor := anchor if not .anchor else .anchor
					}
				)
			`, args); err != nil {
			return
		} else {
			n++
		}
	}
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

func importChildren(tx geltypes.Tx, GBIF_ID int, tracker *ProgressTracker) error {
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
		return taxon.Rank != "UNRANKED" && taxon.Rank != "VARIETY" && taxon.Status == "ACCEPTED"
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

// Retrieves GBIF data for a single taxon
func fetchTaxon(GBIF_ID int) (taxon TaxonGBIF, err error) {
	taxonURL, _ := url.JoinPath(BASE_URL, fmt.Sprint(GBIF_ID))
	body, err := makeRequest(taxonURL, 0)
	if err != nil {
		err = fmt.Errorf("failed to fetch GBIF record of anchor taxon #%d \n %w", GBIF_ID, err)
		return
	}
	if err = json.Unmarshal(body, &taxon); err != nil {
		err = fmt.Errorf("failed to parse JSON response from GBIF \n %w", err)
		logrus.WithFields(logrus.Fields{"body": body}).Errorf("%s", err)
		return
	}
	taxon.Anchor = true
	return
}

type ImportRequestGBIF struct {
	Key      int  `json:"key" doc:"Target GBIF taxon key"`
	Children bool `json:"children" doc:"Import whole clade, including the taxon descendants"`
}

func ImportTaxonTx(tx geltypes.Tx, request ImportRequestGBIF, monitor func(p *ImportProcess)) (err error) {
	taxon, err := fetchTaxon(request.Key)
	if err != nil {
		return
	}

	tracker := NewProgressTracker(&taxon, monitor)

	parents, err := fetchParents(request.Key)
	if err != nil {
		return tracker.Errorf("failed to fetch parent taxa of %s[%d] from GBIF\n%w", taxon.Name, taxon.Key, err)
	}

	insert_count, err := upsertTaxa(tx, append(parents, taxon))
	if err != nil {
		return tracker.Errorf("failed to insert a parent of taxon %s[%d] \n%w",
			taxon.Name, taxon.Key, err)
	}
	tracker.Progress(insert_count)

	if request.Children {
		err := importChildren(tx, request.Key, tracker)
		if err != nil {
			return tracker.Errorf("failed to import children of taxon %s[%d]\b%w", taxon.Name, taxon.Key, err)
		}
	}

	tracker.Terminate()
	return nil
}

func ImportTaxon(db *gel.Client, request ImportRequestGBIF, monitor func(p *ImportProcess)) (err error) {
	return db.Tx(context.Background(), func(ctx context.Context, tx geltypes.Tx) error {
		return ImportTaxonTx(tx, request, monitor)
	})
}
