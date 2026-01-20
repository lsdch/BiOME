package api

import (
	"context"
	"net/http"

	"github.com/caltechlibrary/crossrefapi"
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

func init() {

	crossrefAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/crossref").
			WithTags([]string{"References"})
	}

	referencesAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/references").
			WithTags([]string{"References"})
	}

	router.RegisterSpec(
		crossrefAPI,
		"CrossRef",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "Retrieve article infos from DOI",
		},
		func(ctx context.Context, i *SearchDoiInput) (*controllers.BodyResponse[*crossrefapi.Works], error) {
			ref, err := crossref.RetrieveDOI(i.DOI)
			if err != nil {
				return nil, err
			}
			return &controllers.BodyResponse[*crossrefapi.Works]{Body: ref}, err
		},
	)

	router.RegisterSpec(
		crossrefAPI,
		"CrossRefBibSearch",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Retrieve article infos from query string",
		},
		BibSearch,
	)

	router.RegisterSpec(
		referencesAPI,
		"ListArticles",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List articles",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](references.ListArticles),
	)

	router.RegisterSpec(
		referencesAPI,
		"CreateArticle",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create article",
		},
		controllers.CreateHandler[references.ArticleInput],
	)

	router.RegisterSpec(
		referencesAPI,
		"UpdateArticle",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update article",
		},
		controllers.UpdateByCodeHandler[references.ArticleUpdate],
	)

	router.RegisterSpec(
		referencesAPI,
		"DeleteArticle",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete article",
		},
		controllers.DeleteByCodeHandler(references.DeleteArticle),
	)
}
