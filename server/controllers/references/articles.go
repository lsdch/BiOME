package references

import (
	"darco/proto/controllers"
	"darco/proto/models/references"
	"darco/proto/resolvers"
	"darco/proto/router"
	"darco/proto/services/crossref"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type SearchDoiInput struct {
	resolvers.AuthRequired
	DOI string `query:"doi" required:"true"`
}

func (i SearchDoiInput) Identifier() string {
	return i.DOI
}

func RegisterRoutes(r router.Router) {

	huma.Register(r.API,
		huma.Operation{
			OperationID: "crossref",
			Path:        "/crossref",
			Method:      http.MethodGet,
			Summary:     "Retrieve article infos from DOI",
			Tags:        []string{"References"},
		},
		controllers.GetHandler[*SearchDoiInput](crossref.RetrieveDOI),
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
