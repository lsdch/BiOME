package taxonomy

import (
	"context"
	"fmt"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/db"
	"github.com/lsdch/biome/models/taxonomy"
	_ "github.com/lsdch/biome/models/validations"
	"github.com/lsdch/biome/resolvers"

	"github.com/danielgtaylor/huma/v2"
)

type ListTaxaInput struct {
	taxonomy.ListFilters
}
type ListTaxaOutput struct{ Body []taxonomy.TaxonWithParentRef }

func ListTaxa(ctx context.Context, input *ListTaxaInput) (*ListTaxaOutput, error) {
	taxa, err := taxonomy.ListTaxa(db.Client(), input.ListFilters)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve list of taxa", err)
	}
	return &ListTaxaOutput{Body: taxa}, nil
}

type GetTaxonInput struct{ controllers.NameInput }
type GetTaxonOutput struct{ Body taxonomy.TaxonWithLineage }

func GetTaxon(ctx context.Context, input *GetTaxonInput) (*GetTaxonOutput, error) {
	taxon, err := taxonomy.FindByName(db.Client(), input.Name)
	if db.IsNoData(err) {
		return nil, huma.Error404NotFound(
			fmt.Sprintf("Taxon %s does not exist", input.Name),
		)
	}
	return &GetTaxonOutput{Body: taxon}, err
}

type GetTaxonomyAtRankInput struct {
	resolvers.AuthResolver
	Rank taxonomy.TaxonRank `json:"rank,omitempty" path:"rank"`
}
type GetTaxonomyAtRankOutput struct {
	Body []taxonomy.TaxonomyItem
}

func GetTaxonomyAtRank(ctx context.Context, input *GetTaxonomyAtRankInput) (*GetTaxonomyAtRankOutput, error) {
	var taxonomy, err = taxonomy.GetTaxonomyByRank(input.DB(), input.Rank)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch taxonomy", err)
	}
	return &GetTaxonomyAtRankOutput{Body: taxonomy}, nil
}

type CreateTaxonInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	Body taxonomy.TaxonInput
}

func (i CreateTaxonInput) Item() taxonomy.TaxonInput {
	return i.Body
}
