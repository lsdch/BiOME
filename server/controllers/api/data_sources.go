package api

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/references"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func init() {
	dataSourcesAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/data-sources").
			WithTags([]string{"Data sources"})
	}

	router.RegisterSpec(
		dataSourcesAPI,
		"ListDataSources",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List external data sources",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](references.ListDataSources),
	)

	router.RegisterSpec(
		dataSourcesAPI,
		"CreateDataSource",
		huma.Operation{Path: "/",
			Method:  http.MethodPost,
			Summary: "Register external data source",
		},
		controllers.CreateHandler[references.DataSourceInput],
	)

	router.RegisterSpec(
		dataSourcesAPI,
		"UpdateDataSource",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update external data source",
		},
		controllers.UpdateByCodeHandler[references.DataSourceUpdate],
	)

	router.RegisterSpec(
		dataSourcesAPI,
		"DeleteDataSource",
		huma.Operation{Path: "/{code}",
			Method:  http.MethodDelete,
			Summary: "Delete external data source",
		},
		controllers.DeleteByCodeHandler(references.DeleteDataSources),
	)
}
