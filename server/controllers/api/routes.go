package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lsdch/biome/controllers"
	_ "github.com/lsdch/biome/controllers/api/services"
	_ "github.com/lsdch/biome/controllers/api/settings_api"
	_ "github.com/lsdch/biome/controllers/api/taxonomy"
	_ "github.com/lsdch/biome/controllers/api/users"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
)

func init() {
	baseAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/")
	}
	router.RegisterSpec(
		baseAPI,
		"OccurrencesDateRange",
		huma.Operation{
			Path:    "occurrences-date-range",
			Method:  http.MethodGet,
			Summary: "Get the min and max year for occurrence sampling dates",
		},
		func(ctx context.Context, i *OccurrenceDateRangeInput) (*controllers.BodyTransporter[occurrence.YearRange], error) {
			dateRange, err := occurrence.GetOccurrencesDateRange(i.DB())
			if err != nil {
				return &controllers.BodyTransporter[occurrence.YearRange]{}, controllers.StatusError(err)
			}
			return &controllers.BodyTransporter[occurrence.YearRange]{Body: dateRange}, nil
		},
	)

}

type OccurrenceDateRangeInput struct {
	resolvers.AuthResolver
}
