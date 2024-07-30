package sites

import (
	"darco/proto/controllers"
	"darco/proto/models/location"
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
		}, controllers.ListHandler(location.ListSites))

	router.Register(sites_API, "GetSite",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get site",
			Description: "Get site infos using its code",
		}, controllers.GetByCodeHandler(location.GetSite))

	huma.Register(r.API, huma.Operation{
		Path:    "/access-points",
		Method:  http.MethodGet,
		Summary: "List access points",
		Tags:    sites_API.Tags,
	}, controllers.ListHandler(location.ListAccessPoints))
}
