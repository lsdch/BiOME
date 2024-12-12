package sequences

import (
	"darco/proto/controllers"
	"darco/proto/models/sequences"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGeneRoutes(r router.Router) {
	genesAPI := r.RouteGroup("/genes").
		WithTags([]string{"Sequences"})

	router.Register(genesAPI, "ListGenes",
		huma.Operation{Path: "/",
			Method:  http.MethodGet,
			Summary: "List genes",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](sequences.ListGenes),
	)

	router.Register(genesAPI, "CreateGene",
		huma.Operation{Path: "/",
			Method:  http.MethodPost,
			Summary: "Create gene",
		},
		controllers.CreateHandler[sequences.GeneInput],
	)

	router.Register(genesAPI, "UpdateGene",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update gene",
		},
		controllers.UpdateByCodeHandler[sequences.GeneUpdate],
	)

	router.Register(genesAPI, "DeleteGene",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete gene",
		},
		controllers.DeleteByCodeHandler(sequences.DeleteGene),
	)
}
