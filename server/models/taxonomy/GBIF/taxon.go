package gbif

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/taxonomy"

	"github.com/sirupsen/logrus"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TaxonInnerGBIF struct {
	Key            int32  `json:"key" gel:"GBIF_ID"`
	Parent         string `json:"parent"`
	ParentKey      int32  `json:"parentKey" gel:"parentID"`
	Name           string `json:"canonicalName" gel:"name"`
	ScientificName string `json:"scientificName" gel:"scientific_name"`
	Status         string `json:"taxonomicStatus" gel:"status"`
	Rank           string `json:"rank" gel:"rank"`
	NameType       string `json:"nameType" gel:"name_type"`
}

func (taxon *TaxonInnerGBIF) GetRank() taxonomy.TaxonRank {
	return taxonomy.TaxonRank(cases.Title(language.English).String(taxon.Rank))
}

func (t *TaxonInnerGBIF) normalize() *TaxonInnerGBIF {
	switch t.Status {
	case "ACCEPTED":
		t.Status = "Accepted"
	case "SYNONYM", "HETEROTYPIC_SYNONYM", "HOMOTYPIC_SYNONYM", "PROPARTE_SYNONYM":
		t.Status = "Synonym"
	case "DOUBTFUL":
		t.Status = "Doubtful"
	default:
		t.Status = "Unclassified"
	}
	t.Rank = string(t.GetRank())
	return t
}

type TaxonGBIF struct {
	TaxonInnerGBIF          `json:",inline" gel:"$inline"`
	KingdomKey              int32                        `json:"kingdomKey"`
	PhylumKey               int32                        `json:"phylumKey"`
	ClassKey                int32                        `json:"classKey"`
	OrderKey                int32                        `json:"orderKey"`
	FamilyKey               int32                        `json:"familyKey"`
	GenusKey                int32                        `json:"genusKey"`
	SpeciesKey              int32                        `json:"speciesKey"`
	HigherClassificationMap map[int32]string             `json:"higherClassificationMap"`
	Authorship              models.OptionalInput[string] `json:"authorship,omitempty" gel:"authorship,omitempty"`
	NumDescendants          int                          `json:"numDescendants" gel:"-"`
	Anchor                  bool                         `json:"anchor" gel:"anchor"`
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

func (taxon *TaxonGBIF) normalize() *TaxonGBIF {
	if authorship, isSet := taxon.Authorship.Get(); isSet && authorship == "" {
		taxon.Authorship.Clear()
	}
	taxon.TaxonInnerGBIF.normalize()
	return taxon
}

func UpsertTaxa(tx geltypes.Tx, taxa []TaxonGBIF) (n int, err error) {

	ctx := context.Background()

	for _, taxon := range taxa {
		taxon.normalize()
		logrus.Debugf("Inserting taxon from GBIF %+v", &taxon)
		args, _ := json.Marshal(&taxon)
		if err = tx.Execute(ctx,
			`#edgeql
				with module taxonomy,
					data := <json>$0,
					anchor := <bool>data['anchor'],
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
					anchor := anchor,
					synonym_group := synonym_group,
				}
				unless conflict on .GBIF_ID else (
					update Taxon set {
						anchor := anchor if not .anchor else .anchor
					}
				)
			`, args); err != nil {
			return n, fmt.Errorf("Error inserting taxon %s[%d] %v", taxon.Name, taxon.Key, err)
		} else {
			n++
		}
	}
	return
}
