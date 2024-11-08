package sites

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	sites_API := r.RouteGroup("/sites").WithTags([]string{"Location"})

	router.Register(sites_API, "ListSites",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List sites",
			Description: "List all registered sites",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListSites))

	router.Register(sites_API, "GetSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get site",
			Description: "Get site infos using its code",
		}, controllers.GetByCodeHandler(occurrence.GetSite))

	router.Register(sites_API, "UpdateSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodPatch,
			Summary:     "Update site",
			Description: "Update site infos using its code",
		},
		controllers.UpdateByCodeHandler[occurrence.SiteUpdate],
	)

	huma.Register(r.API, huma.Operation{
		Path:    "/access-points",
		Method:  http.MethodGet,
		Summary: "List access points",
		Tags:    sites_API.Tags,
	}, controllers.ListHandler[*struct {
		resolvers.AuthResolver
	}](occurrence.ListAccessPoints))
}
