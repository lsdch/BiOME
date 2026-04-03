package taxonomy

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/geldata/gel-go/geltypes"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/db/queries"
	"github.com/lsdch/biome/models"
	"github.com/lsdch/biome/models/people"
)

// TaxonCode generates a code from a taxon name, replacing spaces with underscores.
func TaxonCode(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

func TaxonShortCode(name string) string {
	fragments := strings.Split(name, "")
	if len(fragments) > 1 {
		return string(fragments[0][0]) + "_" + strings.Join(fragments[1:], "_")
	}
	return name
}

type TaxonInner struct {
	Name           string      `gel:"name" json:"name" example:"Asellus aquaticus" binding:"required,alpha"`
	ScientificName string      `gel:"scientific_name" json:"scientific_name,omitempty" example:"Asellus aquaticus (Linnaeus, 1758)"`
	Status         TaxonStatus `gel:"status" json:"status" example:"Accepted" binding:"required"`
	Rank           TaxonRank   `gel:"rank" json:"rank" example:"Species" binding:"required"`
}

type TaxonInput struct {
	TaxonInner `gel:"$inline"`
	SynonymOf  models.OptionalInput[string] `gel:"synonym_of" json:"synonym_of,omitempty"`                // name of the taxon this taxon is a synonym of
	Parent     string                       `json:"parent" binding:"required,exist=taxonomy::Taxon.name"` // Parent taxon name
	Authorship models.OptionalInput[string] `json:"authorship,omitempty" example:"(Linnaeus, 1758)"`
	Comment    models.OptionalInput[string] `json:"comment,omitempty" gel:"comment"`
}

type Taxon struct {
	ID            geltypes.UUID          `gel:"id" json:"id" format:"uuid" binding:"required"`
	GBIF_ID       geltypes.OptionalInt32 `gel:"GBIF_ID" json:"GBIF_ID,omitempty" example:"2206247" binding:"numeric"`
	TaxonInner    `gel:"$inline"`
	Authorship    geltypes.OptionalStr `gel:"authorship" json:"authorship,omitempty" example:"(Linnaeus, 1758)"`
	Anchor        bool                 `gel:"anchor" json:"anchor"`
	ChildrenCount int64                `gel:"children_count" json:"children_count"`
	Comment       geltypes.OptionalStr `json:"comment,omitempty" gel:"comment"`
	Meta          people.Meta          `gel:"meta" json:"meta" binding:"required"`
}

type TaxonWithParentRef struct {
	Taxon  `gel:"$inline"`
	Parent geltypes.OptionalStr `gel:"parent_name" json:"parent"`
}

type Lineage struct {
	Kingdom    models.Optional[Taxon] `gel:"kingdom" json:"kingdom,omitempty"`
	Phylum     models.Optional[Taxon] `gel:"phylum" json:"phylum,omitempty"`
	Class      models.Optional[Taxon] `gel:"class" json:"class,omitempty"`
	Order      models.Optional[Taxon] `gel:"order" json:"order,omitempty"`
	Family     models.Optional[Taxon] `gel:"family" json:"family,omitempty"`
	Genus      models.Optional[Taxon] `gel:"genus" json:"genus,omitempty"`
	Species    models.Optional[Taxon] `gel:"species" json:"species,omitempty"`
	Subspecies models.Optional[Taxon] `gel:"subspecies" json:"subspecies,omitempty"`
}

type TaxonWithRelatives struct {
	Taxon    `gel:"$inline"`
	Parent   models.Optional[Taxon] `gel:"parent" json:"parent,omitempty"`
	Synonyms []Taxon                `gel:"synonyms" json:"synonyms,omitempty"`
	Children []Taxon                `gel:"children" json:"children,omitempty"`
}

type TaxonWithLineage struct {
	TaxonWithRelatives `gel:"$inline"`
	Lineage            Lineage `gel:"$inline" json:"lineage"`
}

// TaxonomyItem type is a tree like representation of the taxonomy, or part of it.
type TaxonomyItem struct {
	Taxon    `gel:"$inline"`
	Parent   models.Optional[Taxon] `gel:"parent" json:"parent,omitempty"`
	Synonyms []Taxon                `gel:"synonyms" json:"synonyms,omitempty"`
}

// GetTaxonomyByRank returns a flat list of taxa at a specific rank with parent linkage.
func GetTaxonomyByRank(db geltypes.Executor, rank TaxonRank) ([]TaxonomyItem, error) {
	var taxa []TaxonomyItem
	query := `#edgeql
		with module taxonomy,
		select Taxon { *, meta: {*}, parent: { id, name, rank } }
		filter .rank = <Rank>$0
		order by .name;
	`
	if err := db.Query(context.Background(), query, &taxa, rank); err != nil {
		return nil, err
	}
	return taxa, nil
}

type ListFilters struct {
	Pattern     string      `json:"pattern,omitempty" query:"pattern"`
	Ranks       []TaxonRank `query:"ranks"`
	Status      TaxonStatus `json:"status,omitempty" query:"status"`
	Parent      string      `json:"parent,omitempty" query:"parent"`
	Limit       int64       `json:"limit,omitempty" query:"limit"`
	SampledOnly bool        `json:"sampled_only,omitempty" query:"sampled_only"`
	SynonymOf   string      `json:"synonym_of,omitempty" query:"synonym_of"`
}

func ListTaxa(db geltypes.Executor, filters ListFilters) ([]TaxonWithParentRef, error) {
	var taxa = []TaxonWithParentRef{}
	query := queries.RenderTemplate("list_taxa.tmpl.edgeql", filters)
	err := db.Query(context.Background(), query, &taxa,
		filters.Pattern, filters.Ranks, filters.Status, filters.Parent, filters.SampledOnly, filters.SynonymOf, filters.Limit)
	return taxa, err
}

const TAXON_LINEAGE_SHAPE = `#edgeql
		{ *,
			meta: { * },
			parent : { *, meta: { * } },
			children : { *, meta: { * } },
			kingdom: { *, meta: { * } },
			phylum: { *, meta: { * } },
			class: { *, meta: { * } },
			order: { *, meta: { * } },
			family: { *, meta: { * } },
			genus: { *, meta: { * } },
			species: { *, meta: { * } },
		}`

func FindByID(db geltypes.Executor, id geltypes.UUID) (taxon TaxonWithLineage, err error) {
	query := fmt.Sprintf(
		`select taxonomy::Taxon %s filter .id = <uuid>$0;`,
		TAXON_LINEAGE_SHAPE)
	err = db.QuerySingle(context.Background(), query, &taxon, id)
	return taxon, err
}

func FindByName(db geltypes.Executor, name string) (taxon TaxonWithLineage, err error) {
	query := fmt.Sprintf(
		`select taxonomy::Taxon %s filter .name = <str>$0;`,
		TAXON_LINEAGE_SHAPE)
	err = db.QuerySingle(context.Background(), query, &taxon, name)
	return taxon, err
}

func Delete(db geltypes.Executor, name string) (taxon TaxonWithRelatives, err error) {
	err = db.QuerySingle(context.Background(), `#edgeql
		select (
			delete taxonomy::Taxon filter .name = <str>$0
		) { *, meta: { * }, parent : { * , meta: { * }}, children : { * , meta: { * }} }
	`, &taxon, name)
	return
}

func (taxon TaxonInput) Save(db geltypes.Executor) (created TaxonWithRelatives, err error) {
	args, _ := json.Marshal(taxon)
	err = db.QuerySingle(context.Background(),
		`#edgeql
		with module taxonomy,
		data := <json>$0
		select (
			insert Taxon {
				name := <str>data['name'],
				status := <TaxonStatus>data['status'],
				parent := (
					select detached Taxon filter .name = <str>data['parent']
				),
				synonym_group := (
					select detached Taxon 
					filter .name = <str>json_get(data, 'synonym_of') 
					and .name != <str>json_get(data, 'name')
				).synonym_group,
				rank := <Rank>data['rank'],
				authorship := <str>json_get(data, 'authorship')
			}
		) { *, meta: { * }, parent : { * , meta: { * }}, children : { * , meta: { * }} };
	`, &created, args)
	return created, err
}

type TaxonUpdate struct {
	Name       models.OptionalInput[string]      `gel:"name" json:"name,omitempty"`
	Status     models.OptionalInput[TaxonStatus] `gel:"status" json:"status,omitempty"`
	Authorship models.OptionalNull[string]       `gel:"authorship" json:"authorship,omitempty"`
	Rank       models.OptionalInput[TaxonRank]   `gel:"rank" json:"rank,omitempty"`
	Parent     models.OptionalInput[string]      `gel:"parent" json:"parent,omitempty"` // parent name
	Comment    models.OptionalNull[string]       `json:"comment,omitempty" gel:"comment"`
}

func (u TaxonUpdate) Save(e geltypes.Executor, name string) (updated Taxon, err error) {
	data, _ := json.Marshal(u)
	query := db.UpdateQuery{
		Frame: fmt.Sprintf(`#edgeql
			with item := <json>$1,
			select (update taxonomy::Taxon filter .name = <str>$0 set {
				%s
			}) %s`,
			"%s", TAXON_LINEAGE_SHAPE),
		Mappings: map[string]string{
			"name":       "<str>item['name']",
			"status":     "<taxonomy::TaxonStatus>item['status']",
			"rank":       "<taxonomy::Rank>item['rank']",
			"authorship": "<str>item['authorship']",
			"parent": `#edgeql
				(
					select detached taxonomy::Taxon filter .name = <str>item['parent']
				)`,
		},
	}
	err = e.QuerySingle(context.Background(), query.Query(u), &updated, name, data)
	return
}

func CheckMissingTaxa(db geltypes.Executor, taxa []string) ([]string, error) {
	var missingTaxa []string
	err := db.Query(context.Background(),
		`#edgeql
			with module taxonomy,
			existing := (select Taxon.name union Taxon.scientific_name),
			select distinct (array_unpack(<array<str>>$0) except existing)
		`,
		&missingTaxa, taxa)
	return missingTaxa, err
}
