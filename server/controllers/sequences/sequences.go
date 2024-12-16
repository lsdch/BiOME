package sequences

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	RegisterGeneRoutes(r)
	RegisterSeqDBRoutes(r)

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
