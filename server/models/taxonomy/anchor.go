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

var REQUEST_LIMIT = 1000

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

type ImportProcess struct {
	Name     string    `json:"name"`
	GBIF_ID  int       `json:"GBIF_ID"`
	Expected int       `json:"expected"`
	Imported int       `json:"imported"`
	Rank     string    `json:"rank"`
	Started  time.Time `json:"started"`
	Done     bool      `json:"done"`
}

var jsonDB = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	TagKey:                 "edgedb",
}.Froze()

var baseURL = "https://api.gbif.org/v1/species/"

func makeRequest(strURL string, offset int) (body []byte, err error) {
	URL, err := url.ParseRequestURI(strURL)
	params := url.Values{}
	params.Set("limit", fmt.Sprint(REQUEST_LIMIT))
	params.Set("offset", fmt.Sprint(offset))
	URL.RawQuery = params.Encode()
	strURL = fmt.Sprint(URL)
	response, err := http.Get(strURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
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
	URL, err := url.JoinPath(baseURL, fmt.Sprintf("%d/parents", GBIF_ID))
	if err != nil {
		return nil, err
	}
	return requestTaxa(URL)
}

//go:embed upsert_taxon.edgeql
var upsertTaxonCmd string

func upsertTaxa(tx *edgedb.Tx, taxa []TaxonGBIF) (n int, err error) {
	taxa = funk.Map(taxa, func(taxon TaxonGBIF) TaxonGBIF {
		taxon.normalize()
		return taxon
	}).([]TaxonGBIF)

	ctx := context.Background()
	err = models.DB.Tx(ctx, func(ctx context.Context, tx *edgedb.Tx) (err error) {
		for _, taxon := range taxa {
			log.Debugf("INSERTING %v", &taxon)
			args, _ := jsonDB.Marshal(&taxon)
			if err = tx.Execute(ctx, upsertTaxonCmd, args); err != nil {
				return err
			} else {
				n++
			}
		}
		return
	})
	return n, err
}

type ChildrenGBIF struct {
	Offset       int         `json:"offset"`
	Limit        int         `json:"limit"`
	EndOfRecords bool        `json:"endOfRecords"`
	Results      []TaxonGBIF `json:"results"`
}

func fetchChildren(GBIF_ID int, offset int) (children ChildrenGBIF, err error) {
	URL, err := url.JoinPath(baseURL, fmt.Sprintf("%d/children", GBIF_ID))
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

func importChildren(tx *edgedb.Tx, GBIF_ID int, process *ImportProcess, monitor func(*ImportProcess)) {
	var taxa []TaxonGBIF
	endReached := false
	offset := 0
	for !endReached {
		children, err := fetchChildren(GBIF_ID, offset)
		if err != nil {
			log.Errorf("Failed to fetch data from GBIF \n%s", err)
			return
		}
		taxa = append(taxa, children.Results...)
		endReached = children.EndOfRecords
		offset += REQUEST_LIMIT
	}

	taxa = funk.Filter(taxa, func(taxon TaxonGBIF) bool {
		return taxon.Rank != "UNRANKED" && taxon.Status != "DOUBTFUL"
	}).([]TaxonGBIF)

	if len(taxa) > 0 {
		inserted, _ := upsertTaxa(tx, taxa)
		process.Imported += inserted
		monitor(process)
	}

	for _, taxon := range taxa {
		if taxon.NumDescendants > 0 {
			importChildren(tx, taxon.Key, process, monitor)
		}
	}
}

func ImportTaxon(GBIF_ID int, monitor func(p *ImportProcess)) (err error) {

	taxonURL, _ := url.JoinPath(baseURL, fmt.Sprint(GBIF_ID))
	body, err := makeRequest(taxonURL, 0)
	if err != nil {
		log.Errorf("Failed to fetch GBIF record of anchor taxon #%d \n %s", GBIF_ID, err)
		return
	}
	var taxon TaxonGBIF
	if err = json.Unmarshal(body, &taxon); err != nil {
		log.WithFields(
			log.Fields{"body": body},
		).Errorf("Failed to parse JSON response from GBIF \n %s", err)
		return
	}

	log.Infof("Started import of taxon : %s [GBIF %d]", taxon.Name, taxon.Key)

	process := ImportProcess{
		Name:     taxon.Name,
		GBIF_ID:  taxon.Key,
		Expected: taxon.NumDescendants + 1,
		Imported: 0,
		Rank:     taxon.Rank,
		Started:  time.Now(),
		Done:     false,
	}
	monitor(&process)

	ctx := context.Background()
	err = models.DB.Tx(ctx, func(ctx context.Context, tx *edgedb.Tx) (err error) {

		parents, err := fetchParents(GBIF_ID)
		if err != nil {
			log.Errorf("Failed to fetch parent taxa of %s[%d] from GBIF\n %s",
				taxon.Name, taxon.Key, err)
			return err
		}

		parentCount, err := upsertTaxa(tx, parents)
		if err == nil {
			process.Imported += parentCount
			monitor(&process)
		} else {
			log.Errorf("Failed to insert some parent taxa of %s[%d] \n %s",
				taxon.Name, taxon.Key, err)
		}

		taxon.Anchor = true
		_, err = upsertTaxa(tx, []TaxonGBIF{taxon})
		if err == nil {
			process.Imported++
			monitor(&process)
		} else {
			log.Errorf("Failed to insert anchor taxon %s[%d] \n %s",
				taxon.Name, taxon.Key, err)
		}

		importChildren(tx, GBIF_ID, &process, monitor)

		process.Done = true
		monitor(&process)
		return
	})

	return err
}
