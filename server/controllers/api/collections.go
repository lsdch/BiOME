package api

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
)

func init() {
	collectionsAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/collections").
			WithTags([]string{"References"})
	}

	router.RegisterSpec(
		collectionsAPI,
		"ListCollections",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List collections",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](references.ListCollections),
	)

	router.RegisterSpec(
		collectionsAPI,
		"CreateCollection",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create a new collection",
		},
		controllers.CreateHandler[references.CollectionInput],
	)

	router.RegisterSpec(
		collectionsAPI,
		"UpdateCollection",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update a collection by its code",
		},
		controllers.UpdateHandler[*struct {
			resolvers.AuthRequired
			controllers.CodeInput
			controllers.UpdateInput[references.CollectionUpdate, string, references.Collection]
		}],
	)

	router.RegisterSpec(
		collectionsAPI,
		"DeleteCollection",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete a collection by its code",
		},
		controllers.DeleteHandler[*struct {
			resolvers.AuthResolver
			controllers.CodeInput
		}](references.DeleteCollection),
	)
}
