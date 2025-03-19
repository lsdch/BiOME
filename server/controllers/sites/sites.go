package sites

import (
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

func SitesAPI(r router.Router) router.Group {
	return r.RouteGroup("/sites").WithTags([]string{"Location"})
}

func RegisterRoutes(r router.Router) {

	sites_API := SitesAPI(r)

	router.Register(sites_API, "ListSites",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List sites",
			Description: "List all registered sites",
		}, controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			occurrence.ListSitesOptions
		}](occurrence.ListSites))

	router.Register(sites_API, "GetSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get site",
			Description: "Get site infos using its code",
		}, controllers.GetByCodeHandler(occurrence.GetSite))

	router.Register(sites_API, "CreateSite",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodPost,
			Summary:     "Create site",
			Description: "Create site infos using its code",
		}, controllers.CreateHandler[occurrence.SiteInput])

	router.Register(sites_API, "UpdateSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodPatch,
			Summary:     "Update site",
			Description: "Update site infos using its code",
		},
		controllers.UpdateByCodeHandler[occurrence.SiteUpdate],
	)

	router.Register(sites_API, "ListSiteEvents",
		huma.Operation{
			Path:    "/{code}/events",
			Method:  http.MethodGet,
			Summary: "List site events",
		},
		controllers.GetByCodeHandler(occurrence.ListSiteEvents),
	)

	router.Register(sites_API, "CreateEvent",
		huma.Operation{
			Path:        "/{code}/events",
			Method:      http.MethodPost,
			Summary:     "Create event",
			Description: "Register event on a site identified by its code",
		},
		controllers.UpdateByCodeHandler[occurrence.EventInput],
	)
}
