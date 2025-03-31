package services

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"
	"github.com/lsdch/biome/services/geoapify"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterGeoapifyRoutes(r router.Router) {
	group := r.RouteGroup("/geoapify").
		WithTags([]string{"Services"})

	router.Register(group, "ListGeoapifyUsage",
		huma.Operation{Path: "/usage",
			Method:  http.MethodGet,
			Summary: "List Geoapify usage",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](geoapify.GeoapifyUsageList),
	)

	router.Register(group, "GetGeoapifyStatus",
		huma.Operation{Path: "/status",
			Method:  http.MethodGet,
			Summary: "Get Geoapify API status",
		},
		controllers.FetchHandler[*struct {
			resolvers.AuthResolver
		}](geoapify.GetStatus),
	)

	router.Register(group, "ReverseGeocode",
		huma.Operation{Path: "/reverse-geocode",
			Method:  http.MethodPost,
			Summary: "Reverse geocode coordinates using Geoapify API",
		},
		ReverseGeocode,
	)

}

type ReverseGeocodeInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	Body struct {
		_ struct{} `json:"-" additionalProperties:"true"`
		*occurrence.LatLongCoords
	}
}

type ReverseGeocodeOutput struct {
	Body *geoapify.GeoapifyResult
}

func ReverseGeocode(ctx context.Context, input *ReverseGeocodeInput) (*ReverseGeocodeOutput, error) {
	client, err := geoapify.NewClient()
	if err != nil {
		return nil, huma.Error403Forbidden("Geoapify client was not configured")
	}
	res, err := client.ReverseGeocode(input.DB(), input.Body.ToGeoapify())
	if err != nil {
		return nil, err
	}

	return &ReverseGeocodeOutput{res}, nil
}
