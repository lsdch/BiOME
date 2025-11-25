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

type GetTaxonomyInput struct {
	resolvers.AuthResolver
	taxonomy.TaxonomyQuery
}
type GetTaxonomyOutput struct {
	Body *taxonomy.Taxonomy
}

func GetTaxonomy(ctx context.Context, input *GetTaxonomyInput) (*GetTaxonomyOutput, error) {
	var taxonomy, err = taxonomy.GetTaxonomy(input.DB(), input.TaxonomyQuery)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to fetch taxonomy", err)
	}
	return &GetTaxonomyOutput{Body: taxonomy}, nil
}

type CreateTaxonInput struct {
	resolvers.AccessRestricted[resolvers.Admin]
	Body taxonomy.TaxonInput
}

func (i CreateTaxonInput) Item() taxonomy.TaxonInput {
	return i.Body
}
