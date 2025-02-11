package references

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterDataSourceRoutes(r router.Router) {
	genesAPI := r.RouteGroup("/data-sources").
		WithTags([]string{"Data sources"})

	router.Register(genesAPI, "ListDataSources",
		huma.Operation{Path: "/",
			Method:  http.MethodGet,
			Summary: "List external data sources",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](references.ListDataSources),
	)

	router.Register(genesAPI, "CreateDataSource",
		huma.Operation{Path: "/",
			Method:  http.MethodPost,
			Summary: "Register external data source",
		},
		controllers.CreateHandler[references.DataSourceInput],
	)

	router.Register(genesAPI, "UpdateDataSource",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update external data source",
		},
		controllers.UpdateByCodeHandler[references.DataSourceUpdate],
	)

	router.Register(genesAPI, "DeleteDataSource",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete external data source",
		},
		controllers.DeleteByCodeHandler(references.DeleteDataSources),
	)
}
