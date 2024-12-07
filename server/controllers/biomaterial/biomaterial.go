package biomaterial

import (
	"darco/proto/controllers"
	"darco/proto/models/occurrence"
	"darco/proto/resolvers"
	"darco/proto/router"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(r router.Router) {
	biomat_API := r.RouteGroup("/bio-material").WithTags([]string{"Samples"})

	router.Register(biomat_API, "ListBioMaterial",
		huma.Operation{
			Path:        "/",
			Method:      http.MethodGet,
			Summary:     "List bio-material",
			Description: "Both internal and external",
		}, controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.ListBioMaterials))

	router.Register(biomat_API, "GetBioMaterial",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodGet,
			Summary:     "Get bio-material",
			Description: "Both internal and external",
		}, controllers.GetByCodeHandler(occurrence.GetBioMaterial))

	router.Register(biomat_API, "DeleteBioMaterial",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodDelete,
			Summary:     "Delete bio-material",
			Description: "Delete any (internal/external) bio-material record by its code",
		}, controllers.DeleteByCodeHandler(occurrence.DeleteBioMaterial))
}
