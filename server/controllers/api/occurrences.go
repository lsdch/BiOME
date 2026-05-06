package api

import (
	"context"
	"net/http"

	"github.com/lsdch/biome/controllers"
	"github.com/lsdch/biome/models/occurrence"
	"github.com/lsdch/biome/resolvers"
	"github.com/lsdch/biome/router"

	"github.com/danielgtaylor/huma/v2"
)

type RegisterOccurrenceOutput struct {
	Body occurrence.BaseOccurrence[occurrence.SamplingOutline]
}

func init() {
	occurAPI := func(r *router.Router) router.Group {
		return r.RouteGroup("/occurrences").
			WithTags([]string{"Occurrences"})
	}

	router.RegisterSpec(
		occurAPI,
		"ListOccurrences",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodGet,
			Summary: "List occurrences",
		},
		controllers.ListHandlerWithOpts[*struct {
			resolvers.AuthResolver
			occurrence.ListOccurrencesOptions
		}](occurrence.ListOccurrences),
	)

	router.RegisterSpec(
		occurAPI,
		"GetOccurrence",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodGet,
			Summary: "Get occurrence",
		},
		controllers.GetByCodeHandler(occurrence.GetOccurrence))

	router.RegisterSpec(
		occurAPI,
		"CreateOccurrence",
		huma.Operation{
			Path:    "/",
			Method:  http.MethodPost,
			Summary: "Create occurrence",
		},
		CreateOccurrence)

	router.RegisterSpec(
		occurAPI,
		"UpdateOccurrence",
		huma.Operation{
			Path:    "/{code}",
			Method:  http.MethodPatch,
			Summary: "Update occurrence",
		},
		controllers.UpdateByCodeHandler[occurrence.OccurrenceUpdate])

	router.RegisterSpec(
		occurAPI,
		"DeleteOccurrence",
		huma.Operation{
			Path:        "/{code}",
			Method:      http.MethodDelete,
			Summary:     "Delete occurrence",
			Description: "Delete an occurrence record by its code",
		},
		controllers.DeleteByCodeHandler(occurrence.DeleteOccurrence))

	router.RegisterSpec(
		occurAPI,
		"OccurrenceOverview",
		huma.Operation{
			Path:    "/overview",
			Method:  http.MethodGet,
			Summary: "Occurrences overview",
		},
		controllers.ListHandler[*struct {
			resolvers.AuthResolver
		}](occurrence.OccurrenceOverview),
	)

	router.RegisterSpec(
		occurAPI,
		"OccurrencesBySite",
		huma.Operation{
			Path:    "/by-site",
			Method:  http.MethodPost,
			Summary: "Occurrences by site",
		},
		controllers.ListHandlerWithOpts[*OccurrencesBySiteInput](occurrence.OccurrencesBySite),
	)
}

type OccurrencesBySiteInput struct {
	resolvers.AuthResolver
	Body occurrence.OccurrencesBySiteOptions `json:"options"`
}

func (input *OccurrencesBySiteInput) Options() occurrence.OccurrencesBySiteOptions {
	return input.Body.Options()
}

type CreateOccurrenceInput struct {
	resolvers.AccessRestricted[resolvers.Contributor]
	Body struct {
		Site        occurrence.SiteInput       `json:"site"`
		Sampling    occurrence.SamplingInput   `json:"sampling"`
		Biomaterial occurrence.OccurrenceInput `json:"bio_material"`
	}
}

func CreateOccurrence(ctx context.Context, input *CreateOccurrenceInput) (*RegisterOccurrenceOutput, error) {
	site, err := input.Body.Site.Save(input.DB())
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	sampling, err := input.Body.Sampling.Save(input.DB(), site.Code)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	bioMaterial, err := input.Body.Biomaterial.Save(input.DB(), sampling.Number)
	if err != nil {
		return nil, controllers.StatusError(err)
	}
	return &RegisterOccurrenceOutput{bioMaterial}, nil
}
