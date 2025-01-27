package sequences

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/sequences"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterSeqDBRoutes(r router.Router) {
	genesAPI := r.RouteGroup("/seq-databases").
		WithTags([]string{"Sequences"})

	router.Register(genesAPI, "ListSeqDBs",
		huma.Operation{Path: "/",
			Method:  http.MethodGet,
			Summary: "List external sequence databases",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](sequences.ListSeqDB),
	)

	router.Register(genesAPI, "CreateSeqDB",
		huma.Operation{Path: "/",
			Method:  http.MethodPost,
			Summary: "Create external sequence database",
		},
		controllers.CreateHandler[sequences.SeqDBInput],
	)

	router.Register(genesAPI, "UpdateSeqDB",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update external sequence database",
		},
		controllers.UpdateByCodeHandler[sequences.SeqDBUpdate],
	)

	router.Register(genesAPI, "DeleteSeqDB",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete external sequence database",
		},
		controllers.DeleteByCodeHandler(sequences.DeleteSeqDB),
	)
}
