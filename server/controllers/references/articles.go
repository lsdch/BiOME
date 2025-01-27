package references

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services/crossref"

	"github.com/danielgtaylor/huma/v2"
)

type SearchDoiInput struct {
	resolvers.AuthResolver
	DOI string `query:"doi" required:"true"`
}

func (i SearchDoiInput) Identifier() string {
	return i.DOI
}

type BibSearchInput struct {
	Body string
}
type BibSearchOutput struct {
	Body *crossref.BibSearchResults
}

func BibSearch(ctx context.Context, input *BibSearchInput) (*BibSearchOutput, error) {
	res, err := crossref.BibliographicSearch(input.Body)
	if err != nil {
		return nil, err
	}
	return &BibSearchOutput{Body: res}, nil
}

func RegisterRoutes(r router.Router) {

	huma.Register(r.API,
		huma.Operation{
			OperationID: "CrossRef",
			Path:        "/crossref",
			Method:      http.MethodGet,
			Summary:     "Retrieve article infos from DOI",
			Tags:        []string{"References"},
		},
		controllers.GetHandler[*SearchDoiInput](crossref.RetrieveDOI),
	)

	huma.Register(r.API,
		huma.Operation{
			OperationID: "CrossRefBibSearch",
			Path:        "/crossref",
			Method:      http.MethodPost,
			Summary:     "Retrieve article infos from query string",
			Tags:        []string{"References"},
		},
		BibSearch,
	)

	referencesAPI := r.RouteGroup("/references").
		WithTags([]string{"References"})

	router.Register(referencesAPI, "ListArticles",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List articles",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](references.ListArticles),
	)

	router.Register(referencesAPI, "CreateArticle",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create article",
		},
		controllers.CreateHandler[references.ArticleInput],
	)

	router.Register(referencesAPI, "UpdateArticle",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update article",
		},
		controllers.UpdateByCodeHandler[references.ArticleUpdate],
	)

	router.Register(referencesAPI, "DeleteArticle",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete article",
		},
		controllers.DeleteByCodeHandler(references.DeleteArticle),
	)
}
