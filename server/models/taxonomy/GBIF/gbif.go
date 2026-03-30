package gbif

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models/taxonomy"
	"github.com/sirupsen/logrus"
)

const GBIF_BACKBONE_DATASET_KEY = "d7dddbf4-2cf0-4f39-9b2a-bb099caae36c"

type TaxaByRank map[taxonomy.TaxonRank]map[int32]TaxonGBIF

func newTaxaByRank() TaxaByRank {
	return TaxaByRank{
		taxonomy.Kingdom:    make(map[int32]TaxonGBIF),
		taxonomy.Phylum:     make(map[int32]TaxonGBIF),
		taxonomy.Class:      make(map[int32]TaxonGBIF),
		taxonomy.Order:      make(map[int32]TaxonGBIF),
		taxonomy.Family:     make(map[int32]TaxonGBIF),
		taxonomy.Genus:      make(map[int32]TaxonGBIF),
		taxonomy.Species:    make(map[int32]TaxonGBIF),
		taxonomy.Subspecies: make(map[int32]TaxonGBIF),
	}
}

func (t TaxaByRank) Add(taxon TaxonGBIF) {
	t[taxon.GetRank()][taxon.Key] = taxon
}

func (t TaxaByRank) Contains(key int32, ranks ...taxonomy.TaxonRank) bool {
	if len(ranks) == 0 {
		ranks = taxonomy.TaxonRankValues
	}
	for _, rank := range ranks {
		if _, ok := t[rank][key]; ok {
			return true
		}
	}
	return false
}

func (t TaxaByRank) Count() int {
	count := 0
	for _, rank := range taxonomy.TaxonRankValues {
		count += len(t[rank])
	}
	return count
}

