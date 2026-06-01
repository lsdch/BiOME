package gbif

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/geldata/gel-go/geltypes"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lsdch/biome/db/biomedb"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/taxonomy"

	"github.com/sirupsen/logrus"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TaxonGBIF struct {
	Key                     int32                        `json:"key" gel:"GBIF_ID"`
	Parent                  string                       `json:"parent"`
	ParentKey               int32                        `json:"parentKey" gel:"parentID"`
	Name                    string                       `json:"canonicalName" gel:"name"`
	ScientificName          string                       `json:"scientificName" gel:"scientific_name"`
	Status                  string                       `json:"taxonomicStatus" gel:"status"`
	Rank                    string                       `json:"rank" gel:"rank"`
	NameType                string                       `json:"nameType" gel:"name_type"`
	KingdomKey              int32                        `json:"kingdomKey"`
	PhylumKey               int32                        `json:"phylumKey"`
	ClassKey                int32                        `json:"classKey"`
	OrderKey                int32                        `json:"orderKey"`
	FamilyKey               int32                        `json:"familyKey"`
	GenusKey                int32                        `json:"genusKey"`
	SpeciesKey              int32                        `json:"speciesKey"`
	HigherClassificationMap map[int32]string             `json:"higherClassificationMap"`
	Authorship              models.OptionalInput[string] `json:"authorship,omitempty" gel:"authorship,omitempty"`
	NumDescendants          int32                        `json:"numDescendants" gel:"-"`
	AcceptedKey             int32                        `json:"acceptedKey"`
	AcceptedName            string                       `json:"accepted"`
}

func (taxon *TaxonGBIF) GetParentKeysList() []int32 {
	parents := []int32{}

	addKeys := func(keys ...int32) {
		for _, key := range keys {
			if key == 0 || key == taxon.Key {
				return
			}
			parents = append(parents, key)
		}
	}

	addKeys(
		taxon.KingdomKey,
		taxon.PhylumKey,
		taxon.ClassKey,
		taxon.OrderKey,
		taxon.FamilyKey,
		taxon.GenusKey,
		taxon.SpeciesKey,
	)

	return parents
}

func (taxon *TaxonGBIF) GetRank() taxonomy.TaxonRank {
	return taxonomy.TaxonRank(cases.Title(language.English).String(taxon.Rank))
}

func (taxon *TaxonGBIF) GetStatus() taxonomy.TaxonStatus {
	switch taxon.Status {
	case "Accepted":
		return taxonomy.Accepted
	case "Synonym":
		return taxonomy.Synonym
	case "Doubtful":
		return taxonomy.Doubtful
	default:
		return taxonomy.Unclassified
	}
}

func (taxon *TaxonGBIF) ToStaging() biomedb.InsertGBIFBatchParams {
	return biomedb.InsertGBIFBatchParams{
		Key:              taxon.Key,
		Parent:           taxon.Parent,
		ParentKey:        taxon.ParentKey,
		CanonicalName:    taxon.Name,
		ScientificName:   taxon.ScientificName,
		Status:           string(taxon.GetStatus()),
		Rank:             string(taxon.GetRank()),
		NameType:         taxon.NameType,
		KingdomKey:       taxon.KingdomKey,
		PhylumKey:        taxon.PhylumKey,
		ClassKey:         taxon.ClassKey,
		OrderKey:         taxon.OrderKey,
		FamilyKey:        taxon.FamilyKey,
		GenusKey:         taxon.GenusKey,
		SpeciesKey:       taxon.SpeciesKey,
		HigherTaxonKeys:  slices.Collect(maps.Keys(taxon.HigherClassificationMap)),
		HigherTaxonNames: slices.Collect(maps.Values(taxon.HigherClassificationMap)),
		Authorship:       taxon.Authorship.GetWithDefault(""),
		NumDescendants:   taxon.NumDescendants,
		AcceptedKey:      taxon.AcceptedKey,
		AcceptedName:     taxon.AcceptedName,
	}
}

func (taxon *TaxonGBIF) ToCandidate(inputName string) biomedb.InsertTaxonCandidatesBatchParams {
	var priority int32
	if taxon.Name == inputName || taxon.ScientificName == inputName {
		if taxon.GetStatus() == taxonomy.Accepted {
			// Accepted exact matches
			priority = 100
		} else {
			// Synonyms and other non-accepted exact matches
			priority = 80
		}
	} else {
		// Non-exact matches
		priority = 50
	}

	return biomedb.InsertTaxonCandidatesBatchParams{
		InputName: inputName,
		Name:      taxon.Name,
		Authorship: pgtype.Text{
			String: taxon.Authorship.Value,
			Valid:  taxon.Authorship.IsSet,
		},
		Rank:      biomedb.TaxonRank(taxon.GetRank()),
		Status:    biomedb.TaxonStatus(taxon.GetStatus()),
		Source:    biomedb.TaxonMatchSourceGbif,
		MatchType: biomedb.TaxonMatchTypeExact,
		Priority:  priority,
		GbifID: pgtype.Int4{
			Int32: taxon.Key,
			Valid: true,
		},
	}
}

func (taxon *TaxonGBIF) Normalize() *TaxonGBIF {
	if authorship, isSet := taxon.Authorship.Get(); isSet && authorship == "" {
		taxon.Authorship.Clear()
	}
	taxon.Status = string(taxon.GetStatus())
	taxon.Rank = string(taxon.GetRank())
	return taxon
}

func UpsertTaxa(tx geltypes.Tx, taxa []TaxonGBIF) (n int, err error) {

	ctx := context.Background()

	for _, taxon := range taxa {
		taxon.Normalize()
		logrus.Debugf("Inserting taxon from GBIF %+v", &taxon)
		args, _ := json.Marshal(&taxon)
		if err = tx.Execute(ctx,
			`#edgeql
				with module taxonomy,
					data := <json>$0,
					accepted := (
						select Taxon
						filter .GBIF_ID = <int32>data['acceptedKey']
					),
					synonym_group := (select accepted.synonym_group) ?? (
						update accepted set {
							synonym_group := (insert SynonymGroup {})
						}
					).synonym_group
				insert Taxon {
					name := <str>data['canonicalName'],
					GBIF_ID := <int32>data['key'],
					status := <TaxonStatus>data['taxonomicStatus'],
					parent := (
						select detached Taxon filter .name = <str>data['parent']
					),
					rank := <Rank>data['rank'],
					authorship := <str>data['authorship'],
					synonym_group := synonym_group,
				}
				unless conflict on .GBIF_ID
			`, args); err != nil {
			return n, fmt.Errorf("Error inserting taxon %s[%d] %v", taxon.Name, taxon.Key, err)
		} else {
			n++
		}
	}
	return
}
