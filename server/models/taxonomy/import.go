package taxonomy

import (
	"context"
	"darco/proto/models"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

//go:embed upsert_clade.edgeql
var upsertCladeCmd string

//go:embed upsert_clades.edgeql
var upsertCladesCmd string

type cladeResult struct {
	ID edgedb.UUID `edgedb:"id"`
}

func upsertClades(taxa []TaxonGBIF) (n int, err error) {
	taxa = funk.Map(taxa, func(taxon TaxonGBIF) TaxonGBIF {
		taxon.normalize()
		return taxon
	}).([]TaxonGBIF)

	ctx := context.Background()
	err = models.DB.Tx(ctx, func(ctx context.Context, tx *edgedb.Tx) (err error) {
		for _, taxon := range taxa {
			args, _ := jsonDB.Marshal(&taxon)
			log.Printf("INSERTING %s", args)
			if err = tx.Execute(ctx, upsertCladeCmd, map[string]interface{}{"data": args}); err != nil {
				log.Printf("ERROR : %s", err)
				return err
			} else {
				n++
			}
		}
		return
	})

	// data, _ := jsonDB.Marshal(taxa)
	// log.Printf("%s", data)
	// err = models.DB.Query(context.Background(), upsertCladesCmd, &result, map[string]any{"data": data})
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

func importChildren(GBIF_ID int, process *ImportProcess, monitor func(*ImportProcess)) {
	var taxa []TaxonGBIF
	endReached := false
	offset := 0
	for !endReached {
		children, err := fetchChildren(GBIF_ID, offset)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return
		}
		taxa = append(taxa, children.Results...)
		endReached = children.EndOfRecords
		offset += REQUEST_LIMIT
	}

	taxa = funk.Filter(taxa, func(taxon TaxonGBIF) bool {
		return taxon.Rank != "UNRANKED"
	}).([]TaxonGBIF)

	if len(taxa) > 0 {
		inserted, _ := upsertClades(taxa)
		process.Imported += inserted
		monitor(process)
	}

	for _, taxon := range taxa {
		if taxon.NumDescendants > 0 {
			importChildren(taxon.Key, process, monitor)
		}
	}
}

func ImportTaxon(GBIF_ID int, monitor func(p *ImportProcess)) (result []cladeResult, err error) {

	taxonURL, _ := url.JoinPath(baseURL, fmt.Sprint(GBIF_ID))
	body, err := makeRequest(taxonURL, 0)
	if err != nil {
		log.Printf("Request ERR : %s", err)
	}
	var taxon TaxonGBIF
	if err = json.Unmarshal(body, &taxon); err != nil {
		log.Printf("ERROR : %s", err)
	}

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

	parents, err := fetchParents(GBIF_ID)
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	parentCount, err := upsertClades(parents)
	if err == nil {
		process.Imported += parentCount
		monitor(&process)
	} else {
		log.Printf("%s", err)
	}

	taxon.Anchor = true
	_, err = upsertClades([]TaxonGBIF{taxon})
	if err == nil {
		process.Imported++
		monitor(&process)
	} else {
		log.Printf("ANCHOR TAXON ERROR : %s", err)
	}

	// importChildren(GBIF_ID, &process, monitor)

	process.Done = true
	monitor(&process)

	return result, err

	// for i := 1; i < 10; i++ {
	// 	log.Printf("LOOP")
	// 	time.Sleep(time.Second)
	// 	go monitor(1)
	// }
	// return result, nil
}