func (t TaxaByRank) FetchParentsAndSynonyms(tx geltypes.Executor) error {
	type rankResult struct {
		taxa []TaxonGBIF
		err  error
	}

	resultsChan := make(chan rankResult, len(taxonomy.TaxonRankValues))
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	// Spawn goroutines to fetch parents and synonyms per rank
	for _, rank := range taxonomy.TaxonRankValues {
		wg.Add(1)
		go func(rank taxonomy.TaxonRank) {
			defer wg.Done()
			fetched := make([]TaxonGBIF, 0)

			// Process each taxon for this rank
			for _, taxon := range t[rank] {
				// Fetch accepted taxon if this is a synonym
				if taxon.AcceptedKey != 0 {
					logrus.Infof("Taxon '%s' with key %d is a synonym of taxon with key %d, fetching accepted taxon from GBIF", taxon.Name, taxon.Key, taxon.AcceptedKey)
					if !t.Contains(taxon.AcceptedKey) {
						synonym, err := FetchKey(taxon.AcceptedKey)
						if err != nil {
							resultsChan <- rankResult{err: fmt.Errorf("Error fetching accepted taxon with key %d for synonym '%s': %v", taxon.AcceptedKey, taxon.Name, err)}
							return
						}
						fetched = append(fetched, *synonym)
					}
				}

				// Fetch lineage taxa
				lineageTaxa, err := t.fetchLineageTaxa(tx, taxon)
				if err != nil {
					resultsChan <- rankResult{err: err}
					return
				}
				fetched = append(fetched, lineageTaxa...)
			}

			resultsChan <- rankResult{taxa: fetched}
		}(rank)
	}

	// Wait for all goroutines and close channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect all results and add them sequentially (no concurrent map access)
	for result := range resultsChan {
		if result.err != nil {
			select {
			case errChan <- result.err:
			default:
			}
		}
		for _, taxon := range result.taxa {
			t.Add(taxon)
		}
	}

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// fetchLineageTaxa returns taxa from a taxon's lineage without modifying the map
func (t TaxaByRank) fetchLineageTaxa(tx geltypes.Executor, taxon TaxonGBIF) ([]TaxonGBIF, error) {
	toCheck := make(map[string]int32)
	for key, name := range taxon.HigherClassificationMap {
		if t.Contains(key) || key == taxon.Key {
			continue
		}
		toCheck[name] = key
	}
	if len(toCheck) == 0 {
		return nil, nil
	}

	missing, err := taxonomy.CheckMissingTaxa(tx, slices.Collect(maps.Keys(toCheck)))
	if err != nil {
		return nil, err
	}
	logrus.Debugf(
		"%d taxa for lineage of taxon '%s' already exist in the database",
		len(toCheck), taxon.Name,
	)
	if len(missing) == 0 {
		return nil, nil
	}

	logrus.Debugf("Fetching %d taxa for lineage of taxon '%s' from GBIF", len(missing), taxon.Name)
	taxa := make([]TaxonGBIF, 0, len(missing))
	for _, name := range missing {
		taxon, err := FetchKey(toCheck[name])
		if err != nil {
			return nil, err
		}
		if taxon.Status != "ACCEPTED" {
			logrus.Warnf("Taxon '%s' with key %d in lineage of taxon '%s' is not accepted (status: %s), fetching by name", taxon.Name, taxon.Key, taxon.Name, taxon.Status)
			taxon, err = FetchByName(Client(), SearchParams{Query: name, Status: "ACCEPTED"})
			if err != nil {
				return nil, err
			}
		}
		taxa = append(taxa, *taxon)
	}

	return taxa, nil
}

func (t TaxaByRank) Persist(tx geltypes.Executor) error {
	for _, rank := range taxonomy.TaxonRankValues {
		if _, err := UpsertTaxa(tx, slices.Collect(maps.Values(t[rank]))); err != nil {
			return fmt.Errorf("Error persisting taxa of rank %s: %v", rank, err)
		}
	}
	return nil
}

type NamesSearchResult struct {
	Taxa     TaxaByRank
	NotFound []string
}

func newMultipleNamesResult() NamesSearchResult {

	return NamesSearchResult{
		Taxa: newTaxaByRank(),
	}
}

func (r *NamesSearchResult) Add(taxon TaxonGBIF) {
	if strings.Contains(taxon.Status, "ACCEPTED") || strings.Contains(taxon.Status, "SYNONYM") {
		r.Taxa.Add(taxon)
	} else {
		r.NotFound = append(r.NotFound, taxon.Name)
	}
}

func FetchByName(client *GBIFClient, args SearchParams) (*TaxonGBIF, error) {
	logrus.Debugf("FetchByName received params: %+v", args)
	args.Query = strings.TrimSpace(args.Query)
	args.Limit = 10
	args.DatasetKey = GBIF_BACKBONE_DATASET_KEY
	args.Status = "ACCEPTED"

	switch len(strings.Split(args.Query, " ")) {
	case 1:
		args.Rank = "GENUS"
	case 2:
		args.Rank = "SPECIES"
	default:
		// Chance of an exact match is very low for names with more than 2 parts
	}

	logrus.Debugf("Searching GBIF for taxon name '%s'", args.Query)
	resp, err := client.SearchSpecies(context.Background(), args)
	if err != nil {
		return nil, err
	}

	taxon, found := resp.GetExactMatch(args.Query)
	if found {
		return taxon, nil
	}

	logrus.Debugf("No exact match found for name '%s', searching without rank filter", args.Query)
	args.Rank = ""
	resp, err = client.SearchSpecies(context.Background(), args)
	if err != nil {
		return nil, err
	}

	taxon, found = resp.GetExactMatch(args.Query)
	if found {
		return taxon, nil
	}

	logrus.Debugf("No exact match found for name '%s', including synonyms", args.Query)
	args.Rank = ""
	args.Status = ""
	resp, err = client.SearchSpecies(context.Background(), args)
	if err != nil {
		return nil, err
	}
	taxon, found = resp.GetExactMatch(args.Query)
	if found {
		return taxon, nil
	}
	return nil, nil
}

func FetchNames(client *GBIFClient, names []string, options ...searchParamOption) (NamesSearchResult, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 1)
	results := newMultipleNamesResult()

	for _, name := range names {
		if len(errChan) > 0 {
			break
		}
		wg.Add(1)
		go func(name string) {
			logrus.Debugf("Fetching taxon with name '%s' from GBIF", name)
			defer wg.Done()
			params := SearchParams{Query: name}
			params.Apply(options...)
			taxon, err := FetchByName(client, params)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if taxon == nil {
				results.NotFound = append(results.NotFound, name)
				return
			}
			logrus.Debugf("Fetched taxon '%s' with key %d from GBIF", taxon.Name, taxon.Key)
			results.Add(*taxon)
		}(name)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return results, err
	default:
		return results, nil
	}
}

func FetchKey(key int32) (*TaxonGBIF, error) {
	c := Client()
	return c.GetTaxonByKey(context.Background(), key)
}

func FetchKeys(keys []int32) ([]TaxonGBIF, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]TaxonGBIF, 0, len(keys))
	errChan := make(chan error, 1)

	for _, key := range keys {
		if len(errChan) > 0 {
			break
		}
		wg.Add(1)
		go func(key int32) {
			defer wg.Done()
			resp, err := FetchKey(key)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			mu.Lock()
			results = append(results, *resp)
			mu.Unlock()
		}(key)
	}

	wg.Wait()
	select {
	case err := <-errChan:
		return nil, err
	default:
		return results, nil
	}
}
