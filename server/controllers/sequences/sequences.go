package sequences

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	RegisterGeneRoutes(r)

	seqAPI := r.RouteGroup("/sequences").
		WithTags([]string{"Sequences"})

	router.Register(seqAPI, "ListSequences",
		huma.Operation{Path: "/",
			Method:  http.MethodGet,
			Summary: "List sequences",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSequences),
	)

	router.Register(seqAPI, "GetSequence",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodGet,
			Summary: "Get sequence",
		},
		controllers.GetByCodeHandler(occurrence.GetSequence),
	)

	router.Register(seqAPI, "DeleteSequence",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete sequence",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteSequence),
	)
}
