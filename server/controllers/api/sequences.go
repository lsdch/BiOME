package api

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	seqAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/sequences").
			WithTags([]string{"Sequences"})
	}

	router.RegisterSpec(
		seqAPI,
		"ListSequences",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List sequences",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSequences),
	)

	router.RegisterSpec(
		seqAPI,
		"GetSequence",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodGet,
			Summary: "Get sequence",
		},
		controllers.GetByCodeHandler(occurrence.GetSequence))

	router.RegisterSpec(
		seqAPI,
		"DeleteSequence",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete sequence",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteSequence))
}
